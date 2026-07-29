package yandex

import (
	"fmt"
	"regexp"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

const (
	keyPartAlbumID = "album_id"
	keyPartTrackID = "track_id"
	keyDelimiter   = ":"
)

var (
	numericIDRe = regexp.MustCompile(`^[0-9]+$`)
	keyScheme   = compositekey.NewScheme(
		keyDelimiter,
		compositekey.Part{
			Name: keyPartAlbumID,
			Re:   numericIDRe,
		},
		compositekey.Part{
			Name: keyPartTrackID,
			Re:   numericIDRe,
		},
	)
)

func newTrackKey(albumID, trackID string) (compositekey.Key, error) {
	return keyScheme.NewKey(map[string]string{
		keyPartAlbumID: albumID,
		keyPartTrackID: trackID,
	})
}

func parseTrackKeyParts(id string) (string, string, error) {
	key, err := keyScheme.Parse(id)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse key parts: %w", err)
	}
	albumID, err := keyScheme.Part(key, keyPartAlbumID)
	if err != nil {
		return "", "", err
	}
	trackID, err := keyScheme.Part(key, keyPartTrackID)
	if err != nil {
		return "", "", err
	}
	return albumID, trackID, nil
}
