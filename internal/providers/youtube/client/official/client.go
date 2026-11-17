package official

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/youtube/client"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/duration"
)

type Client struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
}

type ClientOption func(*Client)

func WithAPIURL(url string) ClientOption {
	return func(c *Client) {
		c.apiURL = url
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := Client{
		apiURL:     "https://www.googleapis.com",
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// https://developers.google.com/youtube/v3/docs/videos/list
func (c *Client) FetchVideo(ctx context.Context, id string) (client.Video, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/videos", url.Values{
		"part": {"snippet,contentDetails"},
		"id":   {id},
	})
	if err != nil {
		return client.Video{}, err
	}

	response := videoResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return client.Video{}, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return client.Video{}, client.ErrNotFound
	}

	item := response.Items[0]
	return client.Video{
		Metadata: c.snippetMetadata(item.ID, item.Snippet),
		Duration: duration.ISO8601ToSeconds(item.ContentDetails.Duration),
	}, nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *Client) SearchVideos(ctx context.Context, query string) ([]string, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/search", url.Values{
		"q":               {query},
		"part":            {"snippet"},
		"type":            {"video"},
		"videoCategoryId": {"10"},
		"maxResults":      {strconv.Itoa(client.SearchLimit)},
	})
	if err != nil {
		return nil, err
	}

	response := videoSearchResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	ids := make([]string, 0, len(response.Items))
	for _, item := range response.Items {
		ids = append(ids, item.ID.VideoID)
	}
	return ids, nil
}

// https://developers.google.com/youtube/v3/docs/playlists/list
func (c *Client) FetchPlaylist(ctx context.Context, id string) (client.Metadata, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/playlists", url.Values{
		"part": {"snippet"},
		"id":   {id},
	})
	if err != nil {
		return client.Metadata{}, err
	}

	response := playlistResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return client.Metadata{}, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return client.Metadata{}, client.ErrNotFound
	}

	item := response.Items[0]
	return c.snippetMetadata(item.ID, item.Snippet), nil
}

func (c *Client) FetchPlaylistWithItems(ctx context.Context, id string) (client.Metadata, []string, error) {
	playlist, err := c.FetchPlaylist(ctx, id)
	if err != nil {
		return client.Metadata{}, nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
	}
	items, err := c.fetchPlaylistItems(ctx, playlist.ID)
	if err != nil {
		return client.Metadata{}, nil, fmt.Errorf("failed to get playlist items from youtube: %w", err)
	}
	return playlist, items, nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *Client) SearchPlaylists(ctx context.Context, query string) ([]string, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/search", url.Values{
		"q":          {query},
		"part":       {"snippet"},
		"type":       {"playlist"},
		"maxResults": {strconv.Itoa(client.SearchLimit)},
	})
	if err != nil {
		return nil, err
	}

	response := playlistSearchResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	ids := make([]string, 0, len(response.Items))
	for _, item := range response.Items {
		ids = append(ids, item.ID.PlaylistID)
	}
	return ids, nil
}

// https://developers.google.com/youtube/v3/docs/playlistItems/list
func (c *Client) fetchPlaylistItems(ctx context.Context, id string) ([]string, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/playlistItems", url.Values{
		"part":       {"snippet"},
		"maxResults": {"50"},
		"playlistId": {id},
	})
	if err != nil {
		return nil, err
	}
	response := playlistItemsResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	ids := make([]string, 0, len(response.Items))
	for _, item := range response.Items {
		ids = append(ids, item.Snippet.ResourceID.VideoID)
	}
	return ids, nil
}

func (c *Client) getWithKey(ctx context.Context, path string, values url.Values) ([]byte, error) {
	values.Set("key", c.apiKey)
	u := fmt.Sprintf("%s%s?%s", c.apiURL, path, values.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, c.responseError(response)
	}

	return io.ReadAll(response.Body)
}

func (c *Client) responseError(response *http.Response) error {
	const quotaExceededReason = "quotaExceeded"

	result := fmt.Errorf("non ok http status: %d", response.StatusCode)
	body, err := io.ReadAll(response.Body)
	var errRes errorResponse
	if err != nil || json.Unmarshal(body, &errRes) != nil {
		return result
	}

	reasons := make([]string, 0, len(errRes.Error.Errors))
	for _, detail := range errRes.Error.Errors {
		if detail.Reason != "" {
			reasons = append(reasons, detail.Reason)
		}
	}
	if len(reasons) > 0 {
		result = fmt.Errorf("%w (reason: %s)", result, strings.Join(reasons, ", "))
	}
	if response.StatusCode == http.StatusForbidden && slices.Contains(reasons, quotaExceededReason) {
		return fmt.Errorf("%w: %w", result, release.ErrQuotaExceeded)
	}
	return result
}

func (c *Client) snippetMetadata(id string, s snippet) client.Metadata {
	creator := s.VideoOwnerChannelTitle
	if creator == "" {
		creator = s.ChannelTitle
	}
	return client.Metadata{
		ID:          id,
		Title:       s.Title,
		Creator:     creator,
		Description: s.Description,
		CoverURL:    c.coverURL(s.Thumbnails),
	}
}

func (c *Client) coverURL(t thumbnails) string {
	for _, image := range []thumbnail{t.Maxres, t.Standard, t.High, t.Medium, t.Default} {
		if image.URL != "" {
			return image.URL
		}
	}
	return ""
}
