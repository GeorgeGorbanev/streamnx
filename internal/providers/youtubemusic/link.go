package youtubemusic

import (
	"fmt"
	"regexp"
)

var (
	trackRe       = regexp.MustCompile(`(?:^|[^.[:alnum:]_-])(?:https?://)?music\.youtube\.com/watch\?(?:[^#\s]*&)?v=([a-zA-Z0-9_-]{11})`)
	playlistRe    = regexp.MustCompile(`(?:^|[^.[:alnum:]_-])(?:https?://)?music\.youtube\.com/playlist\?(?:[^#\s]*&)?list=([a-zA-Z0-9_-]+)`)
	browseAlbumRe = regexp.MustCompile(`(?:^|[^.[:alnum:]_-])(?:https?://)?music\.youtube\.com/browse/(MPRE[a-zA-Z0-9_-]+)`)
)

func parseTrackURL(rawURL string) string {
	if matches := trackRe.FindStringSubmatch(rawURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func parseAlbumURL(rawURL string) string {
	if matches := playlistRe.FindStringSubmatch(rawURL); len(matches) > 1 {
		return matches[1]
	}
	if matches := browseAlbumRe.FindStringSubmatch(rawURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func trackURL(id string) string {
	return fmt.Sprintf("https://music.youtube.com/watch?v=%s", id)
}

func albumURL(id string) string {
	if len(id) >= 4 && id[:4] == "MPRE" {
		return fmt.Sprintf("https://music.youtube.com/browse/%s", id)
	}
	return fmt.Sprintf("https://music.youtube.com/playlist?list=%s", id)
}
