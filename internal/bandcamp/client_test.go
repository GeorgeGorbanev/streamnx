package bandcamp

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
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
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqHost, r.Host)
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				w.Write([]byte(tt.respbody))
			}))
			defer srv.Close()

			client := NewHTTPClient(
				WithHTTPTransport(&http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						host, _, err := net.SplitHostPort(addr)
						if err == nil && strings.HasSuffix(host, ".bandcamp.com") {
							d := net.Dialer{Timeout: 5 * time.Second}
							return d.DialContext(ctx, network, strings.TrimPrefix(srv.URL, "https://"))
						}
						d := net.Dialer{Timeout: 5 * time.Second}
						return d.DialContext(ctx, network, addr)
					},
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}),
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
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqHost, r.Host)
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				w.Write([]byte(tt.respbody))
			}))
			defer srv.Close()

			client := NewHTTPClient(
				WithHTTPTransport(&http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						host, _, err := net.SplitHostPort(addr)
						if err == nil && strings.HasSuffix(host, ".bandcamp.com") {
							d := net.Dialer{Timeout: 5 * time.Second}
							return d.DialContext(ctx, network, strings.TrimPrefix(srv.URL, "https://"))
						}
						d := net.Dialer{Timeout: 5 * time.Second}
						return d.DialContext(ctx, network, addr)
					},
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}),
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
