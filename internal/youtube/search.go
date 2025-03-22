package youtube

type SearchResponse struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	ID SearchID `json:"id"`
}

type SearchID struct {
	VideoID    string `json:"videoId"`
	PlaylistID string `json:"playlistId"`
}
