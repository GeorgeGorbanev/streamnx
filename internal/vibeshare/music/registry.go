package music

import (
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/youtube"
)

type Registry struct {
	adapters map[string]Adapter
}

type Adapter interface {
	DetectTrackID(link string) string
	GetTrack(id string) (*Track, error)
	SearchTrack(artistName, trackName string) (*Track, error)

	DetectAlbumID(link string) string
	GetAlbum(id string) (*Album, error)
	SearchAlbum(artistName, albumName string) (*Album, error)
}

type Track struct {
	ID       string
	Title    string
	Artist   string
	URL      string
	Provider *Provider
}

type Album struct {
	ID       string
	Title    string
	Artist   string
	URL      string
	Provider *Provider
}

type RegistryInput struct {
	AppleClient   apple.Client
	SpotifyClient spotify.Client
	YandexClient  yandex.Client
	YoutubeClient youtube.Client
}

func NewRegistry(input *RegistryInput) *Registry {
	return &Registry{
		adapters: map[string]Adapter{
			Apple.Code:   newAppleAdapter(input.AppleClient),
			Spotify.Code: newSpotifyAdapter(input.SpotifyClient),
			Yandex.Code:  newYandexAdapter(input.YandexClient),
			Youtube.Code: newYoutubeAdapter(input.YoutubeClient),
		},
	}
}

func (r *Registry) Adapter(p *Provider) Adapter {
	return r.adapters[p.Code]
}
