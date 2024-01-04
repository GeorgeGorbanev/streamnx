package ymusic

type SearchResult struct {
	Tracks TracksSection `json:"tracks"`
}

type TracksSection struct {
	Results []Track `json:"results"`
}

func (sr *SearchResult) AnyTracksFound() bool {
	return len(sr.Tracks.Results) > 0
}
