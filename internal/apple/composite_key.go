package apple

import (
	"fmt"
	"regexp"
)

const delimiter = "-"

type CompositeKey struct {
	ID         string
	Storefront string
}

var (
	compositeKeyRe = regexp.MustCompile(
		fmt.Sprintf(`^([a-z]{2})%s([0-9]+)$`, delimiter),
	)
)

func (k *CompositeKey) ParseFromTrackURL(trackURL string) {
	matches := AlbumTrackRe.FindStringSubmatch(trackURL)
	if len(matches) == 4 {
		k.Storefront = matches[1]
		k.ID = matches[3]
		return
	}
	matches = SongRe.FindStringSubmatch(trackURL)
	if len(matches) == 3 {
		k.Storefront = matches[1]
		k.ID = matches[2]
	}
}

func (k *CompositeKey) ParseFromAlbumURL(albumURL string) {
	matches := AlbumRe.FindStringSubmatch(albumURL)
	if len(matches) == 3 {
		k.Storefront = matches[1]
		k.ID = matches[2]
	}
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
