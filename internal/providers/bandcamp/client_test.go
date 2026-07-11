package bandcamp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_fetchAlbum(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		artistSlug  string
		respStatus  int
		respbody    string
		want        Entity
		wantReqHost string
		wantReqPath string
		wantErr     error
	}{
		{
			name:       "when found",
			id:         "amber",
			artistSlug: "autechre",
			respStatus: http.StatusOK,
			respbody: `<html>
				<head>
					<script type="application/ld+json">
					{
						"name": "Amber",
						"description": "sample album description",
						"byArtist": {
							"name": "Autechre"
						},
						"publisher": {
							"name": "Warp Records"
						},
						"track": {
							"@type": "ItemList",
							"numberOfItems": 2,
							"itemListElement": [
								{
									"@type": "ListItem",
									"position": 1,
									"item": {
										"@type": "MusicRecording",
										"@id": "https://autechre.bandcamp.com/track/foil",
										"name": "Foil",
										"mainEntityOfPage": "https://autechre.bandcamp.com/track/foil"
									}
								},
								{
									"@type": "ListItem",
									"position": 2,
									"item": {
										"@type": "MusicRecording",
										"@id": "https://autechre.bandcamp.com/track/montreal",
										"name": "Montreal",
										"mainEntityOfPage": "https://autechre.bandcamp.com/track/montreal"
									}
								}
							]
						}
					}
					</script>
				</head>
				<body></body>
			</html>`,
			wantReqHost: "autechre.bandcamp.com",
			wantReqPath: "/album/amber",
			want: Entity{
				Name:        "Amber",
				BandName:    "Autechre",
				CreatorName: "Warp Records",
				Description: "sample album description",
				URL:         "http://autechre.bandcamp.com/album/amber",
				TrackURLs: []string{
					"https://autechre.bandcamp.com/track/foil",
					"https://autechre.bandcamp.com/track/montreal",
				},
			},
		},
		{
			name:        "when not found",
			id:          "notfound",
			artistSlug:  "notfound",
			wantReqHost: "notfound.bandcamp.com",
			wantReqPath: "/album/notfound",
			respStatus:  http.StatusNotFound,
			wantErr:     errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqHost, r.Host)
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respbody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			client := NewClient(
				WithAPIClient(mockSubdomainClient(srv.URL)),
				WithAPI("http", "bandcamp.com"),
			)

			result, err := client.fetchAlbum(t.Context(), tt.artistSlug, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_fetchTrack(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		artistSlug  string
		respStatus  int
		respbody    string
		want        Entity
		wantReqHost string
		wantReqPath string
		wantErr     error
	}{
		{
			name:       "when found",
			id:         "nil",
			artistSlug: "autechre",
			respStatus: http.StatusOK,
			respbody: `<html>
				<head>
					<script type="application/ld+json">
					{
						"name": "Nil",
						"description": "sample track description",
						"byArtist": {
							"name": "Autechre"
						},
						"publisher": {
							"name": "Warp Records"
						}
					}
					</script>
				</head>
				<body></body>
			</html>`,
			wantReqHost: "autechre.bandcamp.com",
			wantReqPath: "/track/nil",
			want: Entity{
				Name:        "Nil",
				BandName:    "Autechre",
				CreatorName: "Warp Records",
				Description: "sample track description",
				URL:         "http://autechre.bandcamp.com/track/nil",
			},
		},
		{
			name:        "when not found",
			id:          "notfound",
			artistSlug:  "notfound",
			wantReqHost: "notfound.bandcamp.com",
			wantReqPath: "/track/notfound",
			respStatus:  http.StatusNotFound,
			wantErr:     errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqHost, r.Host)
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respbody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			client := NewClient(
				WithAPIClient(mockSubdomainClient(srv.URL)),
				WithAPI("http", "bandcamp.com"),
			)

			result, err := client.fetchTrack(t.Context(), tt.artistSlug, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_searchTracks(t *testing.T) {
	tests := []struct {
		name            string
		artist          string
		title           string
		respStatus      int
		respBody        string
		want            []Entity
		wantErr         error
		wantErrContains string
	}{
		{
			name:       "when multiple track candidates found",
			artist:     "Autechre",
			title:      "Nil",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": [
						{
							"name": "Nil",
							"band_name": "Autechre",
							"item_url_root": "https://autechre.bandcamp.com",
							"item_url_path": "https://autechre.bandcamp.com/track/nil"
						},
						{
							"name": "Nil (alternate)",
							"band_name": "Autechre",
							"item_url_root": "https://autechre.bandcamp.com",
							"item_url_path": "https://autechre.bandcamp.com/track/nil-alternate"
						}
					]
				}
			}`,
			want: []Entity{
				{
					Name:     "Nil",
					BandName: "Autechre",
					URL:      "https://autechre.bandcamp.com/track/nil",
				},
				{
					Name:     "Nil (alternate)",
					BandName: "Autechre",
					URL:      "https://autechre.bandcamp.com/track/nil-alternate",
				},
			},
		},
		{
			name:       "when track not found",
			artist:     "Unknown",
			title:      "Track",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": []
				}
			}`,
			want: []Entity{},
		},
		{
			name:            "when api responds with error",
			artist:          "Error",
			title:           "State",
			respStatus:      http.StatusInternalServerError,
			wantErrContains: "unexpected status code: 500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/api/bcsearch_public_api/1/autocomplete_elastic", r.URL.Path)
				require.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Content-Type"))
				var body struct {
					SearchText   string `json:"search_text"`
					SearchFilter string `json:"search_filter"`
					FullPage     bool   `json:"full_page"`
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.Equal(t, tt.artist+" "+tt.title, body.SearchText)
				require.Equal(t, string(trackEntityType), body.SearchFilter)
				require.True(t, body.FullPage)
				w.WriteHeader(tt.respStatus)
				if tt.respBody != "" {
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				}
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewClient(
				WithAPI(srvURL.Scheme, srvURL.Host),
			)

			result, err := client.searchTracks(t.Context(), tt.artist, tt.title)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else if tt.wantErrContains != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_searchAlbums(t *testing.T) {
	tests := []struct {
		name            string
		artist          string
		title           string
		respStatus      int
		respBody        string
		want            []Entity
		wantErr         error
		wantErrContains string
	}{
		{
			name:       "when multiple album candidates found",
			artist:     "Boards of Canada",
			title:      "Music Has The Right To Children",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": [
						{
							"name": "Music Has The Right To Children",
							"band_name": "Boards of Canada",
							"item_url_root": "https://boardsofcanada.bandcamp.com",
							"item_url_path": "https://boardsofcanada.bandcamp.com/album/music-has-the-right-to-children"
						},
						{
							"name": "Hi Scores",
							"band_name": "Boards of Canada",
							"item_url_root": "https://boardsofcanada.bandcamp.com",
							"item_url_path": "https://boardsofcanada.bandcamp.com/album/hi-scores"
						}
					]
				}
			}`,
			want: []Entity{
				{
					Name:     "Music Has The Right To Children",
					BandName: "Boards of Canada",
					URL:      "https://boardsofcanada.bandcamp.com/album/music-has-the-right-to-children",
				},
				{
					Name:     "Hi Scores",
					BandName: "Boards of Canada",
					URL:      "https://boardsofcanada.bandcamp.com/album/hi-scores",
				},
			},
		},
		{
			name:       "when album not found",
			artist:     "Unknown",
			title:      "Album",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": []
				}
			}`,
			want: []Entity{},
		},
		{
			name:            "when api responds with error",
			artist:          "Boards of Canada",
			title:           "Hi Scores",
			respStatus:      http.StatusBadGateway,
			wantErrContains: "unexpected status code: 502",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/api/bcsearch_public_api/1/autocomplete_elastic", r.URL.Path)
				require.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Content-Type"))
				var body struct {
					SearchText   string `json:"search_text"`
					SearchFilter string `json:"search_filter"`
					FullPage     bool   `json:"full_page"`
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.Equal(t, tt.artist+" "+tt.title, body.SearchText)
				require.Equal(t, string(albumEntityType), body.SearchFilter)
				require.True(t, body.FullPage)
				w.WriteHeader(tt.respStatus)
				if tt.respBody != "" {
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				}
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewClient(
				WithAPI(srvURL.Scheme, srvURL.Host),
			)

			result, err := client.searchAlbums(t.Context(), tt.artist, tt.title)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else if tt.wantErrContains != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func mockSubdomainClient(u string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, _, err := net.SplitHostPort(addr)
				if err == nil && strings.HasSuffix(host, ".bandcamp.com") {
					d := net.Dialer{Timeout: 5 * time.Second}
					return d.DialContext(ctx, network, strings.TrimPrefix(u, "http://"))
				}
				d := net.Dialer{Timeout: 5 * time.Second}
				return d.DialContext(ctx, network, addr)
			},
		},
	}
}
