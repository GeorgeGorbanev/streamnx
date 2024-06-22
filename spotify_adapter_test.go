package streaminx

import (
	"testing"

	"github.com/GeorgeGorbanev/streaminx/internal/spotify"
	"github.com/stretchr/testify/require"
)

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

func TestSpotifyAdapter_DetectTrackID(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		expectedError error
	}{
		{
			name:     "Valid URL",
			inputURL: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:          "Invalid URL - Album",
			inputURL:      "https://open.spotify.com/album/3hARuIUZqAIAKSuNvW5dGh",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Empty URL",
			inputURL:      "",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Non-Spotify URL",
			inputURL:      "https://example.com/track/7uv632EkfwYhXoqf8rhYrg",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "URL without ID",
			inputURL:      "https://open.spotify.com/track/",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:     "Valid URL with query",
			inputURL: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?test=123",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:     "Valid URL with intl path",
			inputURL: "https://open.spotify.com/intl-pt/track/2xmQMKTjiOdkdGVgqDzezo",
			expected: "2xmQMKTjiOdkdGVgqDzezo",
		},
		{
			name:     "Valid URL with intl path and query",
			inputURL: "https://open.spotify.com/intl-pt/track/2xmQMKTjiOdkdGVgqDzezo?sample=query",
			expected: "2xmQMKTjiOdkdGVgqDzezo",
		},
		{
			name:     "Valid URL with prefix and suffix",
			inputURL: "prefix https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?test=123 suffix",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newSpotifyAdapter(nil)
			result, err := adapter.DetectTrackID(tt.inputURL)
			require.Equal(t, tt.expected, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSpotifyAdapter_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		expectedError error
	}{
		{
			name:     "Valid URL",
			inputURL: "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:     "Valid URL with intl path",
			inputURL: "https://open.spotify.com/intl-pt/album/7uv632EkfwYhXoqf8rhYrg",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:          "Invalid URL - Track",
			inputURL:      "https://open.spotify.com/track/3hARuIUZqAIAKSuNvW5dGh",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Empty URL",
			inputURL:      "",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Non-Spotify URL",
			inputURL:      "https://example.com/album/7uv632EkfwYhXoqf8rhYrg",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "URL without ID",
			inputURL:      "https://open.spotify.com/album/",
			expected:      "",
			expectedError: IDNotFoundError,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newSpotifyAdapter(nil)
			result, err := adapter.DetectAlbumID(tt.inputURL)
			require.Equal(t, tt.expected, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
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
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/track/sampleID",
				Provider: Spotify,
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
			a := newSpotifyAdapter(&spotifyClientMock{})

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
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/track/sampleID",
				Provider: Spotify,
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
			a := newSpotifyAdapter(&spotifyClientMock{})

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
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/album/sampleID",
				Provider: Spotify,
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
			a := newSpotifyAdapter(&spotifyClientMock{})

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
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/album/sampleID",
				Provider: Spotify,
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
			a := newSpotifyAdapter(&spotifyClientMock{})

			result, err := a.SearchAlbum(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}
