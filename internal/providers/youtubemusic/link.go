package youtubemusic

import (
	"fmt"
	"regexp"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release/compositekey"
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

func parseAlbumURL(rawURL string) (compositekey.Key, bool) {
	if matches := playlistRe.FindStringSubmatch(rawURL); len(matches) > 1 {
		key, err := newAlbumKey(albumIDTypePlaylist, matches[1])
		return key, err == nil
	}
	if matches := browseAlbumRe.FindStringSubmatch(rawURL); len(matches) > 1 {
		key, err := newAlbumKey(albumIDTypeBrowse, matches[1])
		return key, err == nil
	}
	return compositekey.Key{}, false
}

func trackURL(id string) string {
	return fmt.Sprintf("https://music.youtube.com/watch?v=%s", id)
}

func albumURL(idType albumIDType, id string) string {
	if idType == albumIDTypeBrowse {
		return fmt.Sprintf("https://music.youtube.com/browse/%s", id)
	}
	return fmt.Sprintf("https://music.youtube.com/playlist?list=%s", id)
}
