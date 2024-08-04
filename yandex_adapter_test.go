package streamnx

import (
	"context"
	"testing"
	"time"

	"github.com/GeorgeGorbanev/streamnx/internal/yandex"

	"github.com/stretchr/testify/require"
)

type yandexClientMock struct {
	fetchTrack  map[string]*yandex.Track
	fetchAlbum  map[string]*yandex.Album
	searchTrack map[string]map[string]*yandex.Track
	searchAlbum map[string]map[string]*yandex.Album
}

func (c *yandexClientMock) FetchTrack(_ context.Context, id string) (*yandex.Track, error) {
	track, ok := c.fetchTrack[id]
	if !ok {
		return nil, yandex.NotFoundError
	}
	return track, nil
}

func (c *yandexClientMock) SearchTrack(_ context.Context, artistName, trackName string) (*yandex.Track, error) {
	if tracks, ok := c.searchTrack[artistName]; ok {
		track, ok := tracks[trackName]
		if !ok {
			return nil, yandex.NotFoundError
		}
		return track, nil
	}
	return nil, yandex.NotFoundError
}

func (c *yandexClientMock) FetchAlbum(_ context.Context, id string) (*yandex.Album, error) {
	album, ok := c.fetchAlbum[id]
	if !ok {
		return nil, yandex.NotFoundError
	}
	return album, nil
}

func (c *yandexClientMock) SearchAlbum(_ context.Context, artistName, albumName string) (*yandex.Album, error) {
	if albums, ok := c.searchAlbum[artistName]; ok {
		album, ok := albums[albumName]
		if !ok {
			return nil, yandex.NotFoundError
		}
		return album, nil
	}
	return nil, yandex.NotFoundError
}

type translatorMock struct {
	enToRu map[string]string
}

func (t *translatorMock) TranslateEnToRu(_ context.Context, text string) (string, error) {
	return t.enToRu[text], nil
}

func (t *translatorMock) Close() error {
	return nil
}

func TestYandexAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		yandexClientMock yandexClientMock
		expectedTrack    *Entity
		expectedErr      error
	}{
		{
			name: "found ID",
			id:   "42",
			yandexClientMock: yandexClientMock{
				fetchTrack: map[string]*yandex.Track{
					"42": {
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "sample artist"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
				Type:     Track,
			},
		},
		{
			name:             "not found ID",
			id:               "notFoundID",
			yandexClientMock: yandexClientMock{},
			expectedTrack:    nil,
			expectedErr:      EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, nil)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.FetchTrack(ctx, tt.id)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}
		})
	}
}

func TestYandexAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name             string
		artistName       string
		searchName       string
		yandexClientMock yandexClientMock
		translatorMock   translatorMock
		expectedTrack    *Entity
		expectedErr      error
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchTrack: map[string]map[string]*yandex.Track{
					"sample artist": {
						"sample name": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "sample artist"},
							},
							Albums: []yandex.Album{
								{ID: 41},
							},
						},
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
				Type:     Track,
			},
		},
		{
			name:       "found query but artist not matching",
			artistName: "sample artist not matching",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchTrack: map[string]map[string]*yandex.Track{
					"sample artist": {
						"sample artist not matching": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "not matching artist"},
							},
							Albums: []yandex.Album{
								{ID: 41},
							},
						},
					},
				},
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
		{
			name:       "found query matching translit",
			artistName: "sample artist matching translit",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchTrack: map[string]map[string]*yandex.Track{
					"sample artist matching translit": {
						"sample name": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "сампле артист матчинг транслит"},
							},
							Albums: []yandex.Album{
								{ID: 41},
							},
						},
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист матчинг транслит",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
				Type:     Track,
			},
		},
		{
			name:       "found query after translit",
			artistName: "sample artist after translit",
			searchName: "кириллическое название",
			yandexClientMock: yandexClientMock{
				searchTrack: map[string]map[string]*yandex.Track{
					"сампле артист афтер транслит": {
						"кириллическое название": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "сампле артист афтер транслит"},
							},
							Albums: []yandex.Album{
								{ID: 41},
							},
						},
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист афтер транслит",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
				Type:     Track,
			},
		},
		{
			name:       "found query after translation",
			artistName: "translatable artist",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchTrack: map[string]map[string]*yandex.Track{
					"translatable artist": {
						"sample name": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "переведенный артист"},
							},
							Albums: []yandex.Album{
								{ID: 41},
							},
						},
					},
				},
			},
			translatorMock: translatorMock{
				enToRu: map[string]string{
					"translatable artist": "переведенный артист",
				},
			},
			expectedTrack: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "переведенный артист",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
				Type:     Track,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, &tt.translatorMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.SearchTrack(ctx, tt.artistName, tt.searchName)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}
		})
	}
}

func TestYandexAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		yandexClientMock yandexClientMock
		expectedAlbum    *Entity
		expectedErr      error
	}{
		{
			name: "found id",
			id:   "42",
			yandexClientMock: yandexClientMock{
				fetchAlbum: map[string]*yandex.Album{
					"42": {
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "sample artist"},
						},
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
				Type:     Album,
			},
		},
		{
			name:          "not found ID",
			id:            "notFoundID",
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, nil)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.FetchAlbum(ctx, tt.id)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}
		})
	}
}

func TestYandexAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name             string
		artistName       string
		searchName       string
		yandexClientMock yandexClientMock
		translatorMock   translatorMock
		expectedAlbum    *Entity
		expectedErr      error
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchAlbum: map[string]map[string]*yandex.Album{
					"sample artist": {
						"sample name": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "sample artist"},
							},
						},
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
				Type:     Album,
			},
		},
		{
			name:       "found query but artist not matching",
			artistName: "sample artist not matching",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchAlbum: map[string]map[string]*yandex.Album{
					"sample artist": {
						"sample artist not matching": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "not matching artist"},
							},
						},
					},
				},
			},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
		{
			name:       "found query matching translit",
			artistName: "sample artist matching translit",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchAlbum: map[string]map[string]*yandex.Album{
					"sample artist matching translit": {
						"sample name": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "сампле артист матчинг транслит"},
							},
						},
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист матчинг транслит",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
				Type:     Album,
			},
		},
		{
			name:       "found query after translit",
			artistName: "sample artist after translit",
			searchName: "кириллическое название",
			yandexClientMock: yandexClientMock{
				searchAlbum: map[string]map[string]*yandex.Album{
					"сампле артист афтер транслит": {
						"кириллическое название": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "сампле артист афтер транслит"},
							},
						},
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист афтер транслит",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
				Type:     Album,
			},
		},
		{
			name:       "found query after translation",
			artistName: "translatable artist",
			searchName: "sample name",
			yandexClientMock: yandexClientMock{
				searchAlbum: map[string]map[string]*yandex.Album{
					"translatable artist": {
						"sample name": {
							ID:    42,
							Title: "sample name",
							Artists: []yandex.Artist{
								{Name: "переведенный артист"},
							},
						},
					},
				},
			},
			translatorMock: translatorMock{
				enToRu: map[string]string{
					"translatable artist": "переведенный артист",
				},
			},
			expectedAlbum: &Entity{
				ID:       "42",
				Title:    "sample name",
				Artist:   "переведенный артист",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
				Type:     Album,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, &tt.translatorMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.SearchAlbum(ctx, tt.artistName, tt.searchName)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}
		})
	}
}
