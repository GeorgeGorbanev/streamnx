package streamnx

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/bandcamp"
)

func TestBandcampAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *bandcamp.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "autechre:nil",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("FetchTrack", "autechre", "nil").
					Return(&bandcamp.Entity{
						Name:     "Nil",
						BandName: "Autechre",
						URL:      "https://autechre.bandcamp.com/track/nil",
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "autechre:nil",
				Title:    "Nil",
				Artist:   "Autechre",
				URL:      "https://autechre.bandcamp.com/track/nil",
				Provider: Bandcamp,
				Type:     Track,
			},
		},
		{
			name: "not found ID",
			id:   "notfound:notfound",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("FetchTrack", "notfound", "notfound").
					Return(nil, bandcamp.ErrNotFound).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &bandcamp.ClientMock{}
			tt.mockClient(clientMock)

			a := newBandcampAdapter(clientMock)
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

func TestBandcampAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *bandcamp.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artist:     "autechre",
			searchName: "nil",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("SearchTrack", "autechre", "nil").
					Return(&bandcamp.Entity{
						Name:     "Nil",
						BandName: "Autechre",
						URL:      "https://autechre.bandcamp.com/track/nil",
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "autechre:nil",
				Title:    "Nil",
				Artist:   "Autechre",
				URL:      "https://autechre.bandcamp.com/track/nil",
				Provider: Bandcamp,
				Type:     Track,
			},
		},
		{
			name:       "not found query",
			artist:     "not found artist",
			searchName: "not found name",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("SearchTrack", "not found artist", "not found name").
					Return(nil, bandcamp.ErrNotFound).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &bandcamp.ClientMock{}
			tt.mockClient(clientMock)

			a := newBandcampAdapter(clientMock)
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

func TestBandcampAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *bandcamp.ClientMock)
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "autechre:amber",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("FetchAlbum", "autechre", "amber").
					Return(&bandcamp.Entity{
						Name:     "Amber",
						BandName: "Autechre",
						URL:      "https://autechre.bandcamp.com/album/amber",
					}, nil).
					Once()
			},
			expectedAlbum: &Entity{
				ID:       "autechre:amber",
				Title:    "Amber",
				Artist:   "Autechre",
				URL:      "https://autechre.bandcamp.com/album/amber",
				Provider: Bandcamp,
				Type:     Album,
			},
		},
		{
			name: "not found ID",
			id:   "notfound:notfound",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("FetchAlbum", "notfound", "notfound").
					Return(nil, bandcamp.ErrNotFound).
					Once()
			},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &bandcamp.ClientMock{}
			tt.mockClient(clientMock)

			a := newBandcampAdapter(clientMock)
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

func TestBandcampAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *bandcamp.ClientMock)
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artist:     "autechre",
			searchName: "amber",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("SearchAlbum", "autechre", "amber").
					Return(&bandcamp.Entity{
						Name:     "Amber",
						BandName: "Autechre",
						URL:      "https://autechre.bandcamp.com/album/amber",
					}, nil).
					Once()
			},
			expectedAlbum: &Entity{
				ID:       "autechre:amber",
				Title:    "Amber",
				Artist:   "Autechre",
				URL:      "https://autechre.bandcamp.com/album/amber",
				Provider: Bandcamp,
				Type:     Album,
			},
		},
		{
			name:       "not found query",
			artist:     "not found artist",
			searchName: "not found name",
			mockClient: func(m *bandcamp.ClientMock) {
				m.
					On("SearchAlbum", "not found artist", "not found name").
					Return(nil, bandcamp.ErrNotFound).
					Once()
			},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &bandcamp.ClientMock{}
			tt.mockClient(clientMock)

			a := newBandcampAdapter(clientMock)
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

func TestBandcampAdapter_FetchCloak(t *testing.T) {
	a := newBandcampAdapter(nil)
	_, err := a.FetchCloak(t.Context(), "nevermind")
	require.ErrorIs(t, err, UnsupportedEntityTypeError)
}
