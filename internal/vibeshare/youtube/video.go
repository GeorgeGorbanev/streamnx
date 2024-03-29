package youtube

import (
	"fmt"
	"regexp"
)

type Video struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ChannelTitle string `json:"snippet.channelTitle"`
}

var VideoRe = regexp.MustCompile(`(?:youtu\.be/|youtube\.com/watch\?v=)([a-zA-Z0-9_-]{11})`)

func (v *Video) URL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
}
