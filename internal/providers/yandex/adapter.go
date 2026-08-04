package yandex

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/duration"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	fetchTrack(ctx context.Context, id string) (track, error)
	searchTracks(ctx context.Context, query string) ([]searchTrack, error)
	fetchAlbum(ctx context.Context, id string) (album, error)
	searchAlbums(ctx context.Context, query string) ([]searchAlbum, error)
}

func NewAdapter(c adapterClient) *Adapter {
	return &Adapter{client: c}
}

func (a *Adapter) ParseLink(rawURL string) (release.Type, string, bool) {
	if id := a.detectTrackID(rawURL); id != "" {
		return release.TypeTrack, id, true
	}
	if id := parseAlbumLink(rawURL); id != "" {
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

func (a *Adapter) FetchTrack(ctx context.Context, id string) (release.Track, error) {
	albumID, trackID, err := parseTrackKeyParts(id)
	if err != nil {
		return release.Track{}, fmt.Errorf("failed to parse track id: %w", err)
	}

	track, err := a.client.fetchTrack(ctx, trackID)
	switch {
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get track from yandex music: %w", err)
	}

	trackAlbum, ok := a.trackAlbum(track.Albums, albumID)
	if !ok {
		return release.Track{}, fmt.Errorf(
			"yandex music track %q does not belong to album %q: %w",
			trackID,
			albumID,
			release.ErrNotFound,
		)
	}

	return release.Track{
		ID:          id,
		Artist:      a.artistName(track.Artists),
		Title:       track.Title,
		AlbumID:     albumID,
		AlbumTitle:  trackAlbum.Title,
		Provider:    release.Yandex,
		URL:         trackLink(albumID, trackID),
		CoverURL:    a.trackCoverURL(track.CoverURI, trackAlbum),
		Duration:    duration.MsToSeconds(track.DurationMS),
		ReleaseDate: a.releaseDate(trackAlbum.ReleaseDate),
	}, nil
}

func (a *Adapter) FetchTracksByISRC(context.Context, string) ([]release.Track, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	album, err := a.client.fetchAlbum(ctx, id)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get album from yandex music: %w", err)
	}

	trackIDs := make([]string, 0, len(album.Volumes))
	for _, volume := range album.Volumes {
		for _, track := range volume {
			if track.ID == "" {
				continue
			}
			key, err := newTrackKey(strconv.Itoa(album.ID), track.ID)
			if err != nil {
				return release.Album{}, fmt.Errorf("failed to build track key: %w", err)
			}
			id, err := keyScheme.Dump(key)
			if err != nil {
				return release.Album{}, fmt.Errorf("failed to dump track key: %w", err)
			}
			trackIDs = append(trackIDs, id)
		}
	}

	return release.Album{
		ID:          strconv.Itoa(album.ID),
		Artist:      a.artistName(album.Artists),
		Title:       album.Title,
		Label:       a.labelName(album.Labels),
		Provider:    release.Yandex,
		URL:         albumLink(album.ID),
		CoverURL:    a.coverURL(album.CoverURI),
		ReleaseDate: a.releaseDate(album.ReleaseDate),
		TrackIDs:    trackIDs,
	}, nil
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	found, err := a.client.searchTracks(ctx, artist+" – "+title)
	switch {
	case errors.Is(err, errNotFound):
		return []release.SearchTrack{}, nil
	case err != nil:
		return nil, fmt.Errorf("error searching track: %w", err)
	}

	tracks := make([]release.SearchTrack, len(found))
	for i, track := range found {
		if len(track.Albums) == 0 {
			return nil, fmt.Errorf("error adapting track: unexpected yandex track response: missing album")
		}
		trackAlbum := track.Albums[0]
		albumID := strconv.Itoa(trackAlbum.ID)
		trackID := strconv.Itoa(track.ID)
		key, err := newTrackKey(albumID, trackID)
		if err != nil {
			return nil, fmt.Errorf("error building track key: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return nil, fmt.Errorf("error dumping track key: %w", err)
		}

		tracks[i] = release.SearchTrack{
			ID:         id,
			Title:      track.Title,
			Artist:     a.artistName(track.Artists),
			AlbumID:    albumID,
			AlbumTitle: trackAlbum.Title,
			Provider:   release.Yandex,
			URL:        trackLink(albumID, trackID),
			CoverURL:   a.trackCoverURL(track.CoverURI, trackAlbum),
		}
	}
	return tracks, nil
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	found, err := a.client.searchAlbums(ctx, artist+" – "+title)
	switch {
	case errors.Is(err, errNotFound):
		return []release.SearchAlbum{}, nil
	case err != nil:
		return nil, fmt.Errorf("error searching album: %w", err)
	}

	albums := make([]release.SearchAlbum, len(found))
	for i, foundAlbum := range found {
		albums[i] = release.SearchAlbum{
			ID:       strconv.Itoa(foundAlbum.ID),
			Provider: release.Yandex,
			Title:    foundAlbum.Title,
			Artist:   a.artistName(foundAlbum.Artists),
			URL:      albumLink(foundAlbum.ID),
			CoverURL: a.coverURL(foundAlbum.CoverURI),
		}
	}
	return albums, nil
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}

func (a *Adapter) artistName(artists []artist) string {
	if len(artists) == 0 {
		return ""
	}
	return artists[0].Name
}

func (a *Adapter) labelName(labels []label) string {
	if len(labels) == 0 {
		return ""
	}
	return labels[0].Name
}

func (a *Adapter) trackAlbum(albums []albumRef, albumID string) (albumRef, bool) {
	for _, album := range albums {
		if strconv.Itoa(album.ID) == albumID {
			return album, true
		}
	}
	return albumRef{}, false
}

func (a *Adapter) trackCoverURL(trackCoverURI string, album albumRef) string {
	if album.CoverURI != "" {
		return a.coverURL(album.CoverURI)
	}
	return a.coverURL(trackCoverURI)
}

func (a *Adapter) coverURL(coverURI string) string {
	if coverURI == "" {
		return ""
	}

	u := coverURI
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "https://" + u
	}
	return strings.ReplaceAll(u, "%%", "1000x1000")
}

func (a *Adapter) releaseDate(raw string) release.Date {
	parsed, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return release.Date{}
	}
	return release.Date{
		Year:  parsed.Year(),
		Month: parsed.Month(),
		Day:   parsed.Day(),
	}
}
