package streamnx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/translator"
	"github.com/GeorgeGorbanev/streamnx/internal/yandex"
)

func TestYandexAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *yandex.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "42",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("FetchTrack", "42").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "sample artist"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
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
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("FetchTrack", "notFoundID").
					Return(nil, yandex.NotFoundError).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &yandex.ClientMock{}
			tt.mockClient(clientMock)

			a := newYandexAdapter(clientMock, nil)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.FetchTrack(ctx, tt.id)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestYandexAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *yandex.ClientMock)
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name: "found id",
			id:   "42",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("FetchAlbum", "42").
					Return(&yandex.Album{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "sample artist"},
						},
					}, nil).
					Once()
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
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("FetchAlbum", "notFoundID").
					Return(nil, yandex.NotFoundError).
					Once()
			},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &yandex.ClientMock{}
			tt.mockClient(clientMock)

			a := newYandexAdapter(clientMock, nil)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.FetchAlbum(ctx, tt.id)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestYandexAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		searchName     string
		mockClient     func(m *yandex.ClientMock)
		mockTranslator func(m *translator.Mock)
		expectedTrack  *Entity
		expectedErr    error
	}{
		{
			name:       "found query",
			artist:     "sample artist",
			searchName: "sample name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "sample artist – sample name").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "sample artist"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
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
			artist:     "sample artist not matching",
			searchName: "sample name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "sample artist not matching – sample name").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "not matching artist"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedTrack:  nil,
			expectedErr:    EntityNotFoundError,
		},
		{
			name:       "found query matching translit",
			artist:     "sample artist matching translit",
			searchName: "sample name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "sample artist matching translit – sample name").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "сампле артист матчинг транслит"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
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
			artist:     "sample artist after translit",
			searchName: "кириллическое название",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "sample artist after translit – кириллическое название").
					Return(nil, yandex.NotFoundError).
					Once()
				m.
					On("SearchTrack", "сампле артист афтер транслит – кириллическое название").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "сампле артист афтер транслит"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
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
			artist:     "translatable artist",
			searchName: "sample name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "translatable artist – sample name").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "переведенный артист"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {
				m.
					On("TranslateEnToRu", "translatable artist").
					Return("переведенный артист", nil).
					Once()
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
			name:       "not found query",
			artist:     "not found artist",
			searchName: "not found name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "not found artist – not found name").
					Return(nil, yandex.NotFoundError).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedTrack:  nil,
			expectedErr:    EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &yandex.ClientMock{}
			tt.mockClient(clientMock)

			translatorMock := &translator.Mock{}
			tt.mockTranslator(translatorMock)

			a := newYandexAdapter(clientMock, translatorMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.SearchTrack(ctx, tt.artist, tt.searchName)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			clientMock.AssertExpectations(t)
			translatorMock.AssertExpectations(t)
		})
	}
}

func TestYandexAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name           string
		artistName     string
		searchName     string
		mockClient     func(m *yandex.ClientMock)
		mockTranslator func(m *translator.Mock)
		expectedAlbum  *Entity
		expectedErr    error
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "sample artist – sample name").
					Return(&yandex.Album{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "sample artist"},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
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
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "sample artist not matching – sample name").
					Return(&yandex.Album{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "not matching artist"},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedAlbum:  nil,
			expectedErr:    EntityNotFoundError,
		},
		{
			name:       "found query matching translit",
			artistName: "sample artist matching translit",
			searchName: "sample name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "sample artist matching translit – sample name").
					Return(&yandex.Album{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "сампле артист матчинг транслит"},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
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
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "sample artist after translit – кириллическое название").
					Return(nil, yandex.NotFoundError).
					Once()
				m.
					On("SearchAlbum", "сампле артист афтер транслит – кириллическое название").
					Return(&yandex.Album{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "сампле артист афтер транслит"},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
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
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "translatable artist – sample name").
					Return(&yandex.Album{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "переведенный артист"},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {
				m.
					On("TranslateEnToRu", "translatable artist").
					Return("переведенный артист", nil).
					Once()
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
			name:       "not found query",
			artistName: "not found artist",
			searchName: "not found name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "not found artist – not found name").
					Return(nil, yandex.NotFoundError).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedAlbum:  nil,
			expectedErr:    EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &yandex.ClientMock{}
			tt.mockClient(clientMock)

			translatorMock := &translator.Mock{}
			tt.mockTranslator(translatorMock)

			a := newYandexAdapter(clientMock, translatorMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.SearchAlbum(ctx, tt.artistName, tt.searchName)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}

			clientMock.AssertExpectations(t)
			translatorMock.AssertExpectations(t)
		})
	}
}
