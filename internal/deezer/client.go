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
	defaultAPIURL   = "https://api.deezer.com"
	cloakBaseURL    = "https://link.deezer.com/s/"
)

var (
	NotFoundError = errors.New("not found")
)

type Client interface {
	FetchTrack(ctx context.Context, id string) (*Track, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error)
	FetchAlbum(ctx context.Context, id string) (*Album, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error)
	FollowCloak(ctx context.Context, cloakCode string) (string, error)
}

type HTTPClient struct {
	apiURL     string
	apiClient  *http.Client
	cloakClient *http.Client
}

type searchResult struct {
	Data  []*Track `json:"data"`
	Total int      `json:"total"`
}

type albumSearchResult struct {
	Data  []*Album `json:"data"`
	Total int      `json:"total"`
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		apiURL:      defaultAPIURL,
		apiClient:   &http.Client{},
		cloakClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *HTTPClient) FetchTrack(ctx context.Context, id string) (*Track, error) {
	body, err := c.getAPI(ctx, "/track/"+id, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	var track Track
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, fmt.Errorf("failed to parse track response: %w", err)
	}

	// Check if track exists (Deezer returns empty fields for non-existent tracks)
	if track.ID == "" {
		return nil, NotFoundError
	}

	return &track, nil
}

func (c *HTTPClient) SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error) {
	query := url.Values{}
	query.Set("q", fmt.Sprintf("artist:\"%s\" track:\"%s\"", artistName, trackName))

	body, err := c.getAPI(ctx, "/search", query)
	if err != nil {
		return nil, fmt.Errorf("failed to search track: %w", err)
	}

	var result searchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, NotFoundError
	}

	return result.Data[0], nil
}

func (c *HTTPClient) FetchAlbum(ctx context.Context, id string) (*Album, error) {
	body, err := c.getAPI(ctx, "/album/"+id, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to get album: %w", err)
	}

	var album Album
	if err := json.Unmarshal(body, &album); err != nil {
		return nil, fmt.Errorf("failed to parse album response: %w", err)
	}

	// Check if album exists
	if album.ID == "" {
		return nil, NotFoundError
	}

	return &album, nil
}

func (c *HTTPClient) SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error) {
	query := url.Values{}
	query.Set("q", fmt.Sprintf("artist:\"%s\" album:\"%s\"", artistName, albumName))

	body, err := c.getAPI(ctx, "/search/album", query)
	if err != nil {
		return nil, fmt.Errorf("failed to search album: %w", err)
	}

	var result albumSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse album search response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, NotFoundError
	}

	return result.Data[0], nil
}

func (c *HTTPClient) FollowCloak(ctx context.Context, cloakCode string) (string, error) {
	cloakURL := cloakBaseURL + cloakCode
	
	req, err := http.NewRequestWithContext(ctx, "HEAD", cloakURL, nil)
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
	u := fmt.Sprintf("%s%s", c.apiURL, path)
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	
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