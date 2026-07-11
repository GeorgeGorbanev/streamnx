package bandcamp

import (
	"errors"
	"regexp"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

var (
	albumRe = regexp.MustCompile(`^https?://([^.]+)\.bandcamp\.com/album/([^/?#]+)`)
	trackRe = regexp.MustCompile(`^https?://([^.]+)\.bandcamp\.com/track/([^/?#]+)`)
)

func parseTrackLink(trackURL string) (compositekey.Key, error) {
	matches := trackRe.FindStringSubmatch(trackURL)
	if len(matches) != 3 {
		return compositekey.Key{}, errors.New("invalid track url")
	}
	return newKey(matches[1], matches[2])
}

func parseAlbumLink(albumURL string) (compositekey.Key, error) {
	matches := albumRe.FindStringSubmatch(albumURL)
	if len(matches) != 3 {
		return compositekey.Key{}, errors.New("invalid album url")
	}
	return newKey(matches[1], matches[2])
}
