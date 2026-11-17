package official

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/youtube/client"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

const sampleAPIKey = "sampleApiKey"

func TestClientAPIError(t *testing.T) {
	for _, tt := range []struct {
		name           string
		status         int
		body           string
		expectedReason string
		exceeded       bool
	}{
		{"quota exhausted", 403, `{"error":{"errors":[{"reason":"quotaExceeded"}]}}`, "quotaExceeded", true},
		{"later quota reason", 403, `{"error":{"errors":[{"reason":"forbidden"},{"reason":"quotaExceeded"}]}}`, "forbidden, quotaExceeded", true},
		{"access forbidden", 403, `{"error":{"errors":[{"reason":"forbidden"}]}}`, "forbidden", false},
		{"multiple reasons", 403, `{"error":{"errors":[{"reason":"forbidden"},{"reason":"insufficientPermissions"}]}}`, "forbidden, insufficientPermissions", false},
		{"bad request", 400, `{"error":{"errors":[{"reason":"invalidFilters"}]}}`, "invalidFilters", false},
		{"not found", 404, `{"error":{"errors":[{"reason":"videoNotFound"}]}}`, "videoNotFound", false},
		{"empty reason", 403, `{"error":{"errors":[{"reason":""},{"reason":"forbidden"}]}}`, "forbidden", false},
		{"quota mentioned in message", 403, `{"error":{"message":"quotaExceeded"}}`, "", false},
		{"invalid JSON", 403, `{"error":`, "", false},
		{"trailing data", 403, `{"error":{"errors":[{"reason":"quotaExceeded"}]}}invalid`, "", false},
		{"empty response", 403, "", "", false},
		{"server error", 500, `{"error":{"errors":[{"reason":"quotaExceeded"}]}}`, "quotaExceeded", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, err := w.Write([]byte(tt.body))
				require.NoError(t, err)
			}))
			defer server.Close()
			c := NewClient(sampleAPIKey, WithAPIURL(server.URL))

			body, err := c.getWithKey(t.Context(), "/youtube/v3/videos", url.Values{"id": {"video-id"}})

			require.Nil(t, body)
			expected := fmt.Sprintf("non ok http status: %d", tt.status)
			if tt.expectedReason != "" {
				expected += " (reason: " + tt.expectedReason + ")"
			}
			if tt.exceeded {
				require.ErrorIs(t, err, release.ErrQuotaExceeded)
				expected += ": " + release.ErrQuotaExceeded.Error()
			} else {
				require.NotErrorIs(t, err, release.ErrQuotaExceeded)
			}
			require.EqualError(t, err, expected)
		})
	}
}

func TestClient_GetVideo(t *testing.T) {
	tests := []struct {
		name         string
		inputID      string
		responseMock string
		expectedItem client.Video
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
			expectedItem: client.Video{
				Metadata: client.Metadata{
					ID:      "dQw4w9WgXcQ",
					Title:   "Rick Astley - Never Gonna Give You Up (Video)",
					Creator: "RickAstleyVEVO",
				},
				Duration: 213,
			},
		},
		{
			name:    "when video not found",
			inputID: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedErr: client.ErrNotFound,
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

			c := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			item, err := c.FetchVideo(t.Context(), tt.inputID)
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
		expectedItems []string
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
			expectedItems: []string{"dQw4w9WgXcQ", "secondVideoID"},
		},
		{
			name:  "when videos not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedItems: []string{},
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

			c := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			items, err := c.SearchVideos(t.Context(), tt.query)
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
		expectedItem client.Metadata
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
			expectedItem: client.Metadata{
				ID:          "PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6",
				Title:       "Portishead - (1994) Dummy [Full Album]",
				Creator:     "Harry",
				Description: "playlist raw description",
			},
		},
		{
			name:    "when playlist not found",
			inputID: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedErr: client.ErrNotFound,
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

			c := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			item, err := c.FetchPlaylist(t.Context(), tt.inputID)
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
		expectedItems []string
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
			expectedItems: []string{"PLH1JGOJgZ2u2J7bRnfjl-7kDj_vQKTPa6", "secondPlaylistID"},
		},
		{
			name:  "when playlists not found",
			query: "notFoundId",
			responseMock: `{	
				"items": []
			}`,
			expectedItems: []string{},
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

			c := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			items, err := c.SearchPlaylists(t.Context(), tt.query)
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
		expectedItems []string
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
			expectedItems: []string{"dQw4w9WgXcQ"},
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

			c := NewClient(sampleAPIKey, WithAPIURL(apiServerMock.URL))

			items, err := c.fetchPlaylistItems(t.Context(), tt.inputID)

			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedItems, items)
		})
	}
}
