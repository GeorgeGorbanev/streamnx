package deezer

import (
	"fmt"
	"net/url"
	"regexp"
)

type Track struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	Artist Artist    `json:"artist"`
	Album  AlbumInfo `json:"album"`
}

type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist Artist `json:"artist"`
}

type AlbumInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Artist struct {
	Name string `json:"name"`
}

var (
	TrackRe = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/track/(\d+)`)
	AlbumRe = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/album/(\d+)`)
	CloakRe = regexp.MustCompile(`https://link\.deezer\.com/s/([A-Za-z0-9]+)`)
)

func DetectCloakID(url string) string {
	match := CloakRe.FindStringSubmatch(url)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func DetectTrackID(trackURL string) string {
	match := TrackRe.FindStringSubmatch(trackURL)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func DetectAlbumID(albumURL string) string {
	match := AlbumRe.FindStringSubmatch(albumURL)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func DetectCloakDest(cloak string) (string, error) {
	const destParam = "dest"
	u, err := url.Parse(cloak)
	if err != nil {
		return "", fmt.Errorf("failed to parse cloak URL: %w", err)
	}
	val, ok := u.Query()[destParam]
	if !ok || len(val) == 0 {
		return "", fmt.Errorf("failed to find cloak '%s' param in url: %s", destParam, cloak)
	}
	return val[0], nil
}

func (t *Track) URL() string {
	return fmt.Sprintf("https://deezer.com/track/%d", t.ID)
}

func (a *Album) URL() string {
	return fmt.Sprintf("https://deezer.com/album/%d", a.ID)
}
