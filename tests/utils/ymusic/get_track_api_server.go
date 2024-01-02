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
		InvocationInfo: ymusic.InvocationInfo{
			ExecDurationMillis: "11",
			Hostname:           "music-stable-back-vla-24.vla.yp-c.yandex.net",
			ReqID:              "1704222836102466-11015214182353549756",
		},
		Result: []ymusic.Track{
			{
				Albums: []ymusic.Album{
					{
						Artists: []ymusic.Artist{
							{
								Composer: false,
								Cover: ymusic.ArtistCover{
									Prefix: "1cfff9fb.p.67307/",
									Type:   "from-artist-photos",
									URI:    "avatars.yandex.net/get-music-content/63210/1cfff9fb.p.67307/%%",
								},
								Disclaimers: []string{},
								Genres:      []string{},
								ID:          67307,
								Name:        "The Raveonettes",
								Various:     false,
							},
						},
						Available:                true,
						AvailableForMobile:       true,
						AvailableForOptions:      []string{"bookmate"},
						AvailableForPremiumUsers: true,
						AvailablePartially:       false,
						Bests:                    []int{1197793},
						CoverUri:                 "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
						Disclaimers:              []string{},
						Genre:                    "indie",
						ID:                       3192570,
						LikesCount:               30,
						MetaType:                 "music",
						OgImage:                  "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
						Recent:                   false,
						ReleaseDate:              "2008-11-25T00:00:00+03:00",
						Title:                    "Wishing You A Rave Christmas",
						TrackCount:               4,
						TrackPosition:            ymusic.TrackPosition{Index: 2, Volume: 1},
						VeryImportant:            false,
						Year:                     2008,
					},
				},
				Artists: []ymusic.Artist{
					{
						Composer: false,
						Cover: ymusic.ArtistCover{
							Prefix: "1cfff9fb.p.67307/",
							Type:   "from-artist-photos",
							URI:    "avatars.yandex.net/get-music-content/63210/1cfff9fb.p.67307/%%",
						},
						Disclaimers: []string{},
						Genres:      []string{},
						ID:          67307,
						Name:        "The Raveonettes",
						Various:     false,
					},
				},
				Available:                      true,
				AvailableForOptions:            []string{"bookmate"},
				AvailableForPremiumUsers:       true,
				AvailableFullWithoutPermission: false,
				CoverUri:                       "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
				DerivedColors: ymusic.DerivedColors{
					Accent:     "#6C8FC5",
					Average:    "#121D2E",
					MiniPlayer: "#98B0D6",
					WaveText:   "#CCCCCC",
				},
				Disclaimers: []any{},
				DurationMs:  221640,
				Fade: ymusic.Fade{
					InStart:  3,
					InStop:   12,
					OutStart: 205.8,
					OutStop:  207.5,
				},
				FileSize:        0,
				ID:              "1197793",
				LyricsAvailable: false,
				LyricsInfo: ymusic.LyricsInfo{
					HasAvailableSyncLyrics: false,
					HasAvailableTextLyrics: true,
				},
				Major: ymusic.Major{
					ID:   39,
					Name: "ORCHARD",
				},
				OgImage:           "avatars.yandex.net/get-music-content/32236/32877f96.a.3192570-1/%%",
				PreviewDurationMs: 30000,
				R128: ymusic.R128{
					I:  -9.86,
					TP: 0.72,
				},
				RealID:           "1197793",
				RememberPosition: false,
				StorageDir:       "",
				Title:            "Come On Santa",
				TrackSharingFlag: "COVER_ONLY",
				TrackSource:      "OWN",
				Type:             "music",
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
	EmptyGetTrackResponse = ymusic.TrackResponse{
		InvocationInfo: ymusic.InvocationInfo{
			ExecDurationMillis: "11",
			Hostname:           "music-stable-back-vla-24.vla.yp-c.yandex.net",
			ReqID:              "1704222836102466-11015214182353549756",
		},
		Result: []ymusic.Track{},
	}
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
