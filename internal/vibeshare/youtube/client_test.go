package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const sampleAPIKey = "sampleApiKey"

func TestHTTPClient_GetVideo(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		responseMock  string
		expectedVideo *Video
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

			video, err := client.GetVideo(tt.inputID)
			require.NoError(t, err)
			require.Equal(t, tt.expectedVideo, video)
		})
	}
}

func TestHTTPClient_SearchVideo(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		responseMock  string
		expectedVideo *Video
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
			expectedVideo: &Video{
				ID:           "dQw4w9WgXcQ",
				Title:        "Rick Astley - Never Gonna Give You Up (Video)",
				ChannelTitle: "RickAstleyVEVO",
			},
		},
		{
			name:  "when video not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedVideo: nil,
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
				require.Equal(t, "10", r.URL.Query().Get("videoCategoryId"))
				require.Equal(t, "1", r.URL.Query().Get("maxResults"))
				require.Equal(t, "video", r.URL.Query().Get("type"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			video, err := client.SearchVideo(tt.query)
			require.NoError(t, err)
			require.Equal(t, tt.expectedVideo, video)
		})
	}
}

func TestHTTPClient_GetPlaylist(t *testing.T) {
	tests := []struct {
		name             string
		inputID          string
		responseMock     string
		expectedPlaylist *Playlist
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

			playlist, err := client.GetPlaylist(tt.inputID)
			require.NoError(t, err)
			require.Equal(t, tt.expectedPlaylist, playlist)
		})
	}
}

func TestHTTPClient_SearchPlaylist(t *testing.T) {
	tests := []struct {
		name             string
		query            string
		responseMock     string
		expectedPlaylist *Playlist
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
			expectedPlaylist: &Playlist{
				ID:           "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
				Title:        "Portishead - (1994) Dummy [Full Album]",
				ChannelTitle: "Harry",
			},
		},
		{
			name:  "when playlist not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedPlaylist: nil,
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

			video, err := client.SearchPlaylist(tt.query)
			require.NoError(t, err)
			require.Equal(t, tt.expectedPlaylist, video)
		})
	}
}
