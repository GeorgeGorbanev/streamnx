package streaminx

import (
	"fmt"

	"github.com/GeorgeGorbanev/streaminx/internal/spotify"
)

type SpotifyAdapter struct {
	client spotify.Client
}

func newSpotifyAdapter(client spotify.Client) *SpotifyAdapter {
	return &SpotifyAdapter{
		client: client,
	}
}

func (a *SpotifyAdapter) DetectTrackID(trackURL string) (string, error) {
	match := spotify.TrackRe.FindStringSubmatch(trackURL)
	if len(match) < 2 {
		return "", IDNotFoundError
	}
	return match[1], nil
}

func (a *SpotifyAdapter) DetectAlbumID(albumURL string) (string, error) {
	match := spotify.AlbumRe.FindStringSubmatch(albumURL)
	if len(match) < 2 {
		return "", IDNotFoundError
	}
	return match[1], nil
}

func (a *SpotifyAdapter) GetTrack(id string) (*Track, error) {
	track, err := a.client.GetTrack(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get track from spotify: %w", err)
	}
	if track == nil {
		return nil, nil
	}

	return a.adaptTrack(track), nil
}

func (a *SpotifyAdapter) SearchTrack(artistName, trackName string) (*Track, error) {
	track, err := a.client.SearchTrack(artistName, trackName)
	if err != nil {
		return nil, fmt.Errorf("failed to search track on spotify: %w", err)
	}
	if track == nil {
		return nil, nil
	}

	return a.adaptTrack(track), nil
}

func (a *SpotifyAdapter) GetAlbum(id string) (*Album, error) {
	album, err := a.client.GetAlbum(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get album from spotify: %w", err)
	}
	if album == nil {
		return nil, nil
	}

	return a.adaptAlbum(album), nil
}

func (a *SpotifyAdapter) SearchAlbum(artistName, albumName string) (*Album, error) {
	album, err := a.client.SearchAlbum(artistName, albumName)
	if err != nil {
		return nil, fmt.Errorf("failed to search album on spotify: %w", err)
	}
	if album == nil {
		return nil, nil
	}

	return a.adaptAlbum(album), nil
}

func (a *SpotifyAdapter) adaptTrack(track *spotify.Track) *Track {
	return &Track{
		ID:       track.ID,
		Title:    track.Name,
		Artist:   track.Artists[0].Name,
		URL:      track.URL(),
		Provider: Spotify,
	}
}

func (a *SpotifyAdapter) adaptAlbum(album *spotify.Album) *Album {
	return &Album{
		ID:       album.ID,
		Title:    album.Name,
		Artist:   album.Artists[0].Name,
		URL:      album.URL(),
		Provider: Spotify,
	}
}
