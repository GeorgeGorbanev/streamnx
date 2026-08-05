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
	resolveSearchResultURL(ctx context.Context, et entityType, id uint64) (string, error)
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
		ISRC:        track.ISRC,
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

func (a *Adapter) FetchTracksByISRC(context.Context, string) ([]release.Track, error) {
	return nil, release.ErrUnsupportedOperation
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

func (a *Adapter) FetchAlbumsByUPC(context.Context, string) ([]release.Album, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	found, err := a.client.searchTracks(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search track on bandcamp: %w", err)
	}

	tracks := make([]release.SearchTrack, 0, len(found))
	for i, foundTrack := range found {
		track, err := a.adaptSearchTrack(ctx, foundTrack)
		if err != nil {
			return nil, fmt.Errorf("failed to adapt track search result at index %d (%q): %w",
				i, foundTrack.URL, err)
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

func (a *Adapter) adaptSearchTrack(ctx context.Context, track Entity) (release.SearchTrack, error) {
	trackURL, err := a.resolveSearchTrackURL(ctx, track)
	if err != nil {
		return release.SearchTrack{}, err
	}

	key, err := parseTrackLink(trackURL)
	if err != nil {
		return release.SearchTrack{}, fmt.Errorf("failed to parse composite key from track url: %w", err)
	}

	id, err := keyScheme.Dump(key)
	if err != nil {
		return release.SearchTrack{}, fmt.Errorf("failed to parse track id: %w", err)
	}
	return release.SearchTrack{
		ID:          id,
		Title:       track.Name,
		Artist:      track.BandName,
		AlbumTitle:  track.AlbumTitle,
		URL:         trackURL,
		CoverURL:    track.CoverURL,
		Provider:    release.Bandcamp,
		Creator:     a.creatorName(track),
		Description: track.Description,
	}, nil
}

func (a *Adapter) resolveSearchTrackURL(ctx context.Context, track Entity) (string, error) {
	if _, err := parseTrackLink(track.URL); err == nil {
		return track.URL, nil
	}

	canonicalURL, err := a.client.resolveSearchResultURL(ctx, trackEntityType, track.NumericID)
	if err != nil {
		return "", fmt.Errorf("failed to resolve canonical track url: %w", err)
	}
	if _, err := parseTrackLink(canonicalURL); err != nil {
		return "", fmt.Errorf("failed to parse composite key from canonical track url: %w", err)
	}
	return canonicalURL, nil
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

	albums := make([]release.SearchAlbum, 0, len(found))
	for i, foundAlbum := range found {
		album, err := a.adaptSearchAlbum(ctx, foundAlbum)
		if err != nil {
			return nil, fmt.Errorf("failed to adapt album search result at index %d (%q): %w",
				i, foundAlbum.URL, err)
		}
		albums = append(albums, album)
	}
	return albums, nil
}

func (a *Adapter) adaptSearchAlbum(ctx context.Context, album Entity) (release.SearchAlbum, error) {
	albumURL, err := a.resolveSearchAlbumURL(ctx, album)
	if err != nil {
		return release.SearchAlbum{}, err
	}

	key, err := parseAlbumLink(albumURL)
	if err != nil {
		return release.SearchAlbum{}, fmt.Errorf("failed to parse composite key from album url: %w", err)
	}

	id, err := keyScheme.Dump(key)
	if err != nil {
		return release.SearchAlbum{}, fmt.Errorf("failed to parse album id: %w", err)
	}
	return release.SearchAlbum{
		ID:          id,
		Title:       album.Name,
		Artist:      album.BandName,
		URL:         albumURL,
		CoverURL:    album.CoverURL,
		Provider:    release.Bandcamp,
		Description: album.Description,
		Creator:     a.creatorName(album),
	}, nil
}

func (a *Adapter) resolveSearchAlbumURL(ctx context.Context, album Entity) (string, error) {
	if _, err := parseAlbumLink(album.URL); err == nil {
		return album.URL, nil
	}

	canonicalURL, err := a.client.resolveSearchResultURL(ctx, albumEntityType, album.NumericID)
	if err != nil {
		return "", fmt.Errorf("failed to resolve canonical album url: %w", err)
	}
	if _, err := parseAlbumLink(canonicalURL); err != nil {
		return "", fmt.Errorf("failed to parse composite key from canonical album url: %w", err)
	}
	return canonicalURL, nil
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
