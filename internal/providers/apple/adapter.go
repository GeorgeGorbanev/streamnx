package apple

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/duration"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	fetchTrack(ctx context.Context, id, storefront string) (entity, error)
	searchTracks(ctx context.Context, artist, title string) ([]entity, error)
	fetchTracksByISRC(ctx context.Context, isrc string) ([]entity, error)
	fetchAlbum(ctx context.Context, id, storefront string) (entity, error)
	fetchAlbumsByUPC(ctx context.Context, upc string) ([]entity, error)
	searchAlbums(ctx context.Context, artist, title string) ([]entity, error)
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
	dump, err := keyScheme.Dump(key)
	if err != nil {
		return ""
	}
	return dump
}

func (a *Adapter) detectAlbumID(albumURL string) string {
	key, err := parseAlbumLink(albumURL)
	if err != nil {
		return ""
	}
	dump, err := keyScheme.Dump(key)
	if err != nil {
		return ""
	}
	return dump
}

func (a *Adapter) FetchTrack(ctx context.Context, id string) (release.Track, error) {
	trackID, storefront, err := parseKeyParts(id)
	if err != nil {
		return release.Track{}, fmt.Errorf("failed to parse track id: %w", err)
	}

	track, err := a.client.fetchTrack(ctx, trackID, storefront)
	switch {
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get track from apple: %w", err)
	default:
		return a.releaseTrack(track, id, storefront), nil
	}
}

func (a *Adapter) FetchTracksByISRC(ctx context.Context, isrc string) ([]release.Track, error) {
	found, err := a.client.fetchTracksByISRC(ctx, isrc)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tracks by isrc from apple: %w", err)
	}

	tracks := make([]release.Track, len(found))
	for i, track := range found {
		key, err := parseTrackLink(track.Attributes.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse track link: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return nil, fmt.Errorf("failed to dump track key: %w", err)
		}
		_, storefront, err := parseKeyParts(id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse track key: %w", err)
		}
		tracks[i] = a.releaseTrack(track, id, storefront)
	}
	return tracks, nil
}

func (a *Adapter) releaseTrack(track entity, id, storefront string) release.Track {
	return release.Track{
		ID:          id,
		ISRC:        track.Attributes.ISRC,
		Artist:      track.Attributes.ArtistName,
		Title:       track.Attributes.Name,
		AlbumID:     a.trackAlbumID(track, storefront),
		AlbumTitle:  track.Attributes.AlbumName,
		URL:         track.Attributes.URL,
		CoverURL:    a.coverURL(track),
		Duration:    duration.MsToSeconds(track.Attributes.DurationInMillis),
		ReleaseDate: a.releaseDate(track.Attributes.ReleaseDate),
		Provider:    release.Apple,
		Description: a.description(track),
	}
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	albumID, storefront, err := parseKeyParts(id)
	if err != nil {
		return release.Album{}, fmt.Errorf("failed to parse album id: %w", err)
	}

	album, err := a.client.fetchAlbum(ctx, albumID, storefront)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get album from apple: %w", err)
	default:
		return a.releaseAlbum(album, id, storefront)
	}
}

func (a *Adapter) FetchAlbumsByUPC(ctx context.Context, upc string) ([]release.Album, error) {
	found, err := a.client.fetchAlbumsByUPC(ctx, upc)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch albums by upc from apple: %w", err)
	}

	albums := make([]release.Album, len(found))
	for i, album := range found {
		key, err := parseAlbumLink(album.Attributes.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse album link: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return nil, fmt.Errorf("failed to dump album key: %w", err)
		}
		_, storefront, err := parseKeyParts(id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse album key: %w", err)
		}
		albums[i], err = a.releaseAlbum(album, id, storefront)
		if err != nil {
			return nil, err
		}
	}
	return albums, nil
}

func (a *Adapter) releaseAlbum(album entity, id, storefront string) (release.Album, error) {
	tracks := album.Relationships.Tracks.Data
	trackIDs := make([]string, 0, len(tracks))
	for _, track := range tracks {
		var id string
		if track.ID != "" {
			key, err := newKey(storefront, track.ID)
			if err != nil {
				return release.Album{}, fmt.Errorf("failed to build key: %w", err)
			}
			id, err = keyScheme.Dump(key)
			if err != nil {
				return release.Album{}, fmt.Errorf("failed to dump key: %w", err)
			}
		} else {
			key, err := parseTrackLink(track.Attributes.URL)
			if err != nil {
				return release.Album{}, fmt.Errorf("failed to parse track link: %w", err)
			}
			id, err = keyScheme.Dump(key)
			if err != nil {
				return release.Album{}, fmt.Errorf("failed to dump track key: %w", err)
			}
		}
		trackIDs = append(trackIDs, id)
	}

	return release.Album{
		ID:          id,
		UPC:         album.Attributes.UPC,
		Title:       album.Attributes.Name,
		Artist:      album.Attributes.ArtistName,
		Label:       album.Attributes.RecordLabel,
		URL:         album.Attributes.URL,
		CoverURL:    a.coverURL(album),
		ReleaseDate: a.releaseDate(album.Attributes.ReleaseDate),
		Provider:    release.Apple,
		Description: a.description(album),
		TrackIDs:    trackIDs,
	}, nil
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	found, err := a.client.searchTracks(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search track from apple: %w", err)
	}
	return a.searchTracks(found)
}

func (a *Adapter) searchTracks(found []entity) ([]release.SearchTrack, error) {
	tracks := make([]release.SearchTrack, len(found))
	for i, track := range found {
		ck, err := parseTrackLink(track.Attributes.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse track link: %w", err)
		}
		id, err := keyScheme.Dump(ck)
		if err != nil {
			return nil, fmt.Errorf("failed to dump track key: %w", err)
		}
		_, storefront, err := parseKeyParts(id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse track key: %w", err)
		}
		tracks[i] = release.SearchTrack{
			ID:          id,
			ISRC:        track.Attributes.ISRC,
			Artist:      track.Attributes.ArtistName,
			Title:       track.Attributes.Name,
			AlbumID:     a.trackAlbumID(track, storefront),
			AlbumTitle:  track.Attributes.AlbumName,
			URL:         track.Attributes.URL,
			CoverURL:    a.coverURL(track),
			Provider:    release.Apple,
			Description: a.description(track),
		}
	}
	return tracks, nil
}

func (*Adapter) trackAlbumID(track entity, storefront string) string {
	albums := track.Relationships.Albums.Data
	if len(albums) == 0 || albums[0].ID == "" {
		return ""
	}

	key, err := newKey(storefront, albums[0].ID)
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
		return nil, fmt.Errorf("failed to search album from apple: %w", err)
	}

	albums := make([]release.SearchAlbum, len(found))
	for i, album := range found {
		ck, err := parseAlbumLink(album.Attributes.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse album link: %w", err)
		}
		id, err := keyScheme.Dump(ck)
		if err != nil {
			return nil, fmt.Errorf("failed to dump album key: %w", err)
		}
		albums[i] = release.SearchAlbum{
			ID:          id,
			UPC:         album.Attributes.UPC,
			Title:       album.Attributes.Name,
			Artist:      album.Attributes.ArtistName,
			URL:         album.Attributes.URL,
			CoverURL:    a.coverURL(album),
			Provider:    release.Apple,
			Description: a.description(album),
		}
	}
	return albums, nil
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}

func (a *Adapter) description(entity entity) string {
	if entity.Attributes.EditorialNotes.Standard != "" {
		return entity.Attributes.EditorialNotes.Standard
	}
	return entity.Attributes.EditorialNotes.Short
}

var coverSizeReplacer = strings.NewReplacer(
	"{w}", "1200",
	"{h}", "1200",
	"{f}", "jpg",
)

func (*Adapter) coverURL(entity entity) string {
	if entity.Attributes.Artwork.URL == "" {
		return ""
	}

	return coverSizeReplacer.Replace(entity.Attributes.Artwork.URL)
}

func (a *Adapter) releaseDate(raw string) release.Date {
	parsed, err := time.Parse(time.DateOnly, raw)
	if err != nil {
		return release.Date{}
	}
	return release.Date{
		Year:  parsed.Year(),
		Month: parsed.Month(),
		Day:   parsed.Day(),
	}
}
