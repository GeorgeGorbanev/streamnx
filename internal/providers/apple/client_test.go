package apple

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_fetchTrack(t *testing.T) {
	tests := []struct {
		name       string
		trackID    string
		storeFront string
		want       entity
		wantErr    error
	}{
		{
			name:       "when track found",
			trackID:    "foundId",
			storeFront: "us",
			want: entity{
				ID: "foundID",
				Attributes: entityAttributes{
					ArtistName: "sampleArtistName",
					Name:       "sampleTrackName",
					URL:        "sampleURL",
					EditorialNotes: editorialNotes{
						Standard: "sample track editorial notes",
					},
				},
			},
		},
		{
			name:       "when track not found",
			trackID:    "notFoundId",
			storeFront: "nevermind",
			wantErr:    errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var origin string

			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, origin, r.Header.Get("Origin"))

				switch r.URL.Path {
				case "/v1/catalog/us/songs/foundId":
					_, err := w.Write([]byte(`{
					"data":[
						{
							"id":"foundID",
							"attributes": {
								"artistName": "sampleArtistName",
								"editorialNotes": {
									"standard": "sample track editorial notes"
								},
								"name": "sampleTrackName",
								"url": "sampleURL"
							}
						}
					]
				}`))
					require.NoError(t, err)
				case "/v1/catalog/nevermind/songs/notFoundId":
					w.WriteHeader(http.StatusNotFound)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			origin = apiServerMock.URL
			defer apiServerMock.Close()

			client := Client{
				apiURL:       apiServerMock.URL,
				webPlayerURL: apiServerMock.URL,
				token:        "tokenMock",
				httpClient:   &http.Client{},
			}

			result, err := client.fetchTrack(t.Context(), tt.trackID, tt.storeFront)
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
		name    string
		artist  string
		title   string
		want    []entity
		wantErr error
	}{
		{
			name:   "when track candidates found",
			artist: "foundArtistName",
			title:  "foundTrackName",
			want: []entity{
				{
					ID: "firstSongID",
					Attributes: entityAttributes{
						ArtistName: "sampleArtistName",
						Name:       "sampleTrackName",
						URL:        "sampleURL",
						EditorialNotes: editorialNotes{
							Standard: "sample track editorial notes",
						},
					},
				},
				{
					ID: "secondSongID",
					Attributes: entityAttributes{
						ArtistName: "anotherArtistName",
						Name:       "anotherTrackName",
						URL:        "anotherURL",
					},
				},
			},
		},
		{
			name:   "when track candidates are empty",
			artist: "notFoundArtistName",
			title:  "notFoundTrackName",
			want:   []entity{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var origin string

			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, origin, r.Header.Get("Origin"))
				require.Equal(t, "/v1/catalog/us/search", r.URL.Path)

				q := r.URL.Query()
				require.Equal(t, "c", q.Get("art[music-videos:url]"))
				require.Equal(t, "f", q.Get("art[url]"))
				require.Equal(t, "artistUrl", q.Get("extend"))
				require.Equal(t, "url,name,artwork", q.Get("fields[artists]"))
				require.Equal(t, "map", q.Get("format[resources]"))
				require.Equal(t, "artists", q.Get("include[albums]"))
				require.Equal(t, "artists", q.Get("include[music-videos]"))
				require.Equal(t, "artists", q.Get("include[songs]"))
				require.Equal(t, "radio-show", q.Get("include[stations]"))
				require.Equal(t, "en-US", q.Get("l"))
				require.Equal(t, "21", q.Get("limit"))
				require.Equal(t, "autos", q.Get("omit[resource]"))
				require.Equal(t, "web", q.Get("platform"))
				require.Equal(t, "artists", q.Get("relate[albums]"))
				require.Equal(t, "albums", q.Get("relate[songs]"))
				require.Equal(t, "artistName,artistUrl,artwork,contentRating,editorialArtwork,editorialNotes,name,"+
					"playParams,releaseDate,url,trackCount,upc", q.Get("fields[albums]"))
				require.Equal(t, "activities,albums,apple-curators,artists,curators,editorial-items,music-movies,"+
					"music-videos,playlists,record-labels,songs,stations,tv-episodes,uploaded-videos", q.Get("types"))
				require.Equal(t, "lyricHighlights,lyrics,serverBubbles", q.Get("with"))

				var resp string
				if q.Get("term") == "foundArtistName foundTrackName" {
					resp = `{
						"results": {
							"top": {
								"data": [
									{
										"id": "unrelatedAlbumId",
										"type": "albums"
									},
									{
										"id": "firstSongId",
										"type": "songs"
									},
									{
										"id": "secondSongId",
										"type": "songs"
									}
								]		
	 	                     }
						},	
						"resources": {
							"songs": {
								"firstSongId": {
									"id":"firstSongID",
									"attributes": {
										"artistName": "sampleArtistName",
										"editorialNotes": {
											"standard": "sample track editorial notes"
										},
										"name": "sampleTrackName",
										"url": "sampleURL"
									}
								},
								"secondSongId": {
									"id":"secondSongID",
									"attributes": {
										"artistName": "anotherArtistName",
										"name": "anotherTrackName",
										"url": "anotherURL"
									}
								}
							},
							"albums": {
								"unrelatedAlbumId": {
									"id":"unrelatedAlbumID",
									"attributes": {
										"artistName": "albumArtistName",
										"name": "albumName",
										"url": "albumURL"
									}
								}
							}
						}
					}`
				} else {
					resp = `{
						"results": {
							"top": {
								"data": [
									{
										"id": "onlyAlbumId",
										"type": "albums"
									}
								]
							}
						},
						"resources": {
							"albums": {
								"onlyAlbumId": {
									"id":"onlyAlbumID",
									"attributes": {
										"artistName": "albumArtistName",
										"name": "albumName",
										"url": "albumURL"
									}
								}
							}
						}
					}`
				}
				_, err := w.Write([]byte(resp))
				require.NoError(t, err)
			}))
			origin = apiServerMock.URL
			defer apiServerMock.Close()

			client := Client{
				apiURL:       apiServerMock.URL,
				webPlayerURL: apiServerMock.URL,
				token:        "tokenMock",
				httpClient:   &http.Client{},
			}

			result, err := client.searchTracks(t.Context(), tt.artist, tt.title)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_fetchTracksByISRC(t *testing.T) {
	var origin string
	apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
		require.Equal(t, origin, r.Header.Get("Origin"))
		require.Equal(t, "/v1/catalog/us/songs", r.URL.Path)
		require.Equal(t, "GBARL9300135", r.URL.Query().Get("filter[isrc]"))
		require.Equal(t, "albums", r.URL.Query().Get("include"))
		_, err := w.Write([]byte(`{
			"data": [{
				"id": "1559885421",
				"attributes": {
					"artistName": "Rick Astley",
					"isrc": "GBARL9300135",
					"name": "Never Gonna Give You Up"
				}
			}]
		}`))
		require.NoError(t, err)
	}))
	origin = apiServerMock.URL
	defer apiServerMock.Close()

	client := Client{
		apiURL:       apiServerMock.URL,
		webPlayerURL: apiServerMock.URL,
		token:        "tokenMock",
		httpClient:   &http.Client{},
	}

	result, err := client.fetchTracksByISRC(t.Context(), "GBARL9300135")

	require.NoError(t, err)
	require.Equal(t, []entity{{
		ID: "1559885421",
		Attributes: entityAttributes{
			ArtistName: "Rick Astley",
			ISRC:       "GBARL9300135",
			Name:       "Never Gonna Give You Up",
		},
	}}, result)
}

func TestClient_fetchAlbumsByUPC(t *testing.T) {
	var origin string
	apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
		require.Equal(t, origin, r.Header.Get("Origin"))
		require.Equal(t, "/v1/catalog/us/albums", r.URL.Path)
		require.Equal(t, "196006422677", r.URL.Query().Get("filter[upc]"))
		require.Equal(t, "tracks", r.URL.Query().Get("include"))
		_, err := w.Write([]byte(`{
			"data": [{
				"id": "1577580468",
				"attributes": {
					"artistName": "გოგი ცაბაძე",
					"name": "კინომუსიკა: გოგი ცაბაძის შემოქმედება ნაწილი XIII",
					"upc": "196006422677"
				}
			}]
		}`))
		require.NoError(t, err)
	}))
	origin = apiServerMock.URL
	defer apiServerMock.Close()

	client := Client{
		apiURL:       apiServerMock.URL,
		webPlayerURL: apiServerMock.URL,
		token:        "tokenMock",
		httpClient:   &http.Client{},
	}

	result, err := client.fetchAlbumsByUPC(t.Context(), "196006422677")

	require.NoError(t, err)
	require.Equal(t, []entity{{
		ID: "1577580468",
		Attributes: entityAttributes{
			ArtistName: "გოგი ცაბაძე",
			Name:       "კინომუსიკა: გოგი ცაბაძის შემოქმედება ნაწილი XIII",
			UPC:        "196006422677",
		},
	}}, result)
}

func TestClient_fetchAlbum(t *testing.T) {
	tests := []struct {
		name       string
		albumID    string
		storeFront string
		want       entity
		wantErr    error
	}{
		{
			name:       "when album found",
			albumID:    "foundId",
			storeFront: "us",
			want: entity{
				ID: "foundID",
				Attributes: entityAttributes{
					ArtistName: "sampleArtistName",
					Name:       "sampleAlbumName",
					URL:        "sampleURL",
					EditorialNotes: editorialNotes{
						Standard: "sample album editorial notes",
					},
				},
				Relationships: entityRelationships{
					Tracks: tracksRelationship{
						Data: []entity{
							{
								ID: "firstTrackID",
								Attributes: entityAttributes{
									ArtistName: "firstTrackArtistName",
									Name:       "firstTrackName",
									URL:        "https://music.apple.com/us/album/sampleAlbumName/foundID?i=firstTrackID",
								},
							},
							{
								ID: "secondTrackID",
								Attributes: entityAttributes{
									ArtistName: "secondTrackArtistName",
									Name:       "secondTrackName",
									URL:        "https://music.apple.com/us/song/secondTrackName/secondTrackID",
									EditorialNotes: editorialNotes{
										Short: "second track notes",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "when album not found",
			albumID:    "notFoundId",
			storeFront: "nevermind",
			wantErr:    errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var origin string

			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, origin, r.Header.Get("Origin"))

				switch r.URL.Path {
				case "/v1/catalog/us/albums/foundId":
					_, err := w.Write([]byte(`{
					"data":[
						{
							"id":"foundID",
							"attributes": {
								"artistName": "sampleArtistName",
								"editorialNotes": {
									"standard": "sample album editorial notes"
								},
								"name": "sampleAlbumName",
								"url": "sampleURL"
							},
							"relationships": {
								"tracks": {
									"data": [
										{
											"id": "firstTrackID",
											"type": "songs",
											"attributes": {
												"artistName": "firstTrackArtistName",
												"name": "firstTrackName",
												"url": "https://music.apple.com/us/album/sampleAlbumName/foundID?i=firstTrackID"
											}
										},
										{
											"id": "secondTrackID",
											"type": "songs",
											"attributes": {
												"artistName": "secondTrackArtistName",
												"editorialNotes": {
													"short": "second track notes"
												},
												"name": "secondTrackName",
												"url": "https://music.apple.com/us/song/secondTrackName/secondTrackID"
											}
										}
									]
								}
							}
						}
					]
				}`))
					require.NoError(t, err)
				case "/v1/catalog/nevermind/albums/notFoundId":
					w.WriteHeader(http.StatusNotFound)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			origin = apiServerMock.URL
			defer apiServerMock.Close()

			client := Client{
				apiURL:       apiServerMock.URL,
				webPlayerURL: apiServerMock.URL,
				token:        "tokenMock",
				httpClient:   &http.Client{},
			}

			result, err := client.fetchAlbum(t.Context(), tt.albumID, tt.storeFront)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_FetchRejectsEmptyData(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		fetch func(t *testing.T, client *Client) (entity, error)
	}{
		{
			name: "track",
			path: "/v1/catalog/us/songs/emptyDataID",
			fetch: func(t *testing.T, client *Client) (entity, error) {
				return client.fetchTrack(t.Context(), "emptyDataID", "us")
			},
		},
		{
			name: "album",
			path: "/v1/catalog/us/albums/emptyDataID",
			fetch: func(t *testing.T, client *Client) (entity, error) {
				return client.fetchAlbum(t.Context(), "emptyDataID", "us")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.path, r.URL.Path)
				_, err := w.Write([]byte(`{"data":[]}`))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := Client{
				apiURL:       apiServerMock.URL,
				webPlayerURL: apiServerMock.URL,
				token:        "tokenMock",
				httpClient:   &http.Client{},
			}

			result, err := tt.fetch(t, &client)

			require.Zero(t, result)
			require.ErrorContains(t, err, "unexpected response: empty data")
		})
	}
}

func TestClient_FetchRejectsUnexpectedStatus(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		fetch func(t *testing.T, client *Client) (entity, error)
	}{
		{
			name: "track",
			path: "/v1/catalog/us/songs/serverErrorID",
			fetch: func(t *testing.T, client *Client) (entity, error) {
				return client.fetchTrack(t.Context(), "serverErrorID", "us")
			},
		},
		{
			name: "album",
			path: "/v1/catalog/us/albums/serverErrorID",
			fetch: func(t *testing.T, client *Client) (entity, error) {
				return client.fetchAlbum(t.Context(), "serverErrorID", "us")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.path, r.URL.Path)
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"data":[{"id":"should-not-decode"}]}`))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := Client{
				apiURL:       apiServerMock.URL,
				webPlayerURL: apiServerMock.URL,
				token:        "tokenMock",
				httpClient:   &http.Client{},
			}

			result, err := tt.fetch(t, &client)

			require.Zero(t, result)
			require.ErrorContains(t, err, "unexpected status code: 500")
		})
	}
}

func TestClient_searchAlbums(t *testing.T) {
	tests := []struct {
		name       string
		artistName string
		title      string
		want       []entity
		wantErr    error
	}{
		{
			name:       "when album candidates found",
			artistName: "foundArtistName",
			title:      "foundAlbumName",
			want: []entity{
				{
					ID: "firstAlbumID",
					Attributes: entityAttributes{
						ArtistName: "sampleArtistName",
						Name:       "sampleAlbumName",
						URL:        "sampleURL",
						EditorialNotes: editorialNotes{
							Standard: "sample album editorial notes",
						},
					},
				},
				{
					ID: "secondAlbumID",
					Attributes: entityAttributes{
						ArtistName: "anotherArtistName",
						Name:       "anotherAlbumName",
						URL:        "anotherURL",
					},
				},
			},
		},
		{
			name:       "when album candidates are empty",
			artistName: "notFoundArtistName",
			title:      "notFoundAlbumName",
			want:       []entity{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var origin string

			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, origin, r.Header.Get("Origin"))
				require.Equal(t, "/v1/catalog/us/search", r.URL.Path)

				q := r.URL.Query()
				require.Equal(t, "c", q.Get("art[music-videos:url]"))
				require.Equal(t, "f", q.Get("art[url]"))
				require.Equal(t, "artistUrl", q.Get("extend"))
				require.Equal(t, "url,name,artwork", q.Get("fields[artists]"))
				require.Equal(t, "map", q.Get("format[resources]"))
				require.Equal(t, "artists", q.Get("include[albums]"))
				require.Equal(t, "artists", q.Get("include[music-videos]"))
				require.Equal(t, "artists", q.Get("include[songs]"))
				require.Equal(t, "radio-show", q.Get("include[stations]"))
				require.Equal(t, "en-US", q.Get("l"))
				require.Equal(t, "21", q.Get("limit"))
				require.Equal(t, "autos", q.Get("omit[resource]"))
				require.Equal(t, "web", q.Get("platform"))
				require.Equal(t, "artists", q.Get("relate[albums]"))
				require.Equal(t, "albums", q.Get("relate[songs]"))
				require.Equal(t, "artistName,artistUrl,artwork,contentRating,editorialArtwork,editorialNotes,name,"+
					"playParams,releaseDate,url,trackCount,upc", q.Get("fields[albums]"))
				require.Equal(t, "activities,albums,apple-curators,artists,curators,editorial-items,music-movies,"+
					"music-videos,playlists,record-labels,songs,stations,tv-episodes,uploaded-videos", q.Get("types"))
				require.Equal(t, "lyricHighlights,lyrics,serverBubbles", q.Get("with"))

				var resp string
				if q.Get("term") == "foundArtistName foundAlbumName" {
					resp = `{
						"results": {
							"top": {
								"data": [
									{
										"id": "unrelatedSongId",
										"type": "songs"
									},
									{
										"id": "firstAlbumId",
										"type": "albums"
									},
									{
										"id": "secondAlbumId",
										"type": "albums"
									}
								]		
	 	                     }
						},	
						"resources": {
							"albums": {
								"firstAlbumId": {
									"id":"firstAlbumID",
									"attributes": {
										"artistName": "sampleArtistName",
										"editorialNotes": {
											"standard": "sample album editorial notes"
										},
										"name": "sampleAlbumName",
										"url": "sampleURL"
									}
								},
								"secondAlbumId": {
									"id":"secondAlbumID",
									"attributes": {
										"artistName": "anotherArtistName",
										"name": "anotherAlbumName",
										"url": "anotherURL"
									}
								}
							},
							"songs": {
								"unrelatedSongId": {
									"id":"unrelatedSongID",
									"attributes": {
										"artistName": "songArtistName",
										"name": "songName",
										"url": "songURL"
									}
								}
							}
						}
					}`
				} else {
					resp = `{
						"results": {
							"top": {
								"data": [
									{
										"id": "onlySongId",
										"type": "songs"
									}
								]
							}
						},
						"resources": {
							"songs": {
								"onlySongId": {
									"id":"onlySongID",
									"attributes": {
										"artistName": "songArtistName",
										"name": "songName",
										"url": "songURL"
									}
								}
							}
						}
					}`
				}
				_, err := w.Write([]byte(resp))
				require.NoError(t, err)
			}))
			origin = apiServerMock.URL
			defer apiServerMock.Close()

			client := Client{
				apiURL:       apiServerMock.URL,
				webPlayerURL: apiServerMock.URL,
				token:        "tokenMock",
				httpClient:   &http.Client{},
			}

			result, err := client.searchAlbums(t.Context(), tt.artistName, tt.title)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_SearchRejectsUnexpectedStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		search     func(t *testing.T, client *Client) ([]entity, error)
	}{
		{
			name:       "track search 404",
			statusCode: http.StatusNotFound,
			search: func(t *testing.T, client *Client) ([]entity, error) {
				return client.searchTracks(t.Context(), "sample artist", "sample track")
			},
		},
		{
			name:       "track search 502",
			statusCode: http.StatusBadGateway,
			search: func(t *testing.T, client *Client) ([]entity, error) {
				return client.searchTracks(t.Context(), "sample artist", "sample track")
			},
		},
		{
			name:       "album search 404",
			statusCode: http.StatusNotFound,
			search: func(t *testing.T, client *Client) ([]entity, error) {
				return client.searchAlbums(t.Context(), "sample artist", "sample album")
			},
		},
		{
			name:       "album search 502",
			statusCode: http.StatusBadGateway,
			search: func(t *testing.T, client *Client) ([]entity, error) {
				return client.searchAlbums(t.Context(), "sample artist", "sample album")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/v1/catalog/us/search", r.URL.Path)
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(`{
					"results": {
						"top": {
							"data": []
						}
					}
				}`))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := Client{
				apiURL:       apiServerMock.URL,
				webPlayerURL: apiServerMock.URL,
				token:        "tokenMock",
				httpClient:   &http.Client{},
			}

			result, err := tt.search(t, &client)

			require.Nil(t, result)
			require.ErrorContains(t, err, "unexpected status code: "+strconv.Itoa(tt.statusCode))
		})
	}
}

func TestClient_fetchToken(t *testing.T) {
	webPlayerServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, err := w.Write([]byte(`
				<!DOCTYPE html>
				<html>
					<head><script src="/assets/index-Samp1eBund13.js"></script></head>
				</html>
			`))
			require.NoError(t, err)
		case "/assets/index-Samp1eBund13.js":
			_, err := w.Write([]byte(`
				tokenVar = "sampleToken"
				headers.Authorization = ` + "`Bearer ${tokenVar}`",
			))
			require.NoError(t, err)
		default:
			require.Fail(t, "unexpected path: %s", r.URL.Path)
		}
	}))
	defer webPlayerServerMock.Close()

	client := Client{
		httpClient:   &http.Client{},
		webPlayerURL: webPlayerServerMock.URL,
	}

	token, err := client.fetchToken(t.Context())

	require.NoError(t, err)
	require.Equal(t, "sampleToken", token)
}

func TestClient_ConcurrentTokenFetch(t *testing.T) {
	var tokenRequests int64
	var apiRequests int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			atomic.AddInt64(&tokenRequests, 1)
			time.Sleep(10 * time.Millisecond)
			_, err := w.Write([]byte(`
				<!DOCTYPE html>
				<html>
					<head><script src="/assets/index-Samp1eBund13.js"></script></head>
				</html>
			`))
			require.NoError(t, err)
		case "/assets/index-Samp1eBund13.js":
			_, err := w.Write([]byte(`
				tokenVar = "sampleToken"
				headers.Authorization = ` + "`Bearer ${tokenVar}`",
			))
			require.NoError(t, err)
		case "/v1/catalog/us/songs/sampleTrackID":
			atomic.AddInt64(&apiRequests, 1)
			require.Equal(t, "Bearer sampleToken", r.Header.Get("Authorization"))
			_, err := w.Write([]byte(`{"data":[{"id":"sampleTrackID"}]}`))
			require.NoError(t, err)
		default:
			require.Fail(t, "unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithAPIURL(server.URL),
		WithWebPlayerURL(server.URL),
	)

	const concurrency = 10
	errs := make(chan error, concurrency)
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, err := client.fetchTrack(t.Context(), "sampleTrackID", "us")
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
	require.Equal(t, int64(concurrency), atomic.LoadInt64(&apiRequests))
}

func TestClient_RefreshTokenAfterUnauthorized(t *testing.T) {
	var tokenRequests atomic.Int32
	var apiRequests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			tokenRequests.Add(1)
			_, err := w.Write([]byte(`<script src="/assets/index-hash.js"></script>`))
			require.NoError(t, err)
		case "/assets/index-hash.js":
			_, err := w.Write([]byte("token = \"fresh-token\"; headers.Authorization = `Bearer ${token}`"))
			require.NoError(t, err)
		case "/v1/catalog/us/songs/sampleTrackID":
			apiRequests.Add(1)
			if r.Header.Get("Authorization") == "Bearer expired-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			require.Equal(t, "Bearer fresh-token", r.Header.Get("Authorization"))
			_, err := w.Write([]byte(`{"data":[{"id":"sampleTrackID"}]}`))
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithAPIURL(server.URL),
		WithWebPlayerURL(server.URL),
	)
	client.token = "expired-token"

	result, err := client.fetchTrack(t.Context(), "sampleTrackID", "us")

	require.NoError(t, err)
	require.Equal(t, "sampleTrackID", result.ID)
	require.Equal(t, int32(1), tokenRequests.Load())
	require.Equal(t, int32(2), apiRequests.Load())
}

func TestClient_RetriesUnauthorizedOnlyOnce(t *testing.T) {
	var apiRequests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, err := w.Write([]byte(`<script src="/assets/index-hash.js"></script>`))
			require.NoError(t, err)
		case "/assets/index-hash.js":
			_, err := w.Write([]byte("token = \"fresh-token\"; headers.Authorization = `Bearer ${token}`"))
			require.NoError(t, err)
		case "/v1/catalog/us/songs/sampleTrackID":
			apiRequests.Add(1)
			w.WriteHeader(http.StatusUnauthorized)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(
		WithAPIURL(server.URL),
		WithWebPlayerURL(server.URL),
	)
	client.token = "expired-token"

	_, err := client.fetchTrack(t.Context(), "sampleTrackID", "us")

	require.ErrorContains(t, err, "unexpected status code: 401")
	require.Equal(t, int32(2), apiRequests.Load())
}

func TestClient_fetchTokenRejectsUnexpectedStatus(t *testing.T) {
	tests := []struct {
		name            string
		indexStatusCode int
		jsStatusCode    int
		wantErr         string
	}{
		{
			name:            "index page",
			indexStatusCode: http.StatusInternalServerError,
			jsStatusCode:    http.StatusOK,
			wantErr:         "failed to fetch player html page: 500",
		},
		{
			name:            "js bundle",
			indexStatusCode: http.StatusOK,
			jsStatusCode:    http.StatusBadGateway,
			wantErr:         "failed to fetch player js: 502",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webPlayerServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/":
					w.WriteHeader(tt.indexStatusCode)
					_, err := w.Write([]byte(`
						<!DOCTYPE html>
						<html>
							<head><script src="/assets/index-Samp1eBund13.js"></script></head>
						</html>
					`))
					require.NoError(t, err)
				case "/assets/index-Samp1eBund13.js":
					w.WriteHeader(tt.jsStatusCode)
					_, err := w.Write([]byte(`
						tokenVar = "sampleToken" 
						headers.Authorization = ` + "`Bearer ${tokenVar}`",
					))
					require.NoError(t, err)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			defer webPlayerServerMock.Close()

			client := Client{
				httpClient:   &http.Client{},
				webPlayerURL: webPlayerServerMock.URL,
			}

			token, err := client.fetchToken(t.Context())

			require.Empty(t, token)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestClient_parseBundleName(t *testing.T) {
	html := []byte(`
		<!DOCTYPE html>
		<html lang="en">
			<head><script type="module" crossorigin src="/assets/index~d23b7a84.js"></script></head>
			<body>sample body</body>
	 	</html>`)

	result := NewClient().parseBundleName(html)
	require.Equal(t, "/assets/index~d23b7a84.js", result)
}

func TestClient_parseTokenVar(t *testing.T) {
	jsBundle := []byte("sampleJs(); headers.Authorization = `Bearer ${token}`; anotherSampleJs();")
	result := NewClient().parseTokenVar(jsBundle)
	require.Equal(t, "token", result)
}

func TestClient_parseVariableValue(t *testing.T) {
	jsBundle := []byte(`sampleJs(); var var1 = "varValue1"; var var2 = "varValue2"; var var3 = "varValue3"`)
	result, err := NewClient().parseVariableValue(jsBundle, "var2")
	require.NoError(t, err)
	require.Equal(t, "varValue2", result)
}

func TestClient_parseToken(t *testing.T) {
	jsBundle := []byte(`token="sampleToken"; headers.Authorization = ` + "`Bearer ${token}`" + `; anotherSampleJs();`)
	result, err := NewClient().parseToken(jsBundle)
	require.NoError(t, err)
	require.Equal(t, "sampleToken", result)
}

func TestClient_searchQuery(t *testing.T) {
	result := NewClient().searchQuery("sample", "term")

	q, err := url.ParseQuery(result)
	require.NoError(t, err)
	require.Equal(t, "sample term", q.Get("term"))
	require.Equal(t, "c", q.Get("art[music-videos:url]"))
	require.Equal(t, "f", q.Get("art[url]"))
	require.Equal(t, "artistUrl", q.Get("extend"))
	require.Equal(t, "url,name,artwork", q.Get("fields[artists]"))
	require.Equal(t, "map", q.Get("format[resources]"))
	require.Equal(t, "artists", q.Get("include[albums]"))
	require.Equal(t, "artists", q.Get("include[music-videos]"))
	require.Equal(t, "artists", q.Get("include[songs]"))
	require.Equal(t, "radio-show", q.Get("include[stations]"))
	require.Equal(t, "en-US", q.Get("l"))
	require.Equal(t, "21", q.Get("limit"))
	require.Equal(t, "autos", q.Get("omit[resource]"))
	require.Equal(t, "web", q.Get("platform"))
	require.Equal(t, "artists", q.Get("relate[albums]"))
	require.Equal(t, "albums", q.Get("relate[songs]"))
	require.Equal(t, "artistName,artistUrl,artwork,contentRating,editorialArtwork,editorialNotes,name,"+
		"playParams,releaseDate,url,trackCount,upc", q.Get("fields[albums]"))
	require.Equal(t, "activities,albums,apple-curators,artists,curators,editorial-items,music-movies,"+
		"music-videos,playlists,record-labels,songs,stations,tv-episodes,uploaded-videos", q.Get("types"))
	require.Equal(t, "lyricHighlights,lyrics,serverBubbles", q.Get("with"))
}
