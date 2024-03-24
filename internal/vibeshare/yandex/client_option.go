package yandex

type ClientOption func(client *HTTPClient)

func WithAPIURL(url string) ClientOption {
	return func(client *HTTPClient) {
		client.apiURL = url
	}
}
