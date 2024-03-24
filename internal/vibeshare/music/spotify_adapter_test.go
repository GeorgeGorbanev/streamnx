package music

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/stretchr/testify/require"
)

func TestSpotifyAdapter_DetectTrackID(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "Valid URL",
			inputURL: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:     "Invalid URL - Album",
			inputURL: "https://open.spotify.com/album/3hARuIUZqAIAKSuNvW5dGh",
			expected: "",
		},
		{
			name:     "Empty URL",
			inputURL: "",
			expected: "",
		},
		{
			name:     "Non-Spotify URL",
			inputURL: "https://example.com/track/7uv632EkfwYhXoqf8rhYrg",
			expected: "",
		},
		{
			name:     "URL without ID",
			inputURL: "https://open.spotify.com/track/",
			expected: "",
		},
		{
			name:     "Valid URL with query",
			inputURL: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?test=123",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:     "Valid URL with prefix and suffix",
			inputURL: "prefix https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?test=123 suffix",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			adapter := NewSpotifyAdapter(nil)
			result := adapter.DetectTrackID(tc.inputURL)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestSpotifyAdapter_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "Valid URL",
			inputURL: "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:     "Invalid URL - Track",
			inputURL: "https://open.spotify.com/track/3hARuIUZqAIAKSuNvW5dGh",
			expected: "",
		},
		{
			name:     "Empty URL",
			inputURL: "",
			expected: "",
		},
		{
			name:     "Non-Spotify URL",
			inputURL: "https://example.com/album/7uv632EkfwYhXoqf8rhYrg",
			expected: "",
		},
		{
			name:     "URL without ID",
			inputURL: "https://open.spotify.com/album/",
			expected: "",
		},
		{
			name:     "Valid URL with query",
			inputURL: "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg?test=123",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:     "Valid URL with prefix and suffix",
			inputURL: "prefix https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg?test=123 suffix",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			adapter := NewSpotifyAdapter(nil)
			result := adapter.DetectAlbumID(tc.inputURL)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestSpotifyAdapter_GetTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedTrack *Track
	}{
		{
			name: "found ID",
			id:   "sampleID",
			expectedTrack: &Track{
				Title:  "sample name",
				Artist: "sample artist",
				URL:    "https://open.spotify.com/track/sampleID",
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
			a := &SpotifyAdapter{
				client: &spotifyClientMock{},
			}

			result, err := a.GetTrack(tt.id)

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
		expectedTrack *Track
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			expectedTrack: &Track{
				Title:  "sample name",
				Artist: "sample artist",
				URL:    "https://open.spotify.com/track/sampleID",
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
			a := &SpotifyAdapter{
				client: &spotifyClientMock{},
			}

			result, err := a.SearchTrack(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestSpotifyAdapter_GetAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedTrack *Album
	}{
		{
			name: "found ID",
			id:   "sampleID",
			expectedTrack: &Album{
				Title:  "sample name",
				Artist: "sample artist",
				URL:    "https://open.spotify.com/album/sampleID",
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
			a := &SpotifyAdapter{
				client: &spotifyClientMock{},
			}

			result, err := a.GetAlbum(tt.id)

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
		expectedTrack *Album
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			expectedTrack: &Album{
				Title:  "sample name",
				Artist: "sample artist",
				URL:    "https://open.spotify.com/album/sampleID",
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
			a := &SpotifyAdapter{
				client: &spotifyClientMock{},
			}

			result, err := a.SearchAlbum(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

type spotifyClientMock struct{}

func (c *spotifyClientMock) GetTrack(id string) (*spotify.Track, error) {
	if id != "sampleID" {
		return nil, nil
	}
	return &spotify.Track{
		ID:   id,
		Name: "sample name",
		Artists: []spotify.Artist{
			{Name: "sample artist"},
		},
	}, nil
}

func (c *spotifyClientMock) SearchTrack(artistName, trackName string) (*spotify.Track, error) {
	if artistName != "sample artist" || trackName != "sample name" {
		return nil, nil
	}

	return &spotify.Track{
		ID:   "sampleID",
		Name: "sample name",
		Artists: []spotify.Artist{
			{Name: "sample artist"},
		},
	}, nil
}

func (c *spotifyClientMock) GetAlbum(id string) (*spotify.Album, error) {
	if id != "sampleID" {
		return nil, nil
	}
	return &spotify.Album{
		ID:   id,
		Name: "sample name",
		Artists: []spotify.Artist{
			{Name: "sample artist"},
		},
	}, nil
}

func (c *spotifyClientMock) SearchAlbum(artistName, albumName string) (*spotify.Album, error) {
	if artistName != "sample artist" || albumName != "sample name" {
		return nil, nil
	}
	return &spotify.Album{
		ID:   "sampleID",
		Name: "sample name",
		Artists: []spotify.Artist{
			{Name: "sample artist"},
		},
	}, nil
}
