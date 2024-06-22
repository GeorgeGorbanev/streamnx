package streaminx

import (
	"context"
	"testing"

	"github.com/GeorgeGorbanev/streaminx/internal/yandex"
	"github.com/stretchr/testify/require"
)

type yandexClientMock struct {
	getTrack    map[string]*yandex.Track
	getAlbum    map[string]*yandex.Album
	searchTrack map[string]map[string]*yandex.Track
	searchAlbum map[string]map[string]*yandex.Album
}

func (c *yandexClientMock) GetTrack(id string) (*yandex.Track, error) {
	return c.getTrack[id], nil
}

func (c *yandexClientMock) SearchTrack(artistName, trackName string) (*yandex.Track, error) {
	if tracks, ok := c.searchTrack[artistName]; ok {
		return tracks[trackName], nil
	}
	return nil, nil
}

func (c *yandexClientMock) GetAlbum(id string) (*yandex.Album, error) {
	return c.getAlbum[id], nil
}

func (c *yandexClientMock) SearchAlbum(artistName, albumName string) (*yandex.Album, error) {
	if albums, ok := c.searchAlbum[artistName]; ok {
		return albums[albumName], nil
	}
	return nil, nil
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

func TestYandexAdapter_DetectTrackID(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		wantID        string
		expectedError error
	}{
		{
			name:   "Valid Track URL – .com",
			url:    "https://music.yandex.com/album/3192570/track/1197793",
			wantID: "1197793",
		},
		{
			name:   "Valid Track URL – .ru",
			url:    "https://music.yandex.ru/album/3192570/track/1197793",
			wantID: "1197793",
		},
		{
			name:   "Valid Track URL – .by",
			url:    "https://music.yandex.by/album/3192570/track/1197793",
			wantID: "1197793",
		},
		{
			name:   "Valid Track URL – .kz",
			url:    "https://music.yandex.kz/album/3192570/track/1197793",
			wantID: "1197793",
		},
		{
			name:   "Valid Track URL – .uz",
			url:    "https://music.yandex.uz/album/3192570/track/1197793",
			wantID: "1197793",
		},
		{
			name:          "Invalid URL - Missing track ID",
			url:           "https://music.yandex.ru/album/3192570/track/",
			wantID:        "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Invalid URL - Non-numeric track ID",
			url:           "https://music.yandex.ru/album/3192570/track/abc",
			wantID:        "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Invalid URL - Incorrect format",
			url:           "https://example.com/album/3192570/track/1197793",
			wantID:        "",
			expectedError: IDNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newYandexAdapter(nil, nil)
			result, err := adapter.DetectTrackID(tt.url)
			require.Equal(t, tt.wantID, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestYandexAdapter_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		wantID        string
		expectedError error
	}{
		{
			name:   "Valid album URL – .by",
			url:    "https://music.yandex.by/album/1197793",
			wantID: "1197793",
		},
		{
			name:   "Valid album URL – .kz",
			url:    "https://music.yandex.kz/album/1197793",
			wantID: "1197793",
		},
		{
			name:   "Valid album URL – .uz",
			url:    "https://music.yandex.uz/album/1197793",
			wantID: "1197793",
		},
		{
			name:   "Valid album URL – .ru",
			url:    "https://music.yandex.ru/album/1197793",
			wantID: "1197793",
		},
		{
			name:          "Invalid URL - Missing album ID",
			url:           "https://music.yandex.ru/album/",
			wantID:        "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Invalid URL - Non-numeric album ID",
			url:           "https://music.yandex.ru/album/letters",
			wantID:        "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Invalid URL - Incorrect host",
			url:           "https://example.com/album/3192570",
			wantID:        "",
			expectedError: IDNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newYandexAdapter(nil, nil)
			result, err := adapter.DetectAlbumID(tt.url)
			require.Equal(t, tt.wantID, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestYandexAdapter_GetTrack(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		yandexClientMock yandexClientMock
		expectedTrack    *Track
	}{
		{
			name: "found ID",
			id:   "42",
			yandexClientMock: yandexClientMock{
				getTrack: map[string]*yandex.Track{
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
			expectedTrack: &Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
			},
		},
		{
			name:             "not found ID",
			id:               "notFoundID",
			yandexClientMock: yandexClientMock{},
			expectedTrack:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, nil)

			result, err := a.GetTrack(tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
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
		expectedTrack    *Track
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
			expectedTrack: &Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
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
			expectedTrack: &Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист матчинг транслит",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
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
			expectedTrack: &Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист афтер транслит",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
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
			expectedTrack: &Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "переведенный артист",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, &tt.translatorMock)

			result, err := a.SearchTrack(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYandexAdapter_GetAlbum(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		yandexClientMock yandexClientMock
		expectedTrack    *Album
	}{
		{
			name: "found id",
			id:   "42",
			yandexClientMock: yandexClientMock{
				getAlbum: map[string]*yandex.Album{
					"42": {
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "sample artist"},
						},
					},
				},
			},
			expectedTrack: &Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
			},
		},
		{
			name:          "not found ID",
			id:            "notFoundID",
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, nil)

			result, err := a.GetAlbum(tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
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
		expectedAlbum    *Album
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
			expectedAlbum: &Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
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
			expectedAlbum: &Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист матчинг транслит",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
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
			expectedAlbum: &Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "сампле артист афтер транслит",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
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
			expectedAlbum: &Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "переведенный артист",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			expectedAlbum: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYandexAdapter(&tt.yandexClientMock, &tt.translatorMock)

			result, err := a.SearchAlbum(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedAlbum, result)
		})
	}
}
