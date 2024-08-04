package streaminx

import (
	"context"
)

type Adapter interface {
	FetchTrack(ctx context.Context, id string) (*Entity, error)
	SearchTrack(ctx context.Context, artistName, trackName string) (*Entity, error)

	FetchAlbum(ctx context.Context, id string) (*Entity, error)
	SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error)
}
