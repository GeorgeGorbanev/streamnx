package spotify

type ClientOption func(client *HTTPClient)

func WithAuthURL(url string) ClientOption {
	return func(client *HTTPClient) {
		client.authURL = url
	}
}

func WithAPIURL(url string) ClientOption {
	return func(client *HTTPClient) {
		client.apiURL = url
	}
}
