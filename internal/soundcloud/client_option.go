package soundcloud

import "net/http"

type ClientOption func(client *HTTPClient)

func WithAPIURL(url string) ClientOption {
	return func(client *HTTPClient) {
		client.apiURL = url
	}
}

func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *HTTPClient) {
		client.client = c
	}
}

func WithWebURL(host string) ClientOption {
	return func(client *HTTPClient) {
		client.webURL = host
	}
}
