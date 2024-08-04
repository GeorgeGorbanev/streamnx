package streamnx

import (
	"context"
	"errors"
	"fmt"

	"github.com/GeorgeGorbanev/streamnx/internal/apple"
)

type AppleAdapter struct {
	client apple.Client
}

func newAppleAdapter(client apple.Client) *AppleAdapter {
	return &AppleAdapter{
		client: client,
	}
}

func (a *AppleAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	ck := apple.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track id: %w", err)
	}

	track, err := a.client.FetchTrack(ctx, ck.ID, ck.Storefront)
	if err != nil {
		if errors.Is(err, apple.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get track from apple: %w", err)
	}

	res, err := a.adaptTrack(track)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *AppleAdapter) SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error) {
	track, err := a.client.SearchTrack(ctx, artistName, trackName)
	if err != nil {
		if errors.Is(err, apple.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search track from apple: %w", err)
	}
	res, err := a.adaptTrack(track)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *AppleAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	ck := apple.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album id: %w", err)
	}

	album, err := a.client.FetchAlbum(ctx, ck.ID, ck.Storefront)
	if err != nil {
		if errors.Is(err, apple.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get album from apple: %w", err)
	}

	res, err := a.adaptAlbum(album)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *AppleAdapter) SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error) {
	album, err := a.client.SearchAlbum(ctx, artistName, albumName)
	if err != nil {
		if errors.Is(err, apple.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search album from apple: %w", err)
	}
	res, err := a.adaptAlbum(album)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *AppleAdapter) adaptTrack(track *apple.Entity) (*Entity, error) {
	ck := apple.CompositeKey{}
	if err := ck.ParseFromTrackURL(track.Attributes.URL); err != nil {
		return nil, err
	}

	return &Entity{
		ID:       ck.Marshal(),
		Title:    track.Attributes.Name,
		Artist:   track.Attributes.ArtistName,
		URL:      track.Attributes.URL,
		Provider: Apple,
		Type:     Track,
	}, nil
}

func (a *AppleAdapter) adaptAlbum(album *apple.Entity) (*Entity, error) {
	ck := apple.CompositeKey{}
	if err := ck.ParseFromAlbumURL(album.Attributes.URL); err != nil {
		return nil, err
	}

	return &Entity{
		ID:       ck.Marshal(),
		Title:    album.Attributes.Name,
		Artist:   album.Attributes.ArtistName,
		URL:      album.Attributes.URL,
		Provider: Apple,
		Type:     Album,
	}, nil
}
