package ymusic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectAlbumID(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		wantID string
	}{
		{
			name:   "Valid album URL",
			url:    "https://music.yandex.ru/album/1197793",
			wantID: "1197793",
		},
		{
			name:   "Invalid URL - Missing album ID",
			url:    "https://music.yandex.ru/album/",
			wantID: "",
		},
		{
			name:   "Invalid URL - Non-numeric album ID",
			url:    "https://music.yandex.ru/album/letters",
			wantID: "",
		},
		{
			name:   "Invalid URL - Incorrect host",
			url:    "https://example.com/album/3192570",
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectAlbumID(tt.url)
			require.Equal(t, tt.wantID, result)
		})
	}
}

func TestAlbum_URL(t *testing.T) {
	album := Album{ID: 42}
	result := album.URL()
	require.Equal(t, "https://music.yandex.com/album/42", result)
}
