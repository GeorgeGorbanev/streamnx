package apple

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DetectTrackID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid URL with track ID",
			input:    "https://music.apple.com/us/album/song-name/1234567890?i=987654321",
			expected: "us-987654321",
		},
		{
			name:     "valid URL with track ID and th storefront",
			input:    "https://music.apple.com/th/album/song-name/1234567890?i=987654321",
			expected: "th-987654321",
		},
		{
			name:     "valid URL with track ID and invalid iso3611 storefront",
			input:    "https://music.apple.com/invalidstorefront/album/song-name/1234567890?i=987654321",
			expected: "",
		},
		{
			name:     "valid URL without album",
			input:    "https://music.apple.com/us/song/angel/724466660",
			expected: "us-724466660",
		},
		{
			name:     "valid URL without album and invalid storefront",
			input:    "https://music.apple.com/invalidstorefront/song/angel/724466660",
			expected: "",
		},
		{
			name:     "URL without track ID",
			input:    "https://music.apple.com/us/album/song-name/1234567890",
			expected: "",
		},
		{
			name:     "invalid host URL",
			input:    "https://music.orange.com/us/album/song-name/1234567890?i=987654321",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectTrackID(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func Test_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid URL with album ID",
			input:    "https://music.apple.com/us/album/album-name/123456789",
			expected: "us-123456789",
		},
		{
			name:     "valid URL with album ID and gb locale",
			input:    "https://music.apple.com/gb/album/another-album/987654321",
			expected: "gb-987654321",
		},
		{
			name:     "valid URL with album ID and invalid iso3611 storefront",
			input:    "https://music.apple.com/invalidstorefront/album/another-album/987654321",
			expected: "",
		},
		{
			name:     "URL without album ID",
			input:    "https://music.apple.com/us/album/album-name",
			expected: "",
		},
		{
			name:     "invalid host URL",
			input:    "https://music.orange.com/us/album/album-name/123456789",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectAlbumID(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
