package soundcloud

type track struct {
	ID                int64             `json:"id"`
	URN               string            `json:"urn"`
	ArtworkURL        string            `json:"artwork_url"`
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	Permalink         string            `json:"permalink"`
	PermalinkURL      string            `json:"permalink_url"`
	Duration          int               `json:"duration"`
	FullDuration      int               `json:"full_duration"`
	ReleaseDate       string            `json:"release_date"`
	User              user              `json:"user"`
	PublisherMetadata publisherMetadata `json:"publisher_metadata"`
}

type album struct {
	ArtworkURL   string  `json:"artwork_url"`
	LabelName    string  `json:"label_name"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Permalink    string  `json:"permalink"`
	PermalinkURL string  `json:"permalink_url"`
	User         user    `json:"user"`
	Tracks       []track `json:"tracks"`
	ReleaseDate  string  `json:"release_date"`
}

type user struct {
	URN          string `json:"urn"`
	Username     string `json:"username"`
	Permalink    string `json:"permalink"`
	PermalinkURL string `json:"permalink_url"`
}

type publisherMetadata struct {
	Artist     string `json:"artist"`
	AlbumTitle string `json:"album_title"`
}
