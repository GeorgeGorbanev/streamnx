package yandex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlbum_URL(t *testing.T) {
	album := Album{ID: 42}
	result := album.URL()
	require.Equal(t, "https://music.yandex.com/album/42", result)
}
