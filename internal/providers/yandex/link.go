package yandex

import (
	"fmt"
	"regexp"
)

var (
	trackRe = regexp.MustCompile(`https://music\.yandex\.(com|by|kz|ru|uz)/album/\d+/track/(\d+)/?(?:[?#].*)?$`)
	albumRe = regexp.MustCompile(`https://music\.yandex\.(com|by|kz|ru|uz)/album/(\d+)/?(?:[?#].*)?$`)
)

func parseTrackLink(trackURL string) string {
	match := trackRe.FindStringSubmatch(trackURL)
	if len(match) < 3 {
		return ""
	}
	return match[2]
}

func parseAlbumLink(albumURL string) string {
	match := albumRe.FindStringSubmatch(albumURL)
	if len(match) < 3 {
		return ""
	}
	return match[2]
}

func trackLink(albumID int, trackID string) string {
	return fmt.Sprintf("https://music.yandex.com/album/%d/track/%s", albumID, trackID)
}

func albumLink(albumID int) string {
	return fmt.Sprintf("https://music.yandex.com/album/%d", albumID)
}
