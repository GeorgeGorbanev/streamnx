package streaminx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry_Adapter(t *testing.T) {
	spotifyAdapter := newSpotifyAdapter(nil)
	yandexAdapter := newYandexAdapter(nil, nil)
	youtubeAdapter := newYoutubeAdapter(nil)

	registry := Registry{
		adapters: map[string]Adapter{
			Spotify.Code: spotifyAdapter,
			Yandex.Code:  yandexAdapter,
			Youtube.Code: youtubeAdapter,
		},
	}

	require.Equal(t, registry.Adapter(Spotify), spotifyAdapter)
	require.Equal(t, registry.Adapter(Yandex), yandexAdapter)
	require.Equal(t, registry.Adapter(Youtube), youtubeAdapter)
}
