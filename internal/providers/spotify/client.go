package spotify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

type Client struct {
	authURL     string
	apiURL      string
	httpClient  *http.Client
	credentials *Credentials
	token       *token
	tokenMutex  sync.RWMutex
}

type Credentials struct {
	ClientID     string
	ClientSecret string
}

type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	fetchedAt   time.Time
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

func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = c
	}
}

var errNotFound = errors.New("not found")

func NewClient(credentials *Credentials, opts ...ClientOption) *Client {
	const (
		defaultAuthURL = "https://accounts.spotify.com"
		defaultAPIURL  = "https://api.spotify.com"
	)

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

// https://developer.spotify.com/documentation/web-api/reference/get-track
func (c *Client) fetchTrack(ctx context.Context, id string) (track, error) {
	body, err := c.getAPI(ctx, "/v1/tracks/"+id, nil)
	if err != nil {
		return track{}, err
	}

	t := track{}
	if err := json.Unmarshal(body, &t); err != nil {
		return track{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return t, nil
}

// https://developer.spotify.com/documentation/web-api/reference/get-an-album
func (c *Client) fetchAlbum(ctx context.Context, id string) (album, error) {
	body, err := c.getAPI(ctx, "/v1/albums/"+id, nil)
	if err != nil {
		return album{}, err
	}

	a := album{}
	if err := json.Unmarshal(body, &a); err != nil {
		return album{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return a, nil
}

// https://developer.spotify.com/documentation/web-api/reference/search
func (c *Client) searchTracks(ctx context.Context, artist, title string) ([]track, error) {
	return c.searchTracksWithQuery(ctx, fmt.Sprintf("artist:%s track:%s", artist, title))
}

func (c *Client) fetchTracksByISRC(ctx context.Context, isrc string) ([]track, error) {
	return c.searchTracksWithQuery(ctx, "isrc:"+isrc)
}

func (c *Client) searchTracksWithQuery(ctx context.Context, query string) ([]track, error) {
	type searchResult struct {
		Tracks struct {
			Items []track `json:"items"`
		} `json:"tracks"`
	}

	body, err := c.searchAPI(ctx, release.TypeTrack, query)
	if err != nil {
		return nil, err
	}

	sr := searchResult{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return sr.Tracks.Items, nil
}

// https://developer.spotify.com/documentation/web-api/reference/search
func (c *Client) searchAlbums(ctx context.Context, artist, title string) ([]album, error) {
	return c.searchAlbumsWithQuery(ctx, fmt.Sprintf("artist:%s album:%s", artist, title))
}

func (c *Client) fetchAlbumsByUPC(ctx context.Context, upc string) ([]album, error) {
	found, err := c.searchAlbumsWithQuery(ctx, "upc:"+upc)
	if err != nil {
		return nil, err
	}

	albums := make([]album, 0, len(found))
	for _, candidate := range found {
		if candidate.ID == "" {
			continue
		}
		fullAlbum, err := c.fetchAlbum(ctx, candidate.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to hydrate upc search album %q: %w", candidate.ID, err)
		}
		albums = append(albums, fullAlbum)
	}
	return albums, nil
}

func (c *Client) searchAlbumsWithQuery(ctx context.Context, query string) ([]album, error) {
	type searchResult struct {
		Albums struct {
			Items []album `json:"items"`
		} `json:"albums"`
	}

	body, err := c.searchAPI(ctx, release.TypeAlbum, query)
	if err != nil {
		return nil, err
	}

	sr := searchResult{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return sr.Albums.Items, nil
}

// https://developer.spotify.com/documentation/web-api/tutorials/client-credentials-flow
func (c *Client) fetchToken(ctx context.Context) (*token, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.authURL+"/api/token",
		bytes.NewBuffer([]byte("grant_type=client_credentials")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	cred := fmt.Sprintf("%s:%s", c.credentials.ClientID, c.credentials.ClientSecret)
	encodedCred := base64.StdEncoding.EncodeToString([]byte(cred))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedCred))

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

func (c *Client) getAPI(ctx context.Context, path string, query url.Values) ([]byte, error) {
	resp, err := c.requestWithAuth(ctx, c.apiURL+path+"?"+query.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		return nil, errNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		er := struct {
			Error struct {
				Status  int    `json:"status"`
				Message string `json:"message"`
			} `json:"error"`
		}{}
		if err := json.Unmarshal(body, &er); err != nil {
			return nil, fmt.Errorf("failed to load error response")
		}
		return nil, fmt.Errorf("unexpected api response: %d %s", er.Error.Status, er.Error.Message)
	}

	return body, nil
}

func (c *Client) searchAPI(ctx context.Context, rt release.Type, query string) ([]byte, error) {
	const defaultSearchLimit = "10"
	return c.getAPI(ctx, "/v1/search", url.Values{
		"q":     []string{query},
		"type":  []string{string(rt)},
		"limit": []string{defaultSearchLimit},
	})
}

func (c *Client) requestWithAuth(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	authHeader, refresh := c.authHeader()
	if refresh {
		if err = c.refreshToken(ctx, false); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		authHeader, _ = c.authHeader()
	}

	req.Header.Set("Authorization", authHeader)

	resp, err := c.httpClient.Do(req)
	if err == nil && resp.StatusCode == http.StatusUnauthorized {
		_ = resp.Body.Close()
		if err = c.refreshToken(ctx, true); err != nil {
			return nil, fmt.Errorf("failed to refresh token after 401: %w", err)
		}
		authHeader, _ = c.authHeader()
		req.Header.Set("Authorization", authHeader)
		resp, err = c.httpClient.Do(req)
	}
	return resp, err
}

func (c *Client) authHeader() (string, bool) {
	c.tokenMutex.RLock()
	defer c.tokenMutex.RUnlock()

	if c.token == nil || c.isTokenExpired() {
		return "", true
	}
	return fmt.Sprintf("Bearer %s", c.token.AccessToken), false
}

func (c *Client) refreshToken(ctx context.Context, force bool) error {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if !force && c.token != nil && !c.isTokenExpired() {
		return nil
	}

	newToken, err := c.fetchToken(ctx)
	if err != nil {
		return err
	}
	c.token = newToken
	return nil
}

func (c *Client) isTokenExpired() bool {
	return time.Since(c.token.fetchedAt) > time.Duration(c.token.ExpiresIn)*time.Second
}
