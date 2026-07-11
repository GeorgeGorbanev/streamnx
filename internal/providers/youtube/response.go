package youtube

type snippet struct {
	Title                  string     `json:"title"`
	ChannelTitle           string     `json:"channelTitle"`
	Description            string     `json:"description"`
	Thumbnails             thumbnails `json:"thumbnails"`
	VideoOwnerChannelTitle string     `json:"videoOwnerChannelTitle"`
	ResourceID             resourceID `json:"resourceId"`
}

type thumbnails struct {
	Default  thumbnail `json:"default"`
	Medium   thumbnail `json:"medium"`
	High     thumbnail `json:"high"`
	Standard thumbnail `json:"standard"`
	Maxres   thumbnail `json:"maxres"`
}

type thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type resourceID struct {
	VideoID string `json:"videoId"`
}

type contentDetails struct {
	Duration string `json:"duration"`
}

type videoResponse struct {
	Items []video `json:"items"`
}

type video struct {
	ID             string         `json:"id"`
	Snippet        snippet        `json:"snippet"`
	ContentDetails contentDetails `json:"contentDetails"`
}

type playlistResponse struct {
	Items []playlist `json:"items"`
}

type playlist struct {
	ID      string  `json:"id"`
	Snippet snippet `json:"snippet"`
}

type playlistItemsResponse struct {
	Items []playlistItem `json:"items"`
}

type playlistItem struct {
	ID      string  `json:"id"`
	Snippet snippet `json:"snippet"`
}

type videoSearchResponse struct {
	Items []videoSearchResult `json:"items"`
}

type videoSearchResult struct {
	ID videoSearchID `json:"id"`
}

type videoSearchID struct {
	VideoID string `json:"videoId"`
}

type playlistSearchResponse struct {
	Items []playlistSearchResult `json:"items"`
}

type playlistSearchResult struct {
	ID playlistSearchID `json:"id"`
}

type playlistSearchID struct {
	PlaylistID string `json:"playlistId"`
}
