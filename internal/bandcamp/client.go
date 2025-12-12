package bandcamp

type Client interface {
}

type HTTPClient struct {
}

func NewHTTPClient() Client {
	return &HTTPClient{}
}
