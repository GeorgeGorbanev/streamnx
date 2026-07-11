package youtube

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const sampleAPIKey = "sampleApiKey"

func TestClient_GetVideo(t *testing.T) {
	tests := []struct {
		name         string
		inputID      string
		responseMock string
		expectedItem video
		expectedErr  error
	}{
		{
			name:    "when video found",
			inputID: "dQw4w9WgXcQ",
			responseMock: `{	
				"items": [
					{	
						"id": "dQw4w9WgXcQ",
						"snippet": {	
							"title": "Rick Astley - Never Gonna Give You Up (Video)",	
							"channelTitle": "RickAstleyVEVO"	
						},
						"contentDetails": {
							"duration": "PT3M33S"
						}
					}
				]
			}`,
			expectedItem: video{
				ID: "dQw4w9WgXcQ",
				Snippet: snippet{
					Title:        "Rick Astley - Never Gonna Give You Up (Video)",
					ChannelTitle: "RickAstleyVEVO",
				},
				ContentDetails: contentDetails{
					Duration: "PT3M33S",
				},
			},
		},
		{
			name:    "when video not found",
			inputID: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedErr: errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/youtube/v3/videos", r.URL.Path)
				require.Equal(t, sampleAPIKey, r.URL.Query().Get("key"))
				require.Equal(t, "snippet,contentDetails", r.URL.Query().Get("part"))
				require.Equal(t, tt.inputID, r.URL.Query().Get("id"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			item, err := client.fetchVideo(t.Context(), tt.inputID)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedItem, item)
			}
		})
	}
}

func TestClient_SearchVideo(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		responseMock  string
		expectedItems []videoSearchResult
		expectedErr   error
	}{
		{
			name:  "when videos found",
			query: "rick astley - never gonna give you up",
			responseMock: `{	
				"items": [
					{	
						"id": {
							"videoId": "dQw4w9WgXcQ"
						},
						"snippet": {	
							"title": "Rick Astley - Never Gonna Give You Up (Video)",	
							"channelTitle": "RickAstleyVEVO"	
						}
					},
					{
						"id": {
							"videoId": "secondVideoID"
						},
						"snippet": {
							"title": "Rick Astley - Never Gonna Give You Up (Live)",
							"channelTitle": "RickAstleyVEVO"
						}
					}
				]
			}`,
			expectedItems: []videoSearchResult{
				{
					ID: videoSearchID{
						VideoID: "dQw4w9WgXcQ",
					},
				},
				{
					ID: videoSearchID{
						VideoID: "secondVideoID",
					},
				},
			},
		},
		{
			name:  "when videos not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedItems: []videoSearchResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/youtube/v3/search", r.URL.Path)
				require.Equal(t, sampleAPIKey, r.URL.Query().Get("key"))
				require.Equal(t, tt.query, r.URL.Query().Get("q"))
				require.Equal(t, "10", r.URL.Query().Get("videoCategoryId"))
				require.Equal(t, "10", r.URL.Query().Get("maxResults"))
				require.Equal(t, "video", r.URL.Query().Get("type"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			items, err := client.searchVideos(t.Context(), tt.query)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedItems, items)
			}
		})
	}
}

func TestClient_GetPlaylist(t *testing.T) {
	tests := []struct {
		name         string
		inputID      string
		responseMock string
		expectedItem playlist
		expectedErr  error
	}{
		{
			name:    "when playlist found",
			inputID: "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
			responseMock: `{	
				"items": [
					{	
						"id": "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
						"snippet": {	
							"title": "Portishead - (1994) Dummy [Full Album]",	
							"channelTitle": "Harry",
							"description": "playlist raw description"	
						}
					}
				]
			}`,
			expectedItem: playlist{
				ID: "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
				Snippet: snippet{
					Title:        "Portishead - (1994) Dummy [Full Album]",
					ChannelTitle: "Harry",
					Description:  "playlist raw description",
				},
			},
		},
		{
			name:    "when playlist not found",
			inputID: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedErr: errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/youtube/v3/playlists", r.URL.Path)
				require.Equal(t, sampleAPIKey, r.URL.Query().Get("key"))
				require.Equal(t, "snippet", r.URL.Query().Get("part"))
				require.Equal(t, tt.inputID, r.URL.Query().Get("id"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			item, err := client.fetchPlaylist(t.Context(), tt.inputID)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedItem, item)
			}
		})
	}
}

func TestClient_SearchPlaylist(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		responseMock  string
		expectedItems []playlistSearchResult
		expectedErr   error
	}{
		{
			name:  "when playlists found",
			query: "portishead – dummy",
			responseMock: `{	
				"items": [
					{	
						"id": {
							"playlistId": "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6"
						},
						"snippet": {	
							"title": "Portishead - (1994) Dummy [Full Album]",	
							"channelTitle": "Harry"	
						}
					},
					{
						"id": {
							"playlistId": "secondPlaylistID"
						},
						"snippet": {
							"title": "Portishead - Dummy alternate playlist",
							"channelTitle": "Harry"
						}
					}
				]
			}`,
			expectedItems: []playlistSearchResult{
				{
					ID: playlistSearchID{
						PlaylistID: "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
					},
				},
				{
					ID: playlistSearchID{
						PlaylistID: "secondPlaylistID",
					},
				},
			},
		},
		{
			name:  "when playlists not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedItems: []playlistSearchResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/youtube/v3/search", r.URL.Path)
				require.Equal(t, sampleAPIKey, r.URL.Query().Get("key"))
				require.Equal(t, "snippet", r.URL.Query().Get("part"))
				require.Equal(t, tt.query, r.URL.Query().Get("q"))
				require.Equal(t, "10", r.URL.Query().Get("maxResults"))
				require.Equal(t, "playlist", r.URL.Query().Get("type"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			items, err := client.searchPlaylists(t.Context(), tt.query)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedItems, items)
			}
		})
	}
}

func TestClient_GetPlaylistItems(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		responseMock  string
		responseCode  int
		expectedItems []playlistItem
		expectedError error
	}{
		{
			name:         "when playlist found",
			inputID:      "OLAK5uy_n4xauusTJSj6Mtt4cIuq4KZziSfjABYWU",
			responseCode: http.StatusOK,
			responseMock: `{	
				"items": [
					{	
						"id": "T0xBSzV1eV9uNHhhdXVzVEpTajZNdHQ0Y0l1cTRLWnppU2ZqQUJZV1UuQjcxRUYzNEU1RkQxODA0OQ",
						"snippet": {	
							"title": "Space Oddity",	
							"channelTitle": "YouTube",	
							"description": "Provided to YouTube by Revolver Records\n\nSpace Oddity · David Bowie · David Bowie · David Bowie\n\nSpace Oddity\n\n℗ 2018 Revolver Records\n\nReleased on: 2020-01-01\n\nGenerated by the provider.",
							"videoOwnerChannelTitle": "David Bowie Topic Channel",
							"resourceId": {
								"videoId": "dQw4w9WgXcQ"
							}
						}
					}
				]
			}`,
			expectedError: nil,
			expectedItems: []playlistItem{
				{
					ID: "T0xBSzV1eV9uNHhhdXVzVEpTajZNdHQ0Y0l1cTRLWnppU2ZqQUJZV1UuQjcxRUYzNEU1RkQxODA0OQ",
					Snippet: snippet{
						Title:        "Space Oddity",
						ChannelTitle: "YouTube",
						Description: "Provided to YouTube by Revolver Records\n\nSpace Oddity · David Bowie · David Bowie " +
							"· David Bowie\n\nSpace Oddity\n\n℗ 2018 Revolver Records\n\nReleased on: 2020-01-01\n\n" +
							"Generated by the provider.",
						VideoOwnerChannelTitle: "David Bowie Topic Channel",
						ResourceID: resourceID{
							VideoID: "dQw4w9WgXcQ",
						},
					},
				},
			},
		},
		{
			name:          "when playlist not found",
			inputID:       "notFoundId",
			responseCode:  http.StatusNotFound,
			responseMock:  "nevermind",
			expectedItems: nil,
			expectedError: fmt.Errorf("non ok http status: 404"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/youtube/v3/playlistItems", r.URL.Path)
				require.Equal(t, sampleAPIKey, r.URL.Query().Get("key"))
				require.Equal(t, "snippet", r.URL.Query().Get("part"))
				require.Equal(t, "50", r.URL.Query().Get("maxResults"))
				require.Equal(t, tt.inputID, r.URL.Query().Get("playlistId"))

				w.WriteHeader(tt.responseCode)
				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			items, err := client.fetchPlaylistItems(t.Context(), tt.inputID)

			if tt.expectedError != nil {
				require.Error(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedItems, items)
		})
	}
}
