package apple

import (
	"encoding/json"
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
	GetTrack(id string) (*MusicEntity, error)
	SearchTrack(artistName, trackName string) (*MusicEntity, error)
	GetAlbum(id string) (*MusicEntity, error)
	SearchAlbum(artistName, albumName string) (*MusicEntity, error)
}

type HTTPClient struct {
	apiURL       string
	webPlayerURL string
	token        string
	httpClient   *http.Client
}

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

type searchResponse struct {
	Resources searchResources `json:"resources"`
}

type getResponse struct {
	Data []*MusicEntity `json:"data"`
}

type searchResources struct {
	Songs  map[string]*MusicEntity `json:"songs"`
	Albums map[string]*MusicEntity `json:"albums"`
}

func (c *HTTPClient) GetTrack(id string) (*MusicEntity, error) {
	if c.token == "" {
		token, err := c.fetchToken()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %s", err)
		}
		c.token = token
	}

	req, err := c.newAPIRequest(fmt.Sprintf(`%s/v1/catalog/us/songs/%s`, c.apiURL, id))
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	gr := getResponse{}
	if err := json.Unmarshal(body, &gr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	return gr.Data[0], nil
}

func (c *HTTPClient) SearchTrack(artistName, trackName string) (*MusicEntity, error) {
	if c.token == "" {
		token, err := c.fetchToken()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %s", err)
		}
		c.token = token
	}

	searchURL := fmt.Sprintf(`%s/v1/catalog/us/search?%s`, c.apiURL, searchQuery(artistName+" "+trackName))
	req, err := c.newAPIRequest(searchURL)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	sr := searchResponse{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	for _, song := range sr.Resources.Songs {
		return song, nil
	}
	return nil, nil
}
func (c *HTTPClient) GetAlbum(id string) (*MusicEntity, error) {
	if c.token == "" {
		token, err := c.fetchToken()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %s", err)
		}
		c.token = token
	}

	req, err := c.newAPIRequest(fmt.Sprintf(`%s/v1/catalog/us/albums/%s`, c.apiURL, id))
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	gr := getResponse{}
	if err := json.Unmarshal(body, &gr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	return gr.Data[0], nil
}
func (c *HTTPClient) SearchAlbum(artistName, albumName string) (*MusicEntity, error) {
	if c.token == "" {
		token, err := c.fetchToken()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %s", err)
		}
		c.token = token
	}

	searchURL := fmt.Sprintf(`%s/v1/catalog/us/search?%s`, c.apiURL, searchQuery(artistName+" "+albumName))
	req, err := c.newAPIRequest(searchURL)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	sr := searchResponse{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %s", err)
	}
	for _, album := range sr.Resources.Albums {
		return album, nil
	}
	return nil, nil
}

func (c *HTTPClient) fetchToken() (string, error) {
	webPlayerHTML, err := c.fetchWebPlayerHTML()
	if err != nil {
		return "", fmt.Errorf("failed to fetch index page: %s", err)
	}

	bundleName := parseBundleName(webPlayerHTML)
	if bundleName == "" {
		return "", fmt.Errorf("failed to extract bundle name")
	}

	webPlayerJS, err := c.fetchWebPlayerJS(bundleName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch index js: %s", err)
	}

	token, err := parseToken(webPlayerJS)
	if err != nil {
		return "", fmt.Errorf("failed to extract token: %s", err)
	}

	return token, nil
}

func (c *HTTPClient) fetchWebPlayerHTML() ([]byte, error) {
	response, err := c.httpClient.Get(c.webPlayerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func (c *HTTPClient) fetchWebPlayerJS(bundleName string) ([]byte, error) {
	response, err := c.httpClient.Get(c.webPlayerURL + bundleName)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func (c *HTTPClient) newAPIRequest(reqURL string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Origin", defaulWebPlayerURL)
	return req, nil
}

func searchQuery(term string) string {
	query := url.Values{}

	query.Set("term", term)
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
