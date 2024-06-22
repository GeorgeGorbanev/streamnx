package streaminx

import "fmt"

var (
	IDNotFoundError = fmt.Errorf("invalid id")
)

type Adapter interface {
	DetectTrackID(trackURL string) (string, error)
	GetTrack(id string) (*Track, error)
	SearchTrack(artistName, trackName string) (*Track, error)

	DetectAlbumID(albumURL string) (string, error)
	GetAlbum(id string) (*Album, error)
	SearchAlbum(artistName, albumName string) (*Album, error)
}
