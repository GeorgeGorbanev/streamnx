package bandcamp

import "net/http"

type ClientOption func(client *HTTPClient)

func WithAPIClient(c *http.Client) ClientOption {
	return func(client *HTTPClient) {
		client.apiClient = c
	}
}

func WithAPIHost(host string) ClientOption {
	return func(client *HTTPClient) {
		client.apiHost = host
	}
}

func WithAPIScheme(scheme string) ClientOption {
	return func(client *HTTPClient) {
		client.apiScheme = scheme
	}
}
