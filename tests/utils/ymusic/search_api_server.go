package ymusic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
	"github.com/stretchr/testify/require"
)

var (
	SampleSearchQuery        = "Massive Attack – Angel"
	SampleSearchResponseJSON = `{
		"invocationInfo": {
			"exec-duration-millis": "123",
			"hostname": "music-stable-back-klg-79.klg.yp-c.yandex.net",
			"req-id": "1703854680519912-9917339348184729128"
		},
		"result": {
			"page": 0,
			"perPage": 20,
			"searchRequestId": "music-stable-back-klg-79.klg.yp-c.yandex.net-1703854680519912-9917339348184729128-1703854680540",
			"text": "massive attack – angel",
			"tracks": {
				"order": 0,
				"perPage": 20,
				"results": [
					{
						"albums": [
							{
								"artists": [
									{
										"composer": false,
										"cover": {
											"prefix": "87cf11b2.p.12662/",
											"type": "from-artist-photos",
											"uri": "avatars.yandex.net/get-music-content/103235/87cf11b2.p.12662/%%"
										},
										"disclaimers": [],
										"genres": [],
										"id": 12662,
										"name": "Massive Attack",
										"various": false
									}
								],
								"available": true,
								"availableForMobile": true,
								"availableForOptions": [
									"bookmate"
								],
								"availableForPremiumUsers": true,
								"availablePartially": false,
								"availableRegions": [],
								"bests": [
									354093,
									354094,
									354089
								],
								"coverUri": "avatars.yandex.net/get-music-content/28589/b060361b.a.35627-1/%%",
								"disclaimers": [],
								"genre": "triphopgenre",
								"id": 35627,
								"labels": [
									"Virgin"
								],
								"likesCount": 9931,
								"metaType": "music",
								"ogImage": "avatars.yandex.net/get-music-content/28589/b060361b.a.35627-1/%%",
								"recent": false,
								"regions": [
									"BELARUS",
									"BELARUS_PREMIUM"
								],
								"releaseDate": "1998-01-01T00:00:00+03:00",
								"storageDir": "b060361b.a.35627-1",
								"title": "Mezzanine",
								"trackCount": 11,
								"trackPosition": {
									"index": 1,
									"volume": 1
								},
								"type": "single",
								"veryImportant": false,
								"year": 1998
							}
						],
						"artists": [
							{
								"composer": false,
								"cover": {
									"prefix": "87cf11b2.p.12662/",
									"type": "from-artist-photos",
									"uri": "avatars.yandex.net/get-music-content/103235/87cf11b2.p.12662/%%"
								},
								"disclaimers": [],
								"genres": [],
								"id": 12662,
								"name": "Massive Attack",
								"various": false
							}
						],
						"available": true,
						"availableAsRbt": true,
						"availableForOptions": [
							"bookmate"
						],
						"availableForPremiumUsers": true,
						"availableFullWithoutPermission": false,
						"coverUri": "avatars.yandex.net/get-music-content/28589/b060361b.a.35627-1/%%",
						"derivedColors": {
							"accent": "#999999",
							"average": "#666666",
							"miniPlayer": "#B7B7B7",
							"waveText": "#CCCCCC"
						},
						"disclaimers": [],
						"durationMs": 379500,
						"explicit": false,
						"fade": {
							"inStart": 3.4,
							"inStop": 10.1,
							"outStart": 376.9,
							"outStop": 377.1
						},
						"fileSize": 0,
						"id": 354093,
						"lyricsAvailable": true,
						"lyricsInfo": {
							"hasAvailableSyncLyrics": true,
							"hasAvailableTextLyrics": true
						},
						"major": {
							"id": 1,
							"name": "UNIVERSAL_MUSIC"
						},
						"ogImage": "avatars.yandex.net/get-music-content/28589/b060361b.a.35627-1/%%",
						"previewDurationMs": 30000,
						"r128": {
							"i": -11.59,
							"tp": 0.32
						},
						"realId": "354093",
						"regions": [
							"BELARUS",
							"BELARUS_PREMIUM"
						],
						"rememberPosition": false,
						"storageDir": "",
						"title": "Angel",
						"trackSharingFlag": "VIDEO_ALLOWED",
						"trackSource": "OWN",
						"type": "music"
					}                
				],
				"total": 76
			},
			"type": "track"
		}
	}`
	SampleSearchResponse = ymusic.SearchResponse{
		Result: ymusic.SearchResult{
			Tracks: ymusic.TracksSection{
				Results: []ymusic.Track{
					{
						ID:    354093.0,
						Title: "Angel",
						Albums: []ymusic.Album{
							{
								ID: 35627,
							},
						},
						Artists: []ymusic.Artist{
							{
								ID:   12662,
								Name: "Massive Attack",
							},
						},
					},
				},
			},
		},
	}
	EmptySearchJSON = `
		{
			"invocationInfo": {
				"hostname": "music-stable-back-sas-108.sas.yp-c.yandex.net",
				"req-id": "1703860464774050-4133759126158316578",
				"exec-duration-millis": "66"
			},
			"result": {
				"type": "track",
				"page": 0,
				"perPage": 20,
				"text": "SampleImpossibleQuery",
				"searchRequestId": "music-stable-back-sas-108.sas.yp-c.yandex.net-1703860464774050-4133759126158316578-1703860464783"
			}
		}
	`
	EmptySearchResponse = ymusic.SearchResponse{
		Result: ymusic.SearchResult{
			Tracks: ymusic.TracksSection{},
		},
	}
)

func NewAPISearchServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, r.URL.Path, "/search")
		require.Equal(t, r.URL.Query().Get("page"), "0")
		require.Equal(t, r.URL.Query().Get("type"), "track")

		if r.URL.Query().Get("text") == SampleSearchQuery {
			_, err := w.Write([]byte(SampleSearchResponseJSON))
			require.NoError(t, err)
		} else {
			_, err := w.Write([]byte(EmptySearchJSON))
			require.NoError(t, err)
		}
	}))
}
