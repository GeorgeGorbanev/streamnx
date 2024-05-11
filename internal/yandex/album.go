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

var (
	AlbumRe = regexp.MustCompile(
		fmt.Sprintf(
			`https://music\.yandex\.(%s)/album/(\d+)`, allDomainZonesRe(),
		),
	)
)

func (a *Album) URL() string {
	return fmt.Sprintf("https://music.yandex.%s/album/%d", defaultDomainZone, a.ID)
}
