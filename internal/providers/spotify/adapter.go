package spotify

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/duration"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	fetchTrack(ctx context.Context, id string) (track, error)
	searchTracks(ctx context.Context, artist, title string) ([]track, error)
	fetchTracksByISRC(ctx context.Context, isrc string) ([]track, error)
	fetchAlbum(ctx context.Context, id string) (album, error)
	fetchAlbumsByUPC(ctx context.Context, upc string) ([]album, error)
	searchAlbums(ctx context.Context, artist, title string) ([]album, error)
}

func NewAdapter(c adapterClient) *Adapter {
	return &Adapter{client: c}
}

func (a *Adapter) ParseLink(rawURL string) (release.Type, string, bool) {
	if id := parseTrackLink(rawURL); id != "" {
		return release.TypeTrack, id, true
	}
	if id := parseAlbumLink(rawURL); id != "" {
		return release.TypeAlbum, id, true
	}
	return "", "", false
}

func (a *Adapter) FetchTrack(ctx context.Context, id string) (release.Track, error) {
	track, err := a.client.fetchTrack(ctx, id)
	switch {
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get track from spotify: %w", err)
	default:
		return a.releaseTrack(track), nil
	}
}

func (a *Adapter) FetchTracksByISRC(ctx context.Context, isrc string) ([]release.Track, error) {
	found, err := a.client.fetchTracksByISRC(ctx, isrc)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tracks by isrc on spotify: %w", err)
	}

	tracks := make([]release.Track, len(found))
	for i, track := range found {
		tracks[i] = a.releaseTrack(track)
	}
	return tracks, nil
}

func (a *Adapter) releaseTrack(track track) release.Track {
	return release.Track{
		ID:          track.ID,
		ISRC:        track.ExternalIDs.ISRC,
		Title:       track.Name,
		Artist:      a.artist(track.Artists),
		AlbumID:     track.Album.ID,
		AlbumTitle:  track.Album.Name,
		Provider:    release.Spotify,
		URL:         trackLink(track.ID),
		CoverURL:    a.coverURL(track.Album.Images),
		Duration:    duration.MsToSeconds(track.DurationMS),
		ReleaseDate: a.releaseDate(track.Album.ReleaseDate),
	}
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	album, err := a.client.fetchAlbum(ctx, id)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get album from spotify: %w", err)
	default:
		return a.releaseAlbum(album), nil
	}
}

func (a *Adapter) FetchAlbumsByUPC(ctx context.Context, upc string) ([]release.Album, error) {
	found, err := a.client.fetchAlbumsByUPC(ctx, upc)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch albums by upc on spotify: %w", err)
	}

	albums := make([]release.Album, len(found))
	for i, album := range found {
		albums[i] = a.releaseAlbum(album)
	}
	return albums, nil
}

func (a *Adapter) releaseAlbum(album album) release.Album {
	trackIDs := make([]string, 0, len(album.Tracks.Items))
	for _, track := range album.Tracks.Items {
		if track.ID != "" {
			trackIDs = append(trackIDs, track.ID)
		}
	}

	return release.Album{
		ID:          album.ID,
		UPC:         a.albumUPC(album),
		Title:       album.Name,
		Artist:      a.artist(album.Artists),
		Label:       album.Label,
		Provider:    release.Spotify,
		Creator:     album.Label,
		URL:         albumLink(album.ID),
		CoverURL:    a.coverURL(album.Images),
		ReleaseDate: a.releaseDate(album.ReleaseDate),
		TrackIDs:    trackIDs,
	}
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	found, err := a.client.searchTracks(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search track on spotify: %w", err)
	}
	return a.searchTracks(found), nil
}

func (a *Adapter) searchTracks(found []track) []release.SearchTrack {
	tracks := make([]release.SearchTrack, len(found))
	for i, track := range found {
		tracks[i] = release.SearchTrack{
			ID:         track.ID,
			ISRC:       track.ExternalIDs.ISRC,
			Title:      track.Name,
			Artist:     a.artist(track.Artists),
			AlbumID:    track.Album.ID,
			AlbumTitle: track.Album.Name,
			Provider:   release.Spotify,
			URL:        trackLink(track.ID),
			CoverURL:   a.coverURL(track.Album.Images),
		}
	}
	return tracks
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	found, err := a.client.searchAlbums(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search album on spotify: %w", err)
	}

	albums := make([]release.SearchAlbum, len(found))
	for i, foundAlbum := range found {
		albums[i] = release.SearchAlbum{
			ID:       foundAlbum.ID,
			UPC:      a.albumUPC(foundAlbum),
			Title:    foundAlbum.Name,
			Artist:   a.artist(foundAlbum.Artists),
			Provider: release.Spotify,
			Creator:  foundAlbum.Label,
			URL:      albumLink(foundAlbum.ID),
			CoverURL: a.coverURL(foundAlbum.Images),
		}
	}
	return albums, nil
}

func (*Adapter) albumUPC(album album) string {
	if album.ExternalIDs.UPC != "" {
		return album.ExternalIDs.UPC
	}
	return album.ExternalIDs.EAN
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}

func (a *Adapter) artist(artists []artist) string {
	if len(artists) == 0 {
		return ""
	}
	return artists[0].Name
}

func (a *Adapter) coverURL(images []spotifyImage) string {
	if len(images) == 0 {
		return ""
	}

	best := images[0]
	bestArea := best.Width * best.Height
	for _, image := range images[1:] {
		area := image.Width * image.Height
		if area > bestArea {
			best = image
			bestArea = area
		}
	}
	return best.URL
}

func (a *Adapter) releaseDate(raw string) release.Date {
	var (
		layout                   = time.DateOnly
		monthUnknown, dayUnknown bool
	)

	switch len(raw) {
	case 4:
		layout = "2006"
		monthUnknown = true
		dayUnknown = true
	case 7:
		layout = "2006-01"
		dayUnknown = true
	}

	parsed, err := time.Parse(layout, raw)
	if err != nil {
		return release.Date{}
	}

	date := release.Date{Year: parsed.Year()}
	if !monthUnknown {
		date.Month = parsed.Month()
	}
	if !dayUnknown {
		date.Day = parsed.Day()
	}
	return date
}
