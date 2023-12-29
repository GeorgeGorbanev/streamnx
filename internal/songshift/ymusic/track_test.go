package ymusic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrack_URL(t *testing.T) {
	track := Track{
		ID: 123,
		Albums: []Album{
			{ID: 456},
		},
	}
	require.Equal(t, "https://music.yandex.com/album/456/track/123", track.URL())

}
