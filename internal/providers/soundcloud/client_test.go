package soundcloud

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_fetchTrack(t *testing.T) {
	tests := []struct {
		name        string
		userSlug    string
		trackSlug   string
		respStatus  int
		respBody    string
		want        track
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
					"description":"sample track description",
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
			want: track{
				URN:          "soundcloud:tracks:293",
				Title:        "Flickermood",
				Description:  "sample track description",
				Permalink:    "flickermood",
				PermalinkURL: "https://soundcloud.com/forss/flickermood",
				User: user{
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

			client := NewClient(WithWebURL(srv.URL))

			result, err := client.fetchTrack(t.Context(), tt.userSlug, tt.trackSlug)
			if tt.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestClient_fetchAlbum(t *testing.T) {
	tests := []struct {
		name        string
		userSlug    string
		setSlug     string
		respStatus  int
		respBody    string
		want        album
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
					"description":"sample album description",
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
			want: album{
				Title:        "Soulhack",
				Description:  "sample album description",
				Permalink:    "soulhack",
				PermalinkURL: "https://soundcloud.com/forss/sets/soulhack",
				User: user{
					Username:     "Forss",
					Permalink:    "forss",
					PermalinkURL: "https://soundcloud.com/forss",
				},
				Tracks: []track{
					{
						Title:        "Flickermood",
						Permalink:    "flickermood",
						PermalinkURL: "https://soundcloud.com/forss/flickermood",
					},
					{
						Title:        "Using Splashes",
						Permalink:    "using-splashes",
						PermalinkURL: "https://soundcloud.com/forss/using-splashes",
					},
				},
			},
		},
		{
			name:       "returns placeholder without additional requests",
			userSlug:   "artist",
			setSlug:    "album",
			respStatus: http.StatusOK,
			respBody: `<html><body><script>window.__sc_hydration = [
				{"hydratable":"playlist","data":{
					"title":"Album",
					"permalink":"album",
					"permalink_url":"https://soundcloud.com/artist/sets/album",
					"tracks":[{"kind":"track"}]
				}}
			];</script></body></html>`,
			wantReqPath: "/artist/sets/album",
			want: album{
				Title:        "Album",
				Permalink:    "album",
				PermalinkURL: "https://soundcloud.com/artist/sets/album",
				Tracks:       []track{{}},
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
			var requests atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respBody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			client := NewClient(WithWebURL(srv.URL))

			result, err := client.fetchAlbum(t.Context(), tt.userSlug, tt.setSlug)
			if tt.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}
			require.Equal(t, int32(1), requests.Load())
		})
	}
}

func TestClient_fetchTracksByNumericIDs(t *testing.T) {
	const (
		firstID      int64 = 3_000_000_001
		secondID     int64 = 3_000_000_002
		testClientID       = "test-client-id"
	)

	requestErr := errors.New("sample transport failure")
	tests := []struct {
		name                 string
		ids                  []int64
		wantIDs              string
		responseStatus       int
		responseBody         string
		want                 []track
		wantErrContains      string
		errIs                error
		transportErr         error
		wantClientIDRequests int32
	}{
		{
			name:           "returns provider response order",
			ids:            []int64{firstID, secondID},
			wantIDs:        "3000000001,3000000002",
			responseStatus: http.StatusOK,
			responseBody: `[
				{"id":3000000002,"urn":"soundcloud:tracks:3000000002","title":"Second hydrated","permalink":"second-hydrated","permalink_url":"https://soundcloud.com/artist/second-hydrated","user":{"permalink":"artist"}},
				{"id":3000000001,"urn":"soundcloud:tracks:3000000001","title":"First hydrated","permalink":"first-hydrated","permalink_url":"https://soundcloud.com/artist/first-hydrated","user":{"permalink":"artist"}}
			]`,
			want: []track{
				{
					ID:           secondID,
					URN:          "soundcloud:tracks:3000000002",
					Title:        "Second hydrated",
					Permalink:    "second-hydrated",
					PermalinkURL: "https://soundcloud.com/artist/second-hydrated",
					User:         user{Permalink: "artist"},
				},
				{
					ID:           firstID,
					URN:          "soundcloud:tracks:3000000001",
					Title:        "First hydrated",
					Permalink:    "first-hydrated",
					PermalinkURL: "https://soundcloud.com/artist/first-hydrated",
					User:         user{Permalink: "artist"},
				},
			},
			wantClientIDRequests: 1,
		},
		{
			name:           "returns partial provider response",
			ids:            []int64{200, 201},
			wantIDs:        "200,201",
			responseStatus: http.StatusOK,
			responseBody: `[
				{"id":201,"urn":"soundcloud:tracks:201","title":"Hydrated","permalink":"hydrated","permalink_url":"https://soundcloud.com/artist/hydrated"}
			]`,
			want: []track{
				{
					ID:           201,
					URN:          "soundcloud:tracks:201",
					Title:        "Hydrated",
					Permalink:    "hydrated",
					PermalinkURL: "https://soundcloud.com/artist/hydrated",
				},
			},
			wantClientIDRequests: 1,
		},
		{
			name:            "returns request error",
			ids:             []int64{200},
			wantIDs:         "200",
			transportErr:    requestErr,
			wantErrContains: "failed to fetch tracks by numeric ids",
			errIs:           requestErr,
		},
		{
			name:                 "returns status error",
			ids:                  []int64{200},
			wantIDs:              "200",
			responseStatus:       http.StatusNotFound,
			wantErrContains:      "failed to fetch tracks by numeric ids: not found",
			errIs:                errNotFound,
			wantClientIDRequests: 1,
		},
		{
			name:                 "returns malformed response error",
			ids:                  []int64{200},
			wantIDs:              "200",
			responseStatus:       http.StatusOK,
			responseBody:         `{`,
			wantErrContains:      "failed to unmarshal tracks by numeric ids response",
			wantClientIDRequests: 1,
		},
	}

	require.Greater(t, firstID, int64(1<<31-1))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var clientIDRequests atomic.Int32
			var bulkRequests atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/":
					clientIDRequests.Add(1)
					_, err := w.Write([]byte(`<html><body><script>window.__sc_hydration = [
						{"hydratable":"apiClient","data":{"id":"test-client-id"}}
					];</script></body></html>`))
					require.NoError(t, err)
				case "/tracks":
					bulkRequests.Add(1)
					require.Equal(t, tt.wantIDs, r.URL.Query().Get("ids"))
					require.Equal(t, testClientID, r.URL.Query().Get("client_id"))
					w.WriteHeader(tt.responseStatus)
					_, err := w.Write([]byte(tt.responseBody))
					require.NoError(t, err)
				default:
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
			}))
			defer srv.Close()

			httpClient := srv.Client()
			if tt.transportErr != nil {
				baseTransport := httpClient.Transport
				httpClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
					if r.URL.Path == "/tracks" {
						bulkRequests.Add(1)
						require.Equal(t, tt.wantIDs, r.URL.Query().Get("ids"))
						require.Equal(t, testClientID, r.URL.Query().Get("client_id"))
						return nil, tt.transportErr
					}
					return baseTransport.RoundTrip(r)
				})
			}

			client := NewClient(
				WithWebURL(srv.URL),
				WithAPIURL(srv.URL),
				WithHTTPClient(httpClient),
			)
			if tt.transportErr != nil {
				client.clientID = testClientID
			}
			got, err := client.fetchTracksByNumericIDs(t.Context(), tt.ids)

			if tt.wantErrContains != "" {
				require.Nil(t, got)
				require.ErrorContains(t, err, tt.wantErrContains)
				if tt.errIs != nil {
					require.ErrorIs(t, err, tt.errIs)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
			require.Equal(t, tt.wantClientIDRequests, clientIDRequests.Load())
			require.Equal(t, int32(1), bulkRequests.Load())
		})
	}
}

func TestClient_fetchTrackByURN(t *testing.T) {
	const (
		trackURN     = "soundcloud:tracks:urn-only"
		testClientID = "test-client-id"
	)

	tests := []struct {
		name        string
		status      int
		body        string
		want        track
		wantErr     string
		errIs       error
		wantAPICall int32
	}{
		{
			name:   "found",
			status: http.StatusOK,
			body: `{
				"urn":"soundcloud:tracks:urn-only",
				"title":"URN hydrated",
				"permalink":"urn-hydrated",
				"permalink_url":"https://soundcloud.com/artist/urn-hydrated",
				"user":{"permalink":"artist"}
			}`,
			want: track{
				URN:          trackURN,
				Title:        "URN hydrated",
				Permalink:    "urn-hydrated",
				PermalinkURL: "https://soundcloud.com/artist/urn-hydrated",
				User:         user{Permalink: "artist"},
			},
			wantAPICall: 1,
		},
		{
			name:        "not found",
			status:      http.StatusNotFound,
			wantErr:     `failed to fetch track by urn "soundcloud:tracks:urn-only": not found`,
			errIs:       errNotFound,
			wantAPICall: 1,
		},
		{
			name:        "operational failure",
			status:      http.StatusInternalServerError,
			wantErr:     `failed to fetch track by urn "soundcloud:tracks:urn-only": unexpected status code: 500`,
			wantAPICall: 1,
		},
		{
			name:        "malformed response",
			status:      http.StatusOK,
			body:        `{`,
			wantErr:     `failed to unmarshal track by urn "soundcloud:tracks:urn-only" response: unexpected end of JSON input`,
			wantAPICall: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var apiRequests atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/":
					_, err := w.Write([]byte(`<html><body><script>window.__sc_hydration = [
						{"hydratable":"apiClient","data":{"id":"test-client-id"}}
					];</script></body></html>`))
					require.NoError(t, err)
				case "/tracks/" + trackURN:
					apiRequests.Add(1)
					require.Equal(t, "/tracks/soundcloud%3Atracks%3Aurn-only", r.URL.EscapedPath())
					require.Equal(t, testClientID, r.URL.Query().Get("client_id"))
					w.WriteHeader(tt.status)
					_, err := w.Write([]byte(tt.body))
					require.NoError(t, err)
				default:
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
			}))
			defer srv.Close()

			client := NewClient(
				WithWebURL(srv.URL),
				WithAPIURL(srv.URL),
			)
			got, err := client.fetchTrackByURN(t.Context(), trackURN)

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Zero(t, got)
				if tt.errIs != nil {
					require.ErrorIs(t, err, tt.errIs)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
			require.Equal(t, tt.wantAPICall, apiRequests.Load())
		})
	}
}

func TestClient_searchTracks(t *testing.T) {
	tests := []struct {
		name       string
		respBody   string
		wantTracks []track
	}{
		{
			name: "multiple candidates",
			respBody: `{
				"collection": [{
					"kind":"track",
					"urn":"soundcloud:tracks:1441462279",
					"title":"Nil",
					"description":"first search description",
					"permalink":"nil",
					"permalink_url":"https://soundcloud.com/autechreofficial/nil",
					"user":{
						"username":"Autechre",
						"permalink":"autechreofficial",
						"permalink_url":"https://soundcloud.com/autechreofficial"
					}
				}, {
					"kind":"track",
					"urn":"soundcloud:tracks:1441462280",
					"title":"Nil Alternate",
					"description":"second search description",
					"permalink":"nil-alternate",
					"permalink_url":"https://soundcloud.com/autechreofficial/nil-alternate",
					"user":{
						"username":"Autechre Official",
						"permalink":"autechreofficial",
						"permalink_url":"https://soundcloud.com/autechreofficial"
					}
				}]
			}`,
			wantTracks: []track{
				{
					URN:          "soundcloud:tracks:1441462279",
					Title:        "Nil",
					Description:  "first search description",
					Permalink:    "nil",
					PermalinkURL: "https://soundcloud.com/autechreofficial/nil",
					User: user{
						Username:     "Autechre",
						Permalink:    "autechreofficial",
						PermalinkURL: "https://soundcloud.com/autechreofficial",
					},
				},
				{
					URN:          "soundcloud:tracks:1441462280",
					Title:        "Nil Alternate",
					Description:  "second search description",
					Permalink:    "nil-alternate",
					PermalinkURL: "https://soundcloud.com/autechreofficial/nil-alternate",
					User: user{
						Username:     "Autechre Official",
						Permalink:    "autechreofficial",
						PermalinkURL: "https://soundcloud.com/autechreofficial",
					},
				},
			},
		},
		{
			name:       "empty candidates",
			respBody:   `{"collection":[]}`,
			wantTracks: []track{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				default:
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
			}))
			defer srv.Close()

			client := NewClient(
				WithWebURL(srv.URL),
				WithAPIURL(srv.URL),
			)

			tracks, err := client.searchTracks(t.Context(), "autechre", "nil")
			require.NoError(t, err)
			require.Equal(t, tt.wantTracks, tracks)

			_, err = client.searchTracks(t.Context(), "autechre", "nil")
			require.NoError(t, err)
			require.Equal(t, int32(1), homepageHits.Load())
		})
	}
}

func TestClient_RefreshClientIDAfterUnauthorized(t *testing.T) {
	var clientIDRequests atomic.Int32
	var apiRequests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			clientIDRequests.Add(1)
			_, err := w.Write([]byte(`<script>window.__sc_hydration = [
				{"hydratable":"apiClient","data":{"id":"fresh-client-id"}}
			];</script>`))
			require.NoError(t, err)
		case "/search/tracks":
			apiRequests.Add(1)
			if r.URL.Query().Get("client_id") == "expired-client-id" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			require.Equal(t, "fresh-client-id", r.URL.Query().Get("client_id"))
			_, err := w.Write([]byte(`{"collection":[]}`))
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithAPIURL(server.URL),
		WithWebURL(server.URL),
	)
	client.clientID = "expired-client-id"

	tracks, err := client.searchTracks(t.Context(), "artist", "title")

	require.NoError(t, err)
	require.Empty(t, tracks)
	require.Equal(t, int32(1), clientIDRequests.Load())
	require.Equal(t, int32(2), apiRequests.Load())
}

func TestClient_ConcurrentUnauthorizedRefreshesClientIDOnce(t *testing.T) {
	var clientIDRequests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			clientIDRequests.Add(1)
			time.Sleep(10 * time.Millisecond)
			_, err := w.Write([]byte(`<script>window.__sc_hydration = [
				{"hydratable":"apiClient","data":{"id":"fresh-client-id"}}
			];</script>`))
			require.NoError(t, err)
		case "/search/tracks":
			if r.URL.Query().Get("client_id") == "expired-client-id" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_, err := w.Write([]byte(`{"collection":[]}`))
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithAPIURL(server.URL),
		WithWebURL(server.URL),
	)
	client.clientID = "expired-client-id"

	const concurrency = 10
	errs := make(chan error, concurrency)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for range concurrency {
		go func() {
			defer wg.Done()
			_, err := client.searchTracks(t.Context(), "artist", "title")
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
	require.Equal(t, int32(1), clientIDRequests.Load())
}

func TestClient_RetriesUnauthorizedOnlyOnce(t *testing.T) {
	var apiRequests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, err := w.Write([]byte(`<script>window.__sc_hydration = [
				{"hydratable":"apiClient","data":{"id":"fresh-client-id"}}
			];</script>`))
			require.NoError(t, err)
		case "/search/tracks":
			apiRequests.Add(1)
			w.WriteHeader(http.StatusUnauthorized)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithAPIURL(server.URL),
		WithWebURL(server.URL),
	)
	client.clientID = "expired-client-id"

	_, err := client.searchTracks(t.Context(), "artist", "title")

	require.ErrorContains(t, err, "unexpected status code: 401")
	require.Equal(t, int32(2), apiRequests.Load())
}

func TestClient_searchAlbums(t *testing.T) {
	tests := []struct {
		name       string
		respBody   string
		wantAlbums []album
	}{
		{
			name: "multiple candidates",
			respBody: `{
				"collection": [{
					"kind":"playlist",
					"urn":"soundcloud:playlists:1566833785",
					"title":"Amber",
					"description":"first album search description",
					"permalink":"amber-384961498",
					"permalink_url":"https://soundcloud.com/autechreofficial/sets/amber-384961498",
					"set_type":"album",
					"track_count":11,
					"user":{
						"username":"Autechre",
						"permalink":"autechreofficial",
						"permalink_url":"https://soundcloud.com/autechreofficial"
					}
				}, {
					"kind":"playlist",
					"urn":"soundcloud:playlists:1566833786",
					"title":"Amber Alternate",
					"description":"second album search description",
					"permalink":"amber-alternate",
					"permalink_url":"https://soundcloud.com/autechreofficial/sets/amber-alternate",
					"set_type":"album",
					"track_count":10,
					"user":{
						"username":"Autechre Official",
						"permalink":"autechreofficial",
						"permalink_url":"https://soundcloud.com/autechreofficial"
					}
				}]
			}`,
			wantAlbums: []album{
				{
					Title:        "Amber",
					Description:  "first album search description",
					Permalink:    "amber-384961498",
					PermalinkURL: "https://soundcloud.com/autechreofficial/sets/amber-384961498",
					User: user{
						Username:     "Autechre",
						Permalink:    "autechreofficial",
						PermalinkURL: "https://soundcloud.com/autechreofficial",
					},
				},
				{
					Title:        "Amber Alternate",
					Description:  "second album search description",
					Permalink:    "amber-alternate",
					PermalinkURL: "https://soundcloud.com/autechreofficial/sets/amber-alternate",
					User: user{
						Username:     "Autechre Official",
						Permalink:    "autechreofficial",
						PermalinkURL: "https://soundcloud.com/autechreofficial",
					},
				},
			},
		},
		{
			name:       "empty candidates",
			respBody:   `{"collection":[]}`,
			wantAlbums: []album{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				default:
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
			}))
			defer srv.Close()

			client := NewClient(
				WithWebURL(srv.URL),
				WithAPIURL(srv.URL),
			)

			albums, err := client.searchAlbums(t.Context(), "autechre", "amber")
			require.NoError(t, err)
			require.Equal(t, tt.wantAlbums, albums)
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
