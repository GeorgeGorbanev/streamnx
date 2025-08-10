package streamnx

import (
	"context"
)

type Adapter interface {
	FetchTrack(ctx context.Context, id string) (*Entity, error)
	SearchTrack(ctx context.Context, artist, title string) (*Entity, error)

	FetchAlbum(ctx context.Context, id string) (*Entity, error)
	SearchAlbum(ctx context.Context, artist, title string) (*Entity, error)

	FetchCloak(ctx context.Context, cloakCode string) (*Entity, error)
}
