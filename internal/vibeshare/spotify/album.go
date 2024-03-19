package spotify

import (
	"fmt"
	"regexp"
)

type Album struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

var AlbumRe = regexp.MustCompile(`https://open\.spotify\.com/album/([a-zA-Z0-9]+)(?:\?.*)?`)

func DetectAlbumID(albumURL string) string {
	match := AlbumRe.FindStringSubmatch(albumURL)
	if match == nil || len(match) < 2 {
		return ""
	}
	return match[1]
}

func (a *Album) URL() string {
	return fmt.Sprintf("https://open.spotify.com/album/%s", a.ID)
}
