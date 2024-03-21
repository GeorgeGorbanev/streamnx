package spotify

type ClientOption func(client *Client)

func WithAuthURL(url string) ClientOption {
	return func(client *Client) {
		client.authURL = url
	}
}

func WithAPIURL(url string) ClientOption {
	return func(client *Client) {
		client.apiURL = url
	}
}
