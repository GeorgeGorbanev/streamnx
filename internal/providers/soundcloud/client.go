package soundcloud

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
	client        *http.Client
	apiURL        string
	webURL        string
	clientID      string
	clientIDMutex sync.RWMutex
}

type ClientOption func(client *Client)

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = url
	}
}

func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) {
		client.client = c
	}
}

func WithWebURL(host string) ClientOption {
	return func(client *Client) {
		client.webURL = host
	}
}

func NewClient(opts ...ClientOption) *Client {
	const (
		defaultAPIURL = "https://api-v2.soundcloud.com"
		defaultWebURL = "https://soundcloud.com"
	)
	c := Client{
		client: &http.Client{},
		apiURL: defaultAPIURL,
		webURL: defaultWebURL,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

var (
	errNotFound              = errors.New("not found")
	errHydrationDataNotFound = errors.New("hydration data not found")
	errUnauthorized          = errors.New("unexpected status code: 401")
)

func (c *Client) fetchTrack(ctx context.Context, userSlug, trackSlug string) (track, error) {
	path := fmt.Sprintf("/%s/%s", userSlug, trackSlug)
	html, err := c.getWebHTML(ctx, path)
	if err != nil {
		return track{}, fmt.Errorf("failed to fetch track page: %w", err)
	}

	const hydratableKey = "sound"
	trackJSON, err := c.findHydratableJSON(html, hydratableKey)
	switch {
	case errors.Is(err, errHydrationDataNotFound):
		return track{}, errNotFound
	case err != nil:
		return track{}, fmt.Errorf("failed to find track hydration data: %w", err)
	}

	t := track{}
	if err := json.Unmarshal(trackJSON, &t); err != nil {
		return track{}, fmt.Errorf("failed to unmarshal sound hydration data: %w", err)
	}
	return t, nil
}

func (c *Client) fetchAlbum(ctx context.Context, userSlug, setSlug string) (album, error) {
	path := fmt.Sprintf("/%s/sets/%s", userSlug, setSlug)
	html, err := c.getWebHTML(ctx, path)
	if err != nil {
		return album{}, fmt.Errorf("failed to fetch album page: %w", err)
	}

	const hydratableKey = "playlist"
	albumJSON, err := c.findHydratableJSON(html, hydratableKey)
	switch {
	case errors.Is(err, errHydrationDataNotFound):
		return album{}, errNotFound
	case err != nil:
		return album{}, fmt.Errorf("failed to find album hydration data: %w", err)
	}

	a := album{}
	if err := json.Unmarshal(albumJSON, &a); err != nil {
		return album{}, fmt.Errorf("failed to unmarshal playlist hydration data: %w", err)
	}
	return a, nil
}

func (c *Client) searchTracks(ctx context.Context, artist, title string) ([]track, error) {
	body, err := c.getAPI(ctx, "/search/tracks", url.Values{
		"q": []string{artist + " " + title},
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Collection []track `json:"collection"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response body: %w", err)
	}

	return result.Collection, nil
}

func (c *Client) searchAlbums(ctx context.Context, _, title string) ([]album, error) {
	body, err := c.getAPI(ctx, "/search/albums", url.Values{
		"q": []string{title},
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		Collection []album `json:"collection"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response body: %w", err)
	}

	return result.Collection, nil
}

func (c *Client) getWebHTML(ctx context.Context, path string) ([]byte, error) {
	u := c.webURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, errNotFound
	default:
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return html, nil
}

func (c *Client) getAPI(ctx context.Context, path string, query url.Values) ([]byte, error) {
	clientID, err := c.getClientID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %w", err)
	}

	body, err := c.getAPIWithClientID(ctx, path, query, clientID)
	if !errors.Is(err, errUnauthorized) {
		return body, err
	}

	clientID, err = c.refreshClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh client id: %w", err)
	}

	return c.getAPIWithClientID(ctx, path, query, clientID)
}

func (c *Client) getAPIWithClientID(
	ctx context.Context,
	path string,
	query url.Values,
	clientID string,
) ([]byte, error) {
	query.Set("client_id", clientID)

	u := c.apiURL + path + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, errNotFound
	case http.StatusUnauthorized:
		return nil, errUnauthorized
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (c *Client) getClientID(ctx context.Context) (string, error) {
	c.clientIDMutex.RLock()
	if c.clientID != "" {
		defer c.clientIDMutex.RUnlock()
		return c.clientID, nil
	}
	c.clientIDMutex.RUnlock()

	c.clientIDMutex.Lock()
	defer c.clientIDMutex.Unlock()

	if c.clientID != "" {
		return c.clientID, nil
	}

	clientID, err := c.fetchClientID(ctx)
	if err != nil {
		return "", err
	}
	c.clientID = clientID
	return clientID, nil
}

func (c *Client) refreshClientID(ctx context.Context, rejectedClientID string) (string, error) {
	c.clientIDMutex.Lock()
	defer c.clientIDMutex.Unlock()

	if c.clientID != "" && c.clientID != rejectedClientID {
		return c.clientID, nil
	}

	clientID, err := c.fetchClientID(ctx)
	if err != nil {
		return "", err
	}
	c.clientID = clientID
	return clientID, nil
}

func (c *Client) fetchClientID(ctx context.Context) (string, error) {
	html, err := c.getWebHTML(ctx, "/")
	if err != nil {
		return "", err
	}

	const hydratableKey = "apiClient"
	clientJSON, err := c.findHydratableJSON(html, hydratableKey)
	if err != nil {
		return "", fmt.Errorf("failed to find api client hydration data: %w", err)
	}

	var apiClient struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(clientJSON, &apiClient); err != nil {
		return "", fmt.Errorf("failed to unmarshal api client hydration data: %w", err)
	}

	return apiClient.ID, nil
}

var hydrationRe = regexp.MustCompile(`(?s)<script>\s*window\.__sc_hydration\s*=\s*(\[.*?\]);\s*</script>`)

func (c *Client) findHydratableJSON(html []byte, hydratable string) (json.RawMessage, error) {
	matches := hydrationRe.FindSubmatch(html)
	if len(matches) < 2 {
		return nil, errors.New("hydration script not found in html")
	}

	var items []struct {
		Hydratable string          `json:"hydratable"`
		Data       json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(matches[1], &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hydration json: %w", err)
	}

	for _, item := range items {
		if item.Hydratable == hydratable {
			return item.Data, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", errHydrationDataNotFound, hydratable)
}
