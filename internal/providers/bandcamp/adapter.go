package bandcamp

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
	fetchAlbum(ctx context.Context, artistSlug, id string) (Entity, error)
	fetchTrack(ctx context.Context, artistSlug, id string) (Entity, error)
	searchAlbums(ctx context.Context, artist, title string) ([]Entity, error)
	searchTracks(ctx context.Context, artist, title string) ([]Entity, error)
}

func NewAdapter(c adapterClient) *Adapter {
	return &Adapter{client: c}
}

func (a *Adapter) ParseLink(rawURL string) (release.Type, string, bool) {
	if id := a.detectTrackID(rawURL); id != "" {
		return release.TypeTrack, id, true
	}
	if id := a.detectAlbumID(rawURL); id != "" {
		return release.TypeAlbum, id, true
	}
	return "", "", false
}

func (a *Adapter) detectTrackID(trackURL string) string {
	key, err := parseTrackLink(trackURL)
	if err != nil {
		return ""
	}
	id, err := keyScheme.Dump(key)
	if err != nil {
		return ""
	}
	return id
}

func (a *Adapter) detectAlbumID(albumURL string) string {
	key, err := parseAlbumLink(albumURL)
	if err != nil {
		return ""
	}
	id, err := keyScheme.Dump(key)
	if err != nil {
		return ""
	}
	return id
}

func (a *Adapter) FetchTrack(ctx context.Context, id string) (release.Track, error) {
	trackID, artistSlug, err := parseKeyParts(id)
	if err != nil {
		return release.Track{}, fmt.Errorf("failed to parse track id: %w", err)
	}

	track, err := a.client.fetchTrack(ctx, artistSlug, trackID)
	switch {
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get track from bandcamp: %w", err)
	}

	return release.Track{
		ID:          id,
		Artist:      track.BandName,
		Title:       track.Name,
		AlbumID:     a.albumID(track.AlbumURL),
		AlbumTitle:  track.AlbumTitle,
		URL:         track.URL,
		CoverURL:    track.CoverURL,
		Duration:    duration.ISO8601ToSeconds(track.Duration),
		ReleaseDate: a.releaseDate(track.ReleaseDate),
		Provider:    release.Bandcamp,
		Description: track.Description,
		Creator:     a.creatorName(track),
	}, nil
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	albumID, artistSlug, err := parseKeyParts(id)
	if err != nil {
		return release.Album{}, fmt.Errorf("failed to parse album id: %w", err)
	}

	album, err := a.client.fetchAlbum(ctx, artistSlug, albumID)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get album from bandcamp: %w", err)
	}

	trackIDs := make([]string, 0, len(album.TrackURLs))
	for _, trackURL := range album.TrackURLs {
		key, err := parseTrackLink(trackURL)
		if err != nil {
			return release.Album{}, fmt.Errorf("failed to parse track id: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return release.Album{}, fmt.Errorf("failed to parse track id: %w", err)
		}
		trackIDs = append(trackIDs, id)
	}

	return release.Album{
		ID:          id,
		Artist:      album.BandName,
		Title:       album.Name,
		Label:       album.CreatorName,
		URL:         album.URL,
		CoverURL:    album.CoverURL,
		ReleaseDate: a.releaseDate(album.ReleaseDate),
		Provider:    release.Bandcamp,
		Description: album.Description,
		Creator:     a.creatorName(album),
		TrackIDs:    trackIDs,
	}, nil
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	found, err := a.client.searchTracks(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search track on bandcamp: %w", err)
	}

	tracks := make([]release.SearchTrack, len(found))
	for i, track := range found {
		key, err := parseTrackLink(track.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse composite key from url: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse track id: %w", err)
		}
		tracks[i] = release.SearchTrack{
			ID:          id,
			Title:       track.Name,
			Artist:      track.BandName,
			AlbumTitle:  track.AlbumTitle,
			URL:         track.URL,
			CoverURL:    track.CoverURL,
			Provider:    release.Bandcamp,
			Creator:     a.creatorName(track),
			Description: track.Description,
		}
	}
	return tracks, nil
}

func (*Adapter) albumID(albumURL string) string {
	if albumURL == "" {
		return ""
	}
	key, err := parseAlbumLink(albumURL)
	if err != nil {
		return ""
	}
	id, err := keyScheme.Dump(key)
	if err != nil {
		return ""
	}
	return id
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	found, err := a.client.searchAlbums(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search album on bandcamp: %w", err)
	}

	albums := make([]release.SearchAlbum, len(found))
	for i, album := range found {
		key, err := parseAlbumLink(album.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse composite key from url: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse album id: %w", err)
		}
		albums[i] = release.SearchAlbum{
			ID:          id,
			Title:       album.Name,
			Artist:      album.BandName,
			URL:         album.URL,
			CoverURL:    album.CoverURL,
			Provider:    release.Bandcamp,
			Description: album.Description,
			Creator:     a.creatorName(album),
		}
	}
	return albums, nil
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}

func (a *Adapter) creatorName(e Entity) string {
	if e.CreatorName != "" {
		return e.CreatorName
	}
	return e.BandName
}

func (a *Adapter) releaseDate(raw string) release.Date {
	parsed, err := time.Parse("02 Jan 2006 15:04:05 MST", raw)
	if err != nil {
		return release.Date{}
	}
	return release.Date{
		Year:  parsed.Year(),
		Month: parsed.Month(),
		Day:   parsed.Day(),
	}
}
