package youtube

import (
	"fmt"
	"regexp"
)

var (
	videoRe    = regexp.MustCompile(`(?:^|[^.[:alnum:]_-])(?:(?:https?://)?(?:www\.)?youtube\.com/watch\?(?:[^#\s]*&)?v=|(?:https?://)?youtu\.be/)([a-zA-Z0-9_-]{11})`)
	playlistRe = regexp.MustCompile(`(?:^|[^.[:alnum:]_-])(?:(?:https?://)?(?:www\.)?youtube\.com/playlist\?(?:[^#\s]*&)?list=|(?:https?://)?youtu\.be/playlist\?(?:[^#\s]*&)?list=)([a-zA-Z0-9_-]+)`)
)

func parseTrackURL(trackURL string) string {
	if matches := videoRe.FindStringSubmatch(trackURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func parseAlbumURL(albumURL string) string {
	if matches := playlistRe.FindStringSubmatch(albumURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func videoURL(id string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
}

func playlistURL(id string) string {
	return fmt.Sprintf("https://www.youtube.com/playlist?list=%s", id)
}
