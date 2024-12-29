package deezer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrack(t *testing.T) {
	track := Track{
		ID:     "123",
		Title:  "Test Track",
		Artist: "Test Artist",
		URL:    "https://www.deezer.com/track/123",
	}

	require.Equal(t, "123", track.ID)
	require.Equal(t, "Test Track", track.Title)
	require.Equal(t, "Test Artist", track.Artist)
	require.Equal(t, "https://www.deezer.com/track/123", track.URL)
}

func TestAlbum(t *testing.T) {
	album := Album{
		ID:     "456",
		Title:  "Test Album",
		Artist: "Test Artist",
		URL:    "https://www.deezer.com/album/456",
	}

	require.Equal(t, "456", album.ID)
	require.Equal(t, "Test Album", album.Title)
	require.Equal(t, "Test Artist", album.Artist)
	require.Equal(t, "https://www.deezer.com/album/456", album.URL)
}
