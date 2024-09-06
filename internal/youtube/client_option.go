package youtube

import "net/http"

type ClientOption func(client *HTTPClient)

func WithAPIURL(url string) ClientOption {
	return func(client *HTTPClient) {
		client.apiURL = url
	}
}

func WithHTTPTransport(transport *http.Transport) ClientOption {
	return func(client *HTTPClient) {
		client.httpClient.Transport = transport
	}
}
