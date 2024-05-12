package music

import (
	"fmt"

	"github.com/GeorgeGorbanev/vibeshare/internal/apple"
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
	storeFront, id := trackReMatches(trackURL)
	if storeFront == "" || id == "" {
		return "", IDNotFoundError
	}
	if !apple.IsValidStorefront(storeFront) {
		return "", IDNotFoundError
	}

	return id, nil
}

func (a *AppleAdapter) DetectAlbumID(albumURL string) (string, error) {
	matches := apple.AlbumRe.FindStringSubmatch(albumURL)
	if len(matches) < 2 {
		return "", IDNotFoundError
	}

	regionCode := matches[1]
	if !apple.IsValidStorefront(regionCode) {
		return "", IDNotFoundError
	}
	return matches[2], nil
}

func (a *AppleAdapter) GetTrack(id string) (*Track, error) {
	track, err := a.client.GetTrack(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get track from apple: %w", err)
	}
	if track == nil {
		return nil, nil
	}

	return a.adaptTrack(track), nil
}

func (a *AppleAdapter) SearchTrack(artistName, trackName string) (*Track, error) {
	track, err := a.client.SearchTrack(artistName, trackName)
	if err != nil {
		return nil, fmt.Errorf("failed to search track from apple: %w", err)
	}
	if track == nil {
		return nil, nil
	}
	return a.adaptTrack(track), nil
}

func (a *AppleAdapter) GetAlbum(id string) (*Album, error) {
	album, err := a.client.GetAlbum(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get album from apple: %w", err)
	}
	if album == nil {
		return nil, nil
	}

	return a.adaptAlbum(album), nil
}

func (a *AppleAdapter) SearchAlbum(artistName, albumName string) (*Album, error) {
	album, err := a.client.SearchAlbum(artistName, albumName)
	if err != nil {
		return nil, fmt.Errorf("failed to search album from apple: %w", err)
	}
	if album == nil {
		return nil, nil
	}
	return a.adaptAlbum(album), nil
}

func (a *AppleAdapter) adaptTrack(track *apple.MusicEntity) *Track {
	return &Track{
		ID:       track.ID,
		Title:    track.Attributes.Name,
		Artist:   track.Attributes.ArtistName,
		URL:      track.Attributes.URL,
		Provider: Apple,
	}
}

func (a *AppleAdapter) adaptAlbum(album *apple.MusicEntity) *Album {
	return &Album{
		ID:       album.ID,
		Title:    album.Attributes.Name,
		Artist:   album.Attributes.ArtistName,
		URL:      album.Attributes.URL,
		Provider: Apple,
	}
}

func trackReMatches(trackURL string) (string, string) {
	matches := apple.AlbumTrackRe.FindStringSubmatch(trackURL)
	if len(matches) > 3 {
		return matches[1], matches[3]
	}
	matches = apple.SongRe.FindStringSubmatch(trackURL)
	if len(matches) > 2 {
		return matches[1], matches[2]
	}
	return "", ""
}
