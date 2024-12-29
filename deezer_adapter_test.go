package streamnx

import (
	"context"
	"testing"
	"time"

	"github.com/GeorgeGorbanev/streamnx/internal/deezer"
	"github.com/stretchr/testify/require"
)

type deezerClientMock struct {
	fetchTrack  map[string]*deezer.Track
	searchTrack map[string]*deezer.Track
	fetchAlbum  map[string]*deezer.Album
	searchAlbum map[string]*deezer.Album
}

func (c *deezerClientMock) FetchTrack(_ context.Context, id string) (*deezer.Track, error) {
	track, ok := c.fetchTrack[id]
	if !ok {
		return nil, deezer.NotFoundError
	}
	return track, nil
}

func (c *deezerClientMock) SearchTrack(_ context.Context, artistName, trackName string) (*deezer.Track, error) {
	if track, ok := c.searchTrack[artistName+"-"+trackName]; ok {
		return track, nil
	}
	return nil, deezer.NotFoundError
}

func (c *deezerClientMock) FetchAlbum(_ context.Context, id string) (*deezer.Album, error) {
	album, ok := c.fetchAlbum[id]
	if !ok {
		return nil, deezer.NotFoundError
	}
	return album, nil
}

func (c *deezerClientMock) SearchAlbum(_ context.Context, artistName, albumName string) (*deezer.Album, error) {
	if album, ok := c.searchAlbum[artistName+"-"+albumName]; ok {
		return album, nil
	}
	return nil, deezer.NotFoundError
}

func TestDeezerAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		clientMock    *deezerClientMock
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "sampleID",
			clientMock: &deezerClientMock{
				fetchTrack: map[string]*deezer.Track{
					"sampleID": {
						ID:     "sampleID",
						Title:  "sample name",
						Artist: "sample artist",
						URL:    "sampleURL",
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "sampleURL",
				Provider: Deezer,
				Type:     Track,
			},
		},
		{
			name:          "not found ID",
			id:            "notFoundID",
			clientMock:    &deezerClientMock{},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newDeezerAdapter(tt.clientMock)
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

func TestDeezerAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		clientMock    *deezerClientMock
		expectedTrack *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			clientMock: &deezerClientMock{
				searchTrack: map[string]*deezer.Track{
					"sample artist-sample name": {
						ID:     "sampleID",
						Title:  "sample name",
						Artist: "sample artist",
						URL:    "sampleURL",
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "sampleURL",
				Provider: Deezer,
				Type:     Track,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			clientMock:    &deezerClientMock{},
			expectedTrack: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newDeezerAdapter(tt.clientMock)
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

func TestDeezerAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		clientMock    *deezerClientMock
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "sampleID",
			clientMock: &deezerClientMock{
				fetchAlbum: map[string]*deezer.Album{
					"sampleID": {
						ID:     "sampleID",
						Title:  "sample name",
						Artist: "sample artist",
						URL:    "sampleURL",
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "sampleURL",
				Provider: Deezer,
				Type:     Album,
			},
		},
		{
			name:          "not found ID",
			id:            "notFoundID",
			clientMock:    &deezerClientMock{},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newDeezerAdapter(tt.clientMock)
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

func TestDeezerAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		clientMock    *deezerClientMock
		expectedAlbum *Entity
		expectedErr   error
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			clientMock: &deezerClientMock{
				searchAlbum: map[string]*deezer.Album{
					"sample artist-sample name": {
						ID:     "sampleID",
						Title:  "sample name",
						Artist: "sample artist",
						URL:    "sampleURL",
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "sampleURL",
				Provider: Deezer,
				Type:     Album,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			clientMock:    &deezerClientMock{},
			expectedAlbum: nil,
			expectedErr:   EntityNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newDeezerAdapter(tt.clientMock)
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
