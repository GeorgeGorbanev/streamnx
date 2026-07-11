package spotify

import (
	"fmt"
	"regexp"
)

var (
	trackRe = regexp.MustCompile(`https://open\.spotify\.com/(?:[\w-]+/)?track/([a-zA-Z0-9]+)(?:\?.*)?`)
	albumRe = regexp.MustCompile(`https://open\.spotify\.com/(?:[\w-]+/)?album/([a-zA-Z0-9]+)(?:\?.*)?`)
)

func parseTrackLink(trackURL string) string {
	match := trackRe.FindStringSubmatch(trackURL)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func parseAlbumLink(albumURL string) string {
	match := albumRe.FindStringSubmatch(albumURL)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func trackLink(id string) string {
	return fmt.Sprintf("https://open.spotify.com/track/%s", id)
}

func albumLink(id string) string {
	return fmt.Sprintf("https://open.spotify.com/album/%s", id)
}
