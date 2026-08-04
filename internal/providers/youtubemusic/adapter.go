package youtubemusic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	fetchTrack(ctx context.Context, id string) (playerResponse, error)
	fetchWatchNext(ctx context.Context, id string) (playlistPanelVideoRenderer, error)
	fetchAlbum(ctx context.Context, id string) (responsiveHeader, []responsiveListItem, error)
	searchTracks(ctx context.Context, query string) ([]responsiveListItem, error)
	searchAlbums(ctx context.Context, query string) ([]responsiveListItem, error)
}

func NewAdapter(c adapterClient) *Adapter {
	return &Adapter{client: c}
}

func (a *Adapter) ParseLink(rawURL string) (release.Type, string, bool) {
	if id := parseTrackURL(rawURL); id != "" {
		return release.TypeTrack, id, true
	}
	if key, ok := parseAlbumURL(rawURL); ok {
		id, err := dumpAlbumKey(key)
		if err == nil {
			return release.TypeAlbum, id, true
		}
	}
	return "", "", false
}

func (a *Adapter) FetchTrack(ctx context.Context, id string) (release.Track, error) {
	found, err := a.client.fetchTrack(ctx, id)
	switch {
	case errors.Is(err, errLoginRequired):
		return release.Track{}, release.ErrScrapingBlocked
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get track from youtube music: %w", err)
	}
	watchNext, err := a.client.fetchWatchNext(ctx, id)
	if err != nil {
		return release.Track{}, fmt.Errorf("failed to get track metadata from youtube music: %w", err)
	}

	metadata := found.Microformat.Renderer
	artist, albumID, albumTitle := a.metadataFromRuns(watchNext.LongBylineText.Runs)
	albumID = a.browseAlbumID(albumID)
	duration, _ := strconv.Atoi(a.firstNonEmpty(
		found.VideoDetails.LengthSeconds,
		metadata.VideoDetails.DurationSeconds,
	))

	return release.Track{
		ID:         found.VideoDetails.VideoID,
		Title:      a.firstNonEmpty(found.VideoDetails.Title, a.firstText(watchNext.Title)),
		Artist:     a.firstNonEmpty(found.VideoDetails.Author, artist),
		AlbumID:    albumID,
		AlbumTitle: albumTitle,
		URL:        a.firstNonEmpty(metadata.URLCanonical, trackURL(found.VideoDetails.VideoID)),
		CoverURL: a.firstCover(
			found.VideoDetails.Thumbnail.Thumbnails,
			metadata.Thumbnail.Thumbnails,
			watchNext.Thumbnail.Thumbnails,
		),
		Duration:    duration,
		ReleaseDate: release.Date{Year: a.yearFromRuns(watchNext.LongBylineText.Runs)},
		Provider:    release.YoutubeMusic,
		Creator:     metadata.PageOwnerDetails.Name,
		Description: metadata.Description,
	}, nil
}

func (a *Adapter) FetchTracksByISRC(context.Context, string) ([]release.Track, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	idType, rawID, err := parseAlbumID(id)
	if err != nil {
		return release.Album{}, fmt.Errorf("failed to parse youtube music album id: %w", err)
	}

	header, items, err := a.client.fetchAlbum(ctx, rawID)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get album from youtube music: %w", err)
	}

	trackIDs := make([]string, 0, len(items))
	fallbackTitle := ""
	fallbackArtist := ""
	fallbackBrowseID := ""
	fallbackCoverURL := ""
	for _, item := range items {
		if videoID := a.itemVideoID(item); videoID != "" {
			trackIDs = append(trackIDs, videoID)
		}
		artist, albumID, albumTitle := a.itemMetadata(item)
		fallbackArtist = a.firstNonEmpty(fallbackArtist, artist)
		fallbackBrowseID = a.firstNonEmpty(fallbackBrowseID, albumID)
		fallbackTitle = a.firstNonEmpty(fallbackTitle, albumTitle)
		fallbackCoverURL = a.firstNonEmpty(
			fallbackCoverURL,
			a.coverURL(item.Thumbnail.Renderer.Thumbnail.Thumbnails),
		)
	}

	alternativeURL := ""
	switch idType {
	case albumIDTypeBrowse:
		if playlistID := a.headerPlaylistID(header); playlistID != "" {
			alternativeURL = albumURL(albumIDTypePlaylist, playlistID)
		}
	case albumIDTypePlaylist:
		if fallbackBrowseID != "" {
			alternativeURL = albumURL(albumIDTypeBrowse, fallbackBrowseID)
		}
	}

	return release.Album{
		ID:             id,
		Title:          a.firstNonEmpty(a.firstText(header.Title), fallbackTitle),
		Artist:         a.firstNonEmpty(a.firstText(header.StraplineTextOne), fallbackArtist),
		URL:            albumURL(idType, rawID),
		AlternativeURL: alternativeURL,
		CoverURL: a.firstNonEmpty(
			a.coverURL(header.Thumbnail.Renderer.Thumbnail.Thumbnails),
			fallbackCoverURL,
		),
		ReleaseDate: release.Date{Year: a.yearFromRuns(header.Subtitle.Runs)},
		Provider:    release.YoutubeMusic,
		Description: a.joinRuns(header.Description.Shelf.Description.Runs),
		TrackIDs:    trackIDs,
	}, nil
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	items, err := a.client.searchTracks(ctx, a.searchQuery(artist, title))
	if err != nil {
		return nil, fmt.Errorf("failed to search tracks on youtube music: %w", err)
	}

	tracks := make([]release.SearchTrack, 0, len(items))
	for _, item := range items {
		videoID := a.itemVideoID(item)
		if videoID == "" {
			continue
		}
		artistName, albumID, albumTitle := a.itemMetadata(item)
		albumID = a.browseAlbumID(albumID)
		tracks = append(tracks, release.SearchTrack{
			ID:         videoID,
			Title:      a.itemTitle(item),
			Artist:     artistName,
			AlbumID:    albumID,
			AlbumTitle: albumTitle,
			URL:        trackURL(videoID),
			CoverURL:   a.coverURL(item.Thumbnail.Renderer.Thumbnail.Thumbnails),
			Provider:   release.YoutubeMusic,
		})
	}
	return tracks, nil
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	items, err := a.client.searchAlbums(ctx, a.searchQuery(artist, title))
	if err != nil {
		return nil, fmt.Errorf("failed to search albums on youtube music: %w", err)
	}

	albums := make([]release.SearchAlbum, 0, len(items))
	for _, item := range items {
		rawID := item.NavigationEndpoint.BrowseEndpoint.BrowseID
		if rawID == "" {
			continue
		}
		id, err := newAlbumID(albumIDTypeBrowse, rawID)
		if err != nil {
			continue
		}
		artistName, _, _ := a.itemMetadata(item)
		alternativeURL := ""
		if playlistID := a.itemPlaylistID(item); playlistID != "" {
			alternativeURL = albumURL(albumIDTypePlaylist, playlistID)
		}
		albums = append(albums, release.SearchAlbum{
			ID:             id,
			Title:          a.itemTitle(item),
			Artist:         artistName,
			URL:            albumURL(albumIDTypeBrowse, rawID),
			AlternativeURL: alternativeURL,
			CoverURL:       a.coverURL(item.Thumbnail.Renderer.Thumbnail.Thumbnails),
			Provider:       release.YoutubeMusic,
		})
	}
	return albums, nil
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}

func (a *Adapter) searchQuery(artist, title string) string {
	switch {
	case artist == "":
		return title
	case title == "":
		return artist
	default:
		return artist + " – " + title
	}
}

func (a *Adapter) itemPlaylistID(item responsiveListItem) string {
	return item.Overlay.Renderer.Content.PlayButton.
		PlayNavigationEndpoint.WatchPlaylistEndpoint.PlaylistID
}

func (a *Adapter) headerPlaylistID(header responsiveHeader) string {
	for _, button := range header.Buttons {
		if id := button.PlayButton.PlayNavigationEndpoint.WatchPlaylistEndpoint.PlaylistID; id != "" {
			return id
		}
	}
	return ""
}

func (a *Adapter) browseAlbumID(rawID string) string {
	if rawID == "" {
		return ""
	}
	id, err := newAlbumID(albumIDTypeBrowse, rawID)
	if err != nil {
		return ""
	}
	return id
}

func (a *Adapter) itemTitle(item responsiveListItem) string {
	if len(item.FlexColumns) == 0 {
		return ""
	}
	return a.firstText(item.FlexColumns[0].Renderer.Text)
}

func (a *Adapter) itemVideoID(item responsiveListItem) string {
	if item.PlaylistItemData.VideoID != "" {
		return item.PlaylistItemData.VideoID
	}
	for _, column := range item.FlexColumns {
		for _, run := range column.Renderer.Text.Runs {
			if run.NavigationEndpoint.WatchEndpoint.VideoID != "" {
				return run.NavigationEndpoint.WatchEndpoint.VideoID
			}
		}
	}
	return ""
}

func (a *Adapter) itemMetadata(item responsiveListItem) (artist, albumID, albumTitle string) {
	for _, column := range item.FlexColumns {
		foundArtist, foundAlbumID, foundAlbumTitle := a.metadataFromRuns(column.Renderer.Text.Runs)
		if artist == "" {
			artist = foundArtist
		}
		if albumID == "" {
			albumID = foundAlbumID
			albumTitle = foundAlbumTitle
		}
	}
	return artist, albumID, albumTitle
}

func (a *Adapter) metadataFromRuns(runs []run) (artist, albumID, albumTitle string) {
	for _, run := range runs {
		switch run.NavigationEndpoint.BrowseEndpoint.Context.Music.PageType {
		case "MUSIC_PAGE_TYPE_ARTIST":
			if artist == "" {
				artist = run.Text
			}
		case "MUSIC_PAGE_TYPE_ALBUM":
			if albumID == "" {
				albumID = run.NavigationEndpoint.BrowseEndpoint.BrowseID
				albumTitle = run.Text
			}
		}
	}
	return artist, albumID, albumTitle
}

func (a *Adapter) firstText(value text) string {
	if value.SimpleText != "" {
		return value.SimpleText
	}
	if len(value.Runs) == 0 {
		return ""
	}
	return value.Runs[0].Text
}

func (a *Adapter) firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func (a *Adapter) joinRuns(runs []run) string {
	var builder strings.Builder
	for _, run := range runs {
		builder.WriteString(run.Text)
	}
	return builder.String()
}

func (a *Adapter) yearFromRuns(runs []run) int {
	for _, run := range runs {
		if len(run.Text) != 4 {
			continue
		}
		year, err := strconv.Atoi(run.Text)
		if err == nil && year >= 1000 && year <= time.Now().UTC().Year()+1 {
			return year
		}
	}
	return 0
}

func (a *Adapter) firstCover(groups ...[]thumbnail) string {
	for _, thumbnails := range groups {
		if url := a.coverURL(thumbnails); url != "" {
			return url
		}
	}
	return ""
}

func (a *Adapter) coverURL(thumbnails []thumbnail) string {
	bestURL := ""
	bestArea := 0
	for _, thumbnail := range thumbnails {
		area := thumbnail.Width * thumbnail.Height
		if thumbnail.URL != "" && (bestURL == "" || area > bestArea) {
			bestURL = thumbnail.URL
			bestArea = area
		}
	}
	return bestURL
}
