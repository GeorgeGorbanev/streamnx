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
	BrowseEndpoint        browseEndpoint        `json:"browseEndpoint"`
	WatchEndpoint         watchEndpoint         `json:"watchEndpoint"`
	WatchPlaylistEndpoint watchPlaylistEndpoint `json:"watchPlaylistEndpoint"`
}

type browseEndpoint struct {
	BrowseID string                `json:"browseId"`
	Context  browseEndpointContext `json:"browseEndpointContextSupportedConfigs"`
}

type browseEndpointContext struct {
	Music browseEndpointMusic `json:"browseEndpointContextMusicConfig"`
}

type browseEndpointMusic struct {
	PageType string `json:"pageType"`
}

type watchEndpoint struct {
	VideoID string `json:"videoId"`
}

type watchPlaylistEndpoint struct {
	PlaylistID string `json:"playlistId"`
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
	Renderer thumbnailRenderer `json:"musicThumbnailRenderer"`
}

type thumbnailRenderer struct {
	Thumbnail thumbnailList `json:"thumbnail"`
}

type responsiveListItem struct {
	FlexColumns        []flexColumn       `json:"flexColumns"`
	Thumbnail          itemThumbnail      `json:"thumbnail"`
	NavigationEndpoint navigationEndpoint `json:"navigationEndpoint"`
	PlaylistItemData   playlistItemData   `json:"playlistItemData"`
	Overlay            itemOverlay        `json:"overlay"`
}

type itemOverlay struct {
	Renderer itemThumbnailOverlayRenderer `json:"musicItemThumbnailOverlayRenderer"`
}

type itemThumbnailOverlayRenderer struct {
	Content itemOverlayContent `json:"content"`
}

type itemOverlayContent struct {
	PlayButton playButtonRenderer `json:"musicPlayButtonRenderer"`
}

type playButtonRenderer struct {
	PlayNavigationEndpoint navigationEndpoint `json:"playNavigationEndpoint"`
}

type playlistItemData struct {
	VideoID string `json:"videoId"`
}

type flexColumn struct {
	Renderer flexColumnRenderer `json:"musicResponsiveListItemFlexColumnRenderer"`
}

type flexColumnRenderer struct {
	Text text `json:"text"`
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
	Title            text                     `json:"title"`
	Subtitle         text                     `json:"subtitle"`
	StraplineTextOne text                     `json:"straplineTextOne"`
	Thumbnail        itemThumbnail            `json:"thumbnail"`
	Buttons          []responsiveHeaderButton `json:"buttons"`
	Description      struct {
		Shelf struct {
			Description text `json:"description"`
		} `json:"musicDescriptionShelfRenderer"`
	} `json:"description"`
}

type responsiveHeaderButton struct {
	PlayButton playButtonRenderer `json:"musicPlayButtonRenderer"`
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
