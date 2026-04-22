package streamnx

import (
	"github.com/GeorgeGorbanev/streamnx/internal/apple"
	"github.com/GeorgeGorbanev/streamnx/internal/bandcamp"
	"github.com/GeorgeGorbanev/streamnx/internal/deezer"
	"github.com/GeorgeGorbanev/streamnx/internal/soundcloud"
	"github.com/GeorgeGorbanev/streamnx/internal/spotify"
	"github.com/GeorgeGorbanev/streamnx/internal/yandex"
	"github.com/GeorgeGorbanev/streamnx/internal/youtube"
)

var (
	Providers = []*Provider{
		Apple,
		Bandcamp,
		Deezer,
		Spotify,
		Soundcloud,
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
	Bandcamp = &Provider{
		name:          "Bandcamp",
		сode:          "bc",
		trackIDParser: bandcamp.DetectTrackID,
		albumIDParser: bandcamp.DetectAlbumID,
	}
	Deezer = &Provider{
		name:          "Deezer",
		сode:          "dz",
		trackIDParser: deezer.DetectTrackID,
		albumIDParser: deezer.DetectAlbumID,
		cloakIDParser: deezer.DetectCloakID,
	}
	Spotify = &Provider{
		name:          "Spotify",
		сode:          "sf",
		trackIDParser: spotify.DetectTrackID,
		albumIDParser: spotify.DetectAlbumID,
	}
	Soundcloud = &Provider{
		name:          "Soundcloud",
		сode:          "sc",
		trackIDParser: soundcloud.DetectTrackID,
		albumIDParser: soundcloud.DetectAlbumID,
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

	cloakIDParser idParser
	trackIDParser idParser
	albumIDParser idParser
}

type idParser func(url string) (id string)

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Code() string {
	return p.сode
}

func (p *Provider) Regions() []string {
	return p.regions
}

func (p *Provider) DetectCloakEntityID(url string) string {
	if p.cloakIDParser == nil {
		return ""
	}
	return p.cloakIDParser(url)
}

func (p *Provider) DetectTrackID(trackURL string) string {
	return p.trackIDParser(trackURL)
}

func (p *Provider) DetectAlbumID(albumURL string) string {
	return p.albumIDParser(albumURL)
}

func (p *Provider) parseURL(url string) (string, *EntityType) {
	if id := p.DetectTrackID(url); id != "" {
		return id, new(Track)
	}
	if id := p.DetectAlbumID(url); id != "" {
		return id, new(Album)
	}
	if id := p.DetectCloakEntityID(url); id != "" {
		return id, new(Cloak)
	}
	return "", nil
}

func FindProviderByCode(code string) *Provider {
	for _, provider := range Providers {
		if provider.сode == code {
			return provider
		}
	}
	return nil
}
