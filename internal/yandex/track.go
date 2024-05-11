package yandex

import (
	"fmt"
	"regexp"
)

type Track struct {
	Albums  []Album  `json:"albums"`
	Artists []Artist `json:"artists"`
	ID      any      `json:"id"`
	Title   string   `json:"title"`
}

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	TrackRe = regexp.MustCompile(
		fmt.Sprintf(
			`https://music\.yandex\.(%s)/album/\d+/track/(\d+)`, allDomainZonesRe(),
		),
	)
)

func (t *Track) URL() string {
	return fmt.Sprintf("https://music.yandex.%s/album/%d/track/%s", defaultDomainZone, t.Albums[0].ID, t.IDString())
}

func (t *Track) IDString() string {
	switch id := t.ID.(type) {
	case int:
		return fmt.Sprintf("%d", id)
	case string:
		return id
	case float64:
		return fmt.Sprintf("%d", int(id))
	default:
		return ""
	}
}
