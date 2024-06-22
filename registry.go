package streaminx

import (
	"context"
	"fmt"

	"github.com/GeorgeGorbanev/streaminx/internal/apple"
	"github.com/GeorgeGorbanev/streaminx/internal/spotify"
	"github.com/GeorgeGorbanev/streaminx/internal/translator"
	"github.com/GeorgeGorbanev/streaminx/internal/yandex"
	"github.com/GeorgeGorbanev/streaminx/internal/youtube"
)

type Registry struct {
	adapters   map[string]Adapter
	options    registryOptionsMap
	translator translator.Translator
}

type Credentials struct {
	GoogleTranslatorAPIKeyJSON string
	GoogleTranslatorProjectID  string
	YoutubeAPIKey              string
	SpotifyClientID            string
	SpotifyClientSecret        string
}

func NewRegistry(ctx context.Context, credentials Credentials, opts ...RegistryOption) (*Registry, error) {
	registry := Registry{
		adapters: make(map[string]Adapter),
	}
	for _, opt := range opts {
		opt(&registry)
	}

	if registry.translator == nil {
		translatorClient, err := translator.NewGoogleClient(ctx, &translator.GoogleCredentials{
			APIKeyJSON: credentials.GoogleTranslatorAPIKeyJSON,
			ProjectID:  credentials.GoogleTranslatorProjectID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create google translator client: %w", err)
		}
		registry.translator = translatorClient
	}

	if registry.Adapter(Apple) == nil {
		registry.adapters[Apple.Code] = newAppleAdapter(
			apple.NewHTTPClient(registry.options.appleClientOptions...),
		)
	}
	if registry.Adapter(Spotify) == nil {
		registry.adapters[Spotify.Code] = newSpotifyAdapter(
			spotify.NewHTTPClient(
				&spotify.Credentials{
					ClientID:     credentials.SpotifyClientID,
					ClientSecret: credentials.SpotifyClientSecret,
				},
				registry.options.spotifyClientOptions...,
			),
		)
	}
	if registry.Adapter(Yandex) == nil {
		registry.adapters[Yandex.Code] = newYandexAdapter(
			yandex.NewHTTPClient(registry.options.yandexClientOptions...),
			registry.translator,
		)
	}
	if registry.Adapter(Youtube) == nil {
		registry.adapters[Youtube.Code] = newYoutubeAdapter(
			youtube.NewHTTPClient(credentials.YoutubeAPIKey, registry.options.youtubeClientOptions...),
		)
	}

	return &registry, nil
}

func (r *Registry) Adapter(p *Provider) Adapter {
	return r.adapters[p.Code]
}

func (r *Registry) Close() error {
	return r.translator.Close()
}
