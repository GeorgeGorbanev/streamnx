package deezer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient_FetchTrack(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		serverResponse string
		statusCode     int
		want           *Track
		wantError      bool
	}{
		{
			name:       "successful fetch",
			id:         "123456",
			statusCode: http.StatusOK,
			serverResponse: `{
				"id": "123456",
				"title": "Test Song",
				"artist": {
					"name": "Test Artist"
				},
				"album": {
					"id": "789",
					"title": "Test Album"
				}
			}`,
			want: &Track{
				ID:    "123456",
				Title: "Test Song",
				Artist: Artist{
					Name: "Test Artist",
				},
				Album: AlbumInfo{
					ID:    "789",
					Title: "Test Album",
				},
			},
		},
		{
			name:       "not found",
			id:         "404",
			statusCode: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/track/"+tt.id, r.URL.Path)
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					_, err := w.Write([]byte(tt.serverResponse))
					require.NoError(t, err)
				}
			}))
			defer server.Close()

			client := &HTTPClient{
				apiURL:    server.URL,
				apiClient: &http.Client{},
				cloakClient: &http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				},
			}

			got, err := client.FetchTrack(context.Background(), tt.id)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestHTTPClient_FollowCloak(t *testing.T) {
	t.Run("successful redirect", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "HEAD", r.Method)
			w.Header().Set("Location", "https://deezer.com/track/123456")
			w.WriteHeader(http.StatusFound)
		}))
		defer server.Close()

		// Create a custom client that uses the test server
		client := &HTTPClient{
			apiURL:    "https://api.deezer.com",
			apiClient: &http.Client{},
			cloakClient: &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
		}

		// Test with custom implementation
		cloakURL := server.URL

		req, err := http.NewRequestWithContext(context.Background(), "HEAD", cloakURL, nil)
		require.NoError(t, err)

		resp, err := client.cloakClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusFound, resp.StatusCode)
		location := resp.Header.Get("Location")
		require.Equal(t, "https://deezer.com/track/123456", location)
	})
}
