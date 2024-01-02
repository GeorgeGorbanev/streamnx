package spotify

import (
	"fmt"
	"regexp"
)

type Track struct {
	Album            Album             `json:"album"`
	Artists          []Artist          `json:"artists"`
	AvailableMarkets []string          `json:"available_markets"`
	DiscNumber       int               `json:"disc_number"`
	DurationMs       int               `json:"duration_ms"`
	Explicit         bool              `json:"explicit"`
	ExternalIDs      map[string]string `json:"external_ids"`
	ExternalURLs     map[string]string `json:"external_urls"`
	Href             string            `json:"href"`
	ID               string            `json:"id"`
	IsPlayable       bool              `json:"is_playable"`
	LinkedFrom       *Track            `json:"linked_from"`
	Name             string            `json:"name"`
	Popularity       int               `json:"popularity"`
	PreviewURL       string            `json:"preview_url"`
	TrackNumber      int               `json:"track_number"`
	Type             string            `json:"type"`
	URI              string            `json:"uri"`
}

type Album struct {
	AlbumType            string            `json:"album_type"`
	Artists              []Artist          `json:"artists"`
	AvailableMarkets     []string          `json:"available_markets"`
	ExternalURLs         map[string]string `json:"external_urls"`
	Href                 string            `json:"href"`
	ID                   string            `json:"id"`
	Images               []Image           `json:"images"`
	Name                 string            `json:"name"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	Type                 string            `json:"type"`
	URI                  string            `json:"uri"`
}

type Artist struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

var OpenTrackRe = regexp.MustCompile(`https://open\.spotify\.com/track/([a-zA-Z0-9]+)(?:\?.*)?`)

func DetectTrackID(openTrackURL string) string {
	match := OpenTrackRe.FindStringSubmatch(openTrackURL)
	if match == nil || len(match) < 2 {
		return ""
	}
	return match[1]
}

func (t *Track) Title() string {
	return fmt.Sprintf("%s – %s", t.Artists[0].Name, t.Name)
}
