package converter

import (
	"fmt"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
)

type Converter struct {
	adapters map[provider]music.Adapter
}

type Input struct {
	SpotifyClient *spotify.HTTPClient
	YandexClient  *yandex.HTTPClient
}

type provider string

const (
	Yandex  provider = "yandex"
	Spotify provider = "spotify"
)

func NewConverter(input *Input) *Converter {
	spotifyAdapter := music.NewSpotifyAdapter(input.SpotifyClient)
	yandexAdapter := music.NewYandexAdapter(input.YandexClient)

	return &Converter{
		adapters: map[provider]music.Adapter{
			Spotify: spotifyAdapter,
			Yandex:  yandexAdapter,
		},
	}
}

func (c *Converter) ConvertTrack(link string, source, target provider) (string, error) {
	sourceTrackID := c.providerAdapter(source).DetectTrackID(link)
	sourceTrack, err := c.providerAdapter(source).GetTrack(sourceTrackID)
	if err != nil {
		return "", fmt.Errorf("error fetching track: %w", err)
	}
	if sourceTrack == nil {
		return "", nil
	}

	targetTrack, err := c.providerAdapter(target).SearchTrack(sourceTrack.Artist, sourceTrack.Title)
	if err != nil {
		return "", fmt.Errorf("error searching track: %w", err)
	}
	if targetTrack == nil {
		return "", nil
	}

	return targetTrack.URL, nil
}

func (c *Converter) ConvertAlbum(link string, source, target provider) (string, error) {
	sourceAlbumID := c.providerAdapter(source).DetectAlbumID(link)
	sourceAlbum, err := c.providerAdapter(source).GetAlbum(sourceAlbumID)
	if err != nil {
		return "", fmt.Errorf("error fetching album: %w", err)
	}
	if sourceAlbum == nil {
		return "", nil
	}

	targetAlbum, err := c.providerAdapter(target).SearchAlbum(sourceAlbum.Artist, sourceAlbum.Title)
	if err != nil {
		return "", fmt.Errorf("error searching album: %w", err)
	}
	if targetAlbum == nil {
		return "", nil
	}

	return targetAlbum.URL, nil
}

func (с *Converter) providerAdapter(p provider) music.Adapter {
	adapter, ok := с.adapters[p]
	if !ok {
		panic(fmt.Sprintf("adapter for provider %s not found", p))
	}
	return adapter
}
