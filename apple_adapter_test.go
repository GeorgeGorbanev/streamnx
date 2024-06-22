package streaminx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppleAdapter_DetectTrackID(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError error
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
			name:          "valid URL with track ID and invalid iso3611 storefront",
			input:         "https://music.apple.com/invalidstorefront/album/song-name/1234567890?i=987654321",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:     "valid URL without album",
			input:    "https://music.apple.com/us/song/angel/724466660",
			expected: "us-724466660",
		},
		{
			name:          "valid URL without album and invalid storefront",
			input:         "https://music.apple.com/invalidstorefront/song/angel/724466660",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "URL without track ID",
			input:         "https://music.apple.com/us/album/song-name/1234567890",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "invalid host URL",
			input:         "https://music.orange.com/us/album/song-name/1234567890?i=987654321",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "empty string",
			input:         "",
			expected:      "",
			expectedError: IDNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newAppleAdapter(nil)
			result, err := adapter.DetectTrackID(tt.input)
			require.Equal(t, tt.expected, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAppleAdapter_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError error
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
			name:          "valid URL with album ID and invalid iso3611 storefront",
			input:         "https://music.apple.com/invalidstorefront/album/another-album/987654321",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "URL without album ID",
			input:         "https://music.apple.com/us/album/album-name",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "invalid host URL",
			input:         "https://music.orange.com/us/album/album-name/123456789",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "empty string",
			input:         "",
			expected:      "",
			expectedError: IDNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newAppleAdapter(nil)
			result, err := adapter.DetectAlbumID(tt.input)
			require.Equal(t, tt.expected, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
