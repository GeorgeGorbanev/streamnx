package bandcamp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

const (
	keyPartArtistSlug = "artist_slug"
	keyPartID         = "id"
	keyDelimiter      = ":"
)

var keyScheme = compositekey.NewScheme(
	keyDelimiter,
	compositekey.Part{
		Name: keyPartArtistSlug,
		Re:   regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`),
	},
	compositekey.Part{
		Name: keyPartID,
		Re:   regexp.MustCompile(`^[^:\s]+$`),
	},
)

func newKey(artistSlug, id string) (compositekey.Key, error) {
	return keyScheme.NewKey(map[string]string{
		keyPartArtistSlug: strings.ToLower(artistSlug),
		keyPartID:         id,
	})
}

func parseKeyParts(id string) (string, string, error) {
	key, err := keyScheme.Parse(id)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse key parts: %w", err)
	}
	partID, err := keyScheme.Part(key, keyPartID)
	if err != nil {
		return "", "", err
	}
	artistSlug, err := keyScheme.Part(key, keyPartArtistSlug)
	if err != nil {
		return "", "", err
	}
	return partID, artistSlug, nil
}
