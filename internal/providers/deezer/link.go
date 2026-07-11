package deezer

import (
	"fmt"
	"regexp"
)

var (
	trackRe = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/track/(\d+)`)
	albumRe = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/album/(\d+)`)
	cloakRe = regexp.MustCompile(`https://link\.deezer\.com/s/([A-Za-z0-9]+)`)
)

func parseCloakLink(url string) string {
	match := cloakRe.FindStringSubmatch(url)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

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

func trackLink(id int) string {
	return fmt.Sprintf("https://deezer.com/track/%d", id)
}

func albumLink(id int) string {
	return fmt.Sprintf("https://deezer.com/album/%d", id)
}
