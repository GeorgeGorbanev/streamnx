package bandcamp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid URL with album ID",
			input:    "https://autechre.bandcamp.com/album/amber",
			expected: "autechre:amber",
		},
		{
			name:     "URL without album ID",
			input:    "https://autechre.bandcamp.com/album/",
			expected: "",
		},
		{
			name:     "invalid host URL",
			input:    "https://autechre.danbcamp.com/album/amber",
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

func Test_DetectTrackID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid URL with track ID",
			input:    "https://autechre.bandcamp.com/track/nil",
			expected: "autechre:nil",
		},
		{
			name:     "URL without track ID",
			input:    "https://autechre.bandcamp.com/track/",
			expected: "",
		},
		{
			name:     "invalid host URL",
			input:    "https://autechre.danbcamp.com/track/nil",
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
