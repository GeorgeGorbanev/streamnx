package spotify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectAlbumID(t *testing.T) {
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
			result := DetectAlbumID(tc.inputURL)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestAlbum_URL(t *testing.T) {
	album := Album{ID: "sample_id"}
	result := album.URL()
	require.Equal(t, "https://open.spotify.com/album/sample_id", result)
}
