package deezer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient_FetchTrack(t *testing.T) {
	tests := []struct {
		name    string
		trackID string
		want    *Track
		wantErr error
	}{
		{
			name:    "when track found",
			trackID: "foundId",
			want: &Track{
				ID:     "foundID",
				Title:  "sampleTrackName",
				Artist: "sampleArtistName",
				URL:    "sampleURL",
			},
		},
		{
			name:    "when track not found",
			trackID: "notFoundId",
			wantErr: NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "https://api.deezer.com", r.Header.Get("Origin"))

				switch r.URL.Path {
				case "/track/foundId":
					_, err := w.Write([]byte(`{
						"id": "foundID",
						"title": "sampleTrackName",
						"artist": "sampleArtistName",
						"url": "sampleURL"
					}`))
					require.NoError(t, err)
				case "/track/notFoundId":
					w.WriteHeader(http.StatusNotFound)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient("clientID", "clientSecret", WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.FetchTrack(ctx, tt.trackID)
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
		name       string
		artistName string
		trackName  string
		want       *Track
		wantErr    error
	}{
		{
			name:       "when track found",
			artistName: "foundArtistName",
			trackName:  "foundTrackName",
			want: &Track{
				ID:     "foundID",
				Title:  "sampleTrackName",
				Artist: "sampleArtistName",
				URL:    "sampleURL",
			},
		},
		{
			name:       "when track not found",
			artistName: "notFoundArtistName",
			trackName:  "notFoundTrackName",
			wantErr:    NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "https://api.deezer.com", r.Header.Get("Origin"))
				require.Equal(t, "/search", r.URL.Path)

				q := r.URL.Query()
				require.Equal(t, `artist:"`+tt.artistName+`" track:"`+tt.trackName+`"`, q.Get("q"))

				var resp string
				if q.Get("q") == `artist:"foundArtistName" track:"foundTrackName"` {
					resp = `{
						"data": [
							{
								"id": "foundID",
								"title": "sampleTrackName",
								"artist": "sampleArtistName",
								"url": "sampleURL"
							}
						]
					}`
				} else {
					resp = `{"data": []}`
				}
				_, err := w.Write([]byte(resp))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient("clientID", "clientSecret", WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.SearchTrack(ctx, tt.artistName, tt.trackName)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestHTTPClient_FetchAlbum(t *testing.T) {
	tests := []struct {
		name    string
		albumID string
		want    *Album
		wantErr error
	}{
		{
			name:    "when album found",
			albumID: "foundId",
			want: &Album{
				ID:     "foundID",
				Title:  "sampleAlbumName",
				Artist: "sampleArtistName",
				URL:    "sampleURL",
			},
		},
		{
			name:    "when album not found",
			albumID: "notFoundId",
			wantErr: NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "https://api.deezer.com", r.Header.Get("Origin"))

				switch r.URL.Path {
				case "/album/foundId":
					_, err := w.Write([]byte(`{
						"id": "foundID",
						"title": "sampleAlbumName",
						"artist": "sampleArtistName",
						"url": "sampleURL"
					}`))
					require.NoError(t, err)
				case "/album/notFoundId":
					w.WriteHeader(http.StatusNotFound)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient("clientID", "clientSecret", WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.FetchAlbum(ctx, tt.albumID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestHTTPClient_SearchAlbum(t *testing.T) {
	tests := []struct {
		name       string
		artistName string
		albumName  string
		want       *Album
		wantErr    error
	}{
		{
			name:       "when album found",
			artistName: "foundArtistName",
			albumName:  "foundAlbumName",
			want: &Album{
				ID:     "foundID",
				Title:  "sampleAlbumName",
				Artist: "sampleArtistName",
				URL:    "sampleURL",
			},
		},
		{
			name:       "when album not found",
			artistName: "notFoundArtistName",
			albumName:  "notFoundAlbumName",
			wantErr:    NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "https://api.deezer.com", r.Header.Get("Origin"))
				require.Equal(t, "/search", r.URL.Path)

				q := r.URL.Query()
				require.Equal(t, `artist:"`+tt.artistName+`" album:"`+tt.albumName+`"`, q.Get("q"))

				var resp string
				if q.Get("q") == `artist:"foundArtistName" album:"foundAlbumName"` {
					resp = `{
						"data": [
							{
								"id": "foundID",
								"title": "sampleAlbumName",
								"artist": "sampleArtistName",
								"url": "sampleURL"
							}
						]
					}`
				} else {
					resp = `{"data": []}`
				}
				_, err := w.Write([]byte(resp))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient("clientID", "clientSecret", WithAPIURL(apiServerMock.URL))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.SearchAlbum(ctx, tt.artistName, tt.albumName)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}
