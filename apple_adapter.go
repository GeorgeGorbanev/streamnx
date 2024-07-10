package streaminx

import (
	"context"
	"fmt"

	"github.com/GeorgeGorbanev/streaminx/internal/apple"
)

type AppleAdapter struct {
	client apple.Client
}

func newAppleAdapter(client apple.Client) *AppleAdapter {
	return &AppleAdapter{
		client: client,
	}
}

func (a *AppleAdapter) DetectTrackID(trackURL string) (string, error) {
	ck := apple.CompositeKey{}
	ck.ParseFromTrackURL(trackURL)

	if ck.Storefront == "" || ck.ID == "" || !apple.IsValidStorefront(ck.Storefront) {
		return "", IDNotFoundError
	}

	return ck.Marshal(), nil
}

func (a *AppleAdapter) DetectAlbumID(albumURL string) (string, error) {
	ck := apple.CompositeKey{}
	ck.ParseFromAlbumURL(albumURL)

	if ck.Storefront == "" || ck.ID == "" || !apple.IsValidStorefront(ck.Storefront) {
		return "", IDNotFoundError
	}

	return ck.Marshal(), nil
}

func (a *AppleAdapter) GetTrack(ctx context.Context, id string) (*Track, error) {
	ck := apple.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track id: %w", err)
	}

	track, err := a.client.GetTrack(ctx, ck.ID, ck.Storefront)
	if err != nil {
		return nil, fmt.Errorf("failed to get track from apple: %w", err)
	}
	if track == nil {
		return nil, nil
	}

	return a.adaptTrack(track), nil
}

func (a *AppleAdapter) SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error) {
	track, err := a.client.SearchTrack(ctx, artistName, trackName)
	if err != nil {
		return nil, fmt.Errorf("failed to search track from apple: %w", err)
	}
	if track == nil {
		return nil, nil
	}
	return a.adaptTrack(track), nil
}

func (a *AppleAdapter) GetAlbum(ctx context.Context, id string) (*Album, error) {
	ck := apple.CompositeKey{}
	if err := ck.Unmarshal(id); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album id: %w", err)
	}

	album, err := a.client.GetAlbum(ctx, ck.ID, ck.Storefront)
	if err != nil {
		return nil, fmt.Errorf("failed to get album from apple: %w", err)
	}
	if album == nil {
		return nil, nil
	}

	return a.adaptAlbum(album), nil
}

func (a *AppleAdapter) SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error) {
	album, err := a.client.SearchAlbum(ctx, artistName, albumName)
	if err != nil {
		return nil, fmt.Errorf("failed to search album from apple: %w", err)
	}
	if album == nil {
		return nil, nil
	}
	return a.adaptAlbum(album), nil
}

func (a *AppleAdapter) adaptTrack(track *apple.MusicEntity) *Track {
	ck := apple.CompositeKey{}
	ck.ParseFromTrackURL(track.Attributes.URL)

	return &Track{
		ID:       ck.Marshal(),
		Title:    track.Attributes.Name,
		Artist:   track.Attributes.ArtistName,
		URL:      track.Attributes.URL,
		Provider: Apple,
	}
}

func (a *AppleAdapter) adaptAlbum(album *apple.MusicEntity) *Album {
	ck := apple.CompositeKey{}
	ck.ParseFromAlbumURL(album.Attributes.URL)

	return &Album{
		ID:       ck.Marshal(),
		Title:    album.Attributes.Name,
		Artist:   album.Attributes.ArtistName,
		URL:      album.Attributes.URL,
		Provider: Apple,
	}
}
