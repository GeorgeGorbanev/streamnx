package streaminx

const (
	Track EntityType = "track"
	Album EntityType = "album"
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
