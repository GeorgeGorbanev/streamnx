package soundcloud

import (
	"errors"
	"fmt"
	"regexp"
)

type CompositeKey struct {
	ID       string
	UserSlug string
}

const delimiter = ":"

var compositeKeyRe = regexp.MustCompile(
	fmt.Sprintf(`^([^:%s]+)%s([^:%s]+)$`, `\s`, delimiter, `\s`),
)

func (k *CompositeKey) ParseFromTrackURL(url string) error {
	matches := TrackRe.FindStringSubmatch(url)
	if len(matches) != 3 {
		return errors.New("invalid track url")
	}

	k.UserSlug = matches[1]
	k.ID = matches[2]

	return nil
}

func (k *CompositeKey) ParseFromAlbumURL(url string) error {
	matches := AlbumRe.FindStringSubmatch(url)
	if len(matches) != 3 {
		return errors.New("invalid album url")
	}

	k.UserSlug = matches[1]
	k.ID = matches[2]

	return nil
}

func (k *CompositeKey) Marshal() string {
	return k.UserSlug + delimiter + k.ID
}

func (k *CompositeKey) Unmarshal(s string) error {
	matches := compositeKeyRe.FindStringSubmatch(s)
	if len(matches) != 3 {
		return fmt.Errorf("invalid composite key: %s", s)
	}

	k.UserSlug = matches[1]
	k.ID = matches[2]

	return nil
}
