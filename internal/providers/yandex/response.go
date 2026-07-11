package yandex

type track struct {
	ID         string     `json:"id"`
	CoverURI   string     `json:"coverUri"`
	Artists    []artist   `json:"artists"`
	Title      string     `json:"title"`
	DurationMS int        `json:"durationMs"`
	Albums     []albumRef `json:"albums"`
	Error      string     `json:"error"`
}

type searchTrack struct {
	ID       int        `json:"id"`
	CoverURI string     `json:"coverUri"`
	Title    string     `json:"title"`
	Artists  []artist   `json:"artists"`
	Albums   []albumRef `json:"albums"`
}

type album struct {
	ID          int       `json:"id"`
	CoverURI    string    `json:"coverUri"`
	Artists     []artist  `json:"artists"`
	Labels      []label   `json:"labels"`
	Title       string    `json:"title"`
	Volumes     [][]track `json:"volumes"`
	ReleaseDate string    `json:"releaseDate"`
}

type searchAlbum struct {
	ID       int      `json:"id"`
	CoverURI string   `json:"coverUri"`
	Artists  []artist `json:"artists"`
	Title    string   `json:"title"`
}

type albumRef struct {
	ID          int      `json:"id"`
	CoverURI    string   `json:"coverUri"`
	Artists     []artist `json:"artists"`
	Title       string   `json:"title"`
	ReleaseDate string   `json:"releaseDate"`
}

type artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type label struct {
	Name string `json:"name"`
}

type apiError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
