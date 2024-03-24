package yandex

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

func (a *Album) URL() string {
	return fmt.Sprintf("https://music.yandex.com/album/%d", a.ID)
}
