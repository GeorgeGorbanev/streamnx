package apple

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient_FetchTrack(t *testing.T) {
	tests := []struct {
		name       string
		trackID    string
		storeFront string
		want       *Entity
	}{
		{
			name:       "when track found",
			trackID:    "foundId",
			storeFront: "us",
			want: &Entity{
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
			trackID:    "notFoundId",
			storeFront: "nevermind",
			want:       nil,
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
				case "/v1/catalog/nevermind/songs/notFoundId":
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.FetchTrack(ctx, tt.trackID, tt.storeFront)
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
		want       *Entity
	}{
		{
			name:       "when track found",
			artistName: "foundArtistName",
			trackName:  "foundTrackName",
			want: &Entity{
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
					"playParams,releaseDate,url,trackCount", q.Get("fields[albums]"))
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
										"id": "foundId",
										"type": "songs"
									}
								]		
	 	                     }
						},	
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.SearchTrack(ctx, tt.artistName, tt.trackName)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestHTTPClient_FetchAlbum(t *testing.T) {
	tests := []struct {
		name       string
		albumID    string
		storeFront string
		want       *Entity
	}{
		{
			name:       "when album found",
			albumID:    "foundId",
			storeFront: "us",
			want: &Entity{
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
			albumID:    "notFoundId",
			storeFront: "nevermind",
			want:       nil,
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
				case "/v1/catalog/nevermind/albums/notFoundId":
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.FetchAlbum(ctx, tt.albumID, tt.storeFront)
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
		want       *Entity
	}{
		{
			name:       "when album found",
			artistName: "foundArtistName",
			albumName:  "foundAlbumName",
			want: &Entity{
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
					"playParams,releaseDate,url,trackCount", q.Get("fields[albums]"))
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
										"id": "foundId",
										"type": "albums"
									}
								]		
	 	                     }
						},	
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.SearchAlbum(ctx, tt.artistName, tt.albumName)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := client.fetchToken(ctx)

	require.NoError(t, err)
	require.Equal(t, "sampleToken", token)
}

func Test_searchQuery(t *testing.T) {
	sampleTerm := "sample term"
	result := searchQuery(sampleTerm)

	q, err := url.ParseQuery(result)
	require.NoError(t, err)
	require.Equal(t, sampleTerm, q.Get("term"))
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
		"playParams,releaseDate,url,trackCount", q.Get("fields[albums]"))
	require.Equal(t, "activities,albums,apple-curators,artists,curators,editorial-items,music-movies,"+
		"music-videos,playlists,record-labels,songs,stations,tv-episodes,uploaded-videos", q.Get("types"))
	require.Equal(t, "lyricHighlights,lyrics,serverBubbles", q.Get("with"))
}
