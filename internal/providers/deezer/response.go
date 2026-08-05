package deezer

type track struct {
	ID          int       `json:"id"`
	ISRC        string    `json:"isrc"`
	Title       string    `json:"title"`
	Artist      artist    `json:"artist"`
	Album       albumInfo `json:"album"`
	Duration    int       `json:"duration"`
	ReleaseDate string    `json:"release_date"`
}

type album struct {
	ID          int       `json:"id"`
	UPC         string    `json:"upc"`
	Title       string    `json:"title"`
	Label       string    `json:"label"`
	Artist      artist    `json:"artist"`
	Tracks      trackData `json:"tracks"`
	ReleaseDate string    `json:"release_date"`
	coverURLs
}

type albumInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	coverURLs
}

type coverURLs struct {
	Cover       string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	CoverXL     string `json:"cover_xl"`
}

type trackData struct {
	Data []track `json:"data"`
}

type artist struct {
	Name string `json:"name"`
}
