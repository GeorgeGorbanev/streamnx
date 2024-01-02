package ymusic

type SearchResponse struct {
	InvocationInfo InvocationInfo `json:"invocationInfo"`
	Result         SearchResult   `json:"result"`
}

type SearchResult struct {
	Page            int           `json:"page"`
	PerPage         int           `json:"perPage"`
	SearchRequestID string        `json:"searchRequestId"`
	Text            string        `json:"text"`
	Tracks          TracksSection `json:"tracks"`
	Type            string        `json:"type"`
}

type TracksSection struct {
	Order   int     `json:"order"`
	PerPage int     `json:"perPage"`
	Results []Track `json:"results"`
	Total   int     `json:"total"`
}

func (sr *SearchResult) AnyTracksFound() bool {
	return len(sr.Tracks.Results) > 0
}
