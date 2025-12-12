package apple

import (
	"errors"
	"fmt"
	"regexp"
)

type CompositeKey struct {
	ID         string
	Storefront string
}

const delimiter = "-"

var compositeKeyRe = regexp.MustCompile(
	fmt.Sprintf(`^([a-z]{2})%s([0-9]+)$`, delimiter),
)

func (k *CompositeKey) ParseFromTrackURL(url string) error {
	if matches := AlbumTrackRe.FindStringSubmatch(url); len(matches) == 4 {
		if !IsValidStorefront(matches[1]) {
			return fmt.Errorf("invalid storefront: %s", matches[1])
		}
		k.Storefront = matches[1]
		k.ID = matches[3]
		return nil
	}
	if matches := SongRe.FindStringSubmatch(url); len(matches) == 3 {
		if !IsValidStorefront(matches[1]) {
			return fmt.Errorf("invalid storefront: %s", matches[1])
		}
		k.Storefront = matches[1]
		k.ID = matches[2]
		return nil
	}
	return errors.New("invalid track url")
}

func (k *CompositeKey) ParseFromAlbumURL(url string) error {
	matches := AlbumRe.FindStringSubmatch(url)
	if len(matches) != 3 {
		return errors.New("invalid album url")
	}
	if !IsValidStorefront(matches[1]) {
		return fmt.Errorf("invalid storefront: %s", matches[1])
	}

	k.Storefront = matches[1]
	k.ID = matches[2]
	return nil
}

func (k *CompositeKey) Marshal() string {
	return k.Storefront + delimiter + k.ID
}

func (k *CompositeKey) Unmarshal(s string) error {
	matches := compositeKeyRe.FindStringSubmatch(s)
	if len(matches) < 3 {
		return fmt.Errorf("invalid composite key: %s", s)
	}

	k.Storefront = matches[1]
	k.ID = matches[2]

	return nil
}
