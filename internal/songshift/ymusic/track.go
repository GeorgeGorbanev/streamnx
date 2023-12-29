package ymusic

import "fmt"

type Track struct {
	Albums                         []Album       `json:"albums"`
	Artists                        []Artist      `json:"artists"`
	Available                      bool          `json:"available"`
	AvailableAsRbt                 bool          `json:"availableAsRbt"`
	AvailableForOptions            []string      `json:"availableForOptions"`
	AvailableForPremiumUsers       bool          `json:"availableForPremiumUsers"`
	AvailableFullWithoutPermission bool          `json:"availableFullWithoutPermission"`
	CoverUri                       string        `json:"coverUri"`
	DerivedColors                  DerivedColors `json:"derivedColors"`
	Disclaimers                    []any         `json:"disclaimers"`
	DurationMs                     int           `json:"durationMs"`
	Explicit                       bool          `json:"explicit"`
	Fade                           Fade          `json:"fade"`
	FileSize                       int           `json:"fileSize"`
	ID                             int           `json:"id"`
	LyricsAvailable                bool          `json:"lyricsAvailable"`
	LyricsInfo                     LyricsInfo    `json:"lyricsInfo"`
	Major                          Major         `json:"major"`
	OgImage                        string        `json:"ogImage"`
	PreviewDurationMs              int           `json:"previewDurationMs"`
	R128                           R128          `json:"r128"`
	RealID                         string        `json:"realId"`
	Regions                        []string      `json:"regions"`
	RememberPosition               bool          `json:"rememberPosition"`
	StorageDir                     string        `json:"storageDir"`
	Title                          string        `json:"title"`
	TrackSharingFlag               string        `json:"trackSharingFlag"`
	TrackSource                    string        `json:"trackSource"`
	Type                           string        `json:"type"`
}

type Album struct {
	Artists                  []Artist      `json:"artists"`
	Available                bool          `json:"available"`
	AvailableForMobile       bool          `json:"availableForMobile"`
	AvailableForOptions      []string      `json:"availableForOptions"`
	AvailableForPremiumUsers bool          `json:"availableForPremiumUsers"`
	AvailablePartially       bool          `json:"availablePartially"`
	Bests                    []int         `json:"bests"`
	CoverUri                 string        `json:"coverUri"`
	Disclaimers              []any         `json:"disclaimers"`
	Genre                    string        `json:"genre"`
	ID                       int           `json:"id"`
	Labels                   []string      `json:"labels"`
	LikesCount               int           `json:"likesCount"`
	MetaType                 string        `json:"metaType"`
	OgImage                  string        `json:"ogImage"`
	Recent                   bool          `json:"recent"`
	Regions                  []string      `json:"regions"`
	ReleaseDate              string        `json:"releaseDate"`
	StorageDir               string        `json:"storageDir"`
	Title                    string        `json:"title"`
	TrackCount               int           `json:"trackCount"`
	TrackPosition            TrackPosition `json:"trackPosition"`
	Type                     string        `json:"type"`
	VeryImportant            bool          `json:"veryImportant"`
	Year                     int           `json:"year"`
}

type Artist struct {
	Composer    bool        `json:"composer"`
	Cover       ArtistCover `json:"cover"`
	Disclaimers []any       `json:"disclaimers"`
	Genres      []string    `json:"genres"`
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Various     bool        `json:"various"`
}

type ArtistCover struct {
	Prefix string `json:"prefix"`
	Type   string `json:"type"`
	URI    string `json:"uri"`
}

type TrackPosition struct {
	Index  int `json:"index"`
	Volume int `json:"volume"`
}

type DerivedColors struct {
	Accent     string `json:"accent"`
	Average    string `json:"average"`
	MiniPlayer string `json:"miniPlayer"`
	WaveText   string `json:"waveText"`
}

type Fade struct {
	InStart  float64 `json:"inStart"`
	InStop   float64 `json:"inStop"`
	OutStart float64 `json:"outStart"`
	OutStop  float64 `json:"outStop"`
}

type LyricsInfo struct {
	HasAvailableSyncLyrics bool `json:"hasAvailableSyncLyrics"`
	HasAvailableTextLyrics bool `json:"hasAvailableTextLyrics"`
}

type Major struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type R128 struct {
	I  float64 `json:"i"`
	TP float64 `json:"tp"`
}

func (t *Track) URL() string {
	return fmt.Sprintf("https://music.yandex.com/album/%d/track/%d", t.Albums[0].ID, t.ID)
}
