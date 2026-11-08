package unofficial

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/youtube/client"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/duration"
)

type Client struct {
	apiURL     string
	httpClient *http.Client
}

type ClientOption func(*Client)

func WithAPIURL(url string) ClientOption {
	return func(c *Client) {
		c.apiURL = url
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

const webOrigin = "https://www.youtube.com"

var (
	playerJSON     = regexp.MustCompile(`(?:var\s+ytInitialPlayerResponse|window\["ytInitialPlayerResponse"\])\s*=\s*\{`)
	pageJSON       = regexp.MustCompile(`(?:\bytInitialData\b|\bwindow\s*\[\s*(?:"ytInitialData"|'ytInitialData')\s*\])\s*=\s*\{`)
	pageScriptJSON = regexp.MustCompile(`<script(?:\s+[^>]*?)?\s+id\s*=\s*(?:"yt-initial-data"|'yt-initial-data')[^>]*>\s*\{`)
	configJSON     = regexp.MustCompile(`ytcfg\.set\(\s*\{`)
)

func NewClient(opts ...ClientOption) *Client {
	c := Client{apiURL: webOrigin, httpClient: &http.Client{}}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *Client) FetchVideo(ctx context.Context, id string) (client.Video, error) {
	response, err := c.getPage(ctx, "/watch", url.Values{"v": {id}})
	if err != nil {
		return client.Video{}, err
	}

	var player webPlayerResponse
	if err := c.decodeWebObject(response.body, playerJSON, &player); err != nil {
		return client.Video{}, fmt.Errorf("decode watch page: %w", err)
	}
	details := player.Details
	if details.ID != id || details.Title == "" {
		return client.Video{}, c.fetchVideoError(player.Status.Code, player.Status.Reason)
	}

	return client.Video{
		Metadata: client.Metadata{
			ID:          details.ID,
			Title:       details.Title,
			Creator:     details.Author,
			Description: details.Description,
			CoverURL:    c.bestThumbnailURL(details.Thumbnail),
		},
		Duration: duration.ISO8601ToSeconds("PT" + details.Seconds + "S"),
	}, nil
}

func (c *Client) fetchVideoError(status, reason string) error {
	r := strings.ToLower(reason)
	switch {
	case strings.Contains(r, "not a bot"):
		return fmt.Errorf("YouTube requires bot verification (playability %q): %w", status, release.ErrScrapingBlocked)
	case status == "ERROR" && (r == "this video is unavailable" || r == "video unavailable"):
		return client.ErrNotFound
	default:
		return fmt.Errorf("missing video details (playability %q: %s)", status, reason)
	}
}

func (c *Client) SearchVideos(ctx context.Context, term string) ([]string, error) {
	const videosOnlyFilter = "EgIQAfABAQ=="
	sections, err := c.searchSections(ctx, term, videosOnlyFilter)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, client.SearchLimit)
	for _, section := range sections {
		if section.Section == nil {
			continue
		}
		for _, item := range section.Section.Items {
			var id string
			switch {
			case item.Video != nil:
				id = item.Video.ID
			case item.Lockup != nil && item.Lockup.Type == "LOCKUP_CONTENT_TYPE_VIDEO":
				id = item.Lockup.ID
			}
			if id == "" {
				continue
			}
			ids = append(ids, id)
			if len(ids) == client.SearchLimit {
				return ids, nil
			}
		}
	}
	return ids, nil
}

func (c *Client) SearchPlaylists(ctx context.Context, term string) ([]string, error) {
	const playlistsOnlyFilter = "EgIQAw=="
	sections, err := c.searchSections(ctx, term, playlistsOnlyFilter)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, client.SearchLimit)
	for _, section := range sections {
		if section.Section == nil {
			continue
		}
		for _, item := range section.Section.Items {
			var id string
			switch {
			case item.Playlist != nil:
				id = item.Playlist.ID
			case item.Lockup != nil && item.Lockup.Type == "LOCKUP_CONTENT_TYPE_PLAYLIST":
				id = item.Lockup.ID
			}
			if id == "" {
				continue
			}
			ids = append(ids, id)
			if len(ids) == client.SearchLimit {
				return ids, nil
			}
		}
	}
	return ids, nil
}

func (c *Client) searchSections(ctx context.Context, term, filter string) ([]webItem, error) {
	httpResponse, err := c.getPage(ctx, "/results", url.Values{"search_query": {term}, "sp": {filter}})
	if err != nil {
		return nil, err
	}

	var response webSearchResponse
	if err := c.decodeInitialData(httpResponse.body, &response); err != nil {
		return nil, fmt.Errorf("decode search page: %w", err)
	}
	sections := response.Contents.Results.Primary.Sections.Items
	if sections == nil {
		return nil, fmt.Errorf("missing search results section")
	}

	return sections, nil
}

func (c *Client) FetchPlaylist(ctx context.Context, id string) (client.Metadata, error) {
	page, _, err := c.getPlaylistPage(ctx, id)
	if err != nil {
		return client.Metadata{}, err
	}
	return c.playlistMetadata(id, page), nil
}

func (c *Client) FetchPlaylistWithItems(ctx context.Context, id string) (client.Metadata, []string, error) {
	page, body, err := c.getPlaylistPage(ctx, id)
	if err != nil {
		return client.Metadata{}, nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
	}
	items, err := c.playlistItems(ctx, id, page, body)
	if err != nil {
		return client.Metadata{}, nil, fmt.Errorf("failed to get playlist from youtube: failed to get playlist items: %w", err)
	}
	return c.playlistMetadata(id, page), items, nil
}

func (c *Client) playlistMetadata(id string, page webPlaylistResponse) client.Metadata {
	metadata := page.Metadata.Playlist
	header := page.Header.Playlist
	microformat := page.Microformat.Data
	cover := c.bestThumbnailURL(header.Banner.Hero.Thumbnail)
	if cover == "" {
		cover = c.bestThumbnailURL(microformat.Thumbnail)
	}
	owner := c.textValue(header.Owner)
	for _, item := range page.Sidebar.Playlist.Items {
		if name := c.textValue(item.Secondary.Owner.Renderer.Title); name != "" {
			owner = name
			break
		}
	}
	return client.Metadata{
		ID:          id,
		Title:       firstNonEmpty(metadata.Title, c.textValue(header.Title)),
		Description: firstNonEmpty(metadata.Description, c.textValue(header.Description), microformat.Description),
		Creator:     owner,
		CoverURL:    cover,
	}
}

func (c *Client) getPlaylistPage(ctx context.Context, id string) (webPlaylistResponse, []byte, error) {
	var page webPlaylistResponse
	query := url.Values{"list": {id}}
	response, err := c.getPage(ctx, "/playlist", query)
	if errors.Is(err, io.ErrUnexpectedEOF) && ctx.Err() == nil {
		response, err = c.getPage(ctx, "/playlist", query)
	}
	if err != nil {
		return page, nil, fmt.Errorf("load playlist page %q: %w", id, err)
	}
	if err := c.decodeInitialData(response.body, &page); err != nil {
		return page, nil, fmt.Errorf(
			"decode playlist page %q (status=%d, final_host=%q, final_path=%q, content_type=%q, bytes=%d, has_ytInitialData=%t): %w",
			id, response.status, response.host, response.path, response.contentType, len(response.body),
			bytes.Contains(response.body, []byte("ytInitialData")), err,
		)
	}
	for _, alert := range page.Alerts {
		if alert.Renderer.Type != "ERROR" {
			continue
		}
		message := c.textValue(alert.Renderer.Text)
		if message == "The playlist does not exist." || message == "This playlist is private." {
			return page, nil, client.ErrNotFound
		}
		return page, nil, fmt.Errorf("playlist unavailable: %s", message)
	}
	if firstNonEmpty(page.Metadata.Playlist.Title, c.textValue(page.Header.Playlist.Title)) == "" {
		return page, nil, fmt.Errorf("missing playlist metadata")
	}
	return page, response.body, nil
}

func (c *Client) playlistItems(ctx context.Context, id string, page webPlaylistResponse, body []byte) ([]string, error) {
	var entries []webItem
	for _, tab := range page.Contents.Browse.Tabs {
		if !tab.Renderer.Selected && len(page.Contents.Browse.Tabs) > 1 {
			continue
		}
		for _, item := range tab.Renderer.Content.Sections.Items {
			if section := item.Section; section != nil && (section.TargetID == "" || section.TargetID == id) {
				entries = section.Items
				break
			}
		}
		break
	}
	if entries == nil {
		return nil, fmt.Errorf("missing playlist contents")
	}
	for _, entry := range entries {
		if entry.PlaylistList != nil {
			entries = entry.PlaylistList.Items
			break
		}
	}

	items := make([]string, 0, len(entries))
	var token string
	items, token = c.appendPlaylistItems(items, entries)
	if token == "" {
		return items, nil
	}

	var config webConfig
	if err := c.decodeWebObject(body, configJSON, &config); err != nil {
		return nil, fmt.Errorf("decode playlist continuation context: %w", err)
	}
	if len(config.Context) == 0 {
		return nil, fmt.Errorf("missing playlist continuation context")
	}
	seen := make(map[string]bool)
	for token != "" {
		if seen[token] {
			return nil, fmt.Errorf("repeated playlist continuation")
		}
		seen[token] = true
		entries, err := c.browsePlaylist(ctx, id, config.Context, token)
		if err != nil {
			return nil, err
		}
		items, token = c.appendPlaylistItems(items, entries)
	}
	return items, nil
}

func (c *Client) appendPlaylistItems(items []string, entries []webItem) ([]string, string) {
	var token string
	for _, entry := range entries {
		var id string
		switch {
		case entry.PlaylistVideo != nil:
			id = entry.PlaylistVideo.ID
		case entry.Lockup != nil && entry.Lockup.Type == "LOCKUP_CONTENT_TYPE_VIDEO":
			id = entry.Lockup.ID
		}
		if id != "" {
			items = append(items, id)
		}
		if next := firstNonEmpty(entry.Continuation.Endpoint.Command.Token,
			entry.ContinuationView.Command.Endpoint.Command.Token); next != "" {
			token = next
		}
	}
	return items, token
}

func (c *Client) browsePlaylist(ctx context.Context, id string, clientContext json.RawMessage, token string) ([]webItem, error) {
	payload, err := json.Marshal(webBrowseRequest{Context: clientContext, Continuation: token})
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.request(ctx, http.MethodPost, "/youtubei/v1/browse?prettyPrint=false", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var response webBrowseResponse
	if err := json.Unmarshal(httpResponse.body, &response); err != nil {
		return nil, fmt.Errorf("decode playlist continuation: %w", err)
	}
	for _, actions := range [][]webBrowseAction{response.Actions, response.Endpoints} {
		for _, action := range actions {
			if action.Append.Items != nil && (action.Append.TargetID == "" || action.Append.TargetID == id) {
				return action.Append.Items, nil
			}
		}
	}
	return nil, fmt.Errorf("missing playlist continuation items")
}

type webHTTPResponse struct {
	body        []byte
	status      int
	host        string
	path        string
	contentType string
}

func (c *Client) getPage(ctx context.Context, path string, query url.Values) (webHTTPResponse, error) {
	query.Set("hl", "en")
	return c.request(ctx, http.MethodGet, path+"?"+query.Encode(), nil)
}

func (c *Client) request(ctx context.Context, method, path string, payload io.Reader) (webHTTPResponse, error) {
	const maxBody = 16 << 20
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.apiURL, "/")+path, payload)
	if err != nil {
		return webHTTPResponse{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept-Language", "en")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", webOrigin)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return webHTTPResponse{}, err
	}
	defer resp.Body.Close()

	finalURL := req.URL
	if resp.Request != nil {
		finalURL = resp.Request.URL
	}
	switch {
	case finalURL.Hostname() == "consent.youtube.com":
		return webHTTPResponse{}, fmt.Errorf("YouTube consent page (HTTP %d): %w", resp.StatusCode, release.ErrScrapingBlocked)
	case strings.HasPrefix(finalURL.Path, "/sorry/"):
		return webHTTPResponse{}, fmt.Errorf("YouTube verification page (HTTP %d): %w", resp.StatusCode, release.ErrScrapingBlocked)
	case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden:
		return webHTTPResponse{}, fmt.Errorf("YouTube HTTP %d from %q: %w", resp.StatusCode, finalURL.Hostname(), release.ErrScrapingBlocked)
	case resp.StatusCode == http.StatusNotFound && method == http.MethodGet && req.URL.Path != "/results":
		return webHTTPResponse{}, client.ErrNotFound
	case resp.StatusCode != http.StatusOK:
		return webHTTPResponse{}, fmt.Errorf("unexpected YouTube HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return webHTTPResponse{}, err
	}
	if len(body) > maxBody {
		return webHTTPResponse{}, fmt.Errorf("YouTube response exceeds %d bytes", maxBody)
	}
	return webHTTPResponse{
		body:        body,
		status:      resp.StatusCode,
		host:        finalURL.Host,
		path:        finalURL.EscapedPath(),
		contentType: resp.Header.Get("Content-Type"),
	}, nil
}

func (c *Client) decodeInitialData(body []byte, target any) error {
	if pageScriptJSON.Match(body) {
		return c.decodeWebObject(body, pageScriptJSON, target)
	}
	return c.decodeWebObject(body, pageJSON, target)
}

func (c *Client) decodeWebObject(body []byte, marker *regexp.Regexp, target any) error {
	match := marker.FindIndex(body)
	if match == nil {
		return fmt.Errorf("missing embedded YouTube JSON")
	}
	return json.NewDecoder(bytes.NewReader(body[match[1]-1:])).Decode(target)
}

func (c *Client) textValue(t webText) string {
	if t.SimpleText != "" {
		return t.SimpleText
	}
	var text strings.Builder
	for _, run := range t.Runs {
		text.WriteString(run.Text)
	}
	return text.String()
}

func (c *Client) bestThumbnailURL(t webThumbnails) string {
	var best thumbnail
	for _, item := range t.Items {
		if item.URL != "" && (best.URL == "" || item.Width*item.Height >= best.Width*best.Height) {
			best = item
		}
	}
	return best.URL
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
