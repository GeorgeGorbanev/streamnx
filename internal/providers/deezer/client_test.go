package deezer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// https://api.deezer.com/track/3135556
func TestClient_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		responseMock  string
		expectedTrack track
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
			expectedTrack: track{
				ID:    3135556,
				Title: "Harder, Better, Faster, Stronger",
				Artist: artist{
					Name: "Daft Punk",
				},
				Album: albumInfo{
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
			expectedErr: errNotFound,
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

			client := NewClient(
				WithAPIURL(apiServerMock.URL),
			)

			track, err := client.fetchTrack(t.Context(), tt.inputID)
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
func TestClient_searchTracks(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		title          string
		responseMock   string
		expectedTracks []track
		expectedErr    string
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
					},
					{
						"id": 445997983,
						"title": "I Need A Dollar - Live",
						"artist": {
							"name": "Aloe Blacc"
						},
						"album": {
							"id": 44290463,
							"title": "Good Things Live"
						}
					}
				]
			}`,
			expectedTracks: []track{
				{
					ID:    445997982,
					Title: "I Need A Dollar",
					Artist: artist{
						Name: "Aloe Blacc",
					},
					Album: albumInfo{
						ID:    44290462,
						Title: "Good Things",
					},
				},
				{
					ID:    445997983,
					Title: "I Need A Dollar - Live",
					Artist: artist{
						Name: "Aloe Blacc",
					},
					Album: albumInfo{
						ID:    44290463,
						Title: "Good Things Live",
					},
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
			expectedTracks: []track{},
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

			client := NewClient(
				WithAPIURL(apiServerMock.URL),
			)

			tracks, err := client.searchTracks(t.Context(), tt.artist, tt.title)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, tracks)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTracks, tracks)
			}
		})
	}
}

// https://api.deezer.com/album/302127
func TestClient_fetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		responseMock  string
		expectedAlbum album
		expectedErr   error
	}{
		{
			name:    "when album found",
			inputID: "302127",
			responseMock: `{
				"id": 302127,
				"title": "Discovery",
				"label": "Daft Life",
				"artist": {
					"name": "Daft Punk"
				},
				"tracks": {
					"data": [
						{
							"id": 3135556,
							"title": "Harder, Better, Faster, Stronger",
							"artist": {
								"name": "Daft Punk"
							},
							"album": {
								"id": 302127,
								"title": "Discovery"
							}
						},
						{
							"id": 3135557,
							"title": "One More Time",
							"artist": {
								"name": "Daft Punk"
							},
							"album": {
								"id": 302127,
								"title": "Discovery"
							}
						}
					]
				}
			}`,
			expectedAlbum: album{
				ID:    302127,
				Title: "Discovery",
				Label: "Daft Life",
				Artist: artist{
					Name: "Daft Punk",
				},
				Tracks: trackData{
					Data: []track{
						{
							ID:    3135556,
							Title: "Harder, Better, Faster, Stronger",
							Artist: artist{
								Name: "Daft Punk",
							},
							Album: albumInfo{
								ID:    302127,
								Title: "Discovery",
							},
						},
						{
							ID:    3135557,
							Title: "One More Time",
							Artist: artist{
								Name: "Daft Punk",
							},
							Album: albumInfo{
								ID:    302127,
								Title: "Discovery",
							},
						},
					},
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
			expectedErr: errNotFound,
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

			client := NewClient(
				WithAPIURL(apiServerMock.URL),
			)

			album, err := client.fetchAlbum(t.Context(), tt.inputID)
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
func TestClient_searchAlbums(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		title          string
		responseMock   string
		expectedAlbums []album
		expectedErr    string
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
						"label": "Virgin Records",
						"artist": {
							"name": "Massive Attack"
						}
					},
					{
						"id": 302128,
						"title": "Mezzanine (Deluxe)",
						"label": "Virgin Records",
						"artist": {
							"name": "Massive Attack"
						}
					}
				]
			}`,
			expectedAlbums: []album{
				{
					ID:    302127,
					Title: "Mezzanine",
					Label: "Virgin Records",
					Artist: artist{
						Name: "Massive Attack",
					},
				},
				{
					ID:    302128,
					Title: "Mezzanine (Deluxe)",
					Label: "Virgin Records",
					Artist: artist{
						Name: "Massive Attack",
					},
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
			expectedAlbums: []album{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/search/album", r.URL.Path)
				require.Equal(t, fmt.Sprintf(`artist:"%s" album:"%s"`, tt.artist, tt.title), r.URL.Query().Get("q"))

				_, err := w.Write([]byte(tt.responseMock))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(
				WithAPIURL(apiServerMock.URL),
			)

			albums, err := client.searchAlbums(t.Context(), tt.artist, tt.title)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, albums)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbums, albums)
			}
		})
	}
}

func TestClient_FollowCloak(t *testing.T) {
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

			client := NewClient(
				WithCloakBaseURL(cloakServerMock.URL),
				WithAPIURL(apiServerMock.URL),
			)

			location, err := client.followCloak(t.Context(), tt.inputID)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedLocation, location)
			}
		})
	}
}
