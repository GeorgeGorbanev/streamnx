package spotify

import (
	"fmt"
	"regexp"
)

type Track struct {
	Artists []Artist `json:"artists"`
	ID      string   `json:"id"`
	Name    string   `json:"name"`
}

type Artist struct {
	Name string `json:"name"`
}

var TrackRe = regexp.MustCompile(`https://open\.spotify\.com/(?:[\w-]+/)?track/([a-zA-Z0-9]+)(?:\?.*)?`)

func (t *Track) URL() string {
	return fmt.Sprintf("https://open.spotify.com/track/%s", t.ID)
}
