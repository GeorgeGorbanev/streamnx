package deezer

import "regexp"

type Track struct {
}

type Album struct {
}

var (
	TrackRe         = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/track/(\d+)$`)
	AlbumRe         = regexp.MustCompile(`^https://(?:www\.)?deezer\.com(?:/[^/]+)?/album/(\d+)$`)
	UnknownEntityRe = regexp.MustCompile(`https://dzr\.page\.link/([A-Za-z0-9]+)`)
)

func DetectUnknownEntityID(url string) string {
	match := UnknownEntityRe.FindStringSubmatch(url)
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
