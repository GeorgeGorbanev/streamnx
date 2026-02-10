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

func TestHTTPClient_FetchAlbum(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		artistSlug  string
		respStatus  int
		respbody    string
		want        *Entity
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
						"byArtist": {
							"name": "Autechre"
						}
					}
					</script>
				</head>
				<body></body>
			</html>`,
			wantReqHost: "autechre.bandcamp.com",
			wantReqPath: "/album/amber",
			want: &Entity{
				Name:     "Amber",
				BandName: "Autechre",
			},
		},
		{
			name:        "when not found",
			id:          "notfound",
			artistSlug:  "notfound",
			wantReqHost: "notfound.bandcamp.com",
			wantReqPath: "/album/notfound",
			respStatus:  http.StatusNotFound,
			wantErr:     ErrNotFound,
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

			client := NewHTTPClient(
				WithAPIClient(mockSubdomainClient(srv.URL)),
				WithAPIScheme("http"),
			)

			result, err := client.FetchAlbum(t.Context(), tt.artistSlug, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestHTTPClient_FetchTrack(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		artistSlug  string
		respStatus  int
		respbody    string
		want        *Entity
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
						"byArtist": {
							"name": "Autechre"
						}
					}
					</script>
				</head>
				<body></body>
			</html>`,
			wantReqHost: "autechre.bandcamp.com",
			wantReqPath: "/track/nil",
			want: &Entity{
				Name:     "Nil",
				BandName: "Autechre",
			},
		},
		{
			name:        "when not found",
			id:          "notfound",
			artistSlug:  "notfound",
			wantReqHost: "notfound.bandcamp.com",
			wantReqPath: "/track/notfound",
			respStatus:  http.StatusNotFound,
			wantErr:     ErrNotFound,
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

			client := NewHTTPClient(
				WithAPIClient(mockSubdomainClient(srv.URL)),
				WithAPIScheme("http"),
			)

			result, err := client.FetchTrack(t.Context(), tt.artistSlug, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestHTTPClient_SearchTrack(t *testing.T) {
	tests := []struct {
		name            string
		artist          string
		title           string
		respStatus      int
		respBody        string
		want            *Entity
		wantErr         error
		wantErrContains string
	}{
		{
			name:       "when track found",
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
						}
					]
				}
			}`,
			want: &Entity{
				Name:        "Nil",
				BandName:    "Autechre",
				ItemURLPath: "https://autechre.bandcamp.com/track/nil",
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
			wantErr: ErrNotFound,
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
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.Equal(t, tt.artist+" "+tt.title, body.SearchText)
				require.Equal(t, string(trackEntityType), body.SearchFilter)
				w.WriteHeader(tt.respStatus)
				if tt.respBody != "" {
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				}
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewHTTPClient(
				WithAPIHost(srvURL.Host),
				WithAPIScheme(srvURL.Scheme),
			)

			result, err := client.SearchTrack(t.Context(), tt.artist, tt.title)
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

func TestHTTPClient_SearchAlbum(t *testing.T) {
	tests := []struct {
		name            string
		artist          string
		title           string
		respStatus      int
		respBody        string
		want            *Entity
		wantErr         error
		wantErrContains string
	}{
		{
			name:       "when album found",
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
						}
					]
				}
			}`,
			want: &Entity{
				Name:        "Music Has The Right To Children",
				BandName:    "Boards of Canada",
				ItemURLPath: "https://boardsofcanada.bandcamp.com/album/music-has-the-right-to-children",
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
			wantErr: ErrNotFound,
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
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.Equal(t, tt.artist+" "+tt.title, body.SearchText)
				require.Equal(t, string(albumEntityType), body.SearchFilter)
				w.WriteHeader(tt.respStatus)
				if tt.respBody != "" {
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				}
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewHTTPClient(
				WithAPIHost(srvURL.Host),
				WithAPIScheme(srvURL.Scheme),
			)

			result, err := client.SearchAlbum(t.Context(), tt.artist, tt.title)
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
