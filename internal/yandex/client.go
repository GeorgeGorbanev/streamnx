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

type Client interface {
	FetchTrack(ctx context.Context, id string) (*Track, error)
	SearchTrack(ctx context.Context, query string) (*Track, error)
	FetchAlbum(ctx context.Context, id string) (*Album, error)
	SearchAlbum(ctx context.Context, query string) (*Album, error)
}

type HTTPClient struct {
	apiURL     string
	httpClient *http.Client
}

// TODO: rename to ErrFoo
var NotFoundError = errors.New("not found") //nolint:revive

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
	type trackResponse struct {
		Result []Track `json:"result"`
	}

	body, err := c.getAPI(ctx, "/tracks/"+trackID, url.Values{})
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

func (c *HTTPClient) SearchTrack(ctx context.Context, query string) (*Track, error) {
	type searchResponse struct {
		Result struct {
			Tracks struct {
				Results []Track `json:"results"`
			} `json:"tracks"`
		} `json:"result"`
	}

	body, err := c.getAPI(ctx, "/search", url.Values{
		"type": []string{"track"},
		"page": []string{"0"},
		"text": []string{query},
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
	const notFoundAPIErr = "not-found"
	type albumResponse struct {
		Result struct {
			Album
			Error *string `json:"error"`
		} `json:"result"`
	}

	body, err := c.getAPI(ctx, "/albums/"+albumID, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to get api: %s", err)
	}

	ar := albumResponse{}
	if err = json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if ar.Result.Error != nil {
		if *ar.Result.Error == notFoundAPIErr {
			return nil, NotFoundError
		}
		return nil, fmt.Errorf("api error: %s", *ar.Result.Error)
	}

	return &ar.Result.Album, nil
}

func (c *HTTPClient) SearchAlbum(ctx context.Context, query string) (*Album, error) {
	type searchResponse struct {
		Result struct {
			Albums struct {
				Results []Album `json:"results"`
			} `json:"albums"`
		} `json:"result"`
	}

	body, err := c.getAPI(ctx, "/search", url.Values{
		"type": []string{"album"},
		"page": []string{"0"},
		"text": []string{query},
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
