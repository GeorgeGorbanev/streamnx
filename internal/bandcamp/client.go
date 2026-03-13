package bandcamp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Client interface {
	FetchAlbum(ctx context.Context, artistSlug, id string) (*Entity, error)
	FetchTrack(ctx context.Context, artistSlug, id string) (*Entity, error)
	SearchAlbum(ctx context.Context, artist, title string) (*Entity, error)
	SearchTrack(ctx context.Context, artist, title string) (*Entity, error)
}

type HTTPClient struct {
	apiClient *http.Client
	apiHost   string
	apiScheme string
}

type entityLDJSON struct {
	Name     string `json:"name"`
	ByArtist struct {
		Name string `json:"name"`
	} `json:"byArtist"`
}

type searchRequest struct {
	SearchText   string `json:"search_text"`
	SearchFilter string `json:"search_filter"`
	FullPage     bool   `json:"full_page"`
}

type searchResponse struct {
	Auto struct {
		Results []struct {
			Name        string `json:"name"`
			BandName    string `json:"band_name"`
			ItemURLPath string `json:"item_url_path"`
		} `json:"results"`
	} `json:"auto"`
}

var ldJSONRe = regexp.MustCompile(`(?s)<script\s+type=["']application/ld\+json["']\s*>(.*?)</script>`)

var (
	ErrNotFound       = errors.New("entity not found")
	ErrLDJSONNotFound = errors.New("ld+json script not found in html")
)

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		apiClient: &http.Client{},
		apiHost:   "bandcamp.com",
		apiScheme: "https",
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *HTTPClient) FetchAlbum(ctx context.Context, artistSlug, id string) (*Entity, error) {
	u := fmt.Sprintf("%s://%s.%s/album/%s", c.apiScheme, artistSlug, c.apiHost, id)
	return c.fetchEntity(ctx, u)
}

func (c *HTTPClient) FetchTrack(ctx context.Context, artistSlug, id string) (*Entity, error) {
	u := fmt.Sprintf("%s://%s.%s/track/%s", c.apiScheme, artistSlug, c.apiHost, id)
	return c.fetchEntity(ctx, u)
}

func (c *HTTPClient) SearchAlbum(ctx context.Context, artist, title string) (*Entity, error) {
	return c.searchEntity(ctx, albumEntityType, artist, title)
}

func (c *HTTPClient) SearchTrack(ctx context.Context, artist, title string) (*Entity, error) {
	return c.searchEntity(ctx, trackEntityType, artist, title)
}

func (c *HTTPClient) fetchEntity(ctx context.Context, u string) (*Entity, error) {
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
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	matches := ldJSONRe.FindStringSubmatch(string(html))
	if len(matches) < 2 {
		return nil, ErrLDJSONNotFound
	}

	ldjson := entityLDJSON{}
	if err := json.Unmarshal([]byte(matches[1]), &ldjson); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ld+json: %s", err)
	}

	return &Entity{
		Name:     ldjson.Name,
		BandName: ldjson.ByArtist.Name,
		URL:      u,
	}, nil
}

func (c *HTTPClient) searchEntity(ctx context.Context, et entityType, artist, title string) (*Entity, error) {
	reqBody, err := json.Marshal(searchRequest{
		SearchText:   fmt.Sprintf("%s %s", artist, title),
		SearchFilter: string(et),
		FullPage:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %s", err)
	}

	u := fmt.Sprintf("%s://%s/api/bcsearch_public_api/1/autocomplete_elastic", c.apiScheme, c.apiHost)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := c.apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respBody searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %s", err)
	}
	if len(respBody.Auto.Results) == 0 {
		return nil, ErrNotFound
	}
	return &Entity{
		Name:     respBody.Auto.Results[0].Name,
		BandName: respBody.Auto.Results[0].BandName,
		URL:      respBody.Auto.Results[0].ItemURLPath,
	}, nil
}
