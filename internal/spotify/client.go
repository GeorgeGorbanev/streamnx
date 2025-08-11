package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultAuthURL = "https://accounts.spotify.com"
	defaultAPIURL  = "https://api.spotify.com"
)

type Client interface {
	FetchTrack(ctx context.Context, id string) (*Track, error)
	SearchTrack(ctx context.Context, artist, title string) (*Track, error)
	FetchAlbum(ctx context.Context, id string) (*Album, error)
	SearchAlbum(ctx context.Context, artist, title string) (*Album, error)
}

type HTTPClient struct {
	authURL     string
	apiURL      string
	httpClient  *http.Client
	credentials *Credentials
	token       *token
}

// TODO: rename to ErrFoo
var NotFoundError = errors.New("not found") //nolint:revive

func NewHTTPClient(credentials *Credentials, opts ...ClientOption) *HTTPClient {
	c := HTTPClient{
		authURL:     defaultAuthURL,
		apiURL:      defaultAPIURL,
		credentials: credentials,
		httpClient:  &http.Client{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// https://developer.spotify.com/documentation/web-api/reference/get-track
func (c *HTTPClient) FetchTrack(ctx context.Context, id string) (*Track, error) {
	body, err := c.getAPI(ctx, "/v1/tracks/"+id, nil)
	if err != nil {
		return nil, err
	}

	track := Track{}
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &track, nil
}

// https://developer.spotify.com/documentation/web-api/reference/search
func (c *HTTPClient) SearchTrack(ctx context.Context, artist, title string) (*Track, error) {
	type searchResult struct {
		Tracks struct {
			Items []*Track `json:"items"`
		} `json:"tracks"`
	}

	body, err := c.getAPI(ctx, "/v1/search", url.Values{
		"q":     []string{fmt.Sprintf("artist:%s track:%s", artist, title)},
		"type":  []string{"track"},
		"limit": []string{"1"},
	})
	if err != nil {
		return nil, err
	}

	sr := searchResult{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if len(sr.Tracks.Items) == 0 {
		return nil, NotFoundError
	}

	return sr.Tracks.Items[0], nil
}

// https://developer.spotify.com/documentation/web-api/reference/get-an-album
func (c *HTTPClient) FetchAlbum(ctx context.Context, id string) (*Album, error) {
	body, err := c.getAPI(ctx, "/v1/albums/"+id, nil)
	if err != nil {
		return nil, err
	}

	album := Album{}
	if err := json.Unmarshal(body, &album); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &album, nil
}

// https://developer.spotify.com/documentation/web-api/reference/search
func (c *HTTPClient) SearchAlbum(ctx context.Context, artist, title string) (*Album, error) {
	type searchResult struct {
		Albums struct {
			Items []*Album `json:"items"`
		} `json:"albums"`
	}

	body, err := c.getAPI(ctx, "/v1/search", url.Values{
		"q":     []string{fmt.Sprintf("artist:%s album:%s", artist, title)},
		"type":  []string{"album"},
		"limit": []string{"1"},
	})
	if err != nil {
		return nil, err
	}

	sr := searchResult{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if len(sr.Albums.Items) == 0 {
		return nil, NotFoundError
	}

	return sr.Albums.Items[0], nil
}

// https://developer.spotify.com/documentation/web-api/tutorials/client-credentials-flow
func (c *HTTPClient) fetchToken(ctx context.Context) (*token, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.authURL+"/api/token",
		bytes.NewBuffer([]byte("grant_type=client_credentials")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.credentials.authHeader())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	result := token{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	result.fetchedAt = time.Now()
	return &result, nil
}

func (c *HTTPClient) getAPI(ctx context.Context, path string, query url.Values) ([]byte, error) {
	type errorResponse struct {
		Error struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		} `json:"error"`
	}

	resp, err := c.requestWithToken(ctx, c.apiURL+path+"?"+query.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, NotFoundError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		er := errorResponse{}
		if err := json.Unmarshal(body, &er); err != nil {
			return nil, fmt.Errorf("failed to load error response")
		}
		return nil, fmt.Errorf("unexpected API response: %d %s", er.Error.Status, er.Error.Message)
	}

	return body, nil
}

func (c *HTTPClient) requestWithToken(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token == nil || c.token.isExpired() {
		c.token, err = c.fetchToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %w", err)
		}
	}
	req.Header.Set("Authorization", c.token.authHeader())

	resp, err := c.httpClient.Do(req)
	if err == nil && resp.StatusCode == http.StatusUnauthorized {
		c.token, err = c.fetchToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %w", err)
		}
		req.Header.Set("Authorization", c.token.authHeader())
		resp, err = c.httpClient.Do(req)
	}
	return resp, err
}
