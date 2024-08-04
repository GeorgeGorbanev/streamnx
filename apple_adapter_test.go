package streaminx

import (
	"context"
	"testing"
	"time"

	"github.com/GeorgeGorbanev/streaminx/internal/apple"

	"github.com/stretchr/testify/require"
)

type appleClientMock struct {
	fetchTrack  map[string]*apple.Entity
	fetchAlbum  map[string]*apple.Entity
	searchTrack map[string]map[string]*apple.Entity
	searchAlbum map[string]map[string]*apple.Entity
}

func (c *appleClientMock) FetchTrack(_ context.Context, id, storefront string) (*apple.Entity, error) {
	track, ok := c.fetchTrack[storefront+"-"+id]
	if !ok {
		return nil, apple.NotFoundError
	}
	return track, nil
}

func (c *appleClientMock) SearchTrack(_ context.Context, artistName, trackName string) (*apple.Entity, error) {
	if tracks, ok := c.searchTrack[artistName]; ok {
		track, ok := tracks[trackName]
		if !ok {
			return nil, apple.NotFoundError
		}
		return track, nil
	}
	return nil, apple.NotFoundError
}

func (c *appleClientMock) FetchAlbum(_ context.Context, id, storefront string) (*apple.Entity, error) {
	album, ok := c.fetchAlbum[storefront+"-"+id]
	if !ok {
		return nil, apple.NotFoundError
	}
	return album, nil
}

func (c *appleClientMock) SearchAlbum(_ context.Context, artistName, albumName string) (*apple.Entity, error) {
	if albums, ok := c.searchAlbum[artistName]; ok {
		album, ok := albums[albumName]
		if !ok {
			return nil, apple.NotFoundError
		}
		return album, nil
	}
	return nil, apple.NotFoundError
}

func TestAppleAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		clientMock    *appleClientMock
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "ru-123",
			clientMock: &appleClientMock{
				fetchTrack: map[string]*apple.Entity{
					"ru-123": {
						ID: "ru-123",
						Attributes: apple.Attributes{
							ArtistName: "sample artist",
							Name:       "sample name",
							URL:        "https://music.apple.com/ru/album/song-name/1234567890?i=123",
						},
					},
				},
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
			name:          "not found ID",
			id:            "ru-123",
			clientMock:    &appleClientMock{},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newAppleAdapter(tt.clientMock)
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

func TestAppleAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		clientMock    *appleClientMock
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			clientMock: &appleClientMock{
				searchTrack: map[string]map[string]*apple.Entity{
					"sample artist": {
						"sample name": {
							ID: "ru-123",
							Attributes: apple.Attributes{
								ArtistName: "sample artist",
								Name:       "sample name",
								URL:        "https://music.apple.com/ru/album/song-name/1234567890?i=123",
							},
						},
					},
				},
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
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			clientMock:    &appleClientMock{},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newAppleAdapter(tt.clientMock)
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

func TestAppleAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		storefront    string
		clientMock    *appleClientMock
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name:       "found ID",
			id:         "ru-456",
			storefront: "sampleStorefront",
			clientMock: &appleClientMock{
				fetchAlbum: map[string]*apple.Entity{
					"ru-456": {
						ID: "ru-456",
						Attributes: apple.Attributes{
							ArtistName: "sample artist",
							Name:       "sample name",
							URL:        "https://music.apple.com/ru/album/name/456",
						},
					},
				},
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
			name:          "not found ID",
			id:            "ru-456",
			storefront:    "notFoundStorefront",
			clientMock:    &appleClientMock{},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newAppleAdapter(tt.clientMock)
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

func TestAppleAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		clientMock    *appleClientMock
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			clientMock: &appleClientMock{
				searchAlbum: map[string]map[string]*apple.Entity{
					"sample artist": {
						"sample name": {
							ID: "ru-456",
							Attributes: apple.Attributes{
								ArtistName: "sample artist",
								Name:       "sample name",
								URL:        "https://music.apple.com/ru/album/name/456",
							},
						},
					},
				},
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
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			clientMock:    &appleClientMock{},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newAppleAdapter(tt.clientMock)
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
