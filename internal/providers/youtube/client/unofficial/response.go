package unofficial

import "encoding/json"

type webText struct {
	SimpleText string `json:"simpleText"`
	Runs       []struct {
		Text string `json:"text"`
	} `json:"runs"`
}

type webThumbnails struct {
	Items []thumbnail `json:"thumbnails"`
}

type thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type webPlayerResponse struct {
	Status struct {
		Code   string `json:"status"`
		Reason string `json:"reason"`
	} `json:"playabilityStatus"`
	Details struct {
		ID          string        `json:"videoId"`
		Title       string        `json:"title"`
		Author      string        `json:"author"`
		Description string        `json:"shortDescription"`
		Seconds     string        `json:"lengthSeconds"`
		Thumbnail   webThumbnails `json:"thumbnail"`
	} `json:"videoDetails"`
}

type webItem struct {
	Video *struct {
		ID string `json:"videoId"`
	} `json:"videoRenderer"`
	Playlist *struct {
		ID string `json:"playlistId"`
	} `json:"playlistRenderer"`
	PlaylistVideo *struct {
		ID string `json:"videoId"`
	} `json:"playlistVideoRenderer"`
	Lockup *struct {
		ID   string `json:"contentId"`
		Type string `json:"contentType"`
	} `json:"lockupViewModel"`
	Section      *webSection `json:"itemSectionRenderer"`
	PlaylistList *webSection `json:"playlistVideoListRenderer"`
	Continuation struct {
		Endpoint webContinuationEndpoint `json:"continuationEndpoint"`
	} `json:"continuationItemRenderer"`
	ContinuationView struct {
		Command struct {
			Endpoint webContinuationEndpoint `json:"innertubeCommand"`
		} `json:"continuationCommand"`
	} `json:"continuationItemViewModel"`
}

type webContinuationEndpoint struct {
	Command struct {
		Token string `json:"token"`
	} `json:"continuationCommand"`
}

type webSection struct {
	TargetID string    `json:"targetId"`
	Items    []webItem `json:"contents"`
}

type webSearchResponse struct {
	Contents struct {
		Results struct {
			Primary struct {
				Sections webSection `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

type webPlaylistResponse struct {
	Metadata struct {
		Playlist struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"playlistMetadataRenderer"`
	} `json:"metadata"`
	Header struct {
		Playlist struct {
			Title       webText `json:"title"`
			Owner       webText `json:"ownerText"`
			Description webText `json:"descriptionText"`
			Banner      struct {
				Hero struct {
					Thumbnail webThumbnails `json:"thumbnail"`
				} `json:"heroPlaylistThumbnailRenderer"`
			} `json:"playlistHeaderBanner"`
		} `json:"playlistHeaderRenderer"`
	} `json:"header"`
	Microformat struct {
		Data struct {
			Description string        `json:"description"`
			Thumbnail   webThumbnails `json:"thumbnail"`
		} `json:"microformatDataRenderer"`
	} `json:"microformat"`
	Sidebar struct {
		Playlist struct {
			Items []struct {
				Secondary struct {
					Owner struct {
						Renderer struct {
							Title webText `json:"title"`
						} `json:"videoOwnerRenderer"`
					} `json:"videoOwner"`
				} `json:"playlistSidebarSecondaryInfoRenderer"`
			} `json:"items"`
		} `json:"playlistSidebarRenderer"`
	} `json:"sidebar"`
	Contents struct {
		Browse struct {
			Tabs []struct {
				Renderer struct {
					Selected bool `json:"selected"`
					Content  struct {
						Sections webSection `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
	Alerts []struct {
		Renderer struct {
			Type string  `json:"type"`
			Text webText `json:"text"`
		} `json:"alertRenderer"`
	} `json:"alerts"`
}

type webConfig struct {
	Context json.RawMessage `json:"INNERTUBE_CONTEXT"`
}

type webBrowseRequest struct {
	Context      json.RawMessage `json:"context"`
	Continuation string          `json:"continuation"`
}

type webBrowseResponse struct {
	Actions   []webBrowseAction `json:"onResponseReceivedActions"`
	Endpoints []webBrowseAction `json:"onResponseReceivedEndpoints"`
}

type webBrowseAction struct {
	Append struct {
		TargetID string    `json:"targetId"`
		Items    []webItem `json:"continuationItems"`
	} `json:"appendContinuationItemsAction"`
}
