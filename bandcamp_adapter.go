package streamnx

import (
	"context"
	"errors"
	"fmt"

	"github.com/GeorgeGorbanev/streamnx/internal/bandcamp"
)

type BandcampAdapter struct {
	client bandcamp.Client
}

func newBandcampAdapter(c bandcamp.Client) *BandcampAdapter {
	return &BandcampAdapter{
		client: c,
	}
}

func (a *BandcampAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	ck := bandcamp.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track id: %w", err)
	}

	track, err := a.client.FetchTrack(ctx, ck.ArtistSlug, ck.ID)
	if err != nil {
		if errors.Is(err, bandcamp.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get track from bandcamp: %w", err)
	}

	return a.adaptEntity(track, Track, id), nil
}

func (a *BandcampAdapter) SearchTrack(ctx context.Context, artist, title string) (*Entity, error) {
	track, err := a.client.SearchTrack(ctx, artist, title)
	if err != nil {
		if errors.Is(err, bandcamp.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search track on bandcamp: %w", err)
	}

	return a.adaptSearchResult(track, Track)
}

func (a *BandcampAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	ck := bandcamp.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album id: %w", err)
	}

	album, err := a.client.FetchAlbum(ctx, ck.ArtistSlug, ck.ID)
	if err != nil {
		if errors.Is(err, bandcamp.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get album from bandcamp: %w", err)
	}

	return a.adaptEntity(album, Album, id), nil
}

func (a *BandcampAdapter) SearchAlbum(ctx context.Context, artist, title string) (*Entity, error) {
	album, err := a.client.SearchAlbum(ctx, artist, title)
	if err != nil {
		if errors.Is(err, bandcamp.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search album on bandcamp: %w", err)
	}

	return a.adaptSearchResult(album, Album)
}

func (a *BandcampAdapter) FetchCloak(_ context.Context, _ string) (*Entity, error) {
	return nil, UnsupportedEntityTypeError
}

func (a *BandcampAdapter) adaptEntity(e *bandcamp.Entity, et EntityType, id string) *Entity {
	return &Entity{
		ID:       id,
		Type:     et,
		Title:    e.Name,
		Artist:   e.BandName,
		URL:      e.URL,
		Provider: Bandcamp,
	}
}

func (a *BandcampAdapter) adaptSearchResult(e *bandcamp.Entity, et EntityType) (*Entity, error) {
	var (
		ck  bandcamp.CompositeKey
		err error
	)

	switch et {
	case Track:
		err = ck.ParseFromTrackURL(e.URL)
	case Album:
		err = ck.ParseFromAlbumURL(e.URL)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse composite key from url: %w", err)
	}

	return a.adaptEntity(e, et, ck.Marshal()), nil
}
