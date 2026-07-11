package bandcamp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiClient *http.Client
	apiHost   string
	apiScheme string
}

type ClientOption func(client *Client)

func WithAPIClient(c *http.Client) ClientOption {
	return func(client *Client) {
		client.apiClient = c
	}
}

func WithAPI(scheme, host string) ClientOption {
	return func(client *Client) {
		client.apiScheme = scheme
		client.apiHost = host
	}
}

func NewClient(opts ...ClientOption) *Client {
	const (
		defaultAPIHost   = "bandcamp.com"
		defaultAPIScheme = "https"
	)
	c := Client{
		apiClient: &http.Client{},
		apiHost:   defaultAPIHost,
		apiScheme: defaultAPIScheme,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

var errNotFound = errors.New("not found")

func (c *Client) fetchAlbum(ctx context.Context, artistSlug, id string) (Entity, error) {
	u := fmt.Sprintf("%s://%s.%s/album/%s", c.apiScheme, artistSlug, c.apiHost, id)
	return c.fetchEntity(ctx, u)
}

func (c *Client) fetchTrack(ctx context.Context, artistSlug, id string) (Entity, error) {
	u := fmt.Sprintf("%s://%s.%s/track/%s", c.apiScheme, artistSlug, c.apiHost, id)
	return c.fetchEntity(ctx, u)
}

func (c *Client) searchAlbums(ctx context.Context, artist, title string) ([]Entity, error) {
	return c.searchEntity(ctx, albumEntityType, artist, title)
}

func (c *Client) searchTracks(ctx context.Context, artist, title string) ([]Entity, error) {
	return c.searchEntity(ctx, trackEntityType, artist, title)
}

func (c *Client) fetchEntity(ctx context.Context, u string) (Entity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.apiClient.Do(req)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return Entity{}, errNotFound
	default:
		return Entity{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to read response body: %w", err)
	}

	matches := ldJSONRe.FindStringSubmatch(string(html))
	if len(matches) < 2 {
		return Entity{}, errors.New("ld+json script not found in html")
	}

	ldjson := releaseLDJSON{}
	if err := json.Unmarshal([]byte(matches[1]), &ldjson); err != nil {
		return Entity{}, fmt.Errorf("failed to unmarshal ld+json: %w", err)
	}

	return Entity{
		URL:         u,
		Name:        ldjson.Name,
		AlbumTitle:  ldjson.InAlbum.Name,
		AlbumURL:    ldjson.InAlbum.ID,
		BandName:    ldjson.ByArtist.Name,
		CreatorName: ldjson.Publisher.Name,
		Description: ldjson.Description,
		CoverURL:    ldjson.Image,
		Duration:    ldjson.Duration,
		ReleaseDate: ldjson.DatePublished,
		TrackURLs:   trackURLs(ldjson),
	}, nil
}

func trackURLs(ldjson releaseLDJSON) []string {
	items := ldjson.Track.ItemListElement
	if len(items) == 0 {
		return nil
	}

	urls := make([]string, 0, len(items))
	for _, item := range items {
		switch {
		case item.Item.MainEntityOfPage != "":
			urls = append(urls, item.Item.MainEntityOfPage)
		case item.Item.ID != "":
			urls = append(urls, item.Item.ID)
		}
	}
	return urls
}

func (c *Client) searchEntity(ctx context.Context, et entityType, artist, title string) ([]Entity, error) {
	reqBody, err := json.Marshal(searchRequest{
		SearchText:   fmt.Sprintf("%s %s", artist, title),
		SearchFilter: string(et),
		FullPage:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	u := fmt.Sprintf("%s://%s/api/bcsearch_public_api/1/autocomplete_elastic", c.apiScheme, c.apiHost)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := c.apiClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respBody searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	entities := make([]Entity, 0, len(respBody.Auto.Results))
	for _, result := range respBody.Auto.Results {
		entities = append(entities, Entity{
			Name:       result.Name,
			AlbumTitle: result.AlbumName,
			BandName:   result.BandName,
			CoverURL:   result.Img,
			URL:        result.ItemURLPath,
		})
	}

	return entities, nil
}
