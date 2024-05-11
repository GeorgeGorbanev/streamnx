package yandex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
