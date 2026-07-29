package yandex

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

var (
	trackRe = regexp.MustCompile(`https://music\.yandex\.(com|by|kz|ru|uz)/album/(\d+)/track/(\d+)/?(?:[?#].*)?$`)
	albumRe = regexp.MustCompile(`https://music\.yandex\.(com|by|kz|ru|uz)/album/(\d+)/?(?:[?#].*)?$`)
)

func parseTrackLink(trackURL string) (compositekey.Key, error) {
	match := trackRe.FindStringSubmatch(trackURL)
	if len(match) != 4 {
		return compositekey.Key{}, errors.New("invalid track url")
	}
	return newTrackKey(match[2], match[3])
}

func parseAlbumLink(albumURL string) string {
	match := albumRe.FindStringSubmatch(albumURL)
	if len(match) < 3 {
		return ""
	}
	return match[2]
}

func trackLink(albumID, trackID string) string {
	return fmt.Sprintf("https://music.yandex.com/album/%s/track/%s", albumID, trackID)
}

func albumLink(albumID int) string {
	return fmt.Sprintf("https://music.yandex.com/album/%d", albumID)
}
