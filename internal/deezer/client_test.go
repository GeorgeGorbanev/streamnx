package deezer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// https://api.deezer.com/track/3135556
func TestHTTPClient_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		responseMock  string
		expectedTrack *Track
		expectedErr   error
	}{
		{
			name:    "when track found",
			inputID: "3135556",
			responseMock: `{
				"id": 3135556,
				"title": "Harder, Better, Faster, Stronger",
				"artist": {
					"name": "Daft Punk"
				},
				"album": {
					"id": 302127,
					"title": "Discovery"
				}
			}`,
			expectedTrack: &Track{
				ID:    3135556,
				Title: "Harder, Better, Faster, Stronger",
				Artist: Artist{
					Name: "Daft Punk",
				},
				Album: AlbumInfo{
					ID:    302127,
					Title: "Discovery",
				},
			},
		},
		{
			name:    "when track not found",
			inputID: "0",
			responseMock: `{
				"error": {
					"type": "DataException",
					"message": "no data",
					"code": 800
				}
			}`,
			expectedTrack: nil,
			expectedErr:   NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/track/"+tt.inputID, r.URL.Path)

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(
				WithAPIURL(apiServerMock.URL),
			)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			track, err := client.FetchTrack(ctx, tt.inputID)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, track)
			}
		})
	}
}

// https://api.deezer.com/search?q=artist:"aloe blacc" track:"i need a dollar"
func TestHTTPClient_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		title         string
		responseMock  string
		expectedTrack *Track
		expectedErr   error
	}{
		{
			name:   "when track found",
			artist: "Aloe Blacc",
			title:  "I Need A Dollar",
			responseMock: `{
				"data": [
					{
						"id": 445997982,
						"title": "I Need A Dollar",
						"artist": {
							"name": "Aloe Blacc"
						},
						"album": {
							"id": 44290462,
							"title": "Good Things"
						}
					}
				]
			}`,
			expectedTrack: &Track{
				ID:    445997982,
				Title: "I Need A Dollar",
				Artist: Artist{
					Name: "Aloe Blacc",
				},
				Album: AlbumInfo{
					ID:    44290462,
					Title: "Good Things",
				},
			},
		},
		{
			name:   "when track not found",
			artist: "Unknown Artist",
			title:  "Unknown Track",
			responseMock: `{
				"data": []
			}`,
			expectedTrack: nil,
			expectedErr:   NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/search", r.URL.Path)
				require.Equal(t, fmt.Sprintf(`artist:"%s" track:"%s"`, tt.artist, tt.title), r.URL.Query().Get("q"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(
				WithAPIURL(apiServerMock.URL),
			)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			track, err := client.SearchTrack(ctx, tt.artist, tt.title)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, track)
			}
		})
	}
}

// https://api.deezer.com/album/302127
func TestHTTPClient_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		responseMock  string
		expectedAlbum *Album
		expectedErr   error
	}{
		{
			name:    "when album found",
			inputID: "302127",
			responseMock: `{
				"id": 302127,
				"title": "Discovery",
				"artist": {
					"name": "Daft Punk"
				}
			}`,
			expectedAlbum: &Album{
				ID:    302127,
				Title: "Discovery",
				Artist: Artist{
					Name: "Daft Punk",
				},
			},
		},
		{
			name:    "when album not found",
			inputID: "0",
			responseMock: `{
				"error": {
					"type": "DataException",
					"message": "no data",
					"code": 800
				}
			}`,
			expectedAlbum: nil,
			expectedErr:   NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/album/"+tt.inputID, r.URL.Path)

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(
				WithAPIURL(apiServerMock.URL),
			)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			album, err := client.FetchAlbum(ctx, tt.inputID)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, album)
			}
		})
	}
}

// https://api.deezer.com/search?q=artist:"massive attack" album:"mezzanine"
func TestHTTPClient_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		title         string
		responseMock  string
		expectedAlbum *Album
		expectedErr   error
	}{
		{
			name:   "when album found",
			artist: "Massive Attack",
			title:  "Mezzanine",
			responseMock: `{
				"data": [
					{
						"id": 302127,
						"title": "Mezzanine",
						"artist": {
							"name": "Massive Attack"
						}
					}
				]
			}`,
			expectedAlbum: &Album{
				ID:    302127,
				Title: "Mezzanine",
				Artist: Artist{
					Name: "Massive Attack",
				},
			},
		},
		{
			name:   "when album not found",
			artist: "Unknown Artist",
			title:  "Unknown Album",
			responseMock: `{
				"data": []
			}`,
			expectedAlbum: nil,
			expectedErr:   NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/search", r.URL.Path)
				require.Equal(t, fmt.Sprintf(`artist:"%s" album:"%s"`, tt.artist, tt.title), r.URL.Query().Get("q"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(
				WithAPIURL(apiServerMock.URL),
			)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			album, err := client.SearchAlbum(ctx, tt.artist, tt.title)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, album)
			}
		})
	}
}

func TestHTTPClient_FollowCloak(t *testing.T) {
	tests := []struct {
		name             string
		inputID          string
		responseCode     int
		location         string
		expectedLocation string
		expectedErr      bool
	}{
		{
			name:             "when cloak link found",
			inputID:          "abc123",
			responseCode:     http.StatusFound,
			location:         "https://deezer.com/track/123456",
			expectedLocation: "https://deezer.com/track/123456",
			expectedErr:      false,
		},
		{
			name:             "when cloak link not found",
			inputID:          "notfound",
			responseCode:     http.StatusNotFound,
			location:         "",
			expectedLocation: "",
			expectedErr:      true,
		},
		{
			name:             "when no redirect",
			inputID:          "noredirect",
			responseCode:     http.StatusOK,
			location:         "",
			expectedLocation: "",
			expectedErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"id": 123}`))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			cloakServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodHead, r.Method)
				require.Equal(t, "/s/"+tt.inputID, r.URL.Path)

				if tt.location != "" {
					w.Header().Set("Location", tt.location)
				}
				w.WriteHeader(tt.responseCode)
			}))
			defer cloakServerMock.Close()

			client := NewHTTPClient(
				WithCloakBaseURL(cloakServerMock.URL),
				WithAPIURL(apiServerMock.URL),
			)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			location, err := client.FollowCloak(ctx, tt.inputID)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedLocation, location)
			}
		})
	}
}
