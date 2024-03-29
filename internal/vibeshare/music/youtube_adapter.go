package music

import (
	"fmt"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/youtube"
)

type YoutubeAdapter struct {
	client youtube.Client
}

func newYoutubeAdapter(client youtube.Client) *YoutubeAdapter {
	return &YoutubeAdapter{
		client: client,
	}
}

func (a *YoutubeAdapter) DetectTrackID(trackURL string) string {
	if matches := youtube.VideoRe.FindStringSubmatch(trackURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (a *YoutubeAdapter) DetectAlbumID(albumURL string) string {
	if matches := youtube.PlaylistRe.FindStringSubmatch(albumURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (a *YoutubeAdapter) GetTrack(id string) (*Track, error) {
	video, err := a.client.GetVideo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get video from youtube: %w", err)
	}
	if video == nil {
		return nil, nil
	}
	return a.adaptTrack(video), nil
}

func (a *YoutubeAdapter) SearchTrack(artistName, trackName string) (*Track, error) {
	query := fmt.Sprintf("%s – %s", artistName, trackName)
	video, err := a.client.SearchVideo(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search video on youtube: %w", err)
	}
	if video == nil {
		return nil, nil
	}

	return a.adaptTrack(video), nil
}

func (a *YoutubeAdapter) GetAlbum(id string) (*Album, error) {
	album, err := a.client.GetPlaylist(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
	}
	if album == nil {
		return nil, nil
	}
	return a.adaptAlbum(album), nil
}

func (a *YoutubeAdapter) SearchAlbum(artistName, albumName string) (*Album, error) {
	query := fmt.Sprintf("%s – %s", artistName, albumName)
	album, err := a.client.SearchPlaylist(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlist on youtube: %w", err)
	}
	if album == nil {
		return nil, nil
	}

	return a.adaptAlbum(album), nil
}

func (a *YoutubeAdapter) adaptTrack(video *youtube.Video) *Track {
	artist, track := a.parseTitle(video.Title)

	return &Track{
		ID:       video.ID,
		Title:    track,
		Artist:   artist,
		URL:      video.URL(),
		Provider: Youtube,
	}
}

func (a *YoutubeAdapter) adaptAlbum(playlist *youtube.Playlist) *Album {
	artist, album := a.parseTitle(playlist.Title)

	return &Album{
		ID:       playlist.ID,
		Title:    album,
		Artist:   artist,
		URL:      playlist.URL(),
		Provider: Youtube,
	}
}

func (a *YoutubeAdapter) parseTitle(input string) (string, string) {
	splitted := strings.Split(input, " – ")
	if len(splitted) != 2 {
		return input, ""
	}
	return splitted[0], splitted[1]
}
