package streaminx

import (
	"github.com/GeorgeGorbanev/streaminx/internal/apple"
	"github.com/GeorgeGorbanev/streaminx/internal/spotify"
	"github.com/GeorgeGorbanev/streaminx/internal/translator"
	"github.com/GeorgeGorbanev/streaminx/internal/yandex"
	"github.com/GeorgeGorbanev/streaminx/internal/youtube"
)

type RegistryOption func(registry *Registry)

type registryOptionsMap struct {
	appleClientOptions   []apple.ClientOption
	spotifyClientOptions []spotify.ClientOption
	yandexClientOptions  []yandex.ClientOption
	youtubeClientOptions []youtube.ClientOption
}

func WithProviderAdapter(provider *Provider, adapter Adapter) RegistryOption {
	return func(r *Registry) {
		r.adapters[provider.Code] = adapter
	}
}

func WithTranslator(translator translator.Translator) RegistryOption {
	return func(r *Registry) {
		r.translator = translator
	}
}

func WithAppleWebPlayerURL(url string) RegistryOption {
	return func(r *Registry) {
		r.options.appleClientOptions = append(r.options.appleClientOptions, apple.WithWebPlayerURL(url))
	}
}

func WithAppleAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.options.appleClientOptions = append(r.options.appleClientOptions, apple.WithAPIURL(url))
	}
}

func WithSpotifyAuthURL(url string) RegistryOption {
	return func(r *Registry) {
		r.options.spotifyClientOptions = append(r.options.spotifyClientOptions, spotify.WithAuthURL(url))
	}
}

func WithSpotifyAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.options.spotifyClientOptions = append(r.options.spotifyClientOptions, spotify.WithAPIURL(url))
	}
}

func WithYandexAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.options.yandexClientOptions = append(r.options.yandexClientOptions, yandex.WithAPIURL(url))
	}
}

func WithYoutubeAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.options.youtubeClientOptions = append(r.options.youtubeClientOptions, youtube.WithAPIURL(url))
	}
}
