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
	"strings"
	"sync"
)

type Client interface {
	FetchTrack(ctx context.Context, userSlug, trackSlug string) (*Track, error)
	FetchAlbum(ctx context.Context, userSlug, setSlug string) (*Album, error)
	SearchTrack(ctx context.Context, artist, title string) (*Track, error)
	SearchAlbum(ctx context.Context, artist, title string) (*Album, error)
}

type HTTPClient struct {
	apiClient     *http.Client
	apiHost       string
	apiScheme     string
	searchAPIURL  string
	clientID      string
	clientIDMutex sync.RWMutex
}

type hydrationItem struct {
	Hydratable string          `json:"hydratable"`
	Data       json.RawMessage `json:"data"`
}

type apiClientHydration struct {
	ID string `json:"id"`
}

var hydrationRe = regexp.MustCompile(`(?s)<script>\s*window\.__sc_hydration\s*=\s*(\[.*?\]);</script>`)

var (
	NotFoundError            = errors.New("not found")
	ErrHydrationNotFound     = errors.New("hydration script not found in html")
	ErrSoundDataNotFound     = errors.New("sound hydration data not found")
	ErrPlaylistDataNotFound  = errors.New("playlist hydration data not found")
	ErrAPIClientDataNotFound = errors.New("api client hydration data not found")
)

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		apiClient:    &http.Client{},
		apiHost:      "soundcloud.com",
		apiScheme:    "https",
		searchAPIURL: "https://api-v2.soundcloud.com",
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *HTTPClient) FetchTrack(ctx context.Context, userSlug, trackSlug string) (*Track, error) {
	html, err := c.fetchPageHTML(ctx, fmt.Sprintf("%s://%s/%s/%s", c.apiScheme, c.apiHost, userSlug, trackSlug))
	if err != nil {
		return nil, err
	}
	return parseTrackHTML(html)
}

func (c *HTTPClient) FetchAlbum(ctx context.Context, userSlug, setSlug string) (*Album, error) {
	html, err := c.fetchPageHTML(ctx, fmt.Sprintf("%s://%s/%s/sets/%s", c.apiScheme, c.apiHost, userSlug, setSlug))
	if err != nil {
		return nil, err
	}
	return parseAlbumHTML(html)
}

func (c *HTTPClient) SearchTrack(ctx context.Context, artist, title string) (*Track, error) {
	type searchResponse struct {
		Collection []Track `json:"collection"`
	}

	body, err := c.getAPI(ctx, "/search/tracks", url.Values{
		"q": []string{searchQuery(artist, title)},
	})
	if err != nil {
		return nil, err
	}

	result := searchResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response body: %s", err)
	}
	if len(result.Collection) == 0 {
		return nil, NotFoundError
	}

	return &result.Collection[0], nil
}

func (c *HTTPClient) SearchAlbum(ctx context.Context, artist, title string) (*Album, error) {
	type searchResponse struct {
		Collection []Album `json:"collection"`
	}

	body, err := c.getAPI(ctx, "/search/albums", url.Values{
		"q": []string{searchQuery(artist, title)},
	})
	if err != nil {
		return nil, err
	}

	result := searchResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response body: %s", err)
	}
	if len(result.Collection) == 0 {
		return nil, NotFoundError
	}

	return &result.Collection[0], nil
}

func (c *HTTPClient) fetchPageHTML(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	response, err := c.apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, NotFoundError
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
	clientID, err := c.clientIDValue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %s", err)
	}

	if query == nil {
		query = url.Values{}
	}
	query.Set("client_id", clientID)

	u := c.searchAPIURL + path + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
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

func (c *HTTPClient) clientIDValue(ctx context.Context) (string, error) {
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

	html, err := c.fetchPageHTML(ctx, fmt.Sprintf("%s://%s", c.apiScheme, c.apiHost))
	if err != nil {
		return "", err
	}

	clientID, err := parseAPIClientID(html)
	if err != nil {
		return "", err
	}
	c.clientID = clientID

	return c.clientID, nil
}

func parseTrackHTML(html []byte) (*Track, error) {
	items, err := parseHydrationItems(html)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Hydratable != "sound" {
			continue
		}

		track := Track{}
		if err := json.Unmarshal(item.Data, &track); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sound hydration data: %s", err)
		}
		if track.Kind != "track" {
			return nil, NotFoundError
		}

		return &track, nil
	}

	return nil, ErrSoundDataNotFound
}

func parseAlbumHTML(html []byte) (*Album, error) {
	items, err := parseHydrationItems(html)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Hydratable != "playlist" {
			continue
		}

		album := Album{}
		if err := json.Unmarshal(item.Data, &album); err != nil {
			return nil, fmt.Errorf("failed to unmarshal playlist hydration data: %s", err)
		}
		if album.Kind != "playlist" {
			return nil, NotFoundError
		}

		return &album, nil
	}

	return nil, ErrPlaylistDataNotFound
}

func parseHydrationItems(html []byte) ([]hydrationItem, error) {
	matches := hydrationRe.FindSubmatch(html)
	if len(matches) < 2 {
		return nil, ErrHydrationNotFound
	}

	var items []hydrationItem
	if err := json.Unmarshal(matches[1], &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hydration json: %s", err)
	}

	return items, nil
}

func parseAPIClientID(html []byte) (string, error) {
	items, err := parseHydrationItems(html)
	if err != nil {
		return "", err
	}

	for _, item := range items {
		if item.Hydratable != "apiClient" {
			continue
		}

		apiClient := apiClientHydration{}
		if err := json.Unmarshal(item.Data, &apiClient); err != nil {
			return "", fmt.Errorf("failed to unmarshal api client hydration data: %s", err)
		}
		if apiClient.ID == "" {
			return "", ErrAPIClientDataNotFound
		}

		return apiClient.ID, nil
	}

	return "", ErrAPIClientDataNotFound
}

func searchQuery(artist, title string) string {
	return strings.TrimSpace(strings.Join([]string{artist, title}, " "))
}
