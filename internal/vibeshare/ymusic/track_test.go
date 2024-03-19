package ymusic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectTrackID(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		wantID string
	}{
		{
			name:   "Valid Track URL",
			url:    "https://music.yandex.ru/album/3192570/track/1197793",
			wantID: "1197793",
		},
		{
			name:   "Invalid URL - Missing track ID",
			url:    "https://music.yandex.ru/album/3192570/track/",
			wantID: "",
		},
		{
			name:   "Invalid URL - Non-numeric track ID",
			url:    "https://music.yandex.ru/album/3192570/track/abc",
			wantID: "",
		},
		{
			name:   "Invalid URL - Incorrect format",
			url:    "https://example.com/album/3192570/track/1197793",
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectTrackID(tt.url)
			require.Equal(t, tt.wantID, result)
		})
	}
}

func TestTrack_URL(t *testing.T) {
	tests := []struct {
		name  string
		track Track
		want  string
	}{
		{
			name: "track with int ID",
			track: Track{
				ID: 123,
				Albums: []Album{
					{ID: 456},
				},
			},
			want: "https://music.yandex.com/album/456/track/123",
		},
		{
			name: "track with string ID",
			track: Track{
				ID: "123",
				Albums: []Album{
					{ID: 456},
				},
			},
			want: "https://music.yandex.com/album/456/track/123",
		},
		{
			name: "track with float ID",
			track: Track{
				ID: 123.0,
				Albums: []Album{
					{ID: 456},
				},
			},
			want: "https://music.yandex.com/album/456/track/123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.track.URL()
			require.Equal(t, tt.want, result)
		})
	}
}
