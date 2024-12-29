package streamnx

import (
	"context"
)

type Adapter interface {
	FetchTrack(ctx context.Context, id string) (*Entity, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error)

	FetchAlbum(ctx context.Context, id string) (*Entity, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error)
}

type DeezerAdapter struct {
	client DeezerClient
}

func newDeezerAdapter(client DeezerClient) *DeezerAdapter {
	return &DeezerAdapter{
		client: client,
	}
}

func (a *DeezerAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	track, err := a.client.FetchTrack(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Entity{
		ID:       track.ID,
		Title:    track.Title,
		Artist:   track.Artist,
		URL:      track.URL,
		Provider: Deezer,
		Type:     Track,
	}, nil
}

func (a *DeezerAdapter) SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error) {
	track, err := a.client.SearchTrack(ctx, artistName, trackName)
	if err != nil {
		return nil, err
	}
	return &Entity{
		ID:       track.ID,
		Title:    track.Title,
		Artist:   track.Artist,
		URL:      track.URL,
		Provider: Deezer,
		Type:     Track,
	}, nil
}

func (a *DeezerAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	album, err := a.client.FetchAlbum(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Entity{
		ID:       album.ID,
		Title:    album.Title,
		Artist:   album.Artist,
		URL:      album.URL,
		Provider: Deezer,
		Type:     Album,
	}, nil
}

func (a *DeezerAdapter) SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error) {
	album, err := a.client.SearchAlbum(ctx, artistName, albumName)
	if err != nil {
		return nil, err
	}
	return &Entity{
		ID:       album.ID,
		Title:    album.Title,
		Artist:   album.Artist,
		URL:      album.URL,
		Provider: Deezer,
		Type:     Album,
	}, nil
}
