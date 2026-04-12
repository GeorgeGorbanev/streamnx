package soundcloud

import (
	"fmt"
	"regexp"
)

type Track struct {
	Kind         string `json:"kind"`
	URN          string `json:"urn"`
	Title        string `json:"title"`
	Permalink    string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
	User         User   `json:"user"`
}

type Album struct {
	Kind         string  `json:"kind"`
	URN          string  `json:"urn"`
	Title        string  `json:"title"`
	Permalink    string  `json:"permalink"`
	PermalinkURL string  `json:"permalink_url"`
	SetType      string  `json:"set_type"`
	TrackCount   int     `json:"track_count"`
	User         User    `json:"user"`
	Tracks       []Track `json:"tracks"`
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
	if t.PermalinkURL != "" {
		return t.PermalinkURL
	}
	if t.User.Permalink != "" && t.Permalink != "" {
		return fmt.Sprintf("https://soundcloud.com/%s/%s", t.User.Permalink, t.Permalink)
	}
	return ""
}

func (a *Album) URL() string {
	if a.PermalinkURL != "" {
		return a.PermalinkURL
	}
	if a.User.Permalink != "" && a.Permalink != "" {
		return fmt.Sprintf("https://soundcloud.com/%s/sets/%s", a.User.Permalink, a.Permalink)
	}
	return ""
}
