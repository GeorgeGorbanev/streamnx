package soundcloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sync"
)

type Client interface {
	FetchTrack(ctx context.Context, userSlug, trackSlug string) (*Track, error)
	FetchAlbum(ctx context.Context, userSlug, setSlug string) (*Album, error)
	SearchTrack(ctx context.Context, artist, title string) (*Track, error)
	SearchAlbum(ctx context.Context, artist, title string) (*Album, error)
}

type HTTPClient struct {
	client *http.Client

	apiURL string
	webURL string

	clientID      string
	clientIDMutex sync.RWMutex
}

var (
	hydrationRe = regexp.MustCompile(`(?s)<script>\s*window\.__sc_hydration\s*=\s*(\[.*?\]);</script>`)

	ErrNotFound = errors.New("not found")
)

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	const (
		defaultSearchAPIURL = "https://api-v2.soundcloud.com"
		defaultWebURL       = "https://soundcloud.com"
	)
	c := HTTPClient{
		client: &http.Client{},
		apiURL: defaultSearchAPIURL,
		webURL: defaultWebURL,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *HTTPClient) FetchTrack(ctx context.Context, userSlug, trackSlug string) (*Track, error) {
	path := fmt.Sprintf("/%s/%s", userSlug, trackSlug)
	html, err := c.getWebHTML(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch track page: %s", err)
	}

	trackJSON, err := findHydratableJSON(html, "sound")
	if err != nil {
		return nil, fmt.Errorf("failed to find track hydration data: %s", err)
	}

	track := Track{}
	if err := json.Unmarshal(trackJSON, &track); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sound hydration data: %s", err)
	}
	return &track, nil
}

func (c *HTTPClient) FetchAlbum(ctx context.Context, userSlug, setSlug string) (*Album, error) {
	path := fmt.Sprintf("/%s/sets/%s", userSlug, setSlug)
	html, err := c.getWebHTML(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch album page: %s", err)
	}

	albumJSON, err := findHydratableJSON(html, "playlist")
	if err != nil {
		return nil, fmt.Errorf("failed to find album hydration data: %s", err)
	}

	album := Album{}
	if err := json.Unmarshal(albumJSON, &album); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playlist hydration data: %s", err)
	}
	return &album, nil
}

func (c *HTTPClient) SearchTrack(ctx context.Context, artist, title string) (*Track, error) {
	body, err := c.getAPI(ctx, "/search/tracks", url.Values{
		"q": []string{artist + " " + title},
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Collection []Track `json:"collection"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response body: %s", err)
	}
	if len(result.Collection) == 0 {
		return nil, ErrNotFound
	}

	return &result.Collection[0], nil
}

func (c *HTTPClient) SearchAlbum(ctx context.Context, artist, title string) (*Album, error) {
	body, err := c.getAPI(ctx, "/search/albums", url.Values{
		"q": []string{artist + " " + title},
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Collection []Album `json:"collection"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response body: %s", err)
	}
	if len(result.Collection) == 0 {
		return nil, ErrNotFound
	}

	return &result.Collection[0], nil
}

func (c *HTTPClient) getWebHTML(ctx context.Context, path string) ([]byte, error) {
	u := c.webURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	return html, nil
}

func (c *HTTPClient) getAPI(ctx context.Context, path string, query url.Values) ([]byte, error) {
	clientID, err := c.getClientID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %s", err)
	}

	query.Set("client_id", clientID)

	u := c.apiURL + path + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	return body, nil
}

func (c *HTTPClient) getClientID(ctx context.Context) (string, error) {
	c.clientIDMutex.RLock()
	if c.clientID != "" {
		defer c.clientIDMutex.RUnlock()
		return c.clientID, nil
	}
	c.clientIDMutex.RUnlock()

	c.clientIDMutex.Lock()
	defer c.clientIDMutex.Unlock()

	if c.clientID != "" {
		return c.clientID, nil
	}

	html, err := c.getWebHTML(ctx, "/")
	if err != nil {
		return "", err
	}

	clientJSON, err := findHydratableJSON(html, "apiClient")
	if err != nil {
		return "", fmt.Errorf("failed to find api client hydration data: %s", err)
	}

	var apiClient struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(clientJSON, &apiClient); err != nil {
		return "", fmt.Errorf("failed to unmarshal api client hydration data: %s", err)
	}

	c.clientID = apiClient.ID
	return c.clientID, nil
}

func findHydratableJSON(html []byte, hydratable string) (json.RawMessage, error) {
	matches := hydrationRe.FindSubmatch(html)
	if len(matches) < 2 {
		return nil, errors.New("hydration script not found in html")
	}

	var items []struct {
		Hydratable string          `json:"hydratable"`
		Data       json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(matches[1], &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hydration json: %s", err)
	}

	for _, item := range items {
		if item.Hydratable == hydratable {
			return item.Data, nil
		}
	}
	return nil, fmt.Errorf("hydration data for %s not found", hydratable)
}
