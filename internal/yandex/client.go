package yandex

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
	defaultAPIURL = "https://api.music.yandex.net"
)

var (
	NotFoundError = errors.New("not found")
)

type Client interface {
	FetchTrack(ctx context.Context, id string) (*Track, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error)
	FetchAlbum(ctx context.Context, id string) (*Album, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error)
}

type HTTPClient struct {
	apiURL     string
	httpClient *http.Client
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

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		apiURL:     defaultAPIURL,
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func (c *HTTPClient) FetchTrack(ctx context.Context, trackID string) (*Track, error) {
	path := fmt.Sprintf("/tracks/%s", trackID)
	body, err := c.getAPI(ctx, path, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to get api: %s", err)
	}

	tr := trackResponse{}
	if err = json.Unmarshal(body, &tr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if len(tr.Result) < 1 {
		return nil, NotFoundError
	}

	return &tr.Result[0], nil
}

func (c *HTTPClient) SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error) {
	body, err := c.getAPI(ctx, "/search", url.Values{
		"type": []string{"track"},
		"page": []string{"0"},
		"text": []string{fmt.Sprintf("%s – %s", artistName, trackName)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get api: %s", err)
	}

	sr := searchResponse{}
	if err = json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if len(sr.Result.Tracks.Results) == 0 {
		return nil, NotFoundError
	}

	return &sr.Result.Tracks.Results[0], nil
}

func (c *HTTPClient) FetchAlbum(ctx context.Context, albumID string) (*Album, error) {
	path := fmt.Sprintf("/albums/%s", albumID)
	body, err := c.getAPI(ctx, path, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to get api: %s", err)
	}

	ar := albumResponse{}
	if err = json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}
	if ar.Result == nil {
		return nil, NotFoundError
	}

	return ar.Result, nil
}

func (c *HTTPClient) SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error) {
	body, err := c.getAPI(ctx, "/search", url.Values{
		"type": []string{"album"},
		"page": []string{"0"},
		"text": []string{fmt.Sprintf("%s – %s", artistName, albumName)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get api: %s", err)
	}

	sr := searchResponse{}
	if err = json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if len(sr.Result.Albums.Results) == 0 {
		return nil, NotFoundError
	}

	return &sr.Result.Albums.Results[0], nil
}

func (c *HTTPClient) getAPI(ctx context.Context, path string, query url.Values) ([]byte, error) {
	u := fmt.Sprintf("%s%s?%s", c.apiURL, path, query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
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
