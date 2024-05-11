package music

import (
	"github.com/GeorgeGorbanev/vibeshare/internal/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/translator"
	"github.com/GeorgeGorbanev/vibeshare/internal/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/youtube"
)

type Registry struct {
	adapters   map[string]Adapter
	translator translator.Translator
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
	YandexClient  yandex.Client
	YoutubeClient youtube.Client
	SpotifyClient spotify.Client
	Translator    translator.Translator
}

func NewRegistry(input *RegistryInput) *Registry {
	return &Registry{
		adapters: map[string]Adapter{
			Apple.Code:   newAppleAdapter(input.AppleClient),
			Spotify.Code: newSpotifyAdapter(input.SpotifyClient),
			Yandex.Code:  newYandexAdapter(input.YandexClient, input.Translator),
			Youtube.Code: newYoutubeAdapter(input.YoutubeClient),
		},
	}
}

func (r *Registry) Adapter(p *Provider) Adapter {
	return r.adapters[p.Code]
}

func (r *Registry) Close() error {
	return r.translator.Close()
}
