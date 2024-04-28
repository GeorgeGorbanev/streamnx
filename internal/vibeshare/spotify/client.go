package spotify

import (
	"bytes"
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
	GetTrack(id string) (*Track, error)
	SearchTrack(artistName, trackName string) (*Track, error)
	GetAlbum(id string) (*Album, error)
	SearchAlbum(artistName, albumName string) (*Album, error)
}

type HTTPClient struct {
	authURL     string
	apiURL      string
	httpClient  *http.Client
	credentials *Credentials
	token       *token
}

type searchResult struct {
	Tracks tracksSection `json:"tracks"`
	Albums albumsSection `json:"albums"`
}

type tracksSection struct {
	Items []*Track `json:"items"`
}

type albumsSection struct {
	Items []*Album `json:"items"`
}

var InvalidIDError = fmt.Errorf("invalid id")

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
func (c *HTTPClient) GetTrack(id string) (*Track, error) {
	path := fmt.Sprintf("/v1/tracks/%s", id)
	body, err := c.getAPI(path, nil)
	if err != nil {
		if errors.Is(err, InvalidIDError) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	track := Track{}
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &track, nil
}

// https://developer.spotify.com/documentation/web-api/reference/search
func (c *HTTPClient) SearchTrack(artistName, trackName string) (*Track, error) {
	q := fmt.Sprintf("artist:%s track:%s", artistName, trackName)
	body, err := c.getAPI("/v1/search", url.Values{
		"q":     []string{q},
		"type":  []string{"track"},
		"limit": []string{"1"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	sr := searchResult{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if len(sr.Tracks.Items) == 0 {
		return nil, nil
	}

	return sr.Tracks.Items[0], nil
}

// https://developer.spotify.com/documentation/web-api/reference/get-an-album
func (c *HTTPClient) GetAlbum(id string) (*Album, error) {
	path := fmt.Sprintf("/v1/albums/%s", id)
	body, err := c.getAPI(path, nil)
	if err != nil {
		if errors.Is(err, InvalidIDError) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	album := Album{}
	if err := json.Unmarshal(body, &album); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &album, nil
}

// https://developer.spotify.com/documentation/web-api/reference/search
func (c *HTTPClient) SearchAlbum(artistName, albumName string) (*Album, error) {
	q := fmt.Sprintf("artist:%s album:%s", artistName, albumName)
	body, err := c.getAPI("/v1/search", url.Values{
		"q":     []string{q},
		"type":  []string{"album"},
		"limit": []string{"1"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	sr := searchResult{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if len(sr.Albums.Items) == 0 {
		return nil, nil
	}

	return sr.Albums.Items[0], nil
}

func (c *HTTPClient) getAPI(path string, query url.Values) ([]byte, error) {
	u := fmt.Sprintf("%s%s?%s", c.apiURL, path, query.Encode())
	resp, err := c.requestWithToken(u)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.token = nil
		resp, err = c.requestWithToken(u)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, InvalidIDError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// https://developer.spotify.com/documentation/web-api/tutorials/client-credentials-flow
func (c *HTTPClient) fetchToken() (*token, error) {
	url := fmt.Sprintf("%s/api/token", c.authURL)
	form := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest(http.MethodPost, url, form)
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

func (c *HTTPClient) requestWithToken(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token == nil || c.token.isExpired() {
		c.token, err = c.fetchToken()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token: %w", err)
		}
	}
	req.Header.Set("Authorization", c.token.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return resp, nil
}
