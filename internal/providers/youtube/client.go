package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
}

type ClientOption func(client *Client)

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = url
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	const defaultAPIURL = "https://www.googleapis.com"

	c := Client{
		apiKey:     apiKey,
		apiURL:     defaultAPIURL,
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

var errNotFound = errors.New("not found")

// https://developers.google.com/youtube/v3/docs/videos/list
func (c *Client) fetchVideo(ctx context.Context, id string) (video, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/videos", url.Values{
		"part": {"snippet,contentDetails"},
		"id":   {id},
	})
	if err != nil {
		return video{}, err
	}

	response := videoResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return video{}, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return video{}, errNotFound
	}

	return response.Items[0], nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *Client) searchVideos(ctx context.Context, query string) ([]videoSearchResult, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/search", url.Values{
		"q":               {query},
		"part":            {"snippet"},
		"type":            {"video"},
		"videoCategoryId": {"10"},
		"maxResults":      {"10"},
	})
	if err != nil {
		return nil, err
	}

	response := videoSearchResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	return response.Items, nil
}

// https://developers.google.com/youtube/v3/docs/playlists/list
func (c *Client) fetchPlaylist(ctx context.Context, id string) (playlist, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/playlists", url.Values{
		"part": {"snippet"},
		"id":   {id},
	})
	if err != nil {
		return playlist{}, err
	}

	response := playlistResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return playlist{}, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return playlist{}, errNotFound
	}

	return response.Items[0], nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *Client) searchPlaylists(ctx context.Context, query string) ([]playlistSearchResult, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/search", url.Values{
		"q":          {query},
		"part":       {"snippet"},
		"type":       {"playlist"},
		"maxResults": {"10"},
	})
	if err != nil {
		return nil, err
	}

	response := playlistSearchResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	return response.Items, nil
}

// https://developers.google.com/youtube/v3/docs/playlistItems/list
func (c *Client) fetchPlaylistItems(ctx context.Context, id string) ([]playlistItem, error) {
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

	return response.Items, nil
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
		return nil, fmt.Errorf("non ok http status: %d", response.StatusCode)
	}

	return io.ReadAll(response.Body)
}
