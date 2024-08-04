package streamnx

import (
	"context"
	"errors"
	"fmt"

	"github.com/GeorgeGorbanev/streamnx/internal/apple"
	"github.com/GeorgeGorbanev/streamnx/internal/spotify"
	"github.com/GeorgeGorbanev/streamnx/internal/translator"
	"github.com/GeorgeGorbanev/streamnx/internal/yandex"
	"github.com/GeorgeGorbanev/streamnx/internal/youtube"
)

var (
	InvalidProviderError   = errors.New("invalid provider")
	InvalidEntityTypeError = errors.New("invalid entity type")
	EntityNotFoundError    = errors.New("entity not found")
)

type Registry struct {
	adapters      map[string]Adapter
	clientOptions clientOptions
	translator    translator.Translator
}

func NewRegistry(ctx context.Context, cred Credentials, opts ...RegistryOption) (*Registry, error) {
	registry := Registry{
		adapters: map[string]Adapter{},
	}
	for _, opt := range opts {
		opt(&registry)
	}

	if registry.translator == nil {
		translatorClient, err := translator.NewGoogleClient(ctx, cred.google())
		if err != nil {
			return nil, fmt.Errorf("failed to create google translator client: %w", err)
		}
		registry.translator = translatorClient
	}

	if registry.adapter(Apple) == nil {
		client := apple.NewHTTPClient(registry.clientOptions.apple...)
		registry.adapters[Apple.сode] = newAppleAdapter(client)
	}
	if registry.adapter(Spotify) == nil {
		client := spotify.NewHTTPClient(cred.spotify(), registry.clientOptions.spotify...)
		registry.adapters[Spotify.сode] = newSpotifyAdapter(client)
	}
	if registry.adapter(Yandex) == nil {
		client := yandex.NewHTTPClient(registry.clientOptions.yandex...)
		registry.adapters[Yandex.сode] = newYandexAdapter(client, registry.translator)
	}
	if registry.adapter(Youtube) == nil {
		client := youtube.NewHTTPClient(cred.YoutubeAPIKey, registry.clientOptions.youtube...)
		registry.adapters[Youtube.сode] = newYoutubeAdapter(client)
	}

	return &registry, nil
}

func (r *Registry) Close() error {
	return r.translator.Close()
}

func (r *Registry) Fetch(ctx context.Context, p *Provider, et EntityType, id string) (*Entity, error) {
	adapter := r.adapter(p)
	if adapter == nil {
		return nil, InvalidProviderError
	}

	switch et {
	case Track:
		return adapter.FetchTrack(ctx, id)
	case Album:
		return adapter.FetchAlbum(ctx, id)
	default:
		return nil, InvalidEntityTypeError
	}
}

func (r *Registry) Search(ctx context.Context, p *Provider, et EntityType, artist, name string) (*Entity, error) {
	adapter := r.adapter(p)
	if adapter == nil {
		return nil, InvalidProviderError
	}

	switch et {
	case Track:
		return adapter.SearchTrack(ctx, artist, name)
	case Album:
		return adapter.SearchAlbum(ctx, artist, name)
	default:
		return nil, InvalidEntityTypeError
	}
}

func (r *Registry) adapter(p *Provider) Adapter {
	return r.adapters[p.сode]
}
