package apple

type fetchResponse struct {
	Data []entity `json:"data"`
}

type entity struct {
	ID            string              `json:"id"`
	Attributes    entityAttributes    `json:"attributes"`
	Relationships entityRelationships `json:"relationships"`
}

type entityAttributes struct {
	Name             string         `json:"name"`
	ISRC             string         `json:"isrc"`
	UPC              string         `json:"upc"`
	AlbumName        string         `json:"albumName"`
	URL              string         `json:"url"`
	Artwork          artwork        `json:"artwork"`
	ArtistName       string         `json:"artistName"`
	RecordLabel      string         `json:"recordLabel"`
	DurationInMillis int            `json:"durationInMillis"`
	ReleaseDate      string         `json:"releaseDate"`
	EditorialNotes   editorialNotes `json:"editorialNotes"`
}

type artwork struct {
	URL string `json:"url"`
}

type editorialNotes struct {
	Short    string `json:"short"`
	Standard string `json:"standard"`
}

type entityRelationships struct {
	Tracks tracksRelationship `json:"tracks"`
	Albums albumsRelationship `json:"albums"`
}

type tracksRelationship struct {
	Data []entity `json:"data"`
}

type albumsRelationship struct {
	Data []entity `json:"data"`
}

type searchResponse struct {
	Resources searchResources `json:"resources"`
	Results   searchResults   `json:"results"`
}

type searchResources struct {
	Songs  map[string]*entity `json:"songs"`
	Albums map[string]*entity `json:"albums"`
}

type searchResults struct {
	Top searchTop `json:"top"`
}

type searchTop struct {
	Data []searchTopData `json:"data"`
}

type searchTopData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
