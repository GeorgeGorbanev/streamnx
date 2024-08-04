package streaminx

import (
	"context"
	"testing"
	"time"

	"github.com/GeorgeGorbanev/streaminx/internal/spotify"

	"github.com/stretchr/testify/require"
)

type spotifyClientMock struct {
	fetchTrack  map[string]*spotify.Track
	fetchAlbum  map[string]*spotify.Album
	searchTrack map[string]map[string]*spotify.Track
	searchAlbum map[string]map[string]*spotify.Album
}

func (c *spotifyClientMock) FetchTrack(_ context.Context, id string) (*spotify.Track, error) {
	return c.fetchTrack[id], nil
}

func (c *spotifyClientMock) SearchTrack(_ context.Context, artistName, trackName string) (*spotify.Track, error) {
	if tracks, ok := c.searchTrack[artistName]; ok {
		return tracks[trackName], nil
	}
	return nil, nil
}

func (c *spotifyClientMock) FetchAlbum(_ context.Context, id string) (*spotify.Album, error) {
	return c.fetchAlbum[id], nil
}

func (c *spotifyClientMock) SearchAlbum(_ context.Context, artistName, albumName string) (*spotify.Album, error) {
	if albums, ok := c.searchAlbum[artistName]; ok {
		return albums[albumName], nil
	}
	return nil, nil
}

func TestSpotifyAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		clientMock    *spotifyClientMock
		expectedTrack *Entity
	}{
		{
			name: "found ID",
			id:   "sampleID",
			clientMock: &spotifyClientMock{
				fetchTrack: map[string]*spotify.Track{
					"sampleID": {
						ID:   "sampleID",
						Name: "sample name",
						Artists: []spotify.Artist{
							{
								Name: "sample artist",
							},
						},
					},
				},
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
			name:          "not found ID",
			id:            "notFoundID",
			clientMock:    &spotifyClientMock{},
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newSpotifyAdapter(tt.clientMock)
			result, err := a.FetchTrack(ctx, tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestSpotifyAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		clientMock    *spotifyClientMock
		expectedTrack *Entity
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			clientMock: &spotifyClientMock{
				searchTrack: map[string]map[string]*spotify.Track{
					"sample artist": {
						"sample name": {
							ID:   "sampleID",
							Name: "sample name",
							Artists: []spotify.Artist{
								{
									Name: "sample artist",
								},
							},
						},
					},
				},
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
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			clientMock:    &spotifyClientMock{},
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newSpotifyAdapter(tt.clientMock)
			result, err := a.SearchTrack(ctx, tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestSpotifyAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		clientMock    *spotifyClientMock
		expectedTrack *Entity
	}{
		{
			name: "found ID",
			id:   "sampleID",
			clientMock: &spotifyClientMock{
				fetchAlbum: map[string]*spotify.Album{
					"sampleID": {
						ID:   "sampleID",
						Name: "sample name",
						Artists: []spotify.Artist{
							{
								Name: "sample artist",
							},
						},
					},
				},
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
			name:          "not found ID",
			id:            "notFoundID",
			clientMock:    &spotifyClientMock{},
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newSpotifyAdapter(tt.clientMock)
			result, err := a.FetchAlbum(ctx, tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestSpotifyAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		clientMock    *spotifyClientMock
		expectedTrack *Entity
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			clientMock: &spotifyClientMock{
				searchAlbum: map[string]map[string]*spotify.Album{
					"sample artist": {
						"sample name": {
							ID:   "sampleID",
							Name: "sample name",
							Artists: []spotify.Artist{
								{
									Name: "sample artist",
								},
							},
						},
					},
				},
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
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			clientMock:    &spotifyClientMock{},
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a := newSpotifyAdapter(tt.clientMock)
			result, err := a.SearchAlbum(ctx, tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}
