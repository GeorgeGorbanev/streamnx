package streamnx

import (
	"github.com/GeorgeGorbanev/streamnx/internal/apple"
	"github.com/GeorgeGorbanev/streamnx/internal/deezer"
	"github.com/GeorgeGorbanev/streamnx/internal/spotify"
	"github.com/GeorgeGorbanev/streamnx/internal/yandex"
	"github.com/GeorgeGorbanev/streamnx/internal/youtube"
)

var (
	Providers = []*Provider{
		Apple,
		Deezer,
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
	Deezer = &Provider{
		name:            "Deezer",
		сode:            "dz",
		trackIDParser:   deezer.DetectTrackID,
		albumIDParser:   deezer.DetectAlbumID,
		unknownIDParser: deezer.DetectUnknownEntityID,
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

	unknownIDParser idParser
	trackIDParser   idParser
	albumIDParser   idParser
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

func (p *Provider) DetectUnknownEntityID(url string) string {
	if p.unknownIDParser == nil {
		return ""
	}
	return p.unknownIDParser(url)
}

func (p *Provider) DetectTrackID(trackURL string) string {
	return p.trackIDParser(trackURL)
}

func (p *Provider) DetectAlbumID(albumURL string) string {
	return p.albumIDParser(albumURL)
}

func (p *Provider) parseURL(url string) (string, *EntityType) {
	parsers := map[EntityType]idParser{
		Track:   p.DetectTrackID,
		Album:   p.DetectAlbumID,
		Unknown: p.DetectUnknownEntityID,
	}
	for et, detector := range parsers {
		if id := detector(url); id != "" {
			return id, &et
		}
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
