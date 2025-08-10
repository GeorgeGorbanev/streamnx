package streamnx

import "errors"

// TODO: rename to ErrFoo
var UnknownLinkError = errors.New("unknown entity link") //nolint:revive

type Link struct {
	URL        string
	Provider   *Provider
	EntityID   string
	EntityType EntityType
}

func ParseLink(url string) (*Link, error) {
	for _, provider := range Providers {
		if id, et := provider.parseURL(url); id != "" {
			return &Link{
				URL:        url,
				Provider:   provider,
				EntityID:   id,
				EntityType: *et,
			}, nil
		}
	}
	return nil, UnknownLinkError
}
