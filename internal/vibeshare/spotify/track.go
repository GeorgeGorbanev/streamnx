package spotify

import (
	"fmt"
	"regexp"
	"unicode"
)

type Track struct {
	Artists []Artist `json:"artists"`
	ID      string   `json:"id"`
	Name    string   `json:"name"`
}

type Artist struct {
	Name string `json:"name"`
}

var TrackRe = regexp.MustCompile(`https://open\.spotify\.com/track/([a-zA-Z0-9]+)(?:\?.*)?`)

func DetectTrackID(openTrackURL string) string {
	match := TrackRe.FindStringSubmatch(openTrackURL)
	if match == nil || len(match) < 2 {
		return ""
	}
	return match[1]
}

func (t *Track) URL() string {
	return fmt.Sprintf("https://open.spotify.com/track/%s", t.ID)
}

func (t *Track) NameContainsRussianLetters() bool {
	for _, char := range t.Name {
		if unicode.Is(unicode.Cyrillic, char) {
			return true
		}
	}
	return false
}
