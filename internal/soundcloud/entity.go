package soundcloud

import "regexp"

var (
	TrackRe = regexp.MustCompile(`^https?://(?:www\.)?soundcloud\.com/([^/?#]+)/([^/?#]+)(?:\?.*)?$`)
	AlbumRe = regexp.MustCompile(`^https?://(?:www\.)?soundcloud\.com/([^/?#]+)/sets/([^/?#]+)(?:\?.*)?$`)
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
