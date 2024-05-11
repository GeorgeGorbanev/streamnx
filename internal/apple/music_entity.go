package apple

import "regexp"

type MusicEntity struct {
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	ArtistName string `json:"artistName"`
}

var (
	AlbumRe      = regexp.MustCompile(`music\.apple\.com/.*/album/.*/(\d+)`)
	AlbumTrackRe = regexp.MustCompile(`music\.apple\.com/.*/album/.*/(\d+)\?i=(\d+)`)
	SongRe       = regexp.MustCompile(`music\.apple\.com/.*/song/.*/(\d+)`)
	TrackRe      = regexp.MustCompile(AlbumTrackRe.String() + "|" + SongRe.String())
)
