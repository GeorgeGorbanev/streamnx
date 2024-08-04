package yandex

import (
	"fmt"
	"regexp"
)

var (
	TrackRe = regexp.MustCompile(
		fmt.Sprintf(
			`https://music\.yandex\.(%s)/album/\d+/track/(\d+)`, allDomainZonesRe(),
		),
	)
	AlbumRe = regexp.MustCompile(
		fmt.Sprintf(
			`https://music\.yandex\.(%s)/album/(\d+)`, allDomainZonesRe(),
		),
	)
)

type Track struct {
	Albums  []Album  `json:"albums"`
	Artists []Artist `json:"artists"`
	ID      any      `json:"id"`
	Title   string   `json:"title"`
}

type Album struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func DetectTrackID(trackURL string) string {
	match := TrackRe.FindStringSubmatch(trackURL)
	if match == nil || len(match) < 3 {
		return ""
	}
	return match[2]
}

func DetectAlbumID(albumURL string) string {
	match := AlbumRe.FindStringSubmatch(albumURL)
	if match == nil || len(match) < 3 {
		return ""
	}
	return match[2]
}

func (a *Album) URL() string {
	return fmt.Sprintf("https://music.yandex.%s/album/%d", noRegionDomainZone, a.ID)
}

func (t *Track) URL() string {
	return fmt.Sprintf("https://music.yandex.%s/album/%d/track/%s", noRegionDomainZone, t.Albums[0].ID, t.IDString())
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
