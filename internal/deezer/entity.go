package deezer

import (
	"fmt"
	"regexp"
)

type Track struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Artist Artist    `json:"artist"`
	Album  AlbumInfo `json:"album"`
}

type Album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist Artist `json:"artist"`
}

type AlbumInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Artist struct {
	Name string `json:"name"`
}

var (
	TrackRe = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/track/(\d+)$`)
	AlbumRe = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/album/(\d+)$`)
	CloakRe = regexp.MustCompile(`https://dzr\.page\.link/([A-Za-z0-9]+)`)
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

func (t *Track) URL() string {
	return fmt.Sprintf("https://deezer.com/track/%s", t.ID)
}

func (a *Album) URL() string {
	return fmt.Sprintf("https://deezer.com/album/%s", a.ID)
}
