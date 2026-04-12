package soundcloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Client interface {
	FetchTrack(ctx context.Context, userSlug, trackSlug string) (*Track, error)
	FetchAlbum(ctx context.Context, userSlug, setSlug string) (*Album, error)
}

type HTTPClient struct {
	apiClient *http.Client
	apiHost   string
	apiScheme string
}

type hydrationItem struct {
	Hydratable string          `json:"hydratable"`
	Data       json.RawMessage `json:"data"`
}

var hydrationRe = regexp.MustCompile(`(?s)<script>\s*window\.__sc_hydration\s*=\s*(\[.*?\]);</script>`)

var (
	NotFoundError           = errors.New("not found")
	ErrHydrationNotFound    = errors.New("hydration script not found in html")
	ErrSoundDataNotFound    = errors.New("sound hydration data not found")
	ErrPlaylistDataNotFound = errors.New("playlist hydration data not found")
)

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		apiClient: &http.Client{},
		apiHost:   "soundcloud.com",
		apiScheme: "https",
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
