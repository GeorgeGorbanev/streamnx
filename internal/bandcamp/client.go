package bandcamp

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
	FetchAlbum(ctx context.Context, artistSlug, id string) (*Entity, error)
	FetchTrack(ctx context.Context, artistSlug, id string) (*Entity, error)
}

type HTTPClient struct {
	httpClient *http.Client
}

const webHost = "bandcamp.com"

var ldJSONRe = regexp.MustCompile(`(?s)<script\s+type=["']application/ld\+json["']\s*>(.*?)</script>`)

var ErrNotFound = errors.New("entity not found")

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *HTTPClient) FetchAlbum(ctx context.Context, artistSlug, id string) (*Entity, error) {
	u := fmt.Sprintf("https://%s.%s/album/%s", artistSlug, webHost, id)
	return c.fetchEntity(ctx, u)
}

func (c *HTTPClient) FetchTrack(ctx context.Context, artistSlug, id string) (*Entity, error) {
	u := fmt.Sprintf("https://%s.%s/track/%s", artistSlug, webHost, id)
	return c.fetchEntity(ctx, u)
}

func (c *HTTPClient) fetchEntity(ctx context.Context, u string) (*Entity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	response, err := c.httpClient.Do(req)
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
		return nil, errors.New("failed to find ld+json script in html")
	}

	jsonld := struct {
		Name     string `json:"name"`
		ByArtist struct {
			Name string `json:"name"`
		} `json:"byArtist"`
	}{}
	if err := json.Unmarshal([]byte(matches[1]), &jsonld); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ld+json: %s", err)
	}

	return &Entity{
		Name:     jsonld.Name,
		BandName: jsonld.ByArtist.Name,
	}, nil
}
