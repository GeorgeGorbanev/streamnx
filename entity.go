package streamnx

const (
	Track   EntityType = "track"
	Album   EntityType = "album"
	Unknown EntityType = "unknown"
)

type EntityType string

type Entity struct {
	ID       string
	Title    string
	Artist   string
	URL      string
	Provider *Provider
	Type     EntityType
}

func entityFullTitle(artist, title string) string {
	return artist + " – " + title
}
