package deezer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultAPIURL = "https://api.deezer.com"
)

var (
	NotFoundError = errors.New("not found")
)

type Client interface {
	FetchTrack(ctx context.Context, id string) (*Track, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error)
	FetchAlbum(ctx context.Context, id string) (*Album, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error)
}

type HTTPClient struct {
	apiURL     string
	clientID   string
	clientSecret string
	httpClient *http.Client
}

type searchResponse struct {
	Data []searchDataItem `json:"data"`
}

type searchDataItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type getResponse struct {
	Data []json.RawMessage `json:"data"`
}

func NewHTTPClient(clientID, clientSecret string, opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		httpClient:   &http.Client{},
		apiURL:       defaultAPIURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func (c *HTTPClient) FetchTrack(ctx context.Context, id string) (*Track, error) {
	url := fmt.Sprintf(`%s/track/%s`, c.apiURL, id)
	response, err := c.getAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}

	track := Track{}
	if err := json.NewDecoder(response.Body).Decode(&track); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get response: %s", err)
	}
	return &track, nil
}

func (c *HTTPClient) SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error) {
	url := fmt.Sprintf(`%s/search?q=artist:"%s" track:"%s"`, c.apiURL, artistName, trackName)
	response, err := c.getAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	sr := searchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	if len(sr.Data) == 0 {
		return nil, NotFoundError
	}

	track := Track{}
	if err := json.Unmarshal(sr.Data[0], &track); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track data: %s", err)
	}
	return &track, nil
}

func (c *HTTPClient) FetchAlbum(ctx context.Context, id string) (*Album, error) {
	url := fmt.Sprintf(`%s/album/%s`, c.apiURL, id)
	response, err := c.getAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}

	album := Album{}
	if err := json.NewDecoder(response.Body).Decode(&album); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get response: %s", err)
	}
	return &album, nil
}

func (c *HTTPClient) SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error) {
	url := fmt.Sprintf(`%s/search?q=artist:"%s" album:"%s"`, c.apiURL, artistName, albumName)
	response, err := c.getAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	sr := searchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	if len(sr.Data) == 0 {
		return nil, NotFoundError
	}

	album := Album{}
	if err := json.Unmarshal(sr.Data[0], &album); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album data: %s", err)
	}
	return &album, nil
}

func (c *HTTPClient) getAPI(ctx context.Context, reqURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	return c.httpClient.Do(req)
}
