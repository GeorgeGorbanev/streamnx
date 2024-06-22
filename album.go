package streaminx

type Album struct {
	ID       string
	Title    string
	Artist   string
	URL      string
	Provider *Provider
}
