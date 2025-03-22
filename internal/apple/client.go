package apple

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
	defaultAPIURL      = "https://amp-api-edge.music.apple.com"
	defaulWebPlayerURL = "https://music.apple.com"
)

type Client interface {
	FetchTrack(ctx context.Context, id, storefront string) (*Entity, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error)
	FetchAlbum(ctx context.Context, id, storefront string) (*Entity, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error)
}

type HTTPClient struct {
	apiURL       string
	webPlayerURL string
	token        string
	httpClient   *http.Client
}

type fetchResponse struct {
	Data []*Entity `json:"data"`
}

type searchResponse struct {
	Resources struct {
		Songs  map[string]*Entity `json:"songs"`
		Albums map[string]*Entity `json:"albums"`
	} `json:"resources"`
	Results struct {
		Top struct {
			Data []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		} `json:"top"`
	} `json:"results"`
}

var NotFoundError = errors.New("not found")

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		httpClient:   &http.Client{},
		apiURL:       defaultAPIURL,
		webPlayerURL: defaulWebPlayerURL,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *HTTPClient) FetchTrack(ctx context.Context, id, storefront string) (*Entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/%s/songs/%s`, c.apiURL, storefront, id)
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}

	r := fetchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get response: %s", err)
	}
	return r.Data[0], nil
}

func (c *HTTPClient) SearchTrack(ctx context.Context, artist, track string) (*Entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/us/search?%s`, c.apiURL, searchQuery(artist, track))
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	sr := searchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	for _, topResult := range sr.Results.Top.Data {
		if topResult.Type == "songs" {
			return sr.Resources.Songs[topResult.ID], nil
		}
	}
	return nil, NotFoundError
}
func (c *HTTPClient) FetchAlbum(ctx context.Context, id, storefront string) (*Entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/%s/albums/%s`, c.apiURL, storefront, id)
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, NotFoundError
	}
	r := fetchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get response: %s", err)
	}
	return r.Data[0], nil
}
func (c *HTTPClient) SearchAlbum(ctx context.Context, artist, album string) (*Entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/us/search?%s`, c.apiURL, searchQuery(artist, album))
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	sr := searchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	for _, topResult := range sr.Results.Top.Data {
		if topResult.Type == "albums" {
			return sr.Resources.Albums[topResult.ID], nil
		}
	}
	return nil, NotFoundError
}

func (c *HTTPClient) getAPI(ctx context.Context, reqURL string) (*http.Response, error) {
	if c.token == "" {
		token, err := c.fetchToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %s", err)
		}
		c.token = token
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Origin", defaulWebPlayerURL)

	return c.httpClient.Do(req)
}

func (c *HTTPClient) fetchToken(ctx context.Context) (string, error) {
	webPlayerHTML, err := c.fetchWebPlayerHTML(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch index page: %s", err)
	}

	bundleName := parseBundleName(webPlayerHTML)
	if bundleName == "" {
		return "", fmt.Errorf("failed to extract bundle name")
	}

	webPlayerJS, err := c.fetchWebPlayerJS(ctx, bundleName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch index js: %s", err)
	}

	token, err := parseToken(webPlayerJS)
	if err != nil {
		return "", fmt.Errorf("failed to extract token: %s", err)
	}

	return token, nil
}

func (c *HTTPClient) fetchWebPlayerHTML(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.webPlayerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func (c *HTTPClient) fetchWebPlayerJS(ctx context.Context, bundleName string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.webPlayerURL+bundleName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func searchQuery(artist, entity string) string {
	query := url.Values{}

	query.Set("term", artist+" "+entity)
	query.Set("art[music-videos:url]", "c")
	query.Set("art[url]", "f")
	query.Set("extend", "artistUrl")
	query.Set("fields[albums]", "artistName,artistUrl,artwork,contentRating,editorialArtwork,editorialNotes,name,playParams,releaseDate,url,trackCount")
	query.Set("fields[artists]", "url,name,artwork")
	query.Set("format[resources]", "map")
	query.Set("include[albums]", "artists")
	query.Set("include[music-videos]", "artists")
	query.Set("include[songs]", "artists")
	query.Set("include[stations]", "radio-show")
	query.Set("l", "en-US")
	query.Set("limit", "21")
	query.Set("omit[resource]", "autos")
	query.Set("platform", "web")
	query.Set("relate[albums]", "artists")
	query.Set("relate[songs]", "albums")
	query.Set("types", "activities,albums,apple-curators,artists,curators,editorial-items,music-movies,music-videos,playlists,record-labels,songs,stations,tv-episodes,uploaded-videos")
	query.Set("with", "lyricHighlights,lyrics,serverBubbles")

	return query.Encode()
}
