package soundcloud

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

var (
	trackRe = regexp.MustCompile(`^https?://(?:www\.)?soundcloud\.com/([^/?#]+)/([^/?#]+)(?:\?.*)?$`)
	albumRe = regexp.MustCompile(`^https?://(?:www\.)?soundcloud\.com/([^/?#]+)/sets/([^/?#]+)(?:\?.*)?$`)
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

func trackLink(t track) string {
	if t.PermalinkURL != "" {
		return t.PermalinkURL
	}
	return fmt.Sprintf("https://soundcloud.com/%s/%s", t.User.Permalink, t.Permalink)
}

func albumLink(a album) string {
	if a.PermalinkURL != "" {
		return a.PermalinkURL
	}
	return fmt.Sprintf("https://soundcloud.com/%s/sets/%s", a.User.Permalink, a.Permalink)
}
