package streamnx

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/GeorgeGorbanev/streamnx/internal/deezer"
)

type DeezerAdapter struct {
	client deezer.Client
}

func newDeezerAdapter(client deezer.Client) *DeezerAdapter {
	return &DeezerAdapter{
		client: client,
	}
}

func (a *DeezerAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	track, err := a.client.FetchTrack(ctx, id)
	if err != nil {
		if errors.Is(err, deezer.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get track from deezer: %w", err)
	}

	return a.adaptTrack(track), nil
}

func (a *DeezerAdapter) SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error) {
	track, err := a.client.SearchTrack(ctx, artistName, trackName)
	if err != nil {
		if errors.Is(err, deezer.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search track on deezer: %w", err)
	}

	return a.adaptTrack(track), nil
}

func (a *DeezerAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	album, err := a.client.FetchAlbum(ctx, id)
	if err != nil {
		if errors.Is(err, deezer.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get album from deezer: %w", err)
	}

	return a.adaptAlbum(album), nil
}

func (a *DeezerAdapter) SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error) {
	album, err := a.client.SearchAlbum(ctx, artistName, albumName)
	if err != nil {
		if errors.Is(err, deezer.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search album on deezer: %w", err)
	}

	return a.adaptAlbum(album), nil
}

func (a *DeezerAdapter) FetchCloak(ctx context.Context, id string) (*Entity, error) {
	cloak, err := a.client.FollowCloak(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to follow cloak link: %w", err)
	}

	dest, err := deezer.DetectCloakDest(cloak)
	if err != nil {
		return nil, fmt.Errorf("failed to detect cloak destination: %w", err)
	}

	if trackID := deezer.DetectTrackID(dest); trackID != "" {
		return a.FetchTrack(ctx, trackID)
	}
	if albumID := deezer.DetectAlbumID(dest); albumID != "" {
		return a.FetchAlbum(ctx, albumID)
	}

	return nil, fmt.Errorf("%w: cloak dest is not a track or album (%s)", EntityNotFoundError, dest)
}

func (a *DeezerAdapter) adaptTrack(track *deezer.Track) *Entity {
	return &Entity{
		ID:       strconv.Itoa(track.ID),
		Title:    track.Title,
		Artist:   track.Artist.Name,
		URL:      track.URL(),
		Provider: Deezer,
		Type:     Track,
	}
}

func (a *DeezerAdapter) adaptAlbum(album *deezer.Album) *Entity {
	return &Entity{
		ID:       strconv.Itoa(album.ID),
		Title:    album.Title,
		Artist:   album.Artist.Name,
		URL:      album.URL(),
		Provider: Deezer,
		Type:     Album,
	}
}
