package streamnx

import (
	"context"
	"errors"
	"fmt"

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

func (a *DeezerAdapter) FetchCloak(ctx context.Context, cloakCode string) (*Entity, error) {
	resolvedURL, err := a.client.FollowCloak(ctx, cloakCode)
	if err != nil {
		return nil, fmt.Errorf("failed to follow cloak link: %w", err)
	}

	// Parse resolved URL to detect entity type and ID
	if trackID := deezer.DetectTrackID(resolvedURL); trackID != "" {
		return a.FetchTrack(ctx, trackID)
	}
	if albumID := deezer.DetectAlbumID(resolvedURL); albumID != "" {
		return a.FetchAlbum(ctx, albumID)
	}

	return nil, fmt.Errorf("resolved URL is not a track or album: %s", resolvedURL)
}

func (a *DeezerAdapter) adaptTrack(track *deezer.Track) *Entity {
	return &Entity{
		ID:       track.ID,
		Title:    track.Title,
		Artist:   track.Artist.Name,
		URL:      track.URL(),
		Provider: Deezer,
		Type:     Track,
	}
}

func (a *DeezerAdapter) adaptAlbum(album *deezer.Album) *Entity {
	return &Entity{
		ID:       album.ID,
		Title:    album.Title,
		Artist:   album.Artist.Name,
		URL:      album.URL(),
		Provider: Deezer,
		Type:     Album,
	}
}