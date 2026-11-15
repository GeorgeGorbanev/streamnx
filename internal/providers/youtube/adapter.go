package youtube

import (
	"context"
	"errors"
	"fmt"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/youtube/client"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

type Adapter struct {
	client adapterClient
}

type adapterClient interface {
	FetchVideo(ctx context.Context, id string) (client.Video, error)
	FetchPlaylist(ctx context.Context, id string) (client.Metadata, error)
	FetchPlaylistWithItems(ctx context.Context, id string) (client.Metadata, []string, error)
	SearchVideos(ctx context.Context, term string) ([]string, error)
	SearchPlaylists(ctx context.Context, term string) ([]string, error)
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
	video, err := a.client.FetchVideo(ctx, id)
	switch {
	case errors.Is(err, client.ErrNotFound):
		return release.Track{}, release.ErrNotFound
	case err != nil:
		return release.Track{}, fmt.Errorf("failed to get video from youtube: %w", err)
	}

	return release.Track{
		ID:          video.ID,
		Title:       video.Title,
		Provider:    release.Youtube,
		Creator:     video.Creator,
		Description: video.Description,
		CoverURL:    video.CoverURL,
		Duration:    video.Duration,
		URL:         videoURL(video.ID),
	}, nil
}

func (a *Adapter) FetchTracksByISRC(context.Context, string) ([]release.Track, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) FetchAlbum(ctx context.Context, id string) (release.Album, error) {
	album, items, err := a.client.FetchPlaylistWithItems(ctx, id)
	switch {
	case errors.Is(err, client.ErrNotFound):
		return release.Album{}, release.ErrNotFound
	case err != nil:
		return release.Album{}, err
	}

	trackIDs := make([]string, 0, len(items))
	for _, videoID := range items {
		if videoID != "" {
			trackIDs = append(trackIDs, videoID)
		}
	}

	return release.Album{
		ID:          album.ID,
		Title:       album.Title,
		URL:         playlistURL(album.ID),
		Provider:    release.Youtube,
		Creator:     album.Creator,
		Description: album.Description,
		CoverURL:    album.CoverURL,
		TrackIDs:    trackIDs,
	}, nil
}

func (a *Adapter) FetchAlbumsByUPC(context.Context, string) ([]release.Album, error) {
	return nil, release.ErrUnsupportedOperation
}

func (a *Adapter) SearchTracks(ctx context.Context, artist, title string) ([]release.SearchTrack, error) {
	searchResults, err := a.client.SearchVideos(ctx, artist+" – "+title)
	if err != nil {
		return nil, fmt.Errorf("failed to search video on youtube: %w", err)
	}

	tracks := make([]release.SearchTrack, 0, len(searchResults))
	for _, id := range searchResults {
		video, err := a.client.FetchVideo(ctx, id)
		switch {
		case errors.Is(err, client.ErrNotFound):
			continue
		case err != nil:
			return nil, fmt.Errorf("failed to get video from youtube: %w", err)
		}
		tracks = append(tracks, release.SearchTrack{
			ID:          video.ID,
			Title:       video.Title,
			Provider:    release.Youtube,
			Creator:     video.Creator,
			Description: video.Description,
			CoverURL:    video.CoverURL,
			URL:         videoURL(video.ID),
		})
	}
	return tracks, nil
}

func (a *Adapter) SearchAlbums(ctx context.Context, artist, title string) ([]release.SearchAlbum, error) {
	searchResults, err := a.client.SearchPlaylists(ctx, artist+" – "+title)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlist on youtube: %w", err)
	}

	albums := make([]release.SearchAlbum, 0, len(searchResults))
	for _, id := range searchResults {
		playlist, err := a.client.FetchPlaylist(ctx, id)
		switch {
		case errors.Is(err, client.ErrNotFound):
			continue
		case err != nil:
			return nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
		}
		albums = append(albums, release.SearchAlbum{
			ID:          playlist.ID,
			Title:       playlist.Title,
			URL:         playlistURL(playlist.ID),
			Provider:    release.Youtube,
			Creator:     playlist.Creator,
			Description: playlist.Description,
			CoverURL:    playlist.CoverURL,
		})
	}
	return albums, nil
}

func (a *Adapter) Uncloak(context.Context, string) (release.Type, string, error) {
	return "", "", release.ErrUnsupportedOperation
}
