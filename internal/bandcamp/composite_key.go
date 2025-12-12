package bandcamp

import (
	"errors"
	"fmt"
	"regexp"
)

type CompositeKey struct {
	ID        string
	Subdomain string
}

const delimiter = ":"

var compositeKeyRe = regexp.MustCompile(
	fmt.Sprintf(`^([a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)%s([^:\s]+)$`, delimiter),
)

func (k *CompositeKey) ParseFromAlbumURL(url string) error {
	matches := albumRe.FindStringSubmatch(url)
	if len(matches) != 3 {
		return errors.New("invalid album url")
	}
	k.ID = matches[2]
	k.Subdomain = matches[1]
	return nil
}

func (k *CompositeKey) ParseFromTrackURL(url string) error {
	matches := trackRe.FindStringSubmatch(url)
	if len(matches) != 3 {
		return errors.New("invalid track url")
	}
	k.Subdomain = matches[1]
	k.ID = matches[2]
	return nil
}

func (k *CompositeKey) Marshal() string {
	return k.Subdomain + delimiter + k.ID
}

func (k *CompositeKey) Unmarshal(s string) error {
	matches := compositeKeyRe.FindStringSubmatch(s)
	if len(matches) != 3 {
		return fmt.Errorf("invalid composite key: %s", s)
	}
	k.Subdomain = matches[1]
	k.ID = matches[2]
	return nil
}
