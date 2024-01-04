package ymusic

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

type Album struct {
	ID int `json:"id"`
}

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var TrackURLRegExp = regexp.MustCompile(`https://music\.yandex\.ru/album/\d+/track/(\d+)`)

func IsTrackURL(url string) bool {
	return TrackURLRegExp.MatchString(url)
}

func ParseTrackID(trackURL string) string {
	matches := TrackURLRegExp.FindStringSubmatch(trackURL)

	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

func (t *Track) URL() string {
	return fmt.Sprintf("https://music.yandex.com/album/%d/track/%s", t.Albums[0].ID, t.idString())
}

func (t *Track) FullTitle() string {
	return fmt.Sprintf("%s - %s", t.Artists[0].Name, t.Title)
}

func (t *Track) idString() string {
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
