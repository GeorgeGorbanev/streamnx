package music

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppleAdapter_DetectTrackID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid URL with track ID",
			input:    "https://music.apple.com/us/album/song-name/1234567890?i=987654321",
			expected: "987654321",
		},
		{
			name:     "valid URL with track ID and th locale",
			input:    "https://music.apple.com/th/album/song-name/1234567890?i=987654321",
			expected: "987654321",
		},
		{
			name:     "valid URL without album",
			input:    "https://music.apple.com/us/song/angel/724466660",
			expected: "724466660",
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
			adapter := newAppleAdapter(nil)
			result := adapter.DetectTrackID(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAppleAdapter_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid URL with album ID",
			input:    "https://music.apple.com/us/album/album-name/123456789",
			expected: "123456789",
		},
		{
			name:     "valid URL with album ID and gb locale",
			input:    "https://music.apple.com/gb/album/another-album/987654321",
			expected: "987654321",
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
			adapter := newAppleAdapter(nil)
			result := adapter.DetectAlbumID(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
