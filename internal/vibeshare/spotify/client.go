package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	authURL     string
	apiURL      string
	httpClient  *http.Client
	credentials *Credentials
}

type searchResult struct {
	Tracks tracksSection `json:"tracks"`
}

type tracksSection struct {
	Items []*Track `json:"items"`
}

type ClientOption func(client *Client)

func WithAuthURL(url string) ClientOption {
	return func(client *Client) {
		client.authURL = url
	}
}

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = url
	}
}

const (
	defaultAuthURL = "https://accounts.spotify.com"
	defaultAPIURL  = "https://api.spotify.com"
)

func NewClient(credentials *Credentials, opts ...ClientOption) *Client {
	c := Client{
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

// https://developer.spotify.com/documentation/web-api/tutorials/client-credentials-flow
func (c *Client) FetchToken() (*Token, error) {
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

	result := Token{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &result, nil
}

// https://developer.spotify.com/documentation/web-api/reference/get-track
func (c *Client) GetTrack(id string) (*Track, error) {
	url := fmt.Sprintf("%s/v1/tracks/%s", c.apiURL, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token, err := c.FetchToken()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token: %w", err)
	}

	req.Header.Set("Authorization", token.authHeader())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	track := Track{}
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &track, nil
}

// https://developer.spotify.com/documentation/web-api/reference/search
func (c *Client) SearchTrack(artistName, trackName string) (*Track, error) {
	query := url.QueryEscape(fmt.Sprintf("artist:%s track:%s", artistName, trackName))
	u := fmt.Sprintf("%s/v1/search?q=%s&type=track&limit=1", c.apiURL, query)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token, err := c.FetchToken()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token: %w", err)
	}

	req.Header.Set("Authorization", token.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
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
