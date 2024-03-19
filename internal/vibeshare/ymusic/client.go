package ymusic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	apiURL string
}

type trackResponse struct {
	Result []Track `json:"result"`
}

type albumResponse struct {
	Result *Album `json:"result"`
}

type searchResponse struct {
	Result searchResult `json:"result"`
}

type searchResult struct {
	Tracks tracksSection `json:"tracks"`
	Albums albumsSection `json:"albums"`
}

type tracksSection struct {
	Results []Track `json:"results"`
}

type albumsSection struct {
	Results []Album `json:"results"`
}

type ClientOption func(client *Client)

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = url
	}
}

const defaultAPIURL = "https://api.music.yandex.net"

func NewClient(opts ...ClientOption) *Client {
	c := Client{
		apiURL: defaultAPIURL,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func (c *Client) GetTrack(trackID string) (*Track, error) {
	u := fmt.Sprintf("%s/tracks/%s", c.apiURL, trackID)
	response, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	tr := trackResponse{}
	if err = json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if len(tr.Result) < 1 {
		return nil, nil
	}

	return &tr.Result[0], nil
}

func (c *Client) SearchTrack(artistName, trackName string) (*Track, error) {
	u := fmt.Sprintf("%s/search?type=track&page=0&text=", c.apiURL)
	encodedQuery := url.QueryEscape(fmt.Sprintf("%s – %s", artistName, trackName))
	fullUrl := u + encodedQuery

	response, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	sr := searchResponse{}
	if err = json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if len(sr.Result.Tracks.Results) == 0 {
		return nil, nil
	}

	return &sr.Result.Tracks.Results[0], nil
}

func (c *Client) GetAlbum(albumID string) (*Album, error) {
	u := fmt.Sprintf("%s/albums/%s", c.apiURL, albumID)
	response, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	ar := albumResponse{}
	if err = json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	return ar.Result, nil
}

func (c *Client) SearchAlbum(artistName, albumName string) (*Album, error) {
	u := fmt.Sprintf("%s/search?type=album&page=0&text=", c.apiURL)
	encodedQuery := url.QueryEscape(fmt.Sprintf("%s – %s", artistName, albumName))
	fullUrl := u + encodedQuery

	response, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	sr := searchResponse{}
	if err = json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if len(sr.Result.Albums.Results) == 0 {
		return nil, nil
	}

	return &sr.Result.Albums.Results[0], nil
}

func (ar *albumResponse) UnmarshalJSON(data []byte) error {
	parsedResponse := map[string]map[string]any{}
	if err := json.Unmarshal(data, &parsedResponse); err != nil {
		return fmt.Errorf("failed to unmarshal album response: %s", err)
	}

	result, hasResult := parsedResponse["result"]
	if !hasResult {
		return fmt.Errorf("response does not contain result field")
	}

	apiError, hasAPIError := result["error"]
	if hasAPIError {
		apiErrorString, ok := apiError.(string)
		if !ok {
			return fmt.Errorf("api error is not a string")
		}
		if apiErrorString == "not-found" {
			return nil
		}
		return fmt.Errorf("api error: %s", apiErrorString)
	}

	albumJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal album result: %s", err)
	}

	if err = json.Unmarshal(albumJSON, &ar.Result); err != nil {
		return fmt.Errorf("failed to unmarshal album: %s", err)
	}

	return nil
}
