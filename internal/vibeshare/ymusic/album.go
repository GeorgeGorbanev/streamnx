package ymusic

import (
	"fmt"
	"regexp"
)

type Album struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Artists []Artist `json:"artists"`
}

var AlbumRe = regexp.MustCompile(`https://music\.yandex\.(ru|com)/album/(\d+)`)

func DetectAlbumID(trackURL string) string {
	matches := AlbumRe.FindStringSubmatch(trackURL)

	if len(matches) < 3 {
		return ""
	}

	return matches[2]
}

func (a *Album) URL() string {
	return fmt.Sprintf("https://music.yandex.com/album/%d", a.ID)
}
