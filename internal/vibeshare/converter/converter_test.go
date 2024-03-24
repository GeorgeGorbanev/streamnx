package converter

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/stretchr/testify/require"
)

func TestConverter_ConvertTrack(t *testing.T) {
	tests := []struct {
		name         string
		inputLink    string
		expectedLink string
		source       provider
		target       provider
	}{
		{
			name:         "when spotify track link given",
			inputLink:    "sampleSpotifyTrackLink",
			expectedLink: "sample yandex track url",
			source:       Spotify,
			target:       Yandex,
		},
		{
			name:         "when yandex track link given",
			inputLink:    "sampleYandexTrackLink",
			expectedLink: "sample spotify track url",
			source:       Yandex,
			target:       Spotify,
		},
		{
			name:         "when not found spotify track link given",
			inputLink:    "notFoundLink",
			expectedLink: "",
			source:       Spotify,
			target:       Yandex,
		},
		{
			name:         "when not found yandex track link given",
			inputLink:    "notFoundLink",
			expectedLink: "",
			source:       Yandex,
			target:       Spotify,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Converter{
				adapters: map[provider]music.Adapter{
					Spotify: &spotifyAdapterMock{},
					Yandex:  &yandexAdapterMock{},
				},
			}

			result, err := c.ConvertTrack(tt.inputLink, tt.source, tt.target)

			require.NoError(t, err)
			require.Equal(t, tt.expectedLink, result)
		})
	}
}

func TestConverter_ConvertAlbum(t *testing.T) {
	tests := []struct {
		name         string
		inputLink    string
		expectedLink string
		source       provider
		target       provider
	}{
		{
			name:         "when spotify album link given",
			inputLink:    "sampleSpotifyAlbumLink",
			expectedLink: "sample yandex album url",
			source:       Spotify,
			target:       Yandex,
		},
		{
			name:         "when yandex album link given",
			inputLink:    "sampleYandexAlbumLink",
			expectedLink: "sample spotify album url",
			source:       Yandex,
			target:       Spotify,
		},
		{
			name:         "when not found spotify album link given",
			inputLink:    "notFoundLink",
			expectedLink: "",
			source:       Spotify,
			target:       Yandex,
		},
		{
			name:         "when not found yandex album link given",
			inputLink:    "notFoundLink",
			expectedLink: "",
			source:       Yandex,
			target:       Spotify,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Converter{
				adapters: map[provider]music.Adapter{
					Spotify: &spotifyAdapterMock{},
					Yandex:  &yandexAdapterMock{},
				},
			}

			result, err := c.ConvertAlbum(tt.inputLink, tt.source, tt.target)

			require.NoError(t, err)
			require.Equal(t, tt.expectedLink, result)
		})
	}
}

type spotifyAdapterMock struct{}

func (a *spotifyAdapterMock) DetectTrackID(link string) string {
	if link == "sampleSpotifyTrackLink" {
		return "sampleSpotifyTrackId"
	}
	return ""
}

func (a *spotifyAdapterMock) GetTrack(id string) (*music.Track, error) {
	if id == "sampleSpotifyTrackId" {
		return &music.Track{
			Title:  "sample spotify title",
			Artist: "sample spotify artist",
		}, nil
	}
	return nil, nil
}

func (a *spotifyAdapterMock) SearchTrack(artistName, trackName string) (*music.Track, error) {
	if artistName == "sample yandex artist" && trackName == "sample yandex title" {
		return &music.Track{
			Title:  "sample spotify title",
			Artist: "sample spotify artist",
			URL:    "sample spotify track url",
		}, nil
	}
	return nil, nil
}

func (a *spotifyAdapterMock) DetectAlbumID(link string) string {
	if link == "sampleSpotifyAlbumLink" {
		return "sampleSpotifyAlbumId"
	}
	return ""
}

func (a *spotifyAdapterMock) GetAlbum(id string) (*music.Album, error) {
	if id == "sampleSpotifyAlbumId" {
		return &music.Album{
			Title:  "sample spotify title",
			Artist: "sample spotify artist",
		}, nil
	}
	return nil, nil

}

func (a *spotifyAdapterMock) SearchAlbum(artistName, albumName string) (*music.Album, error) {
	if artistName == "sample yandex artist" && albumName == "sample yandex title" {
		return &music.Album{
			Title:  "sample spotify title",
			Artist: "sample spotify artist",
			URL:    "sample spotify album url",
		}, nil
	}
	return nil, nil
}

type yandexAdapterMock struct{}

func (a *yandexAdapterMock) DetectTrackID(link string) string {
	if link == "sampleYandexTrackLink" {
		return "sampleYandexTrackId"
	}
	return ""
}

func (a *yandexAdapterMock) GetTrack(id string) (*music.Track, error) {
	if id == "sampleYandexTrackId" {
		return &music.Track{
			Title:  "sample yandex title",
			Artist: "sample yandex artist",
		}, nil
	}
	return nil, nil
}

func (a *yandexAdapterMock) SearchTrack(artistName, trackName string) (*music.Track, error) {
	if artistName == "sample spotify artist" && trackName == "sample spotify title" {
		return &music.Track{
			Title:  "sample yandex title",
			Artist: "sample yandex artist",
			URL:    "sample yandex track url",
		}, nil
	}
	return nil, nil
}

func (a *yandexAdapterMock) DetectAlbumID(link string) string {
	if link == "sampleYandexAlbumLink" {
		return "sampleYandexAlbumId"
	}
	return ""
}

func (a *yandexAdapterMock) GetAlbum(id string) (*music.Album, error) {
	if id == "sampleYandexAlbumId" {
		return &music.Album{
			Title:  "sample yandex title",
			Artist: "sample yandex artist",
		}, nil
	}
	return nil, nil

}

func (a *yandexAdapterMock) SearchAlbum(artistName, albumName string) (*music.Album, error) {
	if artistName == "sample spotify artist" && albumName == "sample spotify title" {
		return &music.Album{
			Title:  "sample yandex title",
			Artist: "sample yandex artist",
			URL:    "sample yandex album url",
		}, nil
	}
	return nil, nil
}
