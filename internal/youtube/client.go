package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultAPIURL = "https://www.googleapis.com"

type Client interface {
	GetVideo(ctx context.Context, id string) (*Video, error)
	SearchVideo(ctx context.Context, term string) (*Video, error)
	GetPlaylist(ctx context.Context, id string) (*Playlist, error)
	SearchPlaylist(ctx context.Context, term string) (*Playlist, error)
	GetPlaylistItems(ctx context.Context, id string) ([]Video, error)
}

type HTTPClient struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
}

type getSnippetResponse struct {
	Items []*getSnippetItem `json:"items"`
}

type getSnippetItem struct {
	ID      string   `json:"id"`
	Snippet *snippet `json:"snippet"`
}

type searchSnippetResponse struct {
	Items []*searchSnippetItem `json:"items"`
}

type searchSnippetItem struct {
	ID      *searchSnippetID `json:"id"`
	Snippet *snippet         `json:"snippet"`
}

type searchSnippetID struct {
	VideoID    string `json:"videoId"`
	PlaylistID string `json:"playlistId"`
}

type getPlaylistItemsResponse struct {
	Items []*getSnippetItem `json:"items"`
}

func NewHTTPClient(apiKey string, opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		apiKey:     apiKey,
		apiURL:     defaultAPIURL,
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// https://developers.google.com/youtube/v3/docs/videos/list
func (c *HTTPClient) GetVideo(ctx context.Context, id string) (*Video, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/videos", url.Values{
		"part": {"snippet"},
		"id":   {id},
	})
	if err != nil {
		return nil, err
	}

	response := getSnippetResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return nil, nil
	}

	return &Video{
		ID:           id,
		Title:        response.Items[0].Snippet.Title,
		ChannelTitle: response.Items[0].Snippet.ownerChannelTitle(),
		Description:  response.Items[0].Snippet.Description,
	}, nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *HTTPClient) SearchVideo(ctx context.Context, query string) (*Video, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/search", url.Values{
		"q":               {query},
		"part":            {"snippet"},
		"type":            {"video"},
		"videoCategoryId": {"10"},
		"maxResults":      {"1"},
	})
	if err != nil {
		return nil, err
	}

	response := searchSnippetResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, nil
	}

	item := response.Items[0]
	return &Video{
		ID:           item.ID.VideoID,
		Title:        item.Snippet.Title,
		ChannelTitle: item.Snippet.ownerChannelTitle(),
	}, nil
}

// https://developers.google.com/youtube/v3/docs/playlists/list
func (c *HTTPClient) GetPlaylist(ctx context.Context, id string) (*Playlist, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/playlists", url.Values{
		"part": {"snippet"},
		"id":   {id},
	})
	if err != nil {
		return nil, err
	}

	response := getSnippetResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return nil, nil
	}

	item := response.Items[0]
	return &Playlist{
		ID:           item.ID,
		Title:        item.Snippet.Title,
		ChannelTitle: item.Snippet.ownerChannelTitle(),
	}, nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *HTTPClient) SearchPlaylist(ctx context.Context, query string) (*Playlist, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/search", url.Values{
		"q":          {query},
		"part":       {"snippet"},
		"type":       {"playlist"},
		"maxResults": {"1"},
	})
	if err != nil {
		return nil, err
	}

	response := searchSnippetResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return nil, nil
	}

	return &Playlist{
		ID:           response.Items[0].ID.PlaylistID,
		Title:        response.Items[0].Snippet.Title,
		ChannelTitle: response.Items[0].Snippet.ownerChannelTitle(),
	}, nil
}

// https://developers.google.com/youtube/v3/docs/playlistItems/list
func (c *HTTPClient) GetPlaylistItems(ctx context.Context, id string) ([]Video, error) {
	body, err := c.getWithKey(ctx, "/youtube/v3/playlistItems", url.Values{
		"part":       {"snippet"},
		"playlistId": {id},
	})
	if err != nil {
		return nil, err
	}
	response := getPlaylistItemsResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}

	videos := make([]Video, 0, len(response.Items))
	for _, item := range response.Items {
		videos = append(videos, Video{
			ID:           item.ID,
			Title:        item.Snippet.Title,
			ChannelTitle: item.Snippet.VideoOwnerChannelTitle,
			Description:  item.Snippet.Description,
		})
	}

	return videos, nil
}

func (c *HTTPClient) getWithKey(ctx context.Context, path string, values url.Values) ([]byte, error) {
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
