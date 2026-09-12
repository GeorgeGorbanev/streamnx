package apple

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

type Client struct {
	apiURL       string
	webPlayerURL string
	token        string
	tokenMutex   sync.RWMutex
	httpClient   *http.Client
}

type ClientOption func(*Client)

func WithWebPlayerURL(url string) ClientOption {
	return func(client *Client) {
		client.webPlayerURL = url
	}
}

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = url
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func NewClient(opts ...ClientOption) *Client {
	const (
		defaultAPIURL       = "https://amp-api-edge.music.apple.com"
		defaultWebPlayerURL = "https://music.apple.com"
	)

	c := Client{
		httpClient:   &http.Client{},
		apiURL:       defaultAPIURL,
		webPlayerURL: defaultWebPlayerURL,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

var (
	errNotFound     = errors.New("not found")
	errUnauthorized = errors.New("unexpected status code: 401")
)

func (c *Client) fetchTrack(ctx context.Context, id, storefront string) (entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/%s/songs/%s`, c.apiURL, storefront, id)
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return entity{}, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return entity{}, errNotFound
	default:
		return entity{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	r := fetchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return entity{}, fmt.Errorf("failed to unmarshal get response: %w", err)
	}
	if len(r.Data) == 0 {
		return entity{}, fmt.Errorf("unexpected response: empty data")
	}
	return r.Data[0], nil
}

func (c *Client) fetchAlbum(ctx context.Context, id, storefront string) (entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/%s/albums/%s`, c.apiURL, storefront, id)
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return entity{}, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return entity{}, errNotFound
	default:
		return entity{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	r := fetchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return entity{}, fmt.Errorf("failed to unmarshal get response: %w", err)
	}
	if len(r.Data) == 0 {
		return entity{}, fmt.Errorf("unexpected response: empty data")
	}
	return r.Data[0], nil
}

// https://developer.apple.com/documentation/applemusicapi/get-multiple-catalog-albums-by-upc
func (c *Client) fetchAlbumsByUPC(ctx context.Context, upc string) ([]entity, error) {
	query := url.Values{
		"filter[upc]": {upc},
		"include":     {"tracks"},
	}
	u := fmt.Sprintf(`%s/v1/catalog/us/albums?%s`, c.apiURL, query.Encode())
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return []entity{}, nil
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	r := fetchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal upc fetch response: %w", err)
	}
	return r.Data, nil
}

func (c *Client) searchTracks(ctx context.Context, artist, track string) ([]entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/us/search?%s`, c.apiURL, c.searchQuery(artist, track))
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	sr := searchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}
	results := make([]entity, 0, len(sr.Results.Top.Data))
	for _, topResult := range sr.Results.Top.Data {
		if topResult.Type == "songs" {
			result, ok := sr.Resources.Songs[topResult.ID]
			if ok && result != nil {
				results = append(results, *result)
			}
		}
	}
	return results, nil
}

// https://developer.apple.com/documentation/applemusicapi/get-multiple-catalog-songs-by-isrc
func (c *Client) fetchTracksByISRC(ctx context.Context, isrc string) ([]entity, error) {
	query := url.Values{
		"filter[isrc]": {isrc},
		"include":      {"albums"},
	}
	u := fmt.Sprintf(`%s/v1/catalog/us/songs?%s`, c.apiURL, query.Encode())
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return []entity{}, nil
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	r := fetchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal isrc fetch response: %w", err)
	}
	return r.Data, nil
}

func (c *Client) searchAlbums(ctx context.Context, artist, album string) ([]entity, error) {
	u := fmt.Sprintf(`%s/v1/catalog/us/search?%s`, c.apiURL, c.searchQuery(artist, album))
	response, err := c.getAPI(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	sr := searchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}
	results := make([]entity, 0, len(sr.Results.Top.Data))
	for _, topResult := range sr.Results.Top.Data {
		if topResult.Type == "albums" {
			result, ok := sr.Resources.Albums[topResult.ID]
			if ok && result != nil {
				results = append(results, *result)
			}
		}
	}
	return results, nil
}

func (c *Client) searchQuery(artist, title string) string {
	return url.Values{
		"term":                  {artist + " " + title},
		"art[music-videos:url]": {"c"},
		"art[url]":              {"f"},
		"extend":                {"artistUrl"},
		"fields[artists]":       {"url,name,artwork"},
		"format[resources]":     {"map"},
		"include[albums]":       {"artists"},
		"include[music-videos]": {"artists"},
		"include[songs]":        {"artists"},
		"include[stations]":     {"radio-show"},
		"l":                     {"en-US"},
		"limit":                 {"21"},
		"omit[resource]":        {"autos"},
		"platform":              {"web"},
		"relate[albums]":        {"artists"},
		"relate[songs]":         {"albums"},
		"with":                  {"lyricHighlights,lyrics,serverBubbles"},
		"fields[albums]": {
			"artistName,artistUrl,artwork,contentRating,editorialArtwork," +
				"editorialNotes,name,playParams,releaseDate,url,trackCount,upc",
		},
		"types": {
			"activities,albums,apple-curators,artists,curators," +
				"editorial-items,music-movies,music-videos,playlists," +
				"record-labels,songs,stations,tv-episodes,uploaded-videos",
		},
	}.Encode()
}

func (c *Client) getAPI(ctx context.Context, reqURL string) (*http.Response, error) {
	token, err := c.authToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token: %w", err)
	}

	response, err := c.getAPIWithToken(ctx, reqURL, token)
	if !errors.Is(err, errUnauthorized) {
		return response, err
	}

	token, err = c.refreshAuthToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return c.getAPIWithToken(ctx, reqURL, token)
}

func (c *Client) getAPIWithToken(ctx context.Context, reqURL, token string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Origin", c.webPlayerURL)

	response, err := c.httpClient.Do(req)
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to send request: %w", err)
	case response.StatusCode == http.StatusUnauthorized:
		_ = response.Body.Close()
		return nil, errUnauthorized
	}
	return response, nil
}

func (c *Client) authToken(ctx context.Context) (string, error) {
	c.tokenMutex.RLock()
	if c.token != "" {
		token := c.token
		c.tokenMutex.RUnlock()
		return token, nil
	}
	c.tokenMutex.RUnlock()

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if c.token != "" {
		return c.token, nil
	}

	token, err := c.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	c.token = token
	return token, nil
}

func (c *Client) refreshAuthToken(ctx context.Context, rejectedToken string) (string, error) {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if c.token != "" && c.token != rejectedToken {
		return c.token, nil
	}

	token, err := c.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	c.token = token
	return token, nil
}

func (c *Client) fetchToken(ctx context.Context) (string, error) {
	webPlayerHTML, err := c.fetchWebPlayerHTML(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch index page: %w", err)
	}

	bundleName := c.parseBundleName(webPlayerHTML)
	if bundleName == "" {
		return "", fmt.Errorf("failed to extract bundle name")
	}

	webPlayerJS, err := c.fetchWebPlayerJS(ctx, bundleName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch index js: %w", err)
	}

	token, err := c.parseToken(webPlayerJS)
	if err != nil {
		return "", fmt.Errorf("failed to extract token: %w", err)
	}

	return token, nil
}

var tokenBundleRe = regexp.MustCompile(`src="(/assets/index[^"]+\.js)"`)

func (c *Client) parseBundleName(html []byte) string {
	matches := tokenBundleRe.FindSubmatch(html)
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

func (c *Client) fetchWebPlayerHTML(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.webPlayerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch player html page: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func (c *Client) parseToken(jsBundle []byte) (string, error) {
	tokenVar := c.parseTokenVar(jsBundle)
	if tokenVar == "" {
		return "", fmt.Errorf("token variable not found")
	}

	tokenVarValue, err := c.parseVariableValue(jsBundle, tokenVar)
	if err != nil {
		return "", fmt.Errorf("failed to find token variable value: %w", err)
	}

	return tokenVarValue, err
}

var tokenVarRe = regexp.MustCompile(`headers\.Authorization\s*=\s*` + "`Bearer \\${([a-zA-Z_$][a-zA-Z0-9_$]*)}`")

func (c *Client) parseTokenVar(jsBundle []byte) string {
	matches := tokenVarRe.FindSubmatch(jsBundle)
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

func (c *Client) parseVariableValue(jsBundle []byte, variable string) (string, error) {
	re, err := regexp.Compile(fmt.Sprintf(`(?:^|[^a-zA-Z0-9_$])%s\s*=\s*"([^"]+)"`, regexp.QuoteMeta(variable)))
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}
	matches := re.FindSubmatch(jsBundle)
	if len(matches) < 2 {
		return "", fmt.Errorf("value of variable %s not found", variable)
	}
	return string(matches[1]), nil
}

func (c *Client) fetchWebPlayerJS(ctx context.Context, bundleName string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.webPlayerURL+bundleName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch player js: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}
