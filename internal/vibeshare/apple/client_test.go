package apple

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient_GetTrack(t *testing.T) {
	tests := []struct {
		name    string
		trackID string
		want    *MusicEntity
	}{
		{
			name:    "when track found",
			trackID: "foundId",
			want: &MusicEntity{
				ID: "foundID",
				Attributes: Attributes{
					ArtistName: "sampleArtistName",
					Name:       "sampleTrackName",
					URL:        "sampleURL",
				},
			},
		},
		{
			name:    "when track not found",
			trackID: "notFoundId",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, "https://music.apple.com", r.Header.Get("Origin"))

				switch r.URL.Path {
				case "/v1/catalog/us/songs/foundId":
					_, err := w.Write([]byte(`{
					"data":[
						{
							"id":"foundID",
							"attributes": {
								"artistName": "sampleArtistName",
								"name": "sampleTrackName",
								"url": "sampleURL"
							}
						}
					]
				}`))
					require.NoError(t, err)
				case "/v1/catalog/us/songs/notFoundId":
					w.WriteHeader(http.StatusNotFound)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			defer apiServerMock.Close()

			client := HTTPClient{
				apiURL:     apiServerMock.URL,
				token:      "tokenMock",
				httpClient: &http.Client{},
			}

			result, err := client.GetTrack(tt.trackID)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestHTTPClient_SearchTrack(t *testing.T) {
	tests := []struct {
		name       string
		artistName string
		trackName  string
		want       *MusicEntity
	}{
		{
			name:       "when track found",
			artistName: "foundArtistName",
			trackName:  "foundTrackName",
			want: &MusicEntity{
				ID: "foundID",
				Attributes: Attributes{
					ArtistName: "sampleArtistName",
					Name:       "sampleTrackName",
					URL:        "sampleURL",
				},
			},
		},
		{
			name:       "when track not found",
			artistName: "notFoundArtistName",
			trackName:  "notFoundTrackName",
			want:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, "https://music.apple.com", r.Header.Get("Origin"))
				require.Equal(t, "/v1/catalog/us/search/suggestions", r.URL.Path)

				q := r.URL.Query()
				require.Equal(t, "f", q.Get("art[url]"))
				require.Equal(t, "artwork,name,playParams,url,artistName", q.Get("fields[albums]"))
				require.Equal(t, "url,name,artwork", q.Get("fields[artists]"))
				require.Equal(t, "map", q.Get("format[resources]"))
				require.Equal(t, "terms,topResults", q.Get("kinds"))
				require.Equal(t, "en-US", q.Get("l"))
				require.Equal(t, "5", q.Get("limit[results:terms]"))
				require.Equal(t, "10", q.Get("limit[results:topResults]"))
				require.Equal(t, "autos", q.Get("omit[resource]"))
				require.Equal(t, "web", q.Get("platform"))
				require.Equal(t, "activities,albums,artists,editorial-items,music-movies,music-videos,"+
					"playlists,record-labels,songs,stations,tv-episodes", q.Get("types"))

				var resp string
				if q.Get("term") == "foundArtistName foundTrackName" {
					resp = `{
						"resources": {
							"songs": {
								"foundId": {
									"id":"foundID",
									"attributes": {
										"artistName": "sampleArtistName",
										"name": "sampleTrackName",
										"url": "sampleURL"
									}
								}
							}
						}
					}`
				} else {
					resp = `{}`
				}
				_, err := w.Write([]byte(resp))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := HTTPClient{
				apiURL:     apiServerMock.URL,
				token:      "tokenMock",
				httpClient: &http.Client{},
			}

			result, err := client.SearchTrack(tt.artistName, tt.trackName)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestHTTPClient_GetAlbum(t *testing.T) {
	tests := []struct {
		name    string
		albumID string
		want    *MusicEntity
	}{
		{
			name:    "when album found",
			albumID: "foundId",
			want: &MusicEntity{
				ID: "foundID",
				Attributes: Attributes{
					ArtistName: "sampleArtistName",
					Name:       "sampleAlbumName",
					URL:        "sampleURL",
				},
			},
		},
		{
			name:    "when album not found",
			albumID: "notFoundId",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, "https://music.apple.com", r.Header.Get("Origin"))

				switch r.URL.Path {
				case "/v1/catalog/us/albums/foundId":
					_, err := w.Write([]byte(`{
					"data":[
						{
							"id":"foundID",
							"attributes": {
								"artistName": "sampleArtistName",
								"name": "sampleAlbumName",
								"url": "sampleURL"
							}
						}
					]
				}`))
					require.NoError(t, err)
				case "/v1/catalog/us/albums/notFoundId":
					w.WriteHeader(http.StatusNotFound)
				default:
					require.Fail(t, "unexpected path: %s", r.URL.Path)
				}
			}))
			defer apiServerMock.Close()

			client := HTTPClient{
				apiURL:     apiServerMock.URL,
				token:      "tokenMock",
				httpClient: &http.Client{},
			}

			result, err := client.GetAlbum(tt.albumID)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestHTTPClient_SearchAlbum(t *testing.T) {
	tests := []struct {
		name       string
		artistName string
		albumName  string
		want       *MusicEntity
	}{
		{
			name:       "when album found",
			artistName: "foundArtistName",
			albumName:  "foundAlbumName",
			want: &MusicEntity{
				ID: "foundID",
				Attributes: Attributes{
					ArtistName: "sampleArtistName",
					Name:       "sampleAlbumName",
					URL:        "sampleURL",
				},
			},
		},
		{
			name:       "when album not found",
			artistName: "notFoundArtistName",
			albumName:  "notFoundAlbumName",
			want:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer tokenMock", r.Header.Get("Authorization"))
				require.Equal(t, "https://music.apple.com", r.Header.Get("Origin"))
				require.Equal(t, "/v1/catalog/us/search/suggestions", r.URL.Path)

				q := r.URL.Query()
				require.Equal(t, "f", q.Get("art[url]"))
				require.Equal(t, "artwork,name,playParams,url,artistName", q.Get("fields[albums]"))
				require.Equal(t, "url,name,artwork", q.Get("fields[artists]"))
				require.Equal(t, "map", q.Get("format[resources]"))
				require.Equal(t, "terms,topResults", q.Get("kinds"))
				require.Equal(t, "en-US", q.Get("l"))
				require.Equal(t, "5", q.Get("limit[results:terms]"))
				require.Equal(t, "10", q.Get("limit[results:topResults]"))
				require.Equal(t, "autos", q.Get("omit[resource]"))
				require.Equal(t, "web", q.Get("platform"))
				require.Equal(t, "activities,albums,artists,editorial-items,music-movies,music-videos,"+
					"playlists,record-labels,songs,stations,tv-episodes", q.Get("types"))

				var resp string
				if q.Get("term") == "foundArtistName foundAlbumName" {
					resp = `{
						"resources": {
							"albums": {
								"foundId": {
									"id":"foundID",
									"attributes": {
										"artistName": "sampleArtistName",
										"name": "sampleAlbumName",
										"url": "sampleURL"
									}
								}
							}
						}
					}`
				} else {
					resp = `{}`
				}
				_, err := w.Write([]byte(resp))
				require.NoError(t, err)
			}))
			defer apiServerMock.Close()

			client := HTTPClient{
				apiURL:     apiServerMock.URL,
				token:      "tokenMock",
				httpClient: &http.Client{},
			}

			result, err := client.SearchAlbum(tt.artistName, tt.albumName)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestHTTPClient_fetchToken(t *testing.T) {
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

	client := HTTPClient{
		httpClient:   &http.Client{},
		webPlayerURL: webPlayerServerMock.URL,
	}
	token, err := client.fetchToken()

	require.NoError(t, err)
	require.Equal(t, "sampleToken", token)
}

func Test_searchQuery(t *testing.T) {
	sampleTerm := "sample term"
	result := searchQuery(sampleTerm)

	resultQuery, err := url.ParseQuery(result)
	require.NoError(t, err)
	require.Equal(t, sampleTerm, resultQuery.Get("term"))
	require.Equal(t, "f", resultQuery.Get("art[url]"))
	require.Equal(t, "artwork,name,playParams,url,artistName", resultQuery.Get("fields[albums]"))
	require.Equal(t, "url,name,artwork", resultQuery.Get("fields[artists]"))
	require.Equal(t, "map", resultQuery.Get("format[resources]"))
	require.Equal(t, "terms,topResults", resultQuery.Get("kinds"))
	require.Equal(t, "en-US", resultQuery.Get("l"))
	require.Equal(t, "5", resultQuery.Get("limit[results:terms]"))
	require.Equal(t, "10", resultQuery.Get("limit[results:topResults]"))
	require.Equal(t, "autos", resultQuery.Get("omit[resource]"))
	require.Equal(t, "web", resultQuery.Get("platform"))
	require.Equal(t, "activities,albums,artists,editorial-items,music-movies,music-videos,"+
		"playlists,record-labels,songs,stations,tv-episodes", resultQuery.Get("types"))
}
