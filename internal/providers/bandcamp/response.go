package bandcamp

import "regexp"

type releaseLDJSON struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Duration      string `json:"duration"`
	Image         string `json:"image"`
	DatePublished string `json:"datePublished"`
	InAlbum       struct {
		ID   string `json:"@id"`
		Name string `json:"name"`
	} `json:"inAlbum"`
	ByArtist struct {
		Name string `json:"name"`
	} `json:"byArtist"`
	Publisher struct {
		Name string `json:"name"`
	} `json:"publisher"`
	Track struct {
		ItemListElement []struct {
			Item struct {
				Name             string `json:"name"`
				MainEntityOfPage string `json:"mainEntityOfPage"`
				ID               string `json:"@id"`
			} `json:"item"`
		} `json:"itemListElement"`
	} `json:"track"`
}

type searchRequest struct {
	SearchText   string `json:"search_text"`
	SearchFilter string `json:"search_filter"`
	FullPage     bool   `json:"full_page"`
}

type searchResponse struct {
	Auto struct {
		Results []struct {
			Name        string `json:"name"`
			AlbumName   string `json:"album_name"`
			BandName    string `json:"band_name"`
			Img         string `json:"img"`
			ItemURLPath string `json:"item_url_path"`
		} `json:"results"`
	} `json:"auto"`
}

type Entity struct {
	Name        string
	AlbumTitle  string
	AlbumURL    string
	BandName    string
	CreatorName string
	Description string
	URL         string
	CoverURL    string
	Duration    string
	ReleaseDate string
	TrackURLs   []string
}

type entityType string

const (
	albumEntityType entityType = "a"
	trackEntityType entityType = "t"
)

var ldJSONRe = regexp.MustCompile(
	`(?s)<script\s+type=["']application/ld\+json["']\s*>(.*?)</script>`,
)
