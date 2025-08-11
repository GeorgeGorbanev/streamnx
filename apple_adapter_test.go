package streamnx

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/apple"
)

func TestAppleAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *apple.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "ru-123",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("FetchTrack", "123", "ru").
					Return(&apple.Entity{
						ID: "ru-123",
						Attributes: apple.Attributes{
							ArtistName: "sample artist",
							Name:       "sample name",
							URL:        "https://music.apple.com/ru/album/song-name/1234567890?i=123",
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "ru-123",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.apple.com/ru/album/song-name/1234567890?i=123",
				Provider: Apple,
				Type:     Track,
			},
		},
		{
			name: "not found ID",
			id:   "ru-123",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("FetchTrack", "123", "ru").
					Return(nil, apple.NotFoundError).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &apple.ClientMock{}
			tt.mockClient(clientMock)

			a := newAppleAdapter(clientMock)
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

func TestAppleAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *apple.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artist:     "sample artist",
			searchName: "sample name",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("SearchTrack", "sample artist", "sample name").
					Return(&apple.Entity{
						ID: "ru-123",
						Attributes: apple.Attributes{
							ArtistName: "sample artist",
							Name:       "sample name",
							URL:        "https://music.apple.com/ru/album/song-name/1234567890?i=123",
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "ru-123",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.apple.com/ru/album/song-name/1234567890?i=123",
				Provider: Apple,
				Type:     Track,
			},
		},
		{
			name:       "not found query",
			artist:     "not found artist",
			searchName: "not found name",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("SearchTrack", "not found artist", "not found name").
					Return(nil, apple.NotFoundError).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &apple.ClientMock{}
			tt.mockClient(clientMock)

			a := newAppleAdapter(clientMock)
			result, err := a.SearchTrack(t.Context(), tt.artist, tt.searchName)

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

func TestAppleAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		storefront    string
		mockClient    func(m *apple.ClientMock)
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name:       "found ID",
			id:         "ru-456",
			storefront: "sampleStorefront",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("FetchAlbum", "456", "ru").
					Return(&apple.Entity{
						ID: "ru-456",
						Attributes: apple.Attributes{
							ArtistName: "sample artist",
							Name:       "sample name",
							URL:        "https://music.apple.com/ru/album/name/456",
						},
					}, nil).
					Once()
			},
			expectedAlbum: &Entity{
				ID:       "ru-456",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.apple.com/ru/album/name/456",
				Provider: Apple,
				Type:     Album,
			},
		},
		{
			name:       "not found ID",
			id:         "ru-456",
			storefront: "notFoundStorefront",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("FetchAlbum", "456", "ru").
					Return(nil, apple.NotFoundError).
					Once()
			},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &apple.ClientMock{}
			tt.mockClient(clientMock)

			a := newAppleAdapter(clientMock)
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

func TestAppleAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *apple.ClientMock)
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artist:     "sample artist",
			searchName: "sample name",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("SearchAlbum", "sample artist", "sample name").
					Return(&apple.Entity{
						ID: "ru-456",
						Attributes: apple.Attributes{
							ArtistName: "sample artist",
							Name:       "sample name",
							URL:        "https://music.apple.com/ru/album/name/456",
						},
					}, nil).
					Once()
			},
			expectedAlbum: &Entity{
				ID:       "ru-456",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.apple.com/ru/album/name/456",
				Provider: Apple,
				Type:     Album,
			},
		},
		{
			name:       "not found query",
			artist:     "not found artist",
			searchName: "not found name",
			mockClient: func(m *apple.ClientMock) {
				m.
					On("SearchAlbum", "not found artist", "not found name").
					Return(nil, apple.NotFoundError).
					Once()
			},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &apple.ClientMock{}
			tt.mockClient(clientMock)

			a := newAppleAdapter(clientMock)
			result, err := a.SearchAlbum(t.Context(), tt.artist, tt.searchName)

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
