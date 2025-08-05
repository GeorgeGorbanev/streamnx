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
	cloakBaseURL  = "https://link.deezer.com"
)

type Client interface {
	FetchTrack(ctx context.Context, id string) (*Track, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error)
	FetchAlbum(ctx context.Context, id string) (*Album, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error)
	FollowCloak(ctx context.Context, cloakCode string) (string, error)
}

type HTTPClient struct {
	apiURL       string
	cloakBaseURL string
	apiClient    *http.Client
	cloakClient  *http.Client
}

var (
	NotFoundError = errors.New("not found")
)

func NewHTTPClient(options ...ClientOption) *HTTPClient {
	c := &HTTPClient{
		apiURL:       defaultAPIURL,
		cloakBaseURL: cloakBaseURL,
		apiClient:    &http.Client{},
		cloakClient: &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// https://developers.deezer.com/api/track
func (c *HTTPClient) FetchTrack(ctx context.Context, id string) (*Track, error) {
	body, err := c.getAPI(ctx, "/track/"+id, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}
	var track Track
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, fmt.Errorf("failed to parse track response: %w", err)
	}
	if track.ID == 0 {
		return nil, NotFoundError
	}
	return &track, nil
}

// https://developers.deezer.com/api/search
func (c *HTTPClient) SearchTrack(ctx context.Context, artist, title string) (*Track, error) {
	body, err := c.getAPI(ctx, "/search", url.Values{
		"q": []string{
			fmt.Sprintf(`artist:"%s" track:"%s"`, artist, title),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search track: %w", err)
	}
	var result struct {
		Data []Track `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, NotFoundError
	}
	return &result.Data[0], nil
}

// https://developers.deezer.com/api/album
func (c *HTTPClient) FetchAlbum(ctx context.Context, id string) (*Album, error) {
	body, err := c.getAPI(ctx, "/album/"+id, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to get album: %w", err)
	}
	var album Album
	if err := json.Unmarshal(body, &album); err != nil {
		return nil, fmt.Errorf("failed to parse album response: %w", err)
	}
	if album.ID == 0 {
		return nil, NotFoundError
	}
	return &album, nil
}

// https://developers.deezer.com/api/search/album
func (c *HTTPClient) SearchAlbum(ctx context.Context, artist, title string) (*Album, error) {
	body, err := c.getAPI(ctx, "/search/album", url.Values{
		"q": []string{
			fmt.Sprintf(`artist:"%s" album:"%s"`, artist, title),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search album: %w", err)
	}
	var result struct {
		Data []Album `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse album search response: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, NotFoundError
	}
	return &result.Data[0], nil
}

func (c *HTTPClient) FollowCloak(ctx context.Context, id string) (string, error) {
	u := fmt.Sprintf("%s/s/%s", c.cloakBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "HEAD", u, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.cloakClient.Do(req)
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

func (c *HTTPClient) getAPI(ctx context.Context, path string, query url.Values) ([]byte, error) {
	u := fmt.Sprintf("%s%s?%s", c.apiURL, path, query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
