package youtube

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const sampleAPIKey = "sampleApiKey"

func TestHTTPClient_GetVideo(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		responseMock  string
		expectedVideo *Video
		expectedErr   error
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
						}
					}
				]
			}`,
			expectedVideo: &Video{
				ID:           "dQw4w9WgXcQ",
				Title:        "Rick Astley - Never Gonna Give You Up (Video)",
				ChannelTitle: "RickAstleyVEVO",
			},
		},
		{
			name:    "when video not found",
			inputID: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedVideo: nil,
			expectedErr:   NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/youtube/v3/videos", r.URL.Path)
				require.Equal(t, sampleAPIKey, r.URL.Query().Get("key"))
				require.Equal(t, "snippet", r.URL.Query().Get("part"))
				require.Equal(t, tt.inputID, r.URL.Query().Get("id"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			video, err := client.GetVideo(ctx, tt.inputID)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedVideo, video)
			}
		})
	}
}

func TestHTTPClient_SearchVideo(t *testing.T) {
	tests := []struct {
		name             string
		query            string
		responseMock     string
		expectedResponse *SearchResponse
		expectedErr      error
	}{
		{
			name:  "when video found",
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
					}
				]
			}`,
			expectedResponse: &SearchResponse{
				Items: []SearchItem{
					{
						ID: SearchID{
							VideoID:    "dQw4w9WgXcQ",
							PlaylistID: "",
						},
					},
				},
			},
		},
		{
			name:  "when video not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedResponse: nil,
			expectedErr:      NotFoundError,
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
				require.Equal(t, "1", r.URL.Query().Get("maxResults"))
				require.Equal(t, "video", r.URL.Query().Get("type"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			response, err := client.SearchVideo(ctx, tt.query)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

func TestHTTPClient_GetPlaylist(t *testing.T) {
	tests := []struct {
		name             string
		inputID          string
		responseMock     string
		expectedPlaylist *Playlist
		expectedErr      error
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
							"channelTitle": "Harry"	
						}
					}
				]
			}`,
			expectedPlaylist: &Playlist{
				ID:           "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
				Title:        "Portishead - (1994) Dummy [Full Album]",
				ChannelTitle: "Harry",
			},
		},
		{
			name:    "when playlist not found",
			inputID: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedPlaylist: nil,
			expectedErr:      NotFoundError,
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

			client := NewHTTPClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			playlist, err := client.GetPlaylist(ctx, tt.inputID)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPlaylist, playlist)
			}
		})
	}
}

func TestHTTPClient_SearchPlaylist(t *testing.T) {
	tests := []struct {
		name             string
		query            string
		responseMock     string
		expectedResponse *SearchResponse
		expectedErr      error
	}{
		{
			name:  "when playlist found",
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
					}
				]
			}`,
			expectedResponse: &SearchResponse{
				Items: []SearchItem{
					{
						ID: SearchID{
							VideoID:    "",
							PlaylistID: "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
						},
					},
				},
			},
		},
		{
			name:  "when playlist not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedResponse: nil,
			expectedErr:      NotFoundError,
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
				require.Equal(t, "1", r.URL.Query().Get("maxResults"))
				require.Equal(t, "playlist", r.URL.Query().Get("type"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			response, err := client.SearchPlaylist(ctx, tt.query)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

func TestHTTPClient_GetPlaylistItems(t *testing.T) {
	tests := []struct {
		name           string
		inputID        string
		responseMock   string
		responseCode   int
		expectedVideos []Video
		expectedError  error
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
							"description": "Provided to YouTube by Revolver Records\n\nSpace Oddity · David Bowie · David Bowie · David Bowie\n\nSpace Oddity\n\n℗ 2018 Revolver Records\n\nReleased on: 2020-01-01\n\nAuto-generated by YouTube.",
							"videoOwnerChannelTitle": "David Bowie - Topic"
						}
					}
				]
			}`,
			expectedError: nil,
			expectedVideos: []Video{
				{
					ID:           "T0xBSzV1eV9uNHhhdXVzVEpTajZNdHQ0Y0l1cTRLWnppU2ZqQUJZV1UuQjcxRUYzNEU1RkQxODA0OQ",
					Title:        "Space Oddity",
					ChannelTitle: "David Bowie - Topic",
					Description: "Provided to YouTube by Revolver Records\n\nSpace Oddity · David Bowie · David Bowie " +
						"· David Bowie\n\nSpace Oddity\n\n℗ 2018 Revolver Records\n\nReleased on: 2020-01-01\n\n" +
						"Auto-generated by YouTube.",
				},
			},
		},
		{
			name:           "when playlist not found",
			inputID:        "notFoundId",
			responseCode:   http.StatusNotFound,
			responseMock:   "nevermind",
			expectedVideos: nil,
			expectedError:  fmt.Errorf("non ok http status: 404"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/youtube/v3/playlistItems", r.URL.Path)
				require.Equal(t, sampleAPIKey, r.URL.Query().Get("key"))
				require.Equal(t, "snippet", r.URL.Query().Get("part"))
				require.Equal(t, tt.inputID, r.URL.Query().Get("playlistId"))

				w.WriteHeader(tt.responseCode)
				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			videos, err := client.GetPlaylistItems(ctx, tt.inputID)

			if tt.expectedError != nil {
				require.Error(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedVideos, videos)
		})
	}
}
