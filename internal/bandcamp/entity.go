package bandcamp

import "regexp"

var (
	albumRe = regexp.MustCompile(`^https?://([^.]+)\.bandcamp\.com/album/([^/?#]+)`)
	trackRe = regexp.MustCompile(`^https?://([^.]+)\.bandcamp\.com/track/([^/?#]+)`)
)

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
