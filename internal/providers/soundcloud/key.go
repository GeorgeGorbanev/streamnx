package soundcloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

const (
	keyPartUserSlug = "user_slug"
	keyPartID       = "id"
	keyDelimiter    = ":"
)

var keyScheme = compositekey.NewScheme(
	keyDelimiter,
	compositekey.Part{
		Name: keyPartUserSlug,
		Re:   regexp.MustCompile(`^[^:\s]+$`),
	},
	compositekey.Part{
		Name: keyPartID,
		Re:   regexp.MustCompile(`^[^:\s]+$`),
	},
)

func newKey(userSlug, id string) (compositekey.Key, error) {
	return keyScheme.NewKey(map[string]string{
		keyPartUserSlug: strings.ToLower(userSlug),
		keyPartID:       id,
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
	userSlug, err := keyScheme.Part(key, keyPartUserSlug)
	if err != nil {
		return "", "", err
	}
	return partID, userSlug, nil
}
