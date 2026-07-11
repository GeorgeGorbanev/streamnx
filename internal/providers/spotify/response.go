package spotify

type track struct {
	Album      albumInfo `json:"album"`
	Artists    []artist  `json:"artists"`
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	DurationMS int       `json:"duration_ms"`
}

type album struct {
	ID          string         `json:"id"`
	Label       string         `json:"label"`
	Name        string         `json:"name"`
	Artists     []artist       `json:"artists"`
	Images      []spotifyImage `json:"images"`
	Tracks      albumTracks    `json:"tracks"`
	ReleaseDate string         `json:"release_date"`
}

type albumInfo struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Images      []spotifyImage `json:"images"`
	ReleaseDate string         `json:"release_date"`
}

type spotifyImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type artist struct {
	Name string `json:"name"`
}

type albumTracks struct {
	Items []track `json:"items"`
}
