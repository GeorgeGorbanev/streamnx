package spotify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectTrackID(t *testing.T) {
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
			result := DetectTrackID(tc.inputURL)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestTrack_Title(t *testing.T) {
	track := Track{
		Artists: []Artist{
			{Name: "Artist 1"},
			{Name: "Artist 2"},
			{Name: "Artist 3"},
		},
		Name: "Track Name",
	}

	result := track.Title()

	require.Equal(t, "Artist 1 – Track Name", result)
}
