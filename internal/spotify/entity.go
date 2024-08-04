package spotify

import (
	"fmt"
	"regexp"
)

var (
	TrackRe = regexp.MustCompile(`https://open\.spotify\.com/(?:[\w-]+/)?track/([a-zA-Z0-9]+)(?:\?.*)?`)
	AlbumRe = regexp.MustCompile(`https://open\.spotify\.com/(?:[\w-]+/)?album/([a-zA-Z0-9]+)(?:\?.*)?`)
)

type Track struct {
	Artists []Artist `json:"artists"`
	ID      string   `json:"id"`
	Name    string   `json:"name"`
}

type Album struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

func DetectTrackID(trackURL string) string {
	match := TrackRe.FindStringSubmatch(trackURL)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func DetectAlbumID(albumURL string) string {
	match := AlbumRe.FindStringSubmatch(albumURL)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func (t *Track) URL() string {
	return fmt.Sprintf("https://open.spotify.com/track/%s", t.ID)
}

func (a *Album) URL() string {
	return fmt.Sprintf("https://open.spotify.com/album/%s", a.ID)
}
