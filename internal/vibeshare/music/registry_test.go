package music

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry_Adapter(t *testing.T) {
	spotifyAdapter := newSpotifyAdapter(nil)
	yandexAdapter := newYandexAdapter(nil)
	youtubeAdapter := newYoutubeAdapter(nil)

	registry := Registry{
		adapters: map[Provider]Adapter{
			Spotify: spotifyAdapter,
			Yandex:  yandexAdapter,
			Youtube: youtubeAdapter,
		},
	}

	require.Equal(t, registry.Adapter(Spotify), spotifyAdapter)
	require.Equal(t, registry.Adapter(Yandex), yandexAdapter)
	require.Equal(t, registry.Adapter(Youtube), youtubeAdapter)
}
