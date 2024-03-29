package youtube

import (
	"fmt"
	"regexp"
)

type Playlist struct {
	ID           string
	Title        string
	ChannelTitle string
}

var PlaylistRe = regexp.MustCompile(`(?:youtube\.com/playlist\?list=|youtu\.be/playlist\?list=)([a-zA-Z0-9_-]{18,34})`)

func (p *Playlist) URL() string {
	return fmt.Sprintf("https://www.youtube.com/playlist?list=%s", p.ID)
}
