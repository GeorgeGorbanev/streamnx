package yandex

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

func NewClient(opts ...ClientOption) *Client {
	const defaultAPIURL = "https://api.music.yandex.net"
	c := Client{
		apiURL:     defaultAPIURL,
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

var errNotFound = errors.New("not found")

func (c *Client) fetchTrack(ctx context.Context, trackID string) (track, error) {
	body, err := c.getAPI(ctx, "/tracks/"+trackID, url.Values{})
	if err != nil {
		return track{}, fmt.Errorf("failed to get api: %w", err)
	}

	var resp struct {
		Result []track `json:"result"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return track{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if len(resp.Result) < 1 {
		return track{}, errNotFound
	}
	yandexTrack := resp.Result[0]
	if yandexTrack.Error == "not-found" {
		return track{}, errNotFound
	}
	if yandexTrack.Error != "" {
		return track{}, fmt.Errorf("api error: %s", yandexTrack.Error)
	}
	return yandexTrack, nil
}

func (c *Client) searchTracks(ctx context.Context, query string) ([]searchTrack, error) {
	body, err := c.getAPI(ctx, "/search", url.Values{
		"type": []string{"track"},
		"page": []string{"0"},
		"text": []string{query},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get api: %w", err)
	}

	var resp struct {
		Result struct {
			Tracks struct {
				Results []searchTrack `json:"results"`
			} `json:"tracks"`
		} `json:"result"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if resp.Result.Tracks.Results == nil {
		return []searchTrack{}, nil
	}

	return resp.Result.Tracks.Results, nil
}

func (c *Client) fetchAlbum(ctx context.Context, albumID string) (album, error) {
	body, err := c.getAPI(ctx, "/albums/"+albumID+"/with-tracks", url.Values{})
	if err != nil {
		return album{}, fmt.Errorf("failed to get api: %w", err)
	}

	var resp struct {
		Result album     `json:"result"`
		Error  *apiError `json:"error"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return album{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if resp.Error != nil {
		const notFoundAPIErr = "not-found"
		if resp.Error.Name == notFoundAPIErr {
			return album{}, errNotFound
		}
		return album{}, fmt.Errorf("api error: %s: %s", resp.Error.Name, resp.Error.Message)
	}

	return resp.Result, nil
}

func (c *Client) searchAlbums(ctx context.Context, query string) ([]searchAlbum, error) {
	body, err := c.getAPI(ctx, "/search", url.Values{
		"type": []string{"album"},
		"page": []string{"0"},
		"text": []string{query},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get api: %w", err)
	}

	var resp struct {
		Result struct {
			Albums struct {
				Results []searchAlbum `json:"results"`
			} `json:"albums"`
		} `json:"result"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if resp.Result.Albums.Results == nil {
		return []searchAlbum{}, nil
	}
	return resp.Result.Albums.Results, nil
}

func (c *Client) getAPI(ctx context.Context, path string, query url.Values) ([]byte, error) {
	u := fmt.Sprintf("%s%s?%s", c.apiURL, path, query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected api response status %d (%s)", resp.StatusCode, body)
	}

	return body, nil
}
