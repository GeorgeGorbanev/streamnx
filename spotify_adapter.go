package streamnx

import (
	"context"
	"errors"
	"fmt"

	"github.com/GeorgeGorbanev/streamnx/internal/spotify"
)

type SpotifyAdapter struct {
	client spotify.Client
}

func newSpotifyAdapter(client spotify.Client) *SpotifyAdapter {
	return &SpotifyAdapter{
		client: client,
	}
}

func (a *SpotifyAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	track, err := a.client.FetchTrack(ctx, id)
	if err != nil {
		if errors.Is(err, spotify.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get track from spotify: %w", err)
	}

	return a.adaptTrack(track), nil
}

func (a *SpotifyAdapter) SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error) {
	track, err := a.client.SearchTrack(ctx, artistName, trackName)
	if err != nil {
		if errors.Is(err, spotify.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search track on spotify: %w", err)
	}

	return a.adaptTrack(track), nil
}

func (a *SpotifyAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	album, err := a.client.FetchAlbum(ctx, id)
	if err != nil {
		if errors.Is(err, spotify.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get album from spotify: %w", err)
	}

	return a.adaptAlbum(album), nil
}

func (a *SpotifyAdapter) SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error) {
	album, err := a.client.SearchAlbum(ctx, artistName, albumName)
	if err != nil {
		if errors.Is(err, spotify.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search album on spotify: %w", err)
	}

	return a.adaptAlbum(album), nil
}

func (a *SpotifyAdapter) adaptTrack(track *spotify.Track) *Entity {
	return &Entity{
		ID:       track.ID,
		Title:    track.Name,
		Artist:   track.Artists[0].Name,
		URL:      track.URL(),
		Provider: Spotify,
		Type:     Track,
	}
}

func (a *SpotifyAdapter) adaptAlbum(album *spotify.Album) *Entity {
	return &Entity{
		ID:       album.ID,
		Title:    album.Name,
		Artist:   album.Artists[0].Name,
		URL:      album.URL(),
		Provider: Spotify,
		Type:     Album,
	}
}
