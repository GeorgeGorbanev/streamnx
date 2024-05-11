package spotify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlbum_URL(t *testing.T) {
	album := Album{ID: "sample_id"}
	result := album.URL()
	require.Equal(t, "https://open.spotify.com/album/sample_id", result)
}
