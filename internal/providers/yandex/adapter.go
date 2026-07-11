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
		return release.Track{}, fmt.Errorf("failed to get track from yandex music: %w", err)
	}

	albumID, err := a.trackAlbumID(track.Albums)
	if err != nil {
		return release.Track{}, err
	}

	return release.Track{
		ID:          track.ID,
		Artist:      a.artistName(track.Artists),
		Title:       track.Title,
		AlbumID:     strconv.Itoa(albumID),
		AlbumTitle:  track.Albums[0].Title,
		Provider:    release.Yandex,
		URL:         trackLink(albumID, track.ID),
		CoverURL:    a.trackCoverURL(track.CoverURI, track.Albums),
		Duration:    duration.MsToSeconds(track.DurationMS),
		ReleaseDate: a.releaseDate(track.Albums[0].ReleaseDate),
	}, nil
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
			if track.ID != "" {
				trackIDs = append(trackIDs, track.ID)
			}
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
		albumID, err := a.trackAlbumID(track.Albums)
		if err != nil {
			return nil, fmt.Errorf("error adapting track: %w", err)
		}

		trackID := strconv.Itoa(track.ID)
		tracks[i] = release.SearchTrack{
			ID:         trackID,
			Title:      track.Title,
			Artist:     a.artistName(track.Artists),
			AlbumID:    strconv.Itoa(albumID),
			AlbumTitle: track.Albums[0].Title,
			Provider:   release.Yandex,
			URL:        trackLink(albumID, trackID),
			CoverURL:   a.trackCoverURL(track.CoverURI, track.Albums),
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

func (a *Adapter) trackAlbumID(albums []albumRef) (int, error) {
	if len(albums) == 0 {
		return 0, fmt.Errorf("unexpected yandex track response: missing album")
	}
	return albums[0].ID, nil
}

func (a *Adapter) trackCoverURL(trackCoverURI string, albums []albumRef) string {
	if trackCoverURI != "" {
		return a.coverURL(trackCoverURI)
	}
	if len(albums) == 0 {
		return ""
	}
	return a.coverURL(albums[0].CoverURI)
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
