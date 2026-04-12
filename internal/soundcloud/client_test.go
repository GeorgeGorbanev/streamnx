package soundcloud

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient_FetchTrack(t *testing.T) {
	tests := []struct {
		name        string
		userSlug    string
		trackSlug   string
		respStatus  int
		respBody    string
		want        *Track
		wantReqPath string
		wantErr     string
	}{
		{
			name:       "when found",
			userSlug:   "forss",
			trackSlug:  "flickermood",
			respStatus: http.StatusOK,
			respBody: `<html><body><script>window.__sc_hydration = [
				{"hydratable":"user","data":{"username":"Forss","permalink":"forss","permalink_url":"https://soundcloud.com/forss"}},
				{"hydratable":"sound","data":{
					"kind":"track",
					"urn":"soundcloud:tracks:293",
					"title":"Flickermood",
					"permalink":"flickermood",
					"permalink_url":"https://soundcloud.com/forss/flickermood",
					"user":{
						"username":"Forss",
						"permalink":"forss",
						"permalink_url":"https://soundcloud.com/forss"
					}
				}}
			];</script></body></html>`,
			wantReqPath: "/forss/flickermood",
			want: &Track{
				Title:        "Flickermood",
				Permalink:    "flickermood",
				PermalinkURL: "https://soundcloud.com/forss/flickermood",
				User: User{
					Username:     "Forss",
					Permalink:    "forss",
					PermalinkURL: "https://soundcloud.com/forss",
				},
			},
		},
		{
			name:        "when not found",
			userSlug:    "missing",
			trackSlug:   "track",
			respStatus:  http.StatusNotFound,
			wantReqPath: "/missing/track",
			wantErr:     "failed to fetch track page: not found",
		},
		{
			name:        "when hydration script missing",
			userSlug:    "forss",
			trackSlug:   "flickermood",
			respStatus:  http.StatusOK,
			respBody:    `<html><body>No hydration here</body></html>`,
			wantReqPath: "/forss/flickermood",
			wantErr:     "failed to find track hydration data: hydration script not found in html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respBody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			client := NewHTTPClient(WithWebURL(srv.URL))

			result, err := client.FetchTrack(t.Context(), tt.userSlug, tt.trackSlug)
			if tt.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestHTTPClient_FetchAlbum(t *testing.T) {
	tests := []struct {
		name        string
		userSlug    string
		setSlug     string
		respStatus  int
		respBody    string
		want        *Album
		wantReqPath string
		wantErr     string
	}{
		{
			name:       "when found",
			userSlug:   "forss",
			setSlug:    "soulhack",
			respStatus: http.StatusOK,
			respBody: `<html><body><script>window.__sc_hydration = [
				{"hydratable":"playlist","data":{
					"kind":"playlist",
					"urn":"soundcloud:playlists:123",
					"title":"Soulhack",
					"permalink":"soulhack",
					"permalink_url":"https://soundcloud.com/forss/sets/soulhack",
					"set_type":"album",
					"track_count":2,
					"user":{
						"username":"Forss",
						"permalink":"forss",
						"permalink_url":"https://soundcloud.com/forss"
					},
					"tracks":[
						{"kind":"track","title":"Flickermood","permalink":"flickermood","permalink_url":"https://soundcloud.com/forss/flickermood"},
						{"kind":"track","title":"Using Splashes","permalink":"using-splashes","permalink_url":"https://soundcloud.com/forss/using-splashes"}
					]
				}}
			];</script></body></html>`,
			wantReqPath: "/forss/sets/soulhack",
			want: &Album{
				Title:        "Soulhack",
				Permalink:    "soulhack",
				PermalinkURL: "https://soundcloud.com/forss/sets/soulhack",
				User: User{
					Username:     "Forss",
					Permalink:    "forss",
					PermalinkURL: "https://soundcloud.com/forss",
				},
			},
		},
		{
			name:        "when not found",
			userSlug:    "missing",
			setSlug:     "album",
			respStatus:  http.StatusNotFound,
			wantReqPath: "/missing/sets/album",
			wantErr:     "failed to fetch album page: not found",
		},
		{
			name:        "when hydration script missing",
			userSlug:    "forss",
			setSlug:     "soulhack",
			respStatus:  http.StatusOK,
			respBody:    `<html><body>No hydration here</body></html>`,
			wantReqPath: "/forss/sets/soulhack",
			wantErr:     "failed to find album hydration data: hydration script not found in html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respBody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			client := NewHTTPClient(WithWebURL(srv.URL))

			result, err := client.FetchAlbum(t.Context(), tt.userSlug, tt.setSlug)
			if tt.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestHTTPClient_SearchTrack(t *testing.T) {
	var homepageHits atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			homepageHits.Add(1)
			_, err := w.Write([]byte(`<html><body><script>window.__sc_hydration = [
				{"hydratable":"apiClient","data":{"id":"test-client-id","isExpiring":false}}
			];</script></body></html>`))
			require.NoError(t, err)
		case "/search/tracks":
			require.Equal(t, "autechre nil", r.URL.Query().Get("q"))
			require.Equal(t, "test-client-id", r.URL.Query().Get("client_id"))
			_, err := w.Write([]byte(`{
				"collection": [{
					"kind":"track",
					"urn":"soundcloud:tracks:1441462279",
					"title":"Nil",
					"permalink":"nil",
					"permalink_url":"https://soundcloud.com/autechreofficial/nil",
					"user":{
						"username":"Autechre",
						"permalink":"autechreofficial",
						"permalink_url":"https://soundcloud.com/autechreofficial"
					}
				}]
			}`))
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	client := NewHTTPClient(
		WithWebURL(srv.URL),
		WithAPIURL(srv.URL),
	)

	track, err := client.SearchTrack(t.Context(), "autechre", "nil")
	require.NoError(t, err)
	require.Equal(t, "Nil", track.Title)
	require.Equal(t, "Autechre", track.User.Username)

	_, err = client.SearchTrack(t.Context(), "autechre", "nil")
	require.NoError(t, err)
	require.Equal(t, int32(1), homepageHits.Load())
}

func TestHTTPClient_SearchAlbum(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, err := w.Write([]byte(`<html><body><script>window.__sc_hydration = [
				{"hydratable":"apiClient","data":{"id":"test-client-id","isExpiring":false}}
			];</script></body></html>`))
			require.NoError(t, err)
		case "/search/albums":
			require.Equal(t, "autechre amber", r.URL.Query().Get("q"))
			require.Equal(t, "test-client-id", r.URL.Query().Get("client_id"))
			_, err := w.Write([]byte(`{
				"collection": [{
					"kind":"playlist",
					"urn":"soundcloud:playlists:1566833785",
					"title":"Amber",
					"permalink":"amber-384961498",
					"permalink_url":"https://soundcloud.com/autechreofficial/sets/amber-384961498",
					"set_type":"album",
					"track_count":11,
					"user":{
						"username":"Autechre",
						"permalink":"autechreofficial",
						"permalink_url":"https://soundcloud.com/autechreofficial"
					}
				}]
			}`))
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	client := NewHTTPClient(
		WithWebURL(srv.URL),
		WithAPIURL(srv.URL),
	)

	album, err := client.SearchAlbum(t.Context(), "autechre", "amber")
	require.NoError(t, err)
	require.Equal(t, "Amber", album.Title)
}
