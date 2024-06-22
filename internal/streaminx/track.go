package streaminx

type Track struct {
	ID       string
	Title    string
	Artist   string
	URL      string
	Provider *Provider
}
