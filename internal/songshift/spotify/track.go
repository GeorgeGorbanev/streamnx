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

type TracksSection struct {
	Items []*Track `json:"items"`
}

var OpenTrackRe = regexp.MustCompile(`https://open\.spotify\.com/track/([a-zA-Z0-9]+)(?:\?.*)?`)

func DetectTrackID(openTrackURL string) string {
	match := OpenTrackRe.FindStringSubmatch(openTrackURL)
	if match == nil || len(match) < 2 {
		return ""
	}
	return match[1]
}

func (t *Track) Title() string {
	return fmt.Sprintf("%s – %s", t.Artists[0].Name, t.Name)
}

func (t *Track) URL() string {
	return fmt.Sprintf("https://open.spotify.com/track/%s", t.ID)
}
