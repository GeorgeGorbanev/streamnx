package ymusic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	apiURL string
}

type ClientOption func(client *Client)

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = url
	}
}

const defaultAPIURL = "https://api.music.yandex.net"

func NewClient(opts ...ClientOption) *Client {
	c := Client{
		apiURL: defaultAPIURL,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func (c *Client) Search(query string) (*SearchResponse, error) {
	u := fmt.Sprintf("%s/search?type=track&page=0&text=", c.apiURL)
	encodedQuery := url.QueryEscape(query)
	fullUrl := u + encodedQuery

	response, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	searchResponse := SearchResponse{}
	if err = json.Unmarshal(responseData, &searchResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	return &searchResponse, nil
}

func (c *Client) GetTrack(trackID string) (*Track, error) {
	u := fmt.Sprintf("%s/tracks/%s", c.apiURL, trackID)
	response, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request: %s", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	trackResponse := TrackResponse{}
	if err = json.Unmarshal(responseData, &trackResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if len(trackResponse.Result) < 1 {
		return nil, nil
	}

	return &trackResponse.Result[0], nil
}
