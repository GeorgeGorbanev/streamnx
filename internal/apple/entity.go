package apple

import (
	"regexp"
)

var (
	AlbumRe      = regexp.MustCompile(`music\.apple\.com/(\w+)/album/.*/(\d+)`)
	AlbumTrackRe = regexp.MustCompile(`music\.apple\.com/(\w+)/album/.*/(\d+)\?i=(\d+)`)
	SongRe       = regexp.MustCompile(`music\.apple\.com/(\w+)/song/.*/(\d+)`)
)

type Entity struct {
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	ArtistName string `json:"artistName"`
}

func DetectTrackID(trackURL string) string {
	ck := CompositeKey{}
	if err := ck.ParseFromTrackURL(trackURL); err != nil {
		return ""
	}
	return ck.Marshal()
}

func DetectAlbumID(albumURL string) string {
	ck := CompositeKey{}
	if err := ck.ParseFromAlbumURL(albumURL); err != nil {
		return ""
	}
	return ck.Marshal()
}
