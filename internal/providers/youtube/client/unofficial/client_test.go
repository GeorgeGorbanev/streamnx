package unofficial

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/youtube/client"
	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestClient_SearchVideos(t *testing.T) {
	cards := []string{
		`{"playlistRenderer":{"playlistId":"ignored"}}`,
		`{"videoRenderer":{"videoId":"legacy"}}`,
		`{"lockupViewModel":{"contentType":"LOCKUP_CONTENT_TYPE_VIDEO","contentId":"lockup"}}`,
	}
	for i := 3; i <= 11; i++ {
		cards = append(cards, fmt.Sprintf(`{"videoRenderer":{"videoId":"video%d"}}`, i))
	}
	c := newSearchClient(t, "EgIQAfABAQ==", searchPage(strings.Join(cards, ",")))

	ids, err := c.SearchVideos(t.Context(), "artist & title")

	require.NoError(t, err)
	require.Equal(t, []string{
		"legacy", "lockup", "video3", "video4", "video5",
		"video6", "video7", "video8", "video9", "video10",
	}, ids)
}

func TestClient_SearchPlaylists(t *testing.T) {
	tests := []struct {
		name string
		body string
		ids  []string
		err  string
	}{
		{
			name: "legacy and lockup playlist cards",
			body: searchPage(`
				{"playlistRenderer":{"playlistId":"legacy"}},
				{"lockupViewModel":{"contentType":"LOCKUP_CONTENT_TYPE_VIDEO","contentId":"ignored"}},
				{"lockupViewModel":{"contentType":"LOCKUP_CONTENT_TYPE_PLAYLIST","contentId":"lockup"}}
			`),
			ids: []string{"legacy", "lockup"},
		},
		{
			name: "empty results",
			body: searchPage(""),
			ids:  []string{},
		},
		{
			name: "missing results section",
			body: `<script>var ytInitialData = {"contents":{}};</script>`,
			err:  "missing search results section",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newSearchClient(t, "EgIQAw==", tt.body)

			ids, err := c.SearchPlaylists(t.Context(), "artist & title")

			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				require.Nil(t, ids)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.ids, ids)
		})
	}
}

func TestClient_FetchVideoHTTPStatus(t *testing.T) {
	tests := []struct {
		status int
		err    error
	}{
		{http.StatusNotFound, client.ErrNotFound},
		{http.StatusForbidden, release.ErrScrapingBlocked},
		{http.StatusTooManyRequests, release.ErrScrapingBlocked},
	}
	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/watch", r.URL.Path)
				require.Equal(t, "video-id", r.URL.Query().Get("v"))
				w.WriteHeader(tt.status)
			}))
			t.Cleanup(server.Close)
			c := NewClient(WithAPIURL(server.URL), WithHTTPClient(server.Client()))

			video, err := c.FetchVideo(t.Context(), "video-id")

			require.ErrorIs(t, err, tt.err)
			if tt.err == release.ErrScrapingBlocked {
				require.ErrorContains(t, err, fmt.Sprintf("YouTube HTTP %d", tt.status))
			}
			require.Zero(t, video)
		})
	}
}

func TestClient_FetchVideoInvalidPage(t *testing.T) {
	tests := []struct {
		name string
		body string
		err  error
	}{
		{
			name: "bot challenge",
			body: `<script>var ytInitialPlayerResponse = {
				"playabilityStatus":{"status":"LOGIN_REQUIRED","reason":"Sign in to confirm you’re not a bot"}
			};</script>`,
			err: release.ErrScrapingBlocked,
		},
		{
			name: "missing player JSON",
			body: `<!doctype html><title>Unexpected page</title>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, err := io.WriteString(w, tt.body)
				require.NoError(t, err)
			}))
			t.Cleanup(server.Close)
			c := NewClient(WithAPIURL(server.URL), WithHTTPClient(server.Client()))

			video, err := c.FetchVideo(t.Context(), "video-id")
			require.Zero(t, video)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.ErrorContains(t, err, "YouTube requires bot verification")
			} else {
				require.ErrorContains(t, err, "missing embedded YouTube JSON")
				require.NotErrorIs(t, err, client.ErrNotFound)
			}
		})
	}
}

func TestClient_FetchPlaylistInitialDataFormats(t *testing.T) {
	const data = `{"metadata":{"playlistMetadataRenderer":{"title":"Album"}}}`
	for _, tt := range []struct {
		name string
		html string
	}{
		{name: "plain assignment", html: `<script>ytInitialData = ` + data + `;</script>`},
		{name: "window assignment", html: `<script>window [ 'ytInitialData' ] = ` + data + `;</script>`},
		{name: "dedicated script", html: `<script id="yt-initial-data" type="application/json">` + data + `</script>`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, err := io.WriteString(w, tt.html)
				require.NoError(t, err)
			}))
			t.Cleanup(server.Close)
			c := NewClient(WithAPIURL(server.URL), WithHTTPClient(server.Client()))

			metadata, err := c.FetchPlaylist(t.Context(), "playlist-id")

			require.NoError(t, err)
			require.Equal(t, client.Metadata{ID: "playlist-id", Title: "Album"}, metadata)
		})
	}
}

func TestClient_FetchPlaylistRetriesTruncatedResponse(t *testing.T) {
	const body = `<script>var ytInitialData = {"metadata":{"playlistMetadataRenderer":{"title":"Album"}}};</script>`
	for _, alwaysTruncated := range []bool{false, true} {
		t.Run(fmt.Sprintf("always truncated=%t", alwaysTruncated), func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if attempts.Add(1) == 1 || alwaysTruncated {
					w.Header().Set("Content-Length", strconv.Itoa(len(body)+1))
				}
				_, err := io.WriteString(w, body)
				require.NoError(t, err)
			}))
			t.Cleanup(server.Close)
			c := NewClient(WithAPIURL(server.URL), WithHTTPClient(server.Client()))

			metadata, err := c.FetchPlaylist(t.Context(), "playlist-id")

			require.Equal(t, int32(2), attempts.Load())
			if alwaysTruncated {
				require.ErrorIs(t, err, io.ErrUnexpectedEOF)
				require.ErrorContains(t, err, `load playlist page "playlist-id"`)
				require.Zero(t, metadata)
				return
			}
			require.NoError(t, err)
			require.Equal(t, client.Metadata{ID: "playlist-id", Title: "Album"}, metadata)
		})
	}
}

func TestClient_FetchPlaylistReportsUnexpectedPage(t *testing.T) {
	const body = "<html>unexpected response</html>"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/playlist" {
			http.Redirect(w, r, "/interstitial?token=private-token", http.StatusFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, err := io.WriteString(w, body)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	c := NewClient(WithAPIURL(server.URL), WithHTTPClient(server.Client()))

	metadata, err := c.FetchPlaylist(t.Context(), "playlist-id")

	require.Zero(t, metadata)
	require.ErrorContains(t, err, `decode playlist page "playlist-id" (status=200, final_host=`)
	require.ErrorContains(t, err, fmt.Sprintf(`final_path="/interstitial", content_type="text/html", bytes=%d, has_ytInitialData=false`, len(body)))
	require.ErrorContains(t, err, "missing embedded YouTube JSON")
	require.NotContains(t, err.Error(), "private-token")
	require.NotContains(t, err.Error(), body)
}

func TestClient_FetchPlaylistWithItemsStopsOnRepeatedContinuation(t *testing.T) {
	const playlistID = "playlist-id"
	const token = "next-page"
	const continuation = `{"continuationItemRenderer":{"continuationEndpoint":{"continuationCommand":{"token":"` + token + `"}}}}`
	const page = `
		<script>var ytInitialData = {
			"header":{"playlistHeaderRenderer":{"title":{"simpleText":"Legacy playlist"},"ownerText":{"simpleText":"Owner"}}},
			"contents":{"twoColumnBrowseResultsRenderer":{"tabs":[{"tabRenderer":{"selected":true,"content":{
				"sectionListRenderer":{"contents":[{"itemSectionRenderer":{"contents":[{"playlistVideoListRenderer":{
					"contents":[{"playlistVideoRenderer":{"videoId":"first"}},` + continuation + `]
				}}]}}]}
			}}}]}}
		};</script>
		<script>ytcfg.set({"INNERTUBE_CONTEXT":{"client":{"clientName":"WEB"}}});</script>`
	const nextPage = `{"onResponseReceivedEndpoints":[{"appendContinuationItemsAction":{
		"targetId":"` + playlistID + `","continuationItems":[{"playlistVideoRenderer":{"videoId":"second"}},` + continuation + `]
	}}]}`

	var browseCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/playlist":
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, playlistID, r.URL.Query().Get("list"))
			_, err := io.WriteString(w, page)
			require.NoError(t, err)
		case "/youtubei/v1/browse":
			require.Equal(t, int32(1), browseCalls.Add(1))
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var request webBrowseRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
			require.Equal(t, token, request.Continuation)
			require.JSONEq(t, `{"client":{"clientName":"WEB"}}`, string(request.Context))
			_, err := io.WriteString(w, nextPage)
			require.NoError(t, err)
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	httpClient := server.Client()
	httpClient.Timeout = 5 * time.Second
	c := NewClient(WithAPIURL(server.URL), WithHTTPClient(httpClient))

	playlist, err := c.FetchPlaylist(t.Context(), playlistID)
	require.NoError(t, err)
	require.Equal(t, client.Metadata{ID: playlistID, Title: "Legacy playlist", Creator: "Owner"}, playlist)
	require.Zero(t, browseCalls.Load())

	playlist, ids, err := c.FetchPlaylistWithItems(t.Context(), playlistID)
	require.ErrorContains(t, err, "repeated playlist continuation")
	require.Zero(t, playlist)
	require.Nil(t, ids)
	require.Equal(t, int32(1), browseCalls.Load())
}

func newSearchClient(t *testing.T, filter, body string) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/results", r.URL.Path)
		require.Equal(t, "artist & title", r.URL.Query().Get("search_query"))
		require.Equal(t, filter, r.URL.Query().Get("sp"))
		require.Equal(t, "en", r.URL.Query().Get("hl"))
		_, err := io.WriteString(w, body)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	return NewClient(WithAPIURL(server.URL), WithHTTPClient(server.Client()))
}

func searchPage(cards string) string {
	return `<script>var ytInitialData = {"contents":{"twoColumnSearchResultsRenderer":{"primaryContents":{
		"sectionListRenderer":{"contents":[{"itemSectionRenderer":{"contents":[` + cards + `]}}]}
	}}}};</script>`
}
