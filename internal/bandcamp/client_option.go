package bandcamp

import "net/http"

type ClientOption func(client *HTTPClient)

func WithHTTPTransport(transport *http.Transport) ClientOption {
	return func(client *HTTPClient) {
		client.httpClient.Transport = transport
	}
}
