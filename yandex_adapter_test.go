package streamnx

import (
	"context"
	"errors"
	"testing"

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

			result, err := a.FetchTrack(t.Context(), tt.id)

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

			result, err := a.FetchAlbum(t.Context(), tt.id)

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

			result, err := a.SearchTrack(t.Context(), tt.artist, tt.searchName)

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

			result, err := a.SearchAlbum(t.Context(), tt.artistName, tt.searchName)

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

var errTestServer = errors.New("server error")

func TestYandexAdapter_SearchTrack_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		searchName     string
		mockClient     func(m *yandex.ClientMock)
		mockTranslator func(m *translator.Mock)
		expectedErr    error
	}{
		{
			name:       "first search returns non-NotFoundError, should propagate",
			artist:     "test artist",
			searchName: "test name",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "test artist – test name").
					Return(nil, errTestServer).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedErr:    errTestServer,
		},
		{
			name:       "transliterated search returns non-NotFoundError, should propagate",
			artist:     "test artist",
			searchName: "кириллическое название",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "test artist – кириллическое название").
					Return(nil, yandex.NotFoundError).
					Once()
				m.
					On("SearchTrack", "тест артист – кириллическое название").
					Return(nil, errTestServer).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedErr:    errTestServer,
		},
		{
			name:       "transliterated search finds track but artist doesn't match",
			artist:     "test artist",
			searchName: "кириллическое название",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "test artist – кириллическое название").
					Return(nil, yandex.NotFoundError).
					Once()
				m.
					On("SearchTrack", "тест артист – кириллическое название").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "неподходящий артист"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {
				m.
					On("TranslateEnToRu", "test artist").
					Return("тест артист", nil).
					Once()
			},
			expectedErr: EntityNotFoundError,
		},
		{
			name:       "artist match fails during transliterated search",
			artist:     "test artist",
			searchName: "кириллическое название",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchTrack", "test artist – кириллическое название").
					Return(nil, yandex.NotFoundError).
					Once()
				m.
					On("SearchTrack", "тест артист – кириллическое название").
					Return(&yandex.Track{
						ID:    42,
						Title: "sample name",
						Artists: []yandex.Artist{
							{Name: "подходящий артист"},
						},
						Albums: []yandex.Album{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {
				m.
					On("TranslateEnToRu", "test artist").
					Return("", errTestServer).
					Once()
			},
			expectedErr: errTestServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			clientMock := &yandex.ClientMock{}
			tt.mockClient(clientMock)

			translatorMock := &translator.Mock{}
			tt.mockTranslator(translatorMock)

			adapter := newYandexAdapter(clientMock, translatorMock)

			result, err := adapter.SearchTrack(ctx, tt.artist, tt.searchName)

			require.Nil(t, result)
			require.Error(t, err)
			require.ErrorIs(t, err, tt.expectedErr)

			clientMock.AssertExpectations(t)
			translatorMock.AssertExpectations(t)
		})
	}
}

func TestYandexAdapter_SearchAlbum_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		searchName     string
		mockClient     func(m *yandex.ClientMock)
		mockTranslator func(m *translator.Mock)
		expectedErr    error
	}{
		{
			name:       "first search returns non-NotFoundError, should propagate",
			artist:     "test artist",
			searchName: "test album",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "test artist – test album").
					Return(nil, errTestServer).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedErr:    errTestServer,
		},
		{
			name:       "transliterated search returns non-NotFoundError, should propagate",
			artist:     "test artist",
			searchName: "кириллическое название",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "test artist – кириллическое название").
					Return(nil, yandex.NotFoundError).
					Once()
				m.
					On("SearchAlbum", "тест артист – кириллическое название").
					Return(nil, errTestServer).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {},
			expectedErr:    errTestServer,
		},
		{
			name:       "transliterated search finds album but artist doesn't match",
			artist:     "test artist",
			searchName: "кириллическое название",
			mockClient: func(m *yandex.ClientMock) {
				m.
					On("SearchAlbum", "test artist – кириллическое название").
					Return(nil, yandex.NotFoundError).
					Once()
				m.
					On("SearchAlbum", "тест артист – кириллическое название").
					Return(&yandex.Album{
						ID:    42,
						Title: "sample album",
						Artists: []yandex.Artist{
							{Name: "неподходящий артист"},
						},
					}, nil).
					Once()
			},
			mockTranslator: func(m *translator.Mock) {
				m.
					On("TranslateEnToRu", "test artist").
					Return("тест артист", nil).
					Once()
			},
			expectedErr: EntityNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			clientMock := &yandex.ClientMock{}
			tt.mockClient(clientMock)

			translatorMock := &translator.Mock{}
			tt.mockTranslator(translatorMock)

			adapter := newYandexAdapter(clientMock, translatorMock)

			result, err := adapter.SearchAlbum(ctx, tt.artist, tt.searchName)

			require.Nil(t, result)
			require.Error(t, err)
			require.ErrorIs(t, err, tt.expectedErr)

			clientMock.AssertExpectations(t)
			translatorMock.AssertExpectations(t)
		})
	}
}
