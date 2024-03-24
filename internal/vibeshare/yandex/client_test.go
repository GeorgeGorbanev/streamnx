package yandex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_GetTrack(t *testing.T) {
	tests := []struct {
		name    string
		trackID string
		want    *Track
	}{
		{
			name:    "when track found",
			trackID: "foundID",
			want: &Track{
				ID:    "1",
				Title: "sample title",
				Albums: []Album{
					{
						ID:    2,
						Title: "sample title",
						Artists: []Artist{
							{
								ID:   3,
								Name: "sample artist",
							},
						},
					},
				},
				Artists: []Artist{
					{
						ID:   4,
						Name: "sample artist",
					},
				},
			},
		},
		{
			name:    "when track not found",
			trackID: "notFoundID",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)

				if r.URL.Path == "/tracks/foundID" {
					_, err := w.Write([]byte(`{
						"result": [{
							"id": "1",
							"title": "sample title",
							"albums": [
								{
									"id": 2,
									"title": "sample title",
									"artists": [{"id": 3, "name": "sample artist"}]
								}
							],
							"artists": [{"id": 4, "name": "sample artist" }]
						}]
					}`))
					require.NoError(t, err)
				} else {
					_, err := w.Write([]byte(`{"result": []}`))
					require.NoError(t, err)
				}
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(WithAPIURL(apiServerMock.URL))

			result, err := client.GetTrack(tt.trackID)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestClient_GetAlbum(t *testing.T) {
	tests := []struct {
		name    string
		albumID string
		want    *Album
	}{
		{
			name:    "when track found",
			albumID: "foundID",
			want: &Album{
				ID:    1,
				Title: "Sample Title",
				Artists: []Artist{
					{
						ID:   2,
						Name: "Sample Artist",
					},
				},
			},
		},
		{
			name:    "when track not found",
			albumID: "notFoundID",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)

				if r.URL.Path == "/albums/foundID" {
					_, err := w.Write([]byte(`{
						"result": {
							"id": 1,
							"title": "Sample Title",
							"artists": [{"id": 2, "name": "Sample Artist" }]
						}
					}`))
					require.NoError(t, err)
				} else {
					_, err := w.Write([]byte(`{
						"result": {
							"error": "not-found"
						}
					}`))
					require.NoError(t, err)
				}
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(WithAPIURL(apiServerMock.URL))

			result, err := client.GetAlbum(tt.albumID)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestClient_SearchTrack(t *testing.T) {
	tests := []struct {
		name        string
		queryArtist string
		queryTrack  string
		want        *Track
	}{
		{
			name:        "when track found",
			queryArtist: "Found artist",
			queryTrack:  "Found track",
			want: &Track{
				ID:    "1",
				Title: "sample title",
				Albums: []Album{
					{
						ID:    2,
						Title: "sample title",
						Artists: []Artist{
							{
								ID:   3,
								Name: "sample artist",
							},
						},
					},
				},
				Artists: []Artist{
					{
						ID:   4,
						Name: "sample artist",
					},
				},
			},
		},
		{
			name:        "when track not found",
			queryArtist: "any impossible artist",
			queryTrack:  "any impossible track",
			want:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)

				searchType := r.URL.Query().Get("type")
				require.Equal(t, "track", searchType)

				query := r.URL.Query().Get("text")
				if query == "Found artist – Found track" {
					_, err := w.Write([]byte(`{
						"result": {
							"tracks":{
								"results": [{
									"id": "1",
									"title": "sample title",
									"albums": [
										{
											"id": 2,
											"title": "sample title",
											"artists": [{"id": 3, "name": "sample artist"}]
										}
									],
									"artists": [{"id": 4, "name": "sample artist" }]
								}]
							}
						}
					}`))
					require.NoError(t, err)
				} else {
					_, err := w.Write([]byte(`{"result": {}}`))
					require.NoError(t, err)
				}
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(WithAPIURL(apiServerMock.URL))

			result, err := client.SearchTrack(tt.queryArtist, tt.queryTrack)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestClient_SearchAlbum(t *testing.T) {
	tests := []struct {
		name        string
		queryArtist string
		queryAlbum  string
		want        *Album
	}{
		{
			name:        "when track found",
			queryArtist: "Found artist",
			queryAlbum:  "Found album",
			want: &Album{
				ID:    1,
				Title: "Sample Title",
				Artists: []Artist{
					{
						ID:   2,
						Name: "Sample Artist",
					},
				},
			},
		},
		{
			name:        "when track not found",
			queryArtist: "any impossible artist",
			queryAlbum:  "any impossible album",
			want:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)

				searchType := r.URL.Query().Get("type")
				require.Equal(t, "album", searchType)

				query := r.URL.Query().Get("text")
				if query == "Found artist – Found album" {
					_, err := w.Write([]byte(`{
						"result": {
							"albums":{
								"results": [{
									"id": 1,
									"title": "Sample Title",
									"artists": [{"id": 2, "name": "Sample Artist" }]
								}]
							}
						}
					}`))
					require.NoError(t, err)
				} else {
					_, err := w.Write([]byte(`{"result": {}}`))
					require.NoError(t, err)
				}
			}))
			defer apiServerMock.Close()

			client := NewHTTPClient(WithAPIURL(apiServerMock.URL))

			result, err := client.SearchAlbum(tt.queryArtist, tt.queryAlbum)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}
