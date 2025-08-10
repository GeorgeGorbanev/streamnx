package streamnx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/spotify"
)

func TestSpotifyAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *spotify.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "sampleID",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("FetchTrack", "sampleID").
					Return(&spotify.Track{
						ID:   "sampleID",
						Name: "sample name",
						Artists: []spotify.Artist{
							{
								Name: "sample artist",
							},
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/track/sampleID",
				Provider: Spotify,
				Type:     Track,
			},
		},
		{
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("FetchTrack", "notFoundID").
					Return(nil, spotify.NotFoundError).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			clientMock := &spotify.ClientMock{}
			tt.mockClient(clientMock)

			a := newSpotifyAdapter(clientMock)
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

func TestSpotifyAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *spotify.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artist:     "sample artist",
			searchName: "sample name",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("SearchTrack", "sample artist", "sample name").
					Return(&spotify.Track{
						ID:   "sampleID",
						Name: "sample name",
						Artists: []spotify.Artist{
							{
								Name: "sample artist",
							},
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/track/sampleID",
				Provider: Spotify,
				Type:     Track,
			},
		},
		{
			name:       "not found query",
			artist:     "not found artist",
			searchName: "not found name",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("SearchTrack", "not found artist", "not found name").
					Return(nil, spotify.NotFoundError).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			clientMock := &spotify.ClientMock{}
			tt.mockClient(clientMock)

			a := newSpotifyAdapter(clientMock)
			result, err := a.SearchTrack(ctx, tt.artist, tt.searchName)

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

func TestSpotifyAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *spotify.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "sampleID",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("FetchAlbum", "sampleID").
					Return(&spotify.Album{
						ID:   "sampleID",
						Name: "sample name",
						Artists: []spotify.Artist{
							{
								Name: "sample artist",
							},
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/album/sampleID",
				Provider: Spotify,
				Type:     Album,
			},
		},
		{
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("FetchAlbum", "notFoundID").
					Return(nil, spotify.NotFoundError).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			clientMock := &spotify.ClientMock{}
			tt.mockClient(clientMock)

			a := newSpotifyAdapter(clientMock)
			result, err := a.FetchAlbum(ctx, tt.id)

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

func TestSpotifyAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *spotify.ClientMock)
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artist:     "sample artist",
			searchName: "sample name",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("SearchAlbum", "sample artist", "sample name").
					Return(&spotify.Album{
						ID:   "sampleID",
						Name: "sample name",
						Artists: []spotify.Artist{
							{
								Name: "sample artist",
							},
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/album/sampleID",
				Provider: Spotify,
				Type:     Album,
			},
		},
		{
			name:       "not found query",
			artist:     "not found artist",
			searchName: "not found name",
			mockClient: func(m *spotify.ClientMock) {
				m.
					On("SearchAlbum", "not found artist", "not found name").
					Return(nil, spotify.NotFoundError).
					Once()
			},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			clientMock := &spotify.ClientMock{}
			tt.mockClient(clientMock)

			a := newSpotifyAdapter(clientMock)
			result, err := a.SearchAlbum(ctx, tt.artist, tt.searchName)

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
