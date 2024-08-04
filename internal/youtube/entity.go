package youtube

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	autogenVideoDescriptionSubstring = "Auto-generated by YouTube"
	autogenVideoChannelTitleSuffix   = " - Topic"
	autogenPlaylistTitlePrefix       = "Album - "
)

var (
	VideoRe    = regexp.MustCompile(`(?:youtu\.be/|youtube\.com/watch\?v=)([a-zA-Z0-9_-]{11})`)
	PlaylistRe = regexp.MustCompile(`(?:youtube\.com/playlist\?list=|youtu\.be/playlist\?list=)([a-zA-Z0-9_-]+)`)
)

type Video struct {
	ID           string
	Title        string
	ChannelTitle string
	Description  string
}
type Playlist struct {
	ID           string
	Title        string
	ChannelTitle string
}

func DetectTrackID(trackURL string) string {
	if matches := VideoRe.FindStringSubmatch(trackURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func DetectAlbumID(albumURL string) string {
	if matches := PlaylistRe.FindStringSubmatch(albumURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (v *Video) URL() string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
}

func (v *Video) IsAutogenerated() bool {
	return strings.Contains(v.Description, autogenVideoDescriptionSubstring)
}

func (v *Video) Artist() string {
	if !strings.HasSuffix(v.ChannelTitle, autogenVideoChannelTitleSuffix) {
		return ""
	}
	return strings.TrimSuffix(v.ChannelTitle, autogenVideoChannelTitleSuffix)
}

func (p *Playlist) URL() string {
	return fmt.Sprintf("https://www.youtube.com/playlist?list=%s", p.ID)
}

func (p *Playlist) IsAutogenerated() bool {
	return strings.HasPrefix(p.Title, autogenPlaylistTitlePrefix)
}

func (p *Playlist) Album() string {
	if !p.IsAutogenerated() {
		return p.Title
	}
	return strings.TrimPrefix(p.Title, autogenPlaylistTitlePrefix)
}
