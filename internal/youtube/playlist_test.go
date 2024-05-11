package youtube

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlaylist_URL(t *testing.T) {
	playlist := Playlist{ID: "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj"}
	result := playlist.URL()
	require.Equal(t, "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj", result)
}
