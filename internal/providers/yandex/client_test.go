package yandex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_fetchTrack(t *testing.T) {
	tests := []struct {
		name            string
		trackID         string
		want            track
		wantErr         error
		wantErrContains string
	}{
		{
			name:    "when track found",
			trackID: "foundID",
			want: track{
				ID:    "1",
				Title: "sample title",
				Albums: []albumRef{
					{
						ID:    2,
						Title: "sample title",
						Artists: []artist{
							{
								ID:   3,
								Name: "sample artist",
							},
						},
					},
				},
				Artists: []artist{
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
			wantErr: errNotFound,
		},
		{
			name:            "when http status is unexpected",
			trackID:         "serverError",
			wantErrContains: `unexpected api response status 500 ({"error":{"name":"internal-error","message":"Internal server error"}})`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)

				switch r.URL.Path {
				case "/tracks/foundID":
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
				case "/tracks/serverError":
					w.WriteHeader(http.StatusInternalServerError)
					_, err := w.Write([]byte(`{"error":{"name":"internal-error","message":"Internal server error"}}`))
					require.NoError(t, err)
				default:
					_, err := w.Write([]byte(`{"result": []}`))
					require.NoError(t, err)
				}
			}))
			defer apiServerMock.Close()

			client := NewClient(WithAPIURL(apiServerMock.URL))

			result, err := client.fetchTrack(t.Context(), tt.trackID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else if tt.wantErrContains != "" {
				require.Zero(t, result)
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_fetchAlbum(t *testing.T) {
	tests := []struct {
		name    string
		albumID string
		want    album
		wantErr error
	}{
		{
			name:    "when album found",
			albumID: "foundID",
			want: album{
				ID:    1,
				Title: "Sample Title",
				Artists: []artist{
					{
						ID:   2,
						Name: "Sample Artist",
					},
				},
				Volumes: [][]track{
					{
						{
							ID:    "11",
							Title: "First Track",
							Artists: []artist{
								{
									ID:   2,
									Name: "Sample Artist",
								},
							},
							Albums: []albumRef{
								{
									ID:    1,
									Title: "Sample Title",
								},
							},
						},
						{
							ID:    "12",
							Title: "Second Track",
							Artists: []artist{
								{
									ID:   2,
									Name: "Sample Artist",
								},
							},
							Albums: []albumRef{
								{
									ID:    1,
									Title: "Sample Title",
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "when album not found",
			albumID: "-1231231231",
			wantErr: errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)

				switch r.URL.Path {
				case "/albums/foundID/with-tracks":
					_, err := w.Write([]byte(`{
						"result": {
							"id": 1,
							"title": "Sample Title",
							"artists": [{"id": 2, "name": "Sample Artist" }],
							"volumes": [[
								{
									"id": "11",
									"title": "First Track",
									"artists": [{"id": 2, "name": "Sample Artist"}],
									"albums": [{"id": 1, "title": "Sample Title"}]
								},
								{
									"id": "12",
									"title": "Second Track",
									"artists": [{"id": 2, "name": "Sample Artist"}],
									"albums": [{"id": 1, "title": "Sample Title"}]
								}
							]]
						}
					}`))
					require.NoError(t, err)
				case "/albums/-1231231231/with-tracks":
					_, err := w.Write([]byte(`{
						"invocationInfo": {
							"req-id": "12312-31231",
							"hostname": "music-web-default-production-music-86.klg.yp-c.yandex.net"
						},
						"error": {
							"name": "not-found",
							"message": "Requested action for /albums/-1231231231/with-tracks was not found"
						}
					}`))
					require.NoError(t, err)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			defer apiServerMock.Close()

			client := NewClient(WithAPIURL(apiServerMock.URL))

			result, err := client.fetchAlbum(t.Context(), tt.albumID)
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
		name         string
		query        string
		responseBody string
		want         []searchTrack
	}{
		{
			name:  "multiple candidates",
			query: "multiple candidates query",
			responseBody: `{
				"result": {
					"tracks": {
						"results": [{
							"id": 1,
							"title": "sample title",
							"albums": [{
								"id": 2,
								"title": "sample album",
								"artists": [{"id": 3, "name": "sample album artist"}]
							}],
							"artists": [{"id": 4, "name": "sample artist"}]
						}, {
							"id": 5,
							"title": "second title",
							"albums": [{
								"id": 6,
								"title": "second album",
								"artists": [{"id": 7, "name": "second album artist"}]
							}],
							"artists": [{"id": 8, "name": "second artist"}]
						}]
					}
				}
			}`,
			want: []searchTrack{
				{
					ID:    1,
					Title: "sample title",
					Albums: []albumRef{
						{
							ID:    2,
							Title: "sample album",
							Artists: []artist{
								{
									ID:   3,
									Name: "sample album artist",
								},
							},
						},
					},
					Artists: []artist{
						{
							ID:   4,
							Name: "sample artist",
						},
					},
				},
				{
					ID:    5,
					Title: "second title",
					Albums: []albumRef{
						{
							ID:    6,
							Title: "second album",
							Artists: []artist{
								{
									ID:   7,
									Name: "second album artist",
								},
							},
						},
					},
					Artists: []artist{
						{
							ID:   8,
							Name: "second artist",
						},
					},
				},
			},
		},
		{
			name:         "empty candidates",
			query:        "empty query",
			responseBody: `{"result":{"tracks":{"results":[]}}}`,
			want:         []searchTrack{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "0", r.URL.Query().Get("page"))

				searchType := r.URL.Query().Get("type")
				require.Equal(t, "track", searchType)

				query := r.URL.Query().Get("text")
				require.Equal(t, tt.query, query)

				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(WithAPIURL(apiServerMock.URL))

			result, err := client.searchTracks(t.Context(), tt.query)

			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestClient_searchAlbums(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		responseBody string
		want         []searchAlbum
	}{
		{
			name:  "multiple candidates",
			query: "multiple candidates query",
			responseBody: `{
				"result": {
					"albums": {
						"results": [{
							"id": 1,
							"title": "Sample Title",
							"artists": [{"id": 2, "name": "Sample Artist"}]
						}, {
							"id": 3,
							"title": "Second Title",
							"artists": [{"id": 4, "name": "Second Artist"}]
						}]
					}
				}
			}`,
			want: []searchAlbum{
				{
					ID:    1,
					Title: "Sample Title",
					Artists: []artist{
						{
							ID:   2,
							Name: "Sample Artist",
						},
					},
				},
				{
					ID:    3,
					Title: "Second Title",
					Artists: []artist{
						{
							ID:   4,
							Name: "Second Artist",
						},
					},
				},
			},
		},
		{
			name:         "empty candidates",
			query:        "empty query",
			responseBody: `{"result":{"albums":{"results":[]}}}`,
			want:         []searchAlbum{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "0", r.URL.Query().Get("page"))

				searchType := r.URL.Query().Get("type")
				require.Equal(t, "album", searchType)

				query := r.URL.Query().Get("text")
				require.Equal(t, tt.query, query)

				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := NewClient(WithAPIURL(apiServerMock.URL))

			result, err := client.searchAlbums(t.Context(), tt.query)

			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}
