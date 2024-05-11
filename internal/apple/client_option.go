package apple

type ClientOption func(client *HTTPClient)

func WithWebPlayerURL(url string) ClientOption {
	return func(client *HTTPClient) {
		client.webPlayerURL = url
	}
}

func WithAPIURL(url string) ClientOption {
	return func(client *HTTPClient) {
		client.apiURL = url
	}
}
