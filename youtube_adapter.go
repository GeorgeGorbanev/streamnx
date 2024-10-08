package streamnx

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/GeorgeGorbanev/streamnx/internal/youtube"
)

type YoutubeAdapter struct {
	client youtube.Client
}

var (
	nonTitleContentRe = regexp.MustCompile(`\s*\[.*?\]|\s*\{.*?\}|\s*\(.*?\)`)
	titleSeparators   = []string{" - ", " – ", " — ", "|"}
)

func newYoutubeAdapter(client youtube.Client) *YoutubeAdapter {
	return &YoutubeAdapter{
		client: client,
	}
}
func (a *YoutubeAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	video, err := a.client.GetVideo(ctx, id)
	if err != nil {
		if errors.Is(err, youtube.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get video from youtube: %w", err)
	}
	return a.adaptTrack(video), nil
}

func (a *YoutubeAdapter) SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error) {
	query := entityFullTitle(artistName, trackName)
	search, err := a.client.SearchVideo(ctx, query)
	if err != nil {
		if errors.Is(err, youtube.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search video on youtube: %w", err)
	}

	id := search.Items[0].ID.VideoID
	video, err := a.client.GetVideo(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get video from youtube: %w", err)
	}

	return a.adaptTrack(video), nil
}

func (a *YoutubeAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	album, err := a.client.GetPlaylist(ctx, id)
	if err != nil {
		if errors.Is(err, youtube.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
	}
	return a.adaptAlbum(ctx, album)
}

func (a *YoutubeAdapter) SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error) {
	query := entityFullTitle(artistName, albumName)
	search, err := a.client.SearchPlaylist(ctx, query)
	if err != nil {
		if errors.Is(err, youtube.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to search playlist on youtube: %w", err)
	}

	id := search.Items[0].ID.PlaylistID
	album, err := a.client.GetPlaylist(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
	}

	return a.adaptAlbum(ctx, album)
}

func (a *YoutubeAdapter) adaptTrack(video *youtube.Video) *Entity {
	trackTitle := a.extractTrackTitle(video)
	artist, track := a.cleanAndSplitTitle(trackTitle)

	return &Entity{
		ID:       video.ID,
		Title:    track,
		Artist:   artist,
		URL:      video.URL(),
		Provider: Youtube,
		Type:     Track,
	}
}

func (a *YoutubeAdapter) extractTrackTitle(video *youtube.Video) string {
	if video.IsAutogenerated() {
		return fmt.Sprintf("%s - %s", video.Artist(), video.Title)
	}
	return video.Title
}

func (a *YoutubeAdapter) adaptAlbum(ctx context.Context, playlist *youtube.Playlist) (*Entity, error) {
	albumTitle, err := a.extractAlbumTitle(ctx, playlist)
	if err != nil {
		return nil, fmt.Errorf("failed to extract album title: %w", err)
	}

	artist, album := a.cleanAndSplitTitle(albumTitle)

	return &Entity{
		ID:       playlist.ID,
		Title:    album,
		Artist:   artist,
		URL:      playlist.URL(),
		Provider: Youtube,
		Type:     Album,
	}, nil
}

func (a *YoutubeAdapter) extractAlbumTitle(ctx context.Context, playlist *youtube.Playlist) (string, error) {
	if playlist.IsAutogenerated() {
		videos, err := a.client.GetPlaylistItems(ctx, playlist.ID)
		if err != nil {
			return "", fmt.Errorf("failed to get playlist items from youtube: %w", err)
		}
		if len(videos) == 0 {
			return playlist.Title, nil
		}

		v := videos[0]
		if !v.IsAutogenerated() {
			return playlist.Title, nil
		}

		return fmt.Sprintf("%s - %s", v.Artist(), playlist.Album()), nil
	}
	return playlist.Title, nil
}

func (a *YoutubeAdapter) cleanAndSplitTitle(title string) (artist, entity string) {
	cleanTitle := nonTitleContentRe.ReplaceAllString(title, "")

	for _, sep := range titleSeparators {
		if strings.Contains(cleanTitle, sep) {
			parts := strings.Split(cleanTitle, sep)
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}
	}

	words := strings.Fields(cleanTitle)
	if len(words) > 1 {
		return strings.TrimSpace(words[0]), strings.TrimSpace(strings.Join(words[1:], " "))
	}

	return cleanTitle, ""
}
