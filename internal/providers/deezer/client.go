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

type Client struct {
	apiURL       string
	cloakBaseURL string
	httpClient   *http.Client
}

type ClientOption func(*Client)

func WithAPIURL(url string) ClientOption {
	return func(c *Client) {
		c.apiURL = url
	}
}

func WithCloakBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.cloakBaseURL = url
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

var errNotFound = errors.New("not found")

func NewClient(options ...ClientOption) *Client {
	const (
		defaultAPIURL = "https://api.deezer.com"
		cloakBaseURL  = "https://link.deezer.com"
	)

	c := &Client{
		apiURL:       defaultAPIURL,
		cloakBaseURL: cloakBaseURL,
		httpClient:   &http.Client{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// https://developers.deezer.com/api/track
func (c *Client) fetchTrack(ctx context.Context, id string) (track, error) {
	body, err := c.getAPI(ctx, "/track/"+id, url.Values{})
	if err != nil {
		return track{}, fmt.Errorf("failed to get track: %w", err)
	}
	var t track
	if err := json.Unmarshal(body, &t); err != nil {
		return track{}, fmt.Errorf("failed to parse track response: %w", err)
	}
	if t.ID == 0 {
		return track{}, errNotFound
	}
	return t, nil
}

// https://developers.deezer.com/api/search
func (c *Client) searchTracks(ctx context.Context, artist, title string) ([]track, error) {
	body, err := c.getAPI(ctx, "/search", url.Values{
		"q": []string{
			fmt.Sprintf(`artist:"%s" track:"%s"`, artist, title),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search track: %w", err)
	}
	var result struct {
		Data []track `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}
	return result.Data, nil
}

func (c *Client) fetchTrackByISRC(ctx context.Context, isrc string) (track, error) {
	body, err := c.getAPI(ctx, "/track/isrc:"+isrc, url.Values{})
	if err != nil {
		return track{}, fmt.Errorf("failed to fetch track by isrc: %w", err)
	}
	var result track
	if err := json.Unmarshal(body, &result); err != nil {
		return track{}, fmt.Errorf("failed to parse isrc fetch response: %w", err)
	}
	if result.ID == 0 {
		return track{}, errNotFound
	}
	return result, nil
}

// https://developers.deezer.com/api/album
func (c *Client) fetchAlbum(ctx context.Context, id string) (album, error) {
	body, err := c.getAPI(ctx, "/album/"+id, url.Values{})
	if err != nil {
		return album{}, fmt.Errorf("failed to get album: %w", err)
	}
	var a album
	if err := json.Unmarshal(body, &a); err != nil {
		return album{}, fmt.Errorf("failed to parse album response: %w", err)
	}
	if a.ID == 0 {
		return album{}, errNotFound
	}
	return a, nil
}

// https://developers.deezer.com/api/search/album
func (c *Client) searchAlbums(ctx context.Context, artist, title string) ([]album, error) {
	body, err := c.getAPI(ctx, "/search/album", url.Values{
		"q": []string{
			fmt.Sprintf(`artist:"%s" album:"%s"`, artist, title),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search album: %w", err)
	}
	var result struct {
		Data []album `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse album search response: %w", err)
	}
	return result.Data, nil
}

func (c *Client) followCloak(ctx context.Context, id string) (string, error) {
	u := fmt.Sprintf("%s/s/%s", c.cloakBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "HEAD", u, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.doWithoutRedirect(req)
	if err != nil {
		return "", fmt.Errorf("failed to follow cloak link: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		if location != "" {
			return location, nil
		}
	}

	return "", fmt.Errorf("no redirect found for cloak link")
}

func (c *Client) doWithoutRedirect(req *http.Request) (*http.Response, error) {
	noRedirectClient := *c.httpClient
	noRedirectClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return noRedirectClient.Do(req)
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

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, errNotFound
	default:
		return nil, fmt.Errorf("unexpected api response status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
