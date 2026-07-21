package youtubemusic

type text struct {
	Runs       []run  `json:"runs"`
	SimpleText string `json:"simpleText"`
}

type run struct {
	Text               string             `json:"text"`
	NavigationEndpoint navigationEndpoint `json:"navigationEndpoint"`
}

type navigationEndpoint struct {
	BrowseEndpoint browseEndpoint `json:"browseEndpoint"`
	WatchEndpoint  watchEndpoint  `json:"watchEndpoint"`
}

type browseEndpoint struct {
	BrowseID string `json:"browseId"`
	Context  struct {
		Music struct {
			PageType string `json:"pageType"`
		} `json:"browseEndpointContextMusicConfig"`
	} `json:"browseEndpointContextSupportedConfigs"`
}

type watchEndpoint struct {
	VideoID string `json:"videoId"`
}

type thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type thumbnailList struct {
	Thumbnails []thumbnail `json:"thumbnails"`
}

type itemThumbnail struct {
	Renderer struct {
		Thumbnail thumbnailList `json:"thumbnail"`
	} `json:"musicThumbnailRenderer"`
}

type responsiveListItem struct {
	FlexColumns        []flexColumn       `json:"flexColumns"`
	Thumbnail          itemThumbnail      `json:"thumbnail"`
	NavigationEndpoint navigationEndpoint `json:"navigationEndpoint"`
	PlaylistItemData   struct {
		VideoID string `json:"videoId"`
	} `json:"playlistItemData"`
}

type flexColumn struct {
	Renderer struct {
		Text text `json:"text"`
	} `json:"musicResponsiveListItemFlexColumnRenderer"`
}

type responsiveListItemContainer struct {
	Renderer responsiveListItem `json:"musicResponsiveListItemRenderer"`
}

type musicShelf struct {
	Contents []responsiveListItemContainer `json:"contents"`
}

type sectionList struct {
	Contents []struct {
		MusicShelfRenderer         musicShelf       `json:"musicShelfRenderer"`
		MusicPlaylistShelfRenderer musicShelf       `json:"musicPlaylistShelfRenderer"`
		MusicResponsiveHeader      responsiveHeader `json:"musicResponsiveHeaderRenderer"`
	} `json:"contents"`
}

type tab struct {
	Renderer struct {
		Content struct {
			SectionListRenderer sectionList `json:"sectionListRenderer"`
		} `json:"content"`
	} `json:"tabRenderer"`
}

type searchResponse struct {
	Contents struct {
		Tabbed struct {
			Tabs []tab `json:"tabs"`
		} `json:"tabbedSearchResultsRenderer"`
		SectionList sectionList `json:"sectionListRenderer"`
	} `json:"contents"`
}

type browseResponse struct {
	Contents struct {
		TwoColumn struct {
			Tabs              []tab `json:"tabs"`
			SecondaryContents struct {
				SectionListRenderer sectionList `json:"sectionListRenderer"`
			} `json:"secondaryContents"`
		} `json:"twoColumnBrowseResultsRenderer"`
		SingleColumn struct {
			Tabs []tab `json:"tabs"`
		} `json:"singleColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type responsiveHeader struct {
	Title            text          `json:"title"`
	Subtitle         text          `json:"subtitle"`
	StraplineTextOne text          `json:"straplineTextOne"`
	Thumbnail        itemThumbnail `json:"thumbnail"`
	Description      struct {
		Shelf struct {
			Description text `json:"description"`
		} `json:"musicDescriptionShelfRenderer"`
	} `json:"description"`
}

type playerResponse struct {
	PlayabilityStatus struct {
		Status string `json:"status"`
	} `json:"playabilityStatus"`
	VideoDetails videoDetails `json:"videoDetails"`
	Microformat  struct {
		Renderer struct {
			URLCanonical     string        `json:"urlCanonical"`
			Description      string        `json:"description"`
			Thumbnail        thumbnailList `json:"thumbnail"`
			PageOwnerDetails struct {
				Name string `json:"name"`
			} `json:"pageOwnerDetails"`
			VideoDetails struct {
				DurationSeconds string `json:"durationSeconds"`
			} `json:"videoDetails"`
		} `json:"microformatDataRenderer"`
	} `json:"microformat"`
}

type videoDetails struct {
	VideoID       string        `json:"videoId"`
	Title         string        `json:"title"`
	LengthSeconds string        `json:"lengthSeconds"`
	Author        string        `json:"author"`
	Thumbnail     thumbnailList `json:"thumbnail"`
}

type watchNextResponse struct {
	Contents struct {
		Renderer struct {
			TabbedRenderer struct {
				Renderer struct {
					Tabs []struct {
						Renderer struct {
							Content struct {
								Queue struct {
									Content struct {
										Panel struct {
											Contents []struct {
												Video playlistPanelVideoRenderer `json:"playlistPanelVideoRenderer"`
											} `json:"contents"`
										} `json:"playlistPanelRenderer"`
									} `json:"content"`
								} `json:"musicQueueRenderer"`
							} `json:"content"`
						} `json:"tabRenderer"`
					} `json:"tabs"`
				} `json:"watchNextTabbedResultsRenderer"`
			} `json:"tabbedRenderer"`
		} `json:"singleColumnMusicWatchNextResultsRenderer"`
	} `json:"contents"`
}

type playlistPanelVideoRenderer struct {
	Title          text          `json:"title"`
	LongBylineText text          `json:"longBylineText"`
	Thumbnail      thumbnailList `json:"thumbnail"`
	Selected       bool          `json:"selected"`
	VideoID        string        `json:"videoId"`
}
