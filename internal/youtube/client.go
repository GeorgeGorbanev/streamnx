package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const defaultAPIURL = "https://www.googleapis.com"

type Client interface {
	GetVideo(string) (*Video, error)
	SearchVideo(string) (*Video, error)
	GetPlaylist(string) (*Playlist, error)
	SearchPlaylist(string) (*Playlist, error)
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

type snippet struct {
	Title        string `json:"title"`
	ChannelTitle string `json:"channelTitle"`
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
func (c *HTTPClient) GetVideo(id string) (*Video, error) {
	url := fmt.Sprintf("%s/youtube/v3/videos?part=snippet&id=%s&key=%s", c.apiURL, id, c.apiKey)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok http status: %d", resp.StatusCode)
	}

	response := getSnippetResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return nil, nil
	}

	item := response.Items[0].Snippet
	return &Video{
		ID:           id,
		Title:        item.Title,
		ChannelTitle: item.ChannelTitle,
	}, nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *HTTPClient) SearchVideo(query string) (*Video, error) {
	apiURL := fmt.Sprintf("%s/youtube/v3/search?%s", c.apiURL, url.Values{
		"key":             {c.apiKey},
		"q":               {query},
		"part":            {"snippet"},
		"type":            {"video"},
		"videoCategoryId": {"10"},
		"maxResults":      {"1"},
	}.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	searchResult := searchSnippetResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	if len(searchResult.Items) == 0 {
		return nil, nil
	}

	item := searchResult.Items[0]
	return &Video{
		ID:           item.ID.VideoID,
		Title:        item.Snippet.Title,
		ChannelTitle: item.Snippet.ChannelTitle,
	}, nil
}

// https://developers.google.com/youtube/v3/docs/playlists/list
func (c *HTTPClient) GetPlaylist(id string) (*Playlist, error) {
	url := fmt.Sprintf("%s/youtube/v3/playlists?part=snippet&id=%s&key=%s", c.apiURL, id, c.apiKey)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok http status: %d", resp.StatusCode)
	}

	response := getSnippetResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(response.Items) == 0 {
		return nil, nil
	}

	item := response.Items[0]
	return &Playlist{
		ID:           item.ID,
		Title:        item.Snippet.Title,
		ChannelTitle: item.Snippet.ChannelTitle,
	}, nil
}

// https://developers.google.com/youtube/v3/docs/search/list
func (c *HTTPClient) SearchPlaylist(query string) (*Playlist, error) {
	apiURL := fmt.Sprintf("%s/youtube/v3/search?%s", c.apiURL, url.Values{
		"key":        {c.apiKey},
		"q":          {query},
		"part":       {"snippet"},
		"type":       {"playlist"},
		"maxResults": {"1"},
	}.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok http status: %d", resp.StatusCode)
	}

	searchResponse := searchSnippetResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode api response: %w", err)
	}
	if len(searchResponse.Items) == 0 {
		return nil, nil
	}

	item := searchResponse.Items[0]
	return &Playlist{
		ID:           item.ID.PlaylistID,
		Title:        item.Snippet.Title,
		ChannelTitle: item.Snippet.ChannelTitle,
	}, nil
}
