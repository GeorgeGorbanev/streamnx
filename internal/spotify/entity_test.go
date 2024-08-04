package spotify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrack_URL(t *testing.T) {
	track := Track{ID: "sample_id"}
	result := track.URL()
	require.Equal(t, "https://open.spotify.com/track/sample_id", result)
}

func TestAlbum_URL(t *testing.T) {
	album := Album{ID: "sample_id"}
	result := album.URL()
	require.Equal(t, "https://open.spotify.com/album/sample_id", result)
}

func Test_DetectTrackID(t *testing.T) {
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
			name:     "Invalid URL - Entity",
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
			result := DetectTrackID(tt.inputURL)
			require.Equal(t, tt.expected, result)
		})
	}
}

func Test_DetectAlbumID(t *testing.T) {
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
			name:     "Valid URL with intl path",
			inputURL: "https://open.spotify.com/intl-pt/album/7uv632EkfwYhXoqf8rhYrg",
			expected: "7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:     "Invalid URL - Entity",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectAlbumID(tt.inputURL)
			require.Equal(t, tt.expected, result)
		})
	}
}
