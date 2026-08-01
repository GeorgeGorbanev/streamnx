package youtubemusic

import (
	"fmt"
	"regexp"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
)

type albumIDType string

const (
	albumIDTypeBrowse   albumIDType = "b"
	albumIDTypePlaylist albumIDType = "p"

	albumKeyPartIDType = "id_type"
	albumKeyPartID     = "id"
	albumKeyDelimiter  = ":"
)

var albumKeyScheme = compositekey.NewScheme(
	albumKeyDelimiter,
	compositekey.Part{
		Name: albumKeyPartIDType,
		Re:   regexp.MustCompile(`^(?:b|p)$`),
	},
	compositekey.Part{
		Name: albumKeyPartID,
		Re:   regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
	},
)

func newAlbumKey(idType albumIDType, id string) (compositekey.Key, error) {
	key, err := albumKeyScheme.NewKey(map[string]string{
		albumKeyPartIDType: string(idType),
		albumKeyPartID:     id,
	})
	if err != nil {
		return compositekey.Key{}, fmt.Errorf("failed to create album key: %w", err)
	}
	return key, nil
}

func dumpAlbumKey(key compositekey.Key) (string, error) {
	dump, err := albumKeyScheme.Dump(key)
	if err != nil {
		return "", fmt.Errorf("failed to dump album key: %w", err)
	}
	return dump, nil
}

func newAlbumID(idType albumIDType, id string) (string, error) {
	key, err := newAlbumKey(idType, id)
	if err != nil {
		return "", err
	}
	return dumpAlbumKey(key)
}

func parseAlbumID(id string) (albumIDType, string, error) {
	key, err := albumKeyScheme.Parse(id)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse album key: %w", err)
	}
	idType, err := albumKeyScheme.Part(key, albumKeyPartIDType)
	if err != nil {
		return "", "", err
	}
	rawID, err := albumKeyScheme.Part(key, albumKeyPartID)
	if err != nil {
		return "", "", err
	}
	return albumIDType(idType), rawID, nil
}
