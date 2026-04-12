package soundcloud

import (
	"fmt"
	"regexp"
)

type Track struct {
	Title        string `json:"title"`
	Permalink    string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
	User         User   `json:"user"`
}

type Album struct {
	Title        string `json:"title"`
	Permalink    string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
	User         User   `json:"user"`
}

type User struct {
	URN          string `json:"urn"`
	Username     string `json:"username"`
	Permalink    string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
}

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

func (t *Track) URL() string {
	return fmt.Sprintf("https://soundcloud.com/%s/%s", t.User.Permalink, t.Permalink)
}

func (a *Album) URL() string {
	return fmt.Sprintf("https://soundcloud.com/%s/sets/%s", a.User.Permalink, a.Permalink)
}
