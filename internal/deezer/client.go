package deezer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// https://developers.deezer.com/api/explorer
// https://rapidapi.com/deezerdevs/api/deezer-1

const (
	defaultAPIURL = "https://api.deezer.com"
)

type Client interface {
	FetchTrack(ctx context.Context, id string) (*Track, error)
}

type HTTPClient struct {
	apiURL     string
	httpClient *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		apiURL:     defaultAPIURL,
		httpClient: &http.Client{},
	}
}

func (c *HTTPClient) FetchTrack(ctx context.Context, id string) (*Track, error) {
	//body, err := c.getAPI(ctx, c.apiURL+"/track/"+id, url.Values{})
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get api: %s", err)
	//}

	return nil, nil
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
