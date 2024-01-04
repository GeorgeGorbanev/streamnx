package ymusic

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
	"github.com/stretchr/testify/require"
)

var (
	SampleTrackID              = "1197793"
	SampleGetTrackResponseJSON = `{
		"invocationInfo": {
			"exec-duration-millis": 11,
			"hostname": "music-stable-back-vla-24.vla.yp-c.yandex.net",
			"req-id": "1704222836102466-11015214182353549756"
		},
		"result": [
			{
				"albums": [
					{
						"artists": [
							{
								"composer": false,
								"cover": {
									"prefix": "1cfff9fb.p.67307/",
									"type": "from-artist-photos",
									"uri": "avatars.yandex.net/get-music-content/63210/1cfff9fb.p.67307/%%"
								},
								"disclaimers": [],
								"genres": [],
								"id": 67307,
								"name": "The Raveonettes",
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
						"bests": [
							1197793
						],
						"coverUri": "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
						"disclaimers": [],
						"genre": "indie",
						"id": 3192570,
						"labels": [],
						"likesCount": 30,
						"metaType": "music",
						"ogImage": "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
						"recent": false,
						"releaseDate": "2008-11-25T00:00:00+03:00",
						"title": "Wishing You A Rave Christmas",
						"trackCount": 4,
						"trackPosition": {
							"index": 2,
							"volume": 1
						},
						"veryImportant": false,
						"year": 2008
					}
				],
				"artists": [
					{
						"composer": false,
						"cover": {
							"prefix": "1cfff9fb.p.67307/",
							"type": "from-artist-photos",
							"uri": "avatars.yandex.net/get-music-content/63210/1cfff9fb.p.67307/%%"
						},
						"disclaimers": [],
						"genres": [],
						"id": 67307,
						"name": "The Raveonettes",
						"various": false
					}
				],
				"available": true,
				"availableForOptions": [
					"bookmate"
				],
				"availableForPremiumUsers": true,
				"availableFullWithoutPermission": false,
				"coverUri": "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
				"derivedColors": {
					"accent": "#6C8FC5",
					"average": "#121D2E",
					"miniPlayer": "#98B0D6",
					"waveText": "#CCCCCC"
				},
				"disclaimers": [],
				"durationMs": 221640,
				"fade": {
					"inStart": 3,
					"inStop": 12,
					"outStart": 205.8,
					"outStop": 207.5
				},
				"fileSize": 0,
				"id": "1197793",
				"lyricsAvailable": false,
				"lyricsInfo": {
					"hasAvailableSyncLyrics": false,
					"hasAvailableTextLyrics": true
				},
				"major": {
					"id": 39,
					"name": "ORCHARD"
				},
				"ogImage": "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
				"previewDurationMs": 30000,
				"r128": {
					"i": -9.86,
					"tp": 0.72
				},
				"realId": "1197793",
				"rememberPosition": false,
				"storageDir": "",
				"title": "Come On Santa",
				"trackSharingFlag": "COVER_ONLY",
				"trackSource": "OWN",
				"type": "music"
			}
		]
	}`
	SampleGetTrackResponse = ymusic.TrackResponse{
		Result: []ymusic.Track{
			{
				ID:    "1197793",
				Title: "Come On Santa",
				Albums: []ymusic.Album{
					{
						ID: 3192570,
					},
				},
				Artists: []ymusic.Artist{
					{
						ID:   67307,
						Name: "The Raveonettes",
					},
				},
			},
		},
	}

	EmptyGetTrackResponseJSON = `{
		"invocationInfo": {
			"exec-duration-millis": 11,
			"hostname": "music-stable-back-vla-24.vla.yp-c.yandex.net",
			"req-id": "1704222836102466-11015214182353549756"
		},
		"result": []
	}`
)

func NewAPIGetTrackServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		path := strings.Split(r.URL.Path, "/")
		require.Len(t, path, 3)
		require.Equal(t, "", path[0])
		require.Equal(t, "tracks", path[1])

		if path[2] == SampleTrackID {
			_, err := w.Write([]byte(SampleGetTrackResponseJSON))
			require.NoError(t, err)
		} else {
			_, err := w.Write([]byte(EmptyGetTrackResponseJSON))
			require.NoError(t, err)
		}
	}))
}
