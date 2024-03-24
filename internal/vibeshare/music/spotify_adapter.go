package music

import (
	"fmt"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
)

type SpotifyAdapter struct {
	client spotify.Client
}

func NewSpotifyAdapter(client spotify.Client) *SpotifyAdapter {
	return &SpotifyAdapter{
		client: client,
	}
}

func (a *SpotifyAdapter) DetectTrackID(trackURL string) string {
	match := spotify.TrackRe.FindStringSubmatch(trackURL)
	if match == nil || len(match) < 2 {
		return ""
	}
	return match[1]
}

func (a *SpotifyAdapter) DetectAlbumID(albumURL string) string {
	match := spotify.AlbumRe.FindStringSubmatch(albumURL)
	if match == nil || len(match) < 2 {
		return ""
	}
	return match[1]
}

func (a *SpotifyAdapter) GetTrack(id string) (*Track, error) {
	spotifyTrack, err := a.client.GetTrack(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get track from spotify: %w", err)
	}
	if spotifyTrack == nil {
		return nil, nil
	}

	return &Track{
		Title:  spotifyTrack.Name,
		Artist: spotifyTrack.Artists[0].Name,
		URL:    spotifyTrack.URL(),
	}, nil
}

func (a *SpotifyAdapter) SearchTrack(artistName, trackName string) (*Track, error) {
	spotifyTrack, err := a.client.SearchTrack(artistName, trackName)
	if err != nil {
		return nil, fmt.Errorf("failed to search track on spotify: %w", err)
	}
	if spotifyTrack == nil {
		return nil, nil // Трек не найден
	}

	return &Track{
		Title:  spotifyTrack.Name,
		Artist: spotifyTrack.Artists[0].Name,
		URL:    spotifyTrack.URL(),
	}, nil
}

func (a *SpotifyAdapter) GetAlbum(id string) (*Album, error) {
	spotifyAlbum, err := a.client.GetAlbum(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get album from spotify: %w", err)
	}
	if spotifyAlbum == nil {
		return nil, nil
	}

	return &Album{
		Title:  spotifyAlbum.Name,
		Artist: spotifyAlbum.Artists[0].Name,
		URL:    spotifyAlbum.URL(),
	}, nil
}

func (a *SpotifyAdapter) SearchAlbum(artistName, albumName string) (*Album, error) {
	spotifyAlbum, err := a.client.SearchAlbum(artistName, albumName)
	if err != nil {
		return nil, fmt.Errorf("failed to search album on spotify: %w", err)
	}
	if spotifyAlbum == nil {
		return nil, nil
	}

	return &Album{
		Title:  spotifyAlbum.Name,
		Artist: spotifyAlbum.Artists[0].Name,
		URL:    spotifyAlbum.URL(),
	}, nil
}
