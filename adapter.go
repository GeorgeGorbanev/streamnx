package streaminx

import (
	"context"
	"fmt"
)

var (
	IDNotFoundError = fmt.Errorf("invalid id")
)

type Adapter interface {
	DetectTrackID(trackURL string) (string, error)
	GetTrack(ctx context.Context, id string) (*Track, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Track, error)

	DetectAlbumID(albumURL string) (string, error)
	GetAlbum(ctx context.Context, id string) (*Album, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Album, error)
}
