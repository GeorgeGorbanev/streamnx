package streaminx

import (
	"github.com/GeorgeGorbanev/streaminx/internal/apple"
	"github.com/GeorgeGorbanev/streaminx/internal/spotify"
	"github.com/GeorgeGorbanev/streaminx/internal/translator"
	"github.com/GeorgeGorbanev/streaminx/internal/yandex"
	"github.com/GeorgeGorbanev/streaminx/internal/youtube"
)

type RegistryOption func(registry *Registry)

type clientOptions struct {
	apple   []apple.ClientOption
	spotify []spotify.ClientOption
	yandex  []yandex.ClientOption
	youtube []youtube.ClientOption
}

func WithProviderAdapter(provider *Provider, adapter Adapter) RegistryOption {
	return func(r *Registry) {
		r.adapters[provider.сode] = adapter
	}
}

func WithTranslator(translator translator.Translator) RegistryOption {
	return func(r *Registry) {
		r.translator = translator
	}
}

func WithAppleWebPlayerURL(url string) RegistryOption {
	return func(r *Registry) {
		r.clientOptions.apple = append(r.clientOptions.apple, apple.WithWebPlayerURL(url))
	}
}

func WithAppleAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.clientOptions.apple = append(r.clientOptions.apple, apple.WithAPIURL(url))
	}
}

func WithSpotifyAuthURL(url string) RegistryOption {
	return func(r *Registry) {
		r.clientOptions.spotify = append(r.clientOptions.spotify, spotify.WithAuthURL(url))
	}
}

func WithSpotifyAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.clientOptions.spotify = append(r.clientOptions.spotify, spotify.WithAPIURL(url))
	}
}

func WithYandexAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.clientOptions.yandex = append(r.clientOptions.yandex, yandex.WithAPIURL(url))
	}
}

func WithYoutubeAPIURL(url string) RegistryOption {
	return func(r *Registry) {
		r.clientOptions.youtube = append(r.clientOptions.youtube, youtube.WithAPIURL(url))
	}
}
