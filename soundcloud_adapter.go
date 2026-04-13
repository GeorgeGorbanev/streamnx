package streamnx

import (
	"context"
	"errors"
	"fmt"

	"github.com/GeorgeGorbanev/streamnx/internal/soundcloud"
)

type SoundcloudAdapter struct {
	client soundcloud.Client
}

func newSoundcloudAdapter(client soundcloud.Client) *SoundcloudAdapter {
	return &SoundcloudAdapter{
		client: client,
	}
}

func (a *SoundcloudAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	ck := soundcloud.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track id: %w", err)
	}

	track, err := a.client.FetchTrack(ctx, ck.UserSlug, ck.ID)
	if err != nil {
		if errors.Is(err, soundcloud.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get track from soundcloud: %w", err)
	}

	return a.adaptTrack(track, id), nil
}

func (a *SoundcloudAdapter) SearchTrack(ctx context.Context, artist, title string) (*Entity, error) {
	track, err := a.client.SearchTrack(ctx, artist, title)
	if err != nil {
		if errors.Is(err, soundcloud.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search track on soundcloud: %w", err)
	}

	return a.adaptSearchTrack(track)
}

func (a *SoundcloudAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	ck := soundcloud.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album id: %w", err)
	}

	album, err := a.client.FetchAlbum(ctx, ck.UserSlug, ck.ID)
	if err != nil {
		if errors.Is(err, soundcloud.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get album from soundcloud: %w", err)
	}

	return a.adaptAlbum(album, id), nil
}

func (a *SoundcloudAdapter) SearchAlbum(ctx context.Context, artist, title string) (*Entity, error) {
	album, err := a.client.SearchAlbum(ctx, artist, title)
	if err != nil {
		if errors.Is(err, soundcloud.ErrNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search album on soundcloud: %w", err)
	}

	return a.adaptSearchAlbum(album)
}

func (a *SoundcloudAdapter) FetchCloak(_ context.Context, _ string) (*Entity, error) {
	return nil, UnsupportedEntityTypeError
}

func (a *SoundcloudAdapter) adaptTrack(track *soundcloud.Track, id string) *Entity {
	return &Entity{
		ID:       id,
		Title:    track.Title,
		Artist:   track.User.Username,
		URL:      track.URL(),
		Provider: Soundcloud,
		Type:     Track,
	}
}

func (a *SoundcloudAdapter) adaptSearchTrack(track *soundcloud.Track) (*Entity, error) {
	ck := soundcloud.CompositeKey{}
	if err := ck.ParseFromTrackURL(track.URL()); err != nil {
		return nil, fmt.Errorf("failed to parse composite key from track url: %w", err)
	}

	return a.adaptTrack(track, ck.Marshal()), nil
}

func (a *SoundcloudAdapter) adaptAlbum(album *soundcloud.Album, id string) *Entity {
	return &Entity{
		ID:       id,
		Title:    album.Title,
		Artist:   album.User.Username,
		URL:      album.URL(),
		Provider: Soundcloud,
		Type:     Album,
	}
}

func (a *SoundcloudAdapter) adaptSearchAlbum(album *soundcloud.Album) (*Entity, error) {
	ck := soundcloud.CompositeKey{}
	if err := ck.ParseFromAlbumURL(album.URL()); err != nil {
		return nil, fmt.Errorf("failed to parse composite key from album url: %w", err)
	}

	return a.adaptAlbum(album, ck.Marshal()), nil
}
