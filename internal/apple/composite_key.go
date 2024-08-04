package apple

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	delimiter = "-"
)

var (
	compositeKeyRe = regexp.MustCompile(
		fmt.Sprintf(`^([a-z]{2})%s([0-9]+)$`, delimiter),
	)
	CompositeKeyError = errors.New("invalid composite key")
)

type CompositeKey struct {
	ID         string
	Storefront string
}

func (k *CompositeKey) ParseFromTrackURL(url string) error {
	if matches := AlbumTrackRe.FindStringSubmatch(url); len(matches) == 4 {
		if !IsValidStorefront(matches[1]) {
			return fmt.Errorf("%w (invalid storefront)", CompositeKeyError)
		}
		k.Storefront = matches[1]
		k.ID = matches[3]
		return nil
	}
	if matches := SongRe.FindStringSubmatch(url); len(matches) == 3 {
		if !IsValidStorefront(matches[1]) {
			return fmt.Errorf("%w (invalid storefront)", CompositeKeyError)
		}
		k.Storefront = matches[1]
		k.ID = matches[2]
		return nil
	}
	return fmt.Errorf("%w (not valid url)", CompositeKeyError)
}

func (k *CompositeKey) ParseFromAlbumURL(url string) error {
	matches := AlbumRe.FindStringSubmatch(url)
	if len(matches) != 3 {
		return fmt.Errorf("%w (not valid url)", CompositeKeyError)
	}
	if !IsValidStorefront(matches[1]) {
		return fmt.Errorf("%w (invalid storefront)", CompositeKeyError)
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
