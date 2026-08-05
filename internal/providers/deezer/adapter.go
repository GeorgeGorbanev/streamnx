package deezer

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	fetchTrack(ctx context.Context, id string) (track, error)
	searchTracks(ctx context.Context, artist, title string) ([]track, error)
	fetchTrackByISRC(ctx context.Context, isrc string) (track, error)
	fetchAlbum(ctx context.Context, id string) (album, error)
	fetchAlbumByUPC(ctx context.Context, upc string) (album, error)
	searchAlbums(ctx context.Context, artist, title string) ([]album, error)
	followCloak(ctx context.Context, cloakCode string) (string, error)
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
	if id := parseCloakLink(rawURL); id != "" {
		return release.TypeCloak, id, true
	}
	return "", "", false
}

func (a *Adapter) FetchTrack(ctx context.Context, id string) (release.Track, error) {
	track, err := a.client.fetchTrack(ctx, id)
	switch {
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get track from deezer: %w", err)
	default:
		return a.releaseTrack(track), nil
	}
}

func (a *Adapter) FetchTracksByISRC(ctx context.Context, isrc string) ([]release.Track, error) {
	found, err := a.client.fetchTrackByISRC(ctx, isrc)
	switch {
	case errors.Is(err, errNotFound):
		return []release.Track{}, nil
	case err != nil:
		return nil, fmt.Errorf("failed to fetch track by isrc from deezer: %w", err)
	default:
		return []release.Track{a.releaseTrack(found)}, nil
	}
}

func (a *Adapter) releaseTrack(track track) release.Track {
	return release.Track{
		ID:          strconv.Itoa(track.ID),
		ISRC:        track.ISRC,
		Title:       track.Title,
		Artist:      track.Artist.Name,
		AlbumID:     a.albumID(track.Album.ID),
		AlbumTitle:  track.Album.Title,
		Provider:    release.Deezer,
		URL:         trackLink(track.ID),
		CoverURL:    a.coverURL(track.Album.coverURLs),
		Duration:    track.Duration,
		ReleaseDate: a.releaseDate(track.ReleaseDate),
	}
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	album, err := a.client.fetchAlbum(ctx, id)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get album from deezer: %w", err)
	default:
		return a.releaseAlbum(album), nil
	}
}

func (a *Adapter) FetchAlbumsByUPC(ctx context.Context, upc string) ([]release.Album, error) {
	found, err := a.client.fetchAlbumByUPC(ctx, upc)
	switch {
	case errors.Is(err, errNotFound):
		return []release.Album{}, nil
	case err != nil:
		return nil, fmt.Errorf("failed to fetch album by upc from deezer: %w", err)
	default:
		return []release.Album{a.releaseAlbum(found)}, nil
	}
}

func (a *Adapter) releaseAlbum(album album) release.Album {
	trackIDs := make([]string, 0, len(album.Tracks.Data))
	for _, track := range album.Tracks.Data {
		if track.ID != 0 {
			trackIDs = append(trackIDs, strconv.Itoa(track.ID))
		}
	}

	return release.Album{
		ID:          strconv.Itoa(album.ID),
		UPC:         album.UPC,
		Title:       album.Title,
		Artist:      album.Artist.Name,
		Label:       album.Label,
		Provider:    release.Deezer,
		Creator:     album.Label,
		URL:         albumLink(album.ID),
		CoverURL:    a.coverURL(album.coverURLs),
		ReleaseDate: a.releaseDate(album.ReleaseDate),
		TrackIDs:    trackIDs,
	}
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	found, err := a.client.searchTracks(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search track on deezer: %w", err)
	}
	return a.searchTracks(found), nil
}

func (a *Adapter) searchTracks(found []track) []release.SearchTrack {
	tracks := make([]release.SearchTrack, len(found))
	for i, track := range found {
		tracks[i] = release.SearchTrack{
			ID:         strconv.Itoa(track.ID),
			ISRC:       track.ISRC,
			Artist:     track.Artist.Name,
			Title:      track.Title,
			AlbumID:    a.albumID(track.Album.ID),
			AlbumTitle: track.Album.Title,
			Provider:   release.Deezer,
			URL:        trackLink(track.ID),
			CoverURL:   a.coverURL(track.Album.coverURLs),
		}
	}
	return tracks
}

func (*Adapter) albumID(id int) string {
	if id == 0 {
		return ""
	}
	return strconv.Itoa(id)
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	found, err := a.client.searchAlbums(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search album on deezer: %w", err)
	}

	albums := make([]release.SearchAlbum, len(found))
	for i, album := range found {
		albums[i] = release.SearchAlbum{
			ID:       strconv.Itoa(album.ID),
			UPC:      album.UPC,
			Title:    album.Title,
			Artist:   album.Artist.Name,
			Provider: release.Deezer,
			Creator:  album.Label,
			URL:      albumLink(album.ID),
			CoverURL: a.coverURL(album.coverURLs),
		}
	}
	return albums, nil
}

func (a *Adapter) Uncloak(ctx context.Context, id string) (release.Type, string, error) {
	cloak, err := a.client.followCloak(ctx, id)
	if err != nil {
		return "", "", fmt.Errorf("failed to follow cloak link: %w", err)
	}

	u, err := url.Parse(cloak)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse cloak URL: %w", err)
	}

	const destParam = "dest"
	val, ok := u.Query()[destParam]
	if !ok || len(val) == 0 {
		return "", "", fmt.Errorf("failed to find cloak '%s' param in url: %s", destParam, cloak)
	}

	linkType, entityID, ok := a.ParseLink(val[0])
	if !ok || (linkType != release.TypeTrack && linkType != release.TypeAlbum) {
		return "", "", fmt.Errorf("%w: cloak dest is not a track or album (%s)", release.ErrNotFound, val[0])
	}

	return linkType, entityID, nil
}

func (a *Adapter) coverURL(c coverURLs) string {
	switch {
	case c.CoverXL != "":
		return c.CoverXL
	case c.CoverBig != "":
		return c.CoverBig
	case c.CoverMedium != "":
		return c.CoverMedium
	case c.Cover != "":
		return c.Cover
	case c.CoverSmall != "":
		return c.CoverSmall
	default:
		return ""
	}
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
