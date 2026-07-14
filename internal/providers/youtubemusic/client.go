package youtubemusic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	songsParams  = "EgWKAQIIAWoMEA4QChADEAQQCRAF"
	albumsParams = "EgWKAQIYAWoMEA4QChADEAQQCRAF"
)

type Client struct {
	apiURL     string
	httpClient *http.Client
}

type ClientOption func(client *Client)

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = strings.TrimRight(url, "/")
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

var errNotFound = errors.New("not found")

func NewClient(opts ...ClientOption) *Client {
	const defaultAPIURL = "https://music.youtube.com/youtubei/v1"

	c := Client{
		apiURL:     defaultAPIURL,
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (c *Client) fetchTrack(ctx context.Context, id string) (playerResponse, error) {
	body, err := c.post(ctx, "player", map[string]any{
		"video_id": id,
		"playbackContext": map[string]any{
			"contentPlaybackContext": map[string]any{
				"signatureTimestamp": int(time.Now().UTC().Unix()/86400) - 1,
			},
		},
	})
	if err != nil {
		return playerResponse{}, err
	}

	var player playerResponse
	if err := json.Unmarshal(body, &player); err != nil {
		return playerResponse{}, fmt.Errorf("failed to decode player response: %w", err)
	}
	if player.VideoDetails.VideoID == "" {
		return playerResponse{}, errNotFound
	}
	return player, nil
}

func (c *Client) fetchWatchNext(ctx context.Context, id string) (playlistPanelVideoRenderer, error) {
	body, err := c.post(ctx, "next", map[string]any{"videoId": id})
	if err != nil {
		return playlistPanelVideoRenderer{}, err
	}

	var response watchNextResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return playlistPanelVideoRenderer{}, fmt.Errorf("failed to decode next response: %w", err)
	}

	var selected playlistPanelVideoRenderer
	for _, tab := range response.Contents.Renderer.TabbedRenderer.Renderer.Tabs {
		for _, content := range tab.Renderer.Content.Queue.Content.Panel.Contents {
			track := content.Video
			if track.VideoID == id {
				return track, nil
			}
			if selected.VideoID == "" && track.Selected {
				selected = track
			}
		}
	}
	return selected, nil
}

func (c *Client) searchTracks(ctx context.Context, query string) ([]responsiveListItem, error) {
	return c.search(ctx, query, songsParams)
}

func (c *Client) searchAlbums(ctx context.Context, query string) ([]responsiveListItem, error) {
	return c.search(ctx, query, albumsParams)
}

func (c *Client) search(ctx context.Context, query, params string) ([]responsiveListItem, error) {
	body, err := c.post(ctx, "search", map[string]any{"query": query, "params": params})
	if err != nil {
		return nil, err
	}

	var response searchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	sections := response.Contents.SectionList
	if len(response.Contents.Tabbed.Tabs) > 0 {
		sections = response.Contents.Tabbed.Tabs[0].Renderer.Content.SectionListRenderer
	}
	items := make([]responsiveListItem, 0)
	for _, section := range sections.Contents {
		for _, item := range section.MusicShelfRenderer.Contents {
			items = append(items, item.Renderer)
		}
	}
	return items, nil
}

func (c *Client) fetchAlbum(ctx context.Context, id string) (responsiveHeader, []responsiveListItem, error) {
	browseID := id
	if !strings.HasPrefix(id, "MPRE") && !strings.HasPrefix(id, "VL") {
		browseID = "VL" + id
	}

	body, err := c.post(ctx, "browse", map[string]any{"browseId": browseID})
	if err != nil {
		return responsiveHeader{}, nil, err
	}
	var response browseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return responsiveHeader{}, nil, fmt.Errorf("failed to decode browse response: %w", err)
	}

	var header responsiveHeader
	tabs := response.Contents.TwoColumn.Tabs
	if len(tabs) == 0 {
		tabs = response.Contents.SingleColumn.Tabs
	}
	if len(tabs) > 0 {
		for _, section := range tabs[0].Renderer.Content.SectionListRenderer.Contents {
			if len(section.MusicResponsiveHeader.Title.Runs) > 0 {
				header = section.MusicResponsiveHeader
				break
			}
		}
	}

	items := make([]responsiveListItem, 0)
	for _, section := range response.Contents.TwoColumn.SecondaryContents.SectionListRenderer.Contents {
		shelf := section.MusicShelfRenderer.Contents
		if len(shelf) == 0 {
			shelf = section.MusicPlaylistShelfRenderer.Contents
		}
		for _, item := range shelf {
			items = append(items, item.Renderer)
		}
		if len(items) > 0 {
			break
		}
	}
	if len(header.Title.Runs) == 0 && len(items) == 0 {
		return responsiveHeader{}, nil, errNotFound
	}
	return header, items, nil
}

func (c *Client) post(ctx context.Context, endpoint string, body map[string]any) ([]byte, error) {
	body["context"] = map[string]any{
		"client": map[string]any{
			"clientName":    "WEB_REMIX",
			"clientVersion": "1." + time.Now().UTC().Format("20060102") + ".01.00",
		},
		"user": map[string]any{},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.apiURL+"/"+endpoint+"?alt=json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://music.youtube.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0")
	req.Header.Set("Cookie", "SOCS=CAI")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	switch response.StatusCode {
	case http.StatusOK:
		return responseBody, nil
	case http.StatusNotFound:
		return nil, errNotFound
	default:
		return nil, fmt.Errorf("unexpected api response status %d (%s)", response.StatusCode, responseBody)
	}
}
