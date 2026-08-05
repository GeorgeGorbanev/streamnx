package soundcloud

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/duration"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	fetchTrack(ctx context.Context, userSlug, trackSlug string) (track, error)
	fetchAlbum(ctx context.Context, userSlug, setSlug string) (album, error)
	fetchTracksByNumericIDs(ctx context.Context, ids []int64) ([]track, error)
	fetchTrackByURN(ctx context.Context, urn string) (track, error)
	searchTracks(ctx context.Context, artist, title string) ([]track, error)
	searchAlbums(ctx context.Context, artist, title string) ([]album, error)
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
	trackID, userSlug, err := parseKeyParts(id)
	if err != nil {
		return release.Track{}, fmt.Errorf("failed to parse track id: %w", err)
	}

	track, err := a.client.fetchTrack(ctx, userSlug, trackID)
	switch {
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get track from soundcloud: %w", err)
	}

	return release.Track{
		ID:          id,
		ISRC:        track.PublisherMetadata.ISRC,
		Title:       track.Title,
		Artist:      a.trackArtist(track),
		AlbumTitle:  track.PublisherMetadata.AlbumTitle,
		URL:         trackLink(track),
		CoverURL:    a.coverURL(track.ArtworkURL),
		Duration:    a.duration(track),
		ReleaseDate: a.releaseDate(track.ReleaseDate),
		Provider:    release.Soundcloud,
		Creator:     track.User.Username,
		Description: track.Description,
	}, nil
}

func (a *Adapter) FetchTracksByISRC(context.Context, string) ([]release.Track, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	albumID, userSlug, err := parseKeyParts(id)
	if err != nil {
		return release.Album{}, fmt.Errorf("failed to parse album id: %w", err)
	}

	album, err := a.client.fetchAlbum(ctx, userSlug, albumID)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get album from soundcloud: %w", err)
	}

	if !a.albumTracksComplete(album.Tracks) {
		if err := a.enrichIncompleteAlbumTracks(ctx, &album); err != nil {
			return release.Album{}, fmt.Errorf("failed to enrich incomplete soundcloud album tracks: %w", err)
		}
	}

	trackIDs := make([]string, 0, len(album.Tracks))
	for _, track := range album.Tracks {
		if !a.trackHasPermalink(track) {
			continue
		}
		key, err := parseTrackLink(trackLink(track))
		if err != nil {
			return release.Album{}, fmt.Errorf("failed to parse track link: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return release.Album{}, fmt.Errorf("failed to dump composite key: %w", err)
		}
		trackIDs = append(trackIDs, id)
	}

	return release.Album{
		ID:          id,
		Title:       album.Title,
		Artist:      a.albumArtist(album),
		Label:       album.LabelName,
		URL:         albumLink(album),
		CoverURL:    a.coverURL(album.ArtworkURL),
		ReleaseDate: a.releaseDate(album.ReleaseDate),
		Provider:    release.Soundcloud,
		Creator:     album.User.Username,
		Description: album.Description,
		TrackIDs:    trackIDs,
	}, nil
}

func (a *Adapter) FetchAlbumsByUPC(context.Context, string) ([]release.Album, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) enrichIncompleteAlbumTracks(ctx context.Context, album *album) error {
	numericIDs := make([]int64, 0)
	urns := make([]string, 0)
	seenNumericIDs := make(map[int64]struct{})
	seenURNs := make(map[string]struct{})

	for _, t := range album.Tracks {
		if a.trackHasPermalink(t) {
			continue
		}
		if t.ID != 0 {
			if _, seen := seenNumericIDs[t.ID]; seen {
				continue
			}
			seenNumericIDs[t.ID] = struct{}{}
			numericIDs = append(numericIDs, t.ID)
			continue
		}
		if t.URN == "" {
			continue
		}
		if _, seen := seenURNs[t.URN]; seen {
			continue
		}
		seenURNs[t.URN] = struct{}{}
		urns = append(urns, t.URN)
	}

	hydratedByNumericID := make(map[int64]track, len(numericIDs))
	if len(numericIDs) > 0 {
		hydrated, err := a.client.fetchTracksByNumericIDs(ctx, numericIDs)
		if err != nil {
			return fmt.Errorf("failed to enrich numeric track placeholders: %w", err)
		}
		for _, t := range hydrated {
			if t.ID != 0 {
				hydratedByNumericID[t.ID] = t
			}
		}
	}

	hydratedByURN := make(map[string]track, len(urns))
	for _, urn := range urns {
		hydrated, err := a.client.fetchTrackByURN(ctx, urn)
		switch {
		case errors.Is(err, errNotFound):
			continue
		case err != nil:
			return fmt.Errorf("failed to enrich urn track placeholder %q: %w", urn, err)
		}
		if hydrated.URN != "" {
			hydratedByURN[hydrated.URN] = hydrated
		}
	}

	for i, t := range album.Tracks {
		if a.trackHasPermalink(t) {
			continue
		}
		if hydrated, ok := hydratedByNumericID[t.ID]; t.ID != 0 && ok {
			album.Tracks[i] = hydrated
			continue
		}
		if hydrated, ok := hydratedByURN[t.URN]; t.ID == 0 && ok {
			album.Tracks[i] = hydrated
		}
	}

	return nil
}

func (a *Adapter) albumTracksComplete(tracks []track) bool {
	for _, t := range tracks {
		if !a.trackHasPermalink(t) {
			return false
		}
	}
	return true
}

func (a *Adapter) trackHasPermalink(t track) bool {
	return t.PermalinkURL != "" || (t.User.Permalink != "" && t.Permalink != "")
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	found, err := a.client.searchTracks(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search track on soundcloud: %w", err)
	}

	tracks := make([]release.SearchTrack, len(found))
	for i, foundTrack := range found {
		link := trackLink(foundTrack)
		key, err := parseTrackLink(link)
		if err != nil {
			return nil, fmt.Errorf("failed to parse composite key from track url: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return nil, fmt.Errorf("failed to dump composite key: %w", err)
		}

		tracks[i] = release.SearchTrack{
			ID:          id,
			ISRC:        foundTrack.PublisherMetadata.ISRC,
			Title:       foundTrack.Title,
			Artist:      a.trackArtist(foundTrack),
			AlbumTitle:  foundTrack.PublisherMetadata.AlbumTitle,
			URL:         link,
			CoverURL:    a.coverURL(foundTrack.ArtworkURL),
			Provider:    release.Soundcloud,
			Creator:     foundTrack.User.Username,
			Description: foundTrack.Description,
		}
	}
	return tracks, nil
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	found, err := a.client.searchAlbums(ctx, artist, title)
	if err != nil {
		return nil, fmt.Errorf("failed to search album on soundcloud: %w", err)
	}

	albums := make([]release.SearchAlbum, len(found))
	for i, album := range found {
		link := albumLink(album)
		key, err := parseAlbumLink(link)
		if err != nil {
			return nil, fmt.Errorf("failed to parse composite key from album url: %w", err)
		}
		id, err := keyScheme.Dump(key)
		if err != nil {
			return nil, fmt.Errorf("failed to dump composite key: %w", err)
		}
		albums[i] = release.SearchAlbum{
			ID:          id,
			Title:       album.Title,
			Artist:      a.albumArtist(album),
			URL:         link,
			CoverURL:    a.coverURL(album.ArtworkURL),
			Provider:    release.Soundcloud,
			Creator:     album.User.Username,
			Description: album.Description,
		}
	}
	return albums, nil
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}

func (*Adapter) trackArtist(t track) string {
	if t.PublisherMetadata.Artist != "" {
		return t.PublisherMetadata.Artist
	}
	return t.User.Username
}

func (*Adapter) albumArtist(a album) string {
	for _, t := range a.Tracks {
		if t.PublisherMetadata.Artist != "" {
			return t.PublisherMetadata.Artist
		}
	}
	return a.User.Username
}

func (*Adapter) duration(t track) int {
	var dur int
	if t.FullDuration > 0 {
		dur = t.FullDuration
	} else {
		dur = t.Duration
	}
	return duration.MsToSeconds(dur)
}

var artworkSizeSuffixRe = regexp.MustCompile(`-(?:large|t\d+x\d+|original|crop|small|tiny|mini|badge)(\.[^./?#]+)$`)

func (*Adapter) coverURL(artworkURL string) string {
	if artworkURL == "" {
		return ""
	}
	if !artworkSizeSuffixRe.MatchString(artworkURL) {
		return artworkURL
	}
	return artworkSizeSuffixRe.ReplaceAllString(artworkURL, "-original$1")
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
