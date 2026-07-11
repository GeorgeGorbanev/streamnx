package apple

import (
	"errors"
	"regexp"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

var (
	albumRe      = regexp.MustCompile(`music\.apple\.com/(\w+)/album/.*/(\d+)`)
	albumTrackRe = regexp.MustCompile(`music\.apple\.com/(\w+)/album/.*/(\d+)\?i=(\d+)`)
	songRe       = regexp.MustCompile(`music\.apple\.com/(\w+)/song/.*/(\d+)`)
)

func parseTrackLink(trackURL string) (compositekey.Key, error) {
	if matches := albumTrackRe.FindStringSubmatch(trackURL); len(matches) == 4 {
		return newKey(matches[1], matches[3])
	}
	if matches := songRe.FindStringSubmatch(trackURL); len(matches) == 3 {
		return newKey(matches[1], matches[2])
	}
	return compositekey.Key{}, errors.New("invalid track url")
}

func parseAlbumLink(albumURL string) (compositekey.Key, error) {
	matches := albumRe.FindStringSubmatch(albumURL)
	if len(matches) != 3 {
		return compositekey.Key{}, errors.New("invalid album url")
	}
	return newKey(matches[1], matches[2])
}
