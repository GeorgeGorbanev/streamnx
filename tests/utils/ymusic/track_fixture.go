package ymusic

import (
	"encoding/json"
	"fmt"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
)

type TrackFixture struct {
	ID                string
	Track             *ymusic.Track
	GetResponse       string
	SearchQueryArtist string
	SearchQueryTrack  string
	SearchResponse    string
}

var (
	TrackFixtureMassiveAttackAngel = TrackFixture{
		ID: "354093",
		Track: &ymusic.Track{
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
		SearchQueryArtist: "Massive Attack",
		SearchQueryTrack:  "Angel",
		SearchResponse: `{
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
		}`,
		GetResponse: `{
			"invocationInfo": {
				"exec-duration-millis": 8,
				"hostname": "music-stable-back-vla-41.vla.yp-c.yandex.net",
				"req-id": "1704383317792389-4078481512684815117"
			},
			"result": [
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
							"bests": [],
							"contentWarning": "clean",
							"coverUri": "avatars.yandex.net/get-music-content/42108/96350b7f.a.49408-1/%%",
							"disclaimers": [],
							"genre": "triphopgenre",
							"id": 35627,
							"labels": [
								{
									"id": 14959,
									"name": "Virgin"
								}
							],
							"likesCount": 1797,
							"metaType": "music",
							"ogImage": "avatars.yandex.net/get-music-content/42108/96350b7f.a.49408-1/%%",
							"recent": false,
							"releaseDate": "1998-01-01T00:00:00+03:00",
							"title": "Singles Collection",
							"trackCount": 61,
							"trackPosition": {
								"index": 1,
								"volume": 10
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
					"availableForOptions": [
						"bookmate"
					],
					"availableForPremiumUsers": true,
					"availableFullWithoutPermission": false,
					"clipIds": [
						57091
					],
					"coverUri": "avatars.yandex.net/get-music-content/42108/96350b7f.a.49408-1/%%",
					"derivedColors": {
						"accent": "#9A9897",
						"average": "#787776",
						"miniPlayer": "#B8B7B6",
						"waveText": "#CCCCCC"
					},
					"disclaimers": [],
					"durationMs": 379500,
					"fade": {
						"inStart": 3.4,
						"inStop": 10.1,
						"outStart": 376.9,
						"outStop": 377.1
					},
					"fileSize": 0,
					"id": "354093",
					"lyricsAvailable": true,
					"lyricsInfo": {
						"hasAvailableSyncLyrics": true,
						"hasAvailableTextLyrics": true
					},
					"major": {
						"id": 1,
						"name": "UNIVERSAL_MUSIC"
					},
					"ogImage": "avatars.yandex.net/get-music-content/42108/96350b7f.a.49408-1/%%",
					"previewDurationMs": 30000,
					"r128": {
						"i": -11.59,
						"tp": 0.32
					},
					"realId": "354093",
					"rememberPosition": false,
					"storageDir": "",
					"title": "Angel",
					"trackSharingFlag": "VIDEO_ALLOWED",
					"trackSource": "OWN",
					"type": "music"
				}
			]
		}
		`,
	}
	TrackFixtureDJAmor20Flowers = TrackFixture{
		ID: "110791278",
		Track: &ymusic.Track{
			ID:    110791278,
			Title: "Flowers",
			Albums: []ymusic.Album{
				{
					ID: 24665034,
				},
			},
			Artists: []ymusic.Artist{
				{
					ID:   18375620,
					Name: "DJ Amor 2.0",
				},
			},
		},
		SearchQueryArtist: "Miley Cyrus",
		SearchQueryTrack:  "Flowers",
		SearchResponse: `{
			"invocationInfo": {
				"exec-duration-millis": "79",
				"hostname": "music-stable-back-vla-19.vla.yp-c.yandex.net",
				"req-id": "1704449354601545-12020524072032770870"
			},
			"result": {
				"page": 0,
				"perPage": 20,
				"searchRequestId": "music-stable-back-vla-19.vla.yp-c.yandex.net-1704449354601545-12020524072032770870-1704449354610",
				"text": "Miley Cyrus – Flowers",
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
												"prefix": "f95e53d4.a.24665034-1",
												"type": "from-album-cover",
												"uri": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%"
											},
											"disclaimers": [],
											"genres": [],
											"id": 18375620,
											"name": "DJ Amor 2.0",
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
										110791278
									],
									"coverUri": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%",
									"disclaimers": [],
									"genre": "pop",
									"id": 24665034,
									"labels": [],
									"likesCount": 161,
									"metaType": "music",
									"ogImage": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%",
									"recent": false,
									"regions": [
										"BELARUS",
										"BELARUS_PREMIUM"
									],
									"releaseDate": "2023-01-18T00:00:00+03:00",
									"storageDir": "f95e53d4.a.24665034-1",
									"title": "Spiffy Trackz Vol 3",
									"trackCount": 2,
									"trackPosition": {
										"index": 2,
										"volume": 1
									},
									"type": "single",
									"veryImportant": false,
									"year": 2023
								}
							],
							"albums": [
								{
									"artists": [
										{
											"composer": false,
											"cover": {
												"prefix": "f95e53d4.a.24665034-1",
												"type": "from-album-cover",
												"uri": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%"
											},
											"disclaimers": [],
											"genres": [],
											"id": 18375620,
											"name": "DJ Amor 2.0",
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
										110791278
									],
									"coverUri": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%",
									"disclaimers": [],
									"genre": "pop",
									"id": 24665034,
									"labels": [],
									"likesCount": 161,
									"metaType": "music",
									"ogImage": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%",
									"recent": false,
									"releaseDate": "2023-01-18T00:00:00+03:00",
									"title": "Spiffy Trackz Vol 3",
									"trackCount": 2,
									"trackPosition": {
										"index": 2,
										"volume": 1
									},
									"type": "single",
									"veryImportant": false,
									"year": 2023
								}
							],
							"artists": [
								{
									"composer": false,
									"cover": {
										"prefix": "f95e53d4.a.24665034-1",
										"type": "from-album-cover",
										"uri": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%"
									},
									"disclaimers": [],
									"genres": [],
									"id": 18375620,
									"name": "DJ Amor 2.0",
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
							"coverUri": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%",
							"derivedColors": {
								"accent": "#999999",
								"average": "#828282",
								"miniPlayer": "#B7B7B7",
								"waveText": "#CCCCCC"
							},
							"disclaimers": [],
							"durationMs": 196030,
							"explicit": false,
							"fade": {
								"inStart": 0.4,
								"inStop": 1.5,
								"outStart": 191.9,
								"outStop": 193.1
							},
							"fileSize": 0,
							"id": 110791278,
							"lyricsAvailable": false,
							"lyricsInfo": {
								"hasAvailableSyncLyrics": false,
								"hasAvailableTextLyrics": false
							},
							"major": {
								"id": 399,
								"name": "VOICE_EXPRESS"
							},
							"ogImage": "avatars.yandex.net/get-music-content/7852894/f95e53d4.a.24665034-1/%%",
							"previewDurationMs": 30000,
							"r128": {
								"i": -14.15,
								"tp": 0.48
							},
							"realId": "110791278",
							"regions": [
								"BELARUS",
								"BELARUS_PREMIUM"
							],
							"rememberPosition": false,
							"storageDir": "",
							"title": "Flowers",
							"trackSharingFlag": "COVER_ONLY",
							"trackSource": "OWN",
							"type": "music",
							"version": "Instrumental Tribute Version Originally Performed By Miley Cyrus"
						}                
					],
					"total": 22
				},
				"type": "track"
			}
		}`,
	}
	TrackFixtureNotFound = TrackFixture{
		GetResponse: `{
			"invocationInfo": {
				"exec-duration-millis": 11,
				"hostname": "music-stable-back-vla-24.vla.yp-c.yandex.net",
				"req-id": "1704222836102466-11015214182353549756"
			},
			"result": []
		}`,
		SearchResponse: `{
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
		}`,
	}
)

func (tf TrackFixture) TrackWithIDString() *ymusic.Track {
	tJSON, _ := json.Marshal(tf.Track)
	t := ymusic.Track{}
	_ = json.Unmarshal(tJSON, &t)
	t.ID = t.IDString()
	return &t
}

func (tf TrackFixture) SearchQuery() string {
	return fmt.Sprintf("%s – %s", tf.SearchQueryArtist, tf.SearchQueryTrack)
}
