package streamnx

import (
	"github.com/GeorgeGorbanev/streamnx/internal/apple"
	"github.com/GeorgeGorbanev/streamnx/internal/spotify"
	"github.com/GeorgeGorbanev/streamnx/internal/yandex"
	"github.com/GeorgeGorbanev/streamnx/internal/youtube"
)

var (
	Providers = []*Provider{
		Apple,
		Spotify,
		Yandex,
		Youtube,
	}

	Apple = &Provider{
		name:          "Apple",
		сode:          "ap",
		regions:       apple.ISO3166codes,
		trackIDParser: apple.DetectTrackID,
		albumIDParser: apple.DetectAlbumID,
	}
	Spotify = &Provider{
		name:          "Spotify",
		сode:          "sf",
		trackIDParser: spotify.DetectTrackID,
		albumIDParser: spotify.DetectAlbumID,
	}
	Yandex = &Provider{
		name:          "Yandex",
		сode:          "ya",
		regions:       yandex.Regions,
		trackIDParser: yandex.DetectTrackID,
		albumIDParser: yandex.DetectAlbumID,
	}
	Youtube = &Provider{
		name:          "Youtube",
		сode:          "yt",
		trackIDParser: youtube.DetectTrackID,
		albumIDParser: youtube.DetectAlbumID,
	}
)

type Provider struct {
	name    string
	сode    string
	regions []string

	trackIDParser func(trackURL string) string
	albumIDParser func(albumURL string) string
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Code() string {
	return p.сode
}

func (p *Provider) Regions() []string {
	return p.regions
}

func (p *Provider) DetectTrackID(trackURL string) string {
	return p.trackIDParser(trackURL)
}

func (p *Provider) DetectAlbumID(albumURL string) string {
	return p.albumIDParser(albumURL)
}

func FindProviderByCode(code string) *Provider {
	for _, provider := range Providers {
		if provider.сode == code {
			return provider
		}
	}
	return nil
}
