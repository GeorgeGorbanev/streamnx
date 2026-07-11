package streamnx

import (
	"net/http"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/apple"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/bandcamp"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/deezer"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/soundcloud"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/spotify"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/yandex"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/youtube"
)

type CatalogOption func(catalog *Catalog) error

type AppleOption func(*appleOptions)

type appleOptions struct {
	clientOptions []apple.ClientOption
}

func WithApple(opts ...AppleOption) CatalogOption {
	return func(r *Catalog) error {
		options := appleOptions{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&options)
		}
		client := apple.NewClient(options.clientOptions...)
		adapter := apple.NewAdapter(client)
		return r.register(Apple, adapter)
	}
}

func WithAppleWebPlayerURL(url string) AppleOption {
	return func(cfg *appleOptions) {
		cfg.clientOptions = append(cfg.clientOptions, apple.WithWebPlayerURL(url))
	}
}

func WithAppleAPIURL(url string) AppleOption {
	return func(cfg *appleOptions) {
		cfg.clientOptions = append(cfg.clientOptions, apple.WithAPIURL(url))
	}
}

func WithAppleHTTPClient(client *http.Client) AppleOption {
	return func(cfg *appleOptions) {
		if client == nil {
			return
		}
		cfg.clientOptions = append(cfg.clientOptions, apple.WithHTTPClient(client))
	}
}

type BandcampOption func(*bandcampOptions)

type bandcampOptions struct {
	client []bandcamp.ClientOption
}

func WithBandcamp(opts ...BandcampOption) CatalogOption {
	return func(r *Catalog) error {
		cfg := bandcampOptions{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&cfg)
		}
		return r.register(Bandcamp, bandcamp.NewAdapter(bandcamp.NewClient(cfg.client...)))
	}
}

func WithBandcampAPI(scheme, host string) BandcampOption {
	return func(cfg *bandcampOptions) {
		cfg.client = append(cfg.client, bandcamp.WithAPI(scheme, host))
	}
}

func WithBandcampHTTPClient(client *http.Client) BandcampOption {
	return func(cfg *bandcampOptions) {
		if client == nil {
			return
		}
		cfg.client = append(cfg.client, bandcamp.WithAPIClient(client))
	}
}

type DeezerOption func(*deezerOptions)

type deezerOptions struct {
	client []deezer.ClientOption
}

func WithDeezer(opts ...DeezerOption) CatalogOption {
	return func(r *Catalog) error {
		cfg := deezerOptions{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&cfg)
		}
		return r.register(Deezer, deezer.NewAdapter(deezer.NewClient(cfg.client...)))
	}
}

func WithDeezerAPIURL(url string) DeezerOption {
	return func(cfg *deezerOptions) {
		cfg.client = append(cfg.client, deezer.WithAPIURL(url))
	}
}

func WithDeezerCloakBaseURL(url string) DeezerOption {
	return func(cfg *deezerOptions) {
		cfg.client = append(cfg.client, deezer.WithCloakBaseURL(url))
	}
}

func WithDeezerHTTPClient(client *http.Client) DeezerOption {
	return func(cfg *deezerOptions) {
		if client == nil {
			return
		}
		cfg.client = append(cfg.client, deezer.WithHTTPClient(client))
	}
}

type SpotifyOption func(*spotifyOptions)

type SpotifyCredentials struct {
	ClientID     string
	ClientSecret string
}
type spotifyOptions struct {
	client []spotify.ClientOption
}

func WithSpotify(creds SpotifyCredentials, opts ...SpotifyOption) CatalogOption {
	return func(r *Catalog) error {
		cfg := spotifyOptions{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&cfg)
		}
		client := spotify.NewClient(&spotify.Credentials{
			ClientID:     creds.ClientID,
			ClientSecret: creds.ClientSecret,
		}, cfg.client...)
		return r.register(Spotify, spotify.NewAdapter(client))
	}
}

func WithSpotifyAuthURL(url string) SpotifyOption {
	return func(cfg *spotifyOptions) {
		cfg.client = append(cfg.client, spotify.WithAuthURL(url))
	}
}

func WithSpotifyAPIURL(url string) SpotifyOption {
	return func(cfg *spotifyOptions) {
		cfg.client = append(cfg.client, spotify.WithAPIURL(url))
	}
}

func WithSpotifyHTTPClient(client *http.Client) SpotifyOption {
	return func(cfg *spotifyOptions) {
		if client == nil {
			return
		}
		cfg.client = append(cfg.client, spotify.WithHTTPClient(client))
	}
}

type SoundcloudOption func(*soundcloudOptions)

type soundcloudOptions struct {
	client []soundcloud.ClientOption
}

func WithSoundcloud(opts ...SoundcloudOption) CatalogOption {
	return func(r *Catalog) error {
		cfg := soundcloudOptions{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&cfg)
		}
		return r.register(Soundcloud, soundcloud.NewAdapter(soundcloud.NewClient(cfg.client...)))
	}
}

func WithSoundcloudAPIURL(url string) SoundcloudOption {
	return func(cfg *soundcloudOptions) {
		cfg.client = append(cfg.client, soundcloud.WithAPIURL(url))
	}
}

func WithSoundcloudWebURL(url string) SoundcloudOption {
	return func(cfg *soundcloudOptions) {
		cfg.client = append(cfg.client, soundcloud.WithWebURL(url))
	}
}

func WithSoundcloudHTTPClient(client *http.Client) SoundcloudOption {
	return func(cfg *soundcloudOptions) {
		if client == nil {
			return
		}
		cfg.client = append(cfg.client, soundcloud.WithHTTPClient(client))
	}
}

type YandexOption func(*yandexOptions)

type yandexOptions struct {
	client []yandex.ClientOption
}

func WithYandex(opts ...YandexOption) CatalogOption {
	return func(r *Catalog) error {
		cfg := yandexOptions{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&cfg)
		}
		return r.register(Yandex, yandex.NewAdapter(yandex.NewClient(cfg.client...)))
	}
}

func WithYandexAPIURL(url string) YandexOption {
	return func(cfg *yandexOptions) {
		cfg.client = append(cfg.client, yandex.WithAPIURL(url))
	}
}

func WithYandexHTTPClient(client *http.Client) YandexOption {
	return func(cfg *yandexOptions) {
		if client == nil {
			return
		}
		cfg.client = append(cfg.client, yandex.WithHTTPClient(client))
	}
}

type YoutubeOption func(*youtubeOptions)

type YoutubeCredentials struct {
	APIKey string
}

type youtubeOptions struct {
	client []youtube.ClientOption
}

func WithYoutube(creds YoutubeCredentials, opts ...YoutubeOption) CatalogOption {
	return func(r *Catalog) error {
		cfg := youtubeOptions{}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&cfg)
		}
		return r.register(Youtube, youtube.NewAdapter(youtube.NewClient(creds.APIKey, cfg.client...)))
	}
}

func WithYoutubeAPIURL(url string) YoutubeOption {
	return func(cfg *youtubeOptions) {
		cfg.client = append(cfg.client, youtube.WithAPIURL(url))
	}
}

func WithYoutubeHTTPClient(client *http.Client) YoutubeOption {
	return func(cfg *youtubeOptions) {
		if client == nil {
			return
		}
		cfg.client = append(cfg.client, youtube.WithHTTPClient(client))
	}
}
