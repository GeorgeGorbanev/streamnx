package streamnx

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/isrc"
)

type Catalog struct {
	adapters  map[ReleaseProvider]adapter
	providers []ReleaseProvider
}

type adapter interface {
	ParseLink(rawURL string) (ReleaseType, string, bool)
	FetchTrack(ctx context.Context, id string) (Track, error)
	FetchTracksByISRC(ctx context.Context, isrc string) ([]Track, error)
	FetchAlbum(ctx context.Context, id string) (Album, error)
	SearchTracks(ctx context.Context, artist, title string) ([]SearchTrack, error)
	SearchAlbums(ctx context.Context, artist, title string) ([]SearchAlbum, error)
	Uncloak(ctx context.Context, cloakID string) (ReleaseType, string, error)
}

type Link struct {
	URL         string
	Provider    ReleaseProvider
	ReleaseID   string
	ReleaseType ReleaseType
}

type SearchQuery struct {
	Artist string
	Title  string
}

var (
	ErrInvalidProvider    = errors.New("invalid provider")
	ErrDuplicateProvider  = errors.New("duplicate provider")
	ErrUnknownLink        = errors.New("unknown release link")
	ErrInvalidID          = errors.New("invalid release ID")
	ErrInvalidSearchQuery = errors.New("invalid search query")
)

func NewCatalog(opts ...CatalogOption) (*Catalog, error) {
	catalog := &Catalog{
		adapters: make(map[ReleaseProvider]adapter),
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(catalog); err != nil {
			return nil, err
		}
	}
	return catalog, nil
}

func (r *Catalog) ParseLink(rawURL string) (Link, error) {
	for p, a := range r.adapters {
		releaseType, id, ok := a.ParseLink(rawURL)
		if !ok {
			continue
		}
		return Link{
			URL:         rawURL,
			Provider:    p,
			ReleaseID:   id,
			ReleaseType: releaseType,
		}, nil
	}
	return Link{}, ErrUnknownLink
}

func (r *Catalog) FetchTrack(ctx context.Context, p ReleaseProvider, id string) (Track, error) {
	a, ok := r.adapters[p]
	if !ok {
		return Track{}, ErrInvalidProvider
	}
	if empty(id) {
		return Track{}, ErrInvalidID
	}
	return a.FetchTrack(ctx, id)
}

func (r *Catalog) FetchTracksByISRC(ctx context.Context, p ReleaseProvider, rawISRC string) ([]Track, error) {
	a, ok := r.adapters[p]
	if !ok {
		return nil, ErrInvalidProvider
	}
	normalizedISRC, valid := isrc.Normalize(rawISRC)
	if !valid {
		return nil, fmt.Errorf("%w: invalid isrc", ErrInvalidID)
	}
	tracks, err := a.FetchTracksByISRC(ctx, normalizedISRC)
	if err != nil {
		return nil, err
	}
	return forceNonNil(tracks), nil
}

func (r *Catalog) FetchAlbum(ctx context.Context, p ReleaseProvider, id string) (Album, error) {
	a, ok := r.adapters[p]
	if !ok {
		return Album{}, ErrInvalidProvider
	}
	if empty(id) {
		return Album{}, ErrInvalidID
	}
	return a.FetchAlbum(ctx, id)
}

func (r *Catalog) SearchTracks(ctx context.Context, p ReleaseProvider, query SearchQuery) ([]SearchTrack, error) {
	a, ok := r.adapters[p]
	if !ok {
		return nil, ErrInvalidProvider
	}
	if emptyQuery(query) {
		return nil, ErrInvalidSearchQuery
	}
	results, err := a.SearchTracks(ctx, query.Artist, query.Title)
	if err != nil {
		return nil, err
	}
	return forceNonNil(results), nil
}

func (r *Catalog) SearchAlbums(ctx context.Context, p ReleaseProvider, query SearchQuery) ([]SearchAlbum, error) {
	a, ok := r.adapters[p]
	if !ok {
		return nil, ErrInvalidProvider
	}
	if emptyQuery(query) {
		return nil, ErrInvalidSearchQuery
	}
	results, err := a.SearchAlbums(ctx, query.Artist, query.Title)
	if err != nil {
		return nil, err
	}
	return forceNonNil(results), nil
}

func (r *Catalog) Uncloak(ctx context.Context, p ReleaseProvider, cloakID string) (ReleaseType, string, error) {
	a, ok := r.adapters[p]
	if !ok {
		return "", "", ErrInvalidProvider
	}
	if empty(cloakID) {
		return "", "", ErrInvalidID
	}
	return a.Uncloak(ctx, cloakID)
}

func (r *Catalog) register(provider ReleaseProvider, adapter adapter) error {
	if !slices.Contains(Providers(), provider) {
		return fmt.Errorf("%w: %s", ErrInvalidProvider, provider)
	}
	if adapter == nil {
		return fmt.Errorf("%w: nil adapter for %s", ErrInvalidProvider, provider)
	}
	if _, ok := r.adapters[provider]; ok {
		return fmt.Errorf("%w: %s", ErrDuplicateProvider, provider)
	}
	r.adapters[provider] = adapter
	r.providers = append(r.providers, provider)
	return nil
}

func empty(value string) bool {
	return strings.TrimSpace(value) == ""
}

func emptyQuery(q SearchQuery) bool {
	return empty(q.Artist) && empty(q.Title)
}

func forceNonNil[T any](values []T) []T {
	if values == nil {
		return []T{}
	}
	return values
}
