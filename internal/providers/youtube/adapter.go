package youtube

import (
	"context"
	"errors"
	"fmt"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/duration"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	fetchVideo(ctx context.Context, id string) (video, error)
	fetchPlaylist(ctx context.Context, id string) (playlist, error)
	searchVideos(ctx context.Context, term string) ([]videoSearchResult, error)
	searchPlaylists(ctx context.Context, term string) ([]playlistSearchResult, error)
	fetchPlaylistItems(ctx context.Context, id string) ([]playlistItem, error)
}

func NewAdapter(c adapterClient) *Adapter {
	return &Adapter{client: c}
}

func (a *Adapter) ParseLink(rawURL string) (release.Type, string, bool) {
	if id := parseTrackURL(rawURL); id != "" {
		return release.TypeTrack, id, true
	}
	if id := parseAlbumURL(rawURL); id != "" {
		return release.TypeAlbum, id, true
	}
	return "", "", false
}

func (a *Adapter) FetchTrack(ctx context.Context, id string) (release.Track, error) {
	video, err := a.client.fetchVideo(ctx, id)
	switch {
	case errors.Is(err, errNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get video from youtube: %w", err)
	}

	return release.Track{
		ID:          video.ID,
		Title:       video.Snippet.Title,
		Provider:    release.Youtube,
		Creator:     a.ownerChannelTitle(video.Snippet),
		Description: video.Snippet.Description,
		CoverURL:    coverURL(video.Snippet.Thumbnails),
		Duration:    duration.ISO8601ToSeconds(video.ContentDetails.Duration),
		URL:         videoURL(video.ID),
	}, nil
}

func (a *Adapter) FetchTracksByISRC(context.Context, string) ([]release.Track, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	album, err := a.client.fetchPlaylist(ctx, id)
	switch {
	case errors.Is(err, errNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, fmt.Errorf("failed to get playlist from youtube: %w", err)
	}
	items, err := a.client.fetchPlaylistItems(ctx, album.ID)
	if err != nil {
		return release.Album{}, fmt.Errorf("failed to get playlist items from youtube: %w", err)
	}

	trackIDs := make([]string, 0, len(items))
	for _, item := range items {
		if videoID := item.Snippet.ResourceID.VideoID; videoID != "" {
			trackIDs = append(trackIDs, videoID)
		}
	}

	return release.Album{
		ID:          album.ID,
		Title:       album.Snippet.Title,
		URL:         playlistURL(album.ID),
		Provider:    release.Youtube,
		Creator:     a.ownerChannelTitle(album.Snippet),
		Description: album.Snippet.Description,
		CoverURL:    coverURL(album.Snippet.Thumbnails),
		TrackIDs:    trackIDs,
	}, nil
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	searchResults, err := a.client.searchVideos(ctx, artist+" – "+title)
	if err != nil {
		return nil, fmt.Errorf("failed to search video on youtube: %w", err)
	}

	tracks := make([]release.SearchTrack, len(searchResults))
	for i, searchResult := range searchResults {
		video, err := a.client.fetchVideo(ctx, searchResult.ID.VideoID)
		if err != nil {
			return nil, fmt.Errorf("failed to get video from youtube: %w", err)
		}
		tracks[i] = release.SearchTrack{
			ID:          video.ID,
			Title:       video.Snippet.Title,
			Provider:    release.Youtube,
			Creator:     a.ownerChannelTitle(video.Snippet),
			Description: video.Snippet.Description,
			CoverURL:    coverURL(video.Snippet.Thumbnails),
			URL:         videoURL(video.ID),
		}
	}
	return tracks, nil
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	searchResults, err := a.client.searchPlaylists(ctx, artist+" – "+title)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlist on youtube: %w", err)
	}

	albums := make([]release.SearchAlbum, len(searchResults))
	for i, searchResult := range searchResults {
		playlist, err := a.client.fetchPlaylist(ctx, searchResult.ID.PlaylistID)
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
		}
		albums[i] = release.SearchAlbum{
			ID:          playlist.ID,
			Title:       playlist.Snippet.Title,
			URL:         playlistURL(playlist.ID),
			Provider:    release.Youtube,
			Creator:     a.ownerChannelTitle(playlist.Snippet),
			Description: playlist.Snippet.Description,
			CoverURL:    coverURL(playlist.Snippet.Thumbnails),
		}
	}
	return albums, nil
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}

func (a *Adapter) ownerChannelTitle(s snippet) string {
	if s.VideoOwnerChannelTitle != "" {
		return s.VideoOwnerChannelTitle
	}
	return s.ChannelTitle
}

func coverURL(thumbnails thumbnails) string {
	for _, thumbnail := range []thumbnail{
		thumbnails.Maxres,
		thumbnails.Standard,
		thumbnails.High,
		thumbnails.Medium,
		thumbnails.Default,
	} {
		if thumbnail.URL != "" {
			return thumbnail.URL
		}
	}
	return ""
}
