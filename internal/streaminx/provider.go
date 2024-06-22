package streaminx

import (
	"regexp"

	"github.com/GeorgeGorbanev/vibeshare/internal/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/youtube"
)

type Provider struct {
	Name    string
	Code    string
	TrackRe *regexp.Regexp
	AlbumRe *regexp.Regexp
	Regions []string
}

var (
	Apple = &Provider{
		Name:    "Apple",
		Code:    "ap",
		TrackRe: apple.TrackRe,
		AlbumRe: apple.AlbumRe,
		Regions: apple.ISO3166codes,
	}
	Spotify = &Provider{
		Name:    "Spotify",
		Code:    "sf",
		TrackRe: spotify.TrackRe,
		AlbumRe: spotify.AlbumRe,
	}
	Yandex = &Provider{
		Name:    "Yandex",
		Code:    "ya",
		TrackRe: yandex.TrackRe,
		AlbumRe: yandex.AlbumRe,
		Regions: yandex.Regions,
	}
	Youtube = &Provider{
		Name:    "Youtube",
		Code:    "yt",
		TrackRe: youtube.VideoRe,
		AlbumRe: youtube.PlaylistRe,
	}

	Providers = []*Provider{
		Apple,
		Spotify,
		Yandex,
		Youtube,
	}
)

func FindProviderByCode(code string) *Provider {
	for _, provider := range Providers {
		if provider.Code == code {
			return provider
		}
	}
	return nil
}
