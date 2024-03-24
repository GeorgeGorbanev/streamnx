package music

type Adapter interface {
	DetectTrackID(link string) string
	GetTrack(id string) (*Track, error)
	SearchTrack(artistName, trackName string) (*Track, error)

	DetectAlbumID(link string) string
	GetAlbum(id string) (*Album, error)
	SearchAlbum(artistName, albumName string) (*Album, error)
}

type Track struct {
	Title  string
	Artist string
	URL    string
}

type Album struct {
	Title  string
	Artist string
	URL    string
}
