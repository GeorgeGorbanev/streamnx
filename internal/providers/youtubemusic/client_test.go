package youtubemusic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientRequestsAndDecodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, "https://music.youtube.com", r.Header.Get("Origin"))
		require.Equal(t, "alt=json", r.URL.RawQuery)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		contextValue := body["context"].(map[string]any)
		clientValue := contextValue["client"].(map[string]any)
		require.Equal(t, "WEB_REMIX", clientValue["clientName"])
		require.True(t, strings.HasPrefix(clientValue["clientVersion"].(string), "1."))

		switch r.URL.Path {
		case "/player":
			require.Equal(t, "video-id", body["video_id"])
			_, err := w.Write([]byte(`{
				"playabilityStatus":{"status":"OK"},
				"videoDetails":{"videoId":"video-id","title":"Track","author":"Artist","lengthSeconds":"120"},
				"microformat":{"microformatDataRenderer":{
					"urlCanonical":"https://music.youtube.com/watch?v=video-id",
					"description":"Description",
					"thumbnail":{"thumbnails":[{"url":"player-cover","width":120,"height":120}]},
					"pageOwnerDetails":{"name":"Artist - Topic"},
					"videoDetails":{"durationSeconds":"121"}
				}}
			}`))
			require.NoError(t, err)
		case "/next":
			require.Equal(t, "video-id", body["videoId"])
			_, err := w.Write([]byte(`{
				"contents": {
					"singleColumnMusicWatchNextResultsRenderer": {
						"tabbedRenderer": {
							"watchNextTabbedResultsRenderer": {
								"tabs": [{
									"tabRenderer": {
										"content": {
											"musicQueueRenderer": {
												"content": {
													"playlistPanelRenderer": {
														"contents": [{
															"playlistPanelVideoRenderer": {
																"videoId": "video-id",
																"selected": true
															}
														}]
													}
												}
											}
										}
									}
								}]
							}
						}
					}
				}
			}`))
			require.NoError(t, err)
		case "/search":
			require.Equal(t, "Artist – Title", body["query"])
			require.Contains(t, []any{songsParams, albumsParams}, body["params"])
			_, err := w.Write([]byte(`{
				"contents":{"tabbedSearchResultsRenderer":{"tabs":[{"tabRenderer":{"content":{"sectionListRenderer":{"contents":[
					{"musicShelfRenderer":{"contents":[{"musicResponsiveListItemRenderer":{
						"playlistItemData":{"videoId":"video-id"},
						"thumbnail":{"musicThumbnailRenderer":{"thumbnail":{"thumbnails":[{"url":"search-cover"}]}}},
						"flexColumns":[{"musicResponsiveListItemFlexColumnRenderer":{"text":{"runs":[{
							"text":"Artist",
							"navigationEndpoint":{"browseEndpoint":{
								"browseId":"artist-id",
								"browseEndpointContextSupportedConfigs":{"browseEndpointContextMusicConfig":{"pageType":"MUSIC_PAGE_TYPE_ARTIST"}}
							}}
						}]}}}]
					}}]}}
				]}}}}]}}
			}`))
			require.NoError(t, err)
		case "/browse":
			require.Equal(t, "VLOLAKalbum", body["browseId"])
			_, err := w.Write([]byte(`{
				"contents":{"twoColumnBrowseResultsRenderer":{
					"tabs":[{"tabRenderer":{"content":{"sectionListRenderer":{"contents":[
						{"musicResponsiveHeaderRenderer":{
							"title":{"runs":[{"text":"Album"}]},
							"thumbnail":{"musicThumbnailRenderer":{"thumbnail":{"thumbnails":[{"url":"album-cover"}]}}},
							"description":{"musicDescriptionShelfRenderer":{"description":{"runs":[{"text":"Album description"}]}}}
						}}
					]}}}}],
					"secondaryContents":{"sectionListRenderer":{"contents":[
						{"musicShelfRenderer":{"contents":[{"musicResponsiveListItemRenderer":{"playlistItemData":{"videoId":"video-id"}}}]}}
					]}}
				}}
			}`))
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(WithAPIURL(server.URL))

	track, err := client.fetchTrack(t.Context(), "video-id")
	require.NoError(t, err)
	require.Equal(t, "video-id", track.VideoDetails.VideoID)
	require.Equal(t, "Artist - Topic", track.Microformat.Renderer.PageOwnerDetails.Name)
	require.Equal(t, "player-cover", track.Microformat.Renderer.Thumbnail.Thumbnails[0].URL)

	metadata, err := client.fetchWatchNext(t.Context(), "video-id")
	require.NoError(t, err)
	require.Equal(t, "video-id", metadata.VideoID)

	tracks, err := client.searchTracks(t.Context(), "Artist – Title")
	require.NoError(t, err)
	require.Len(t, tracks, 1)
	require.Equal(t, "search-cover", tracks[0].Thumbnail.Renderer.Thumbnail.Thumbnails[0].URL)
	require.Equal(t, "MUSIC_PAGE_TYPE_ARTIST",
		tracks[0].FlexColumns[0].Renderer.Text.Runs[0].NavigationEndpoint.BrowseEndpoint.Context.Music.PageType)

	albums, err := client.searchAlbums(t.Context(), "Artist – Title")
	require.NoError(t, err)
	require.Len(t, albums, 1)

	header, items, err := client.fetchAlbum(t.Context(), "OLAKalbum")
	require.NoError(t, err)
	require.Equal(t, "Album", header.Title.Runs[0].Text)
	require.Equal(t, "album-cover", header.Thumbnail.Renderer.Thumbnail.Thumbnails[0].URL)
	require.Equal(t, "Album description", header.Description.Shelf.Description.Runs[0].Text)
	require.Len(t, items, 1)
}

func TestClientNotFound(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{name: "http 404", statusCode: http.StatusNotFound, body: `{}`},
		{name: "missing video details", statusCode: http.StatusOK, body: `{"playabilityStatus":{"status":"ERROR"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.body))
				require.NoError(t, err)
			}))
			defer server.Close()

			got, err := NewClient(WithAPIURL(server.URL)).fetchTrack(t.Context(), "missing")

			require.Zero(t, got)
			require.ErrorIs(t, err, errNotFound)
		})
	}
}
