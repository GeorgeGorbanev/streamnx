package streamnx

import "errors"

var (
	UnknownLinkError = errors.New("unknown entity link")
)

type Link struct {
	URL        string
	Provider   *Provider
	EntityID   string
	EntityType EntityType
}

func ParseLink(url string) (*Link, error) {
	for _, provider := range Providers {
		if id := provider.DetectTrackID(url); id != "" {
			return &Link{
				URL:        url,
				Provider:   provider,
				EntityID:   id,
				EntityType: Track,
			}, nil
		}
		if id := provider.DetectAlbumID(url); id != "" {
			return &Link{
				URL:        url,
				Provider:   provider,
				EntityID:   id,
				EntityType: Album,
			}, nil
		}
	}

	return nil, UnknownLinkError
}
