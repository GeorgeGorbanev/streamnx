package deezer

import "net/http"

type ClientOption func(*HTTPClient)

func WithAPIURL(url string) ClientOption {
	return func(c *HTTPClient) {
		c.apiURL = url
	}
}

func WithCloakBaseURL(url string) ClientOption {
	return func(c *HTTPClient) {
		c.cloakBaseURL = url
	}
}

func WithAPIClient(client *http.Client) ClientOption {
	return func(c *HTTPClient) {
		c.apiClient = client
	}
}

func WithCloakClient(client *http.Client) ClientOption {
	return func(c *HTTPClient) {
		c.cloakClient = client
	}
}
