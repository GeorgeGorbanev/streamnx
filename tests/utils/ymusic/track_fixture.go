package ymusic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/ymusic"
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
	TrackFixtureZemfiraIskala = TrackFixture{
		ID: "732401",
		Track: &ymusic.Track{
			ID:    732401.0,
			Title: "ИСКАЛА",
			Artists: []ymusic.Artist{
				{
					ID:   218099,
					Name: "Земфира",
				},
			},
			Albums: []ymusic.Album{
				{
					ID: 81431,
				},
			},
		},
		SearchQueryArtist: "Zemfira",
		SearchQueryTrack:  "ИСКАЛА",
		SearchResponse: `{
			"invocationInfo": {
				"exec-duration-millis": "109",
				"hostname": "music-stable-back-vla-94.vla.yp-c.yandex.net",
				"req-id": "1704481543921044-1749570170531180058"
			},
			"result": {
				"page": 0,
				"perPage": 20,
				"searchRequestId": "music-stable-back-vla-94.vla.yp-c.yandex.net-1704481543921044-1749570170531180058-1704481543924",
				"text": "Zemfira – ИСКАЛА",
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
												"prefix": "e01c1f0f.p.218099/",
												"type": "from-artist-photos",
												"uri": "avatars.yandex.net/get-music-content/4747389/e01c1f0f.p.218099/%%"
											},
											"disclaimers": [],
											"genres": [],
											"id": 218099,
											"name": "Земфира",
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
										732401,
										732400,
										732405
									],
									"coverUri": "avatars.yandex.net/get-music-content/5234929/e893dd2c.a.81431-2/%%",
									"disclaimers": [],
									"genre": "rusrock",
									"id": 81431,
									"labels": [
										"Земфира"
									],
									"likesCount": 60258,
									"metaType": "music",
									"ogImage": "avatars.yandex.net/get-music-content/5234929/e893dd2c.a.81431-2/%%",
									"recent": false,
									"regions": [
										"BELARUS",
										"BELARUS_PREMIUM"
									],
									"storageDir": "e893dd2c.a.81431-2",
									"title": "ПРОСТИ МЕНЯ МОЯ ЛЮБОВЬ",
									"trackCount": 13,
									"trackPosition": {
										"index": 11,
										"volume": 1
									},
									"veryImportant": false,
									"year": 2000
								}
							],
							"albums": [
								{
									"artists": [
										{
											"composer": false,
											"cover": {
												"prefix": "e01c1f0f.p.218099/",
												"type": "from-artist-photos",
												"uri": "avatars.yandex.net/get-music-content/4747389/e01c1f0f.p.218099/%%"
											},
											"disclaimers": [],
											"genres": [],
											"id": 218099,
											"name": "Земфира",
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
										732401,
										732400,
										732405
									],
									"coverUri": "avatars.yandex.net/get-music-content/5234929/e893dd2c.a.81431-2/%%",
									"disclaimers": [],
									"genre": "rusrock",
									"id": 81431,
									"labels": [
										{
											"id": 573394,
											"name": "Земфира"
										}
									],
									"likesCount": 60258,
									"metaType": "music",
									"ogImage": "avatars.yandex.net/get-music-content/5234929/e893dd2c.a.81431-2/%%",
									"recent": false,
									"title": "ПРОСТИ МЕНЯ МОЯ ЛЮБОВЬ",
									"trackCount": 13,
									"trackPosition": {
										"index": 11,
										"volume": 1
									},
									"veryImportant": false,
									"year": 2000
								}
							],
							"artists": [
								{
									"composer": false,
									"cover": {
										"prefix": "e01c1f0f.p.218099/",
										"type": "from-artist-photos",
										"uri": "avatars.yandex.net/get-music-content/4747389/e01c1f0f.p.218099/%%"
									},
									"disclaimers": [],
									"genres": [],
									"id": 218099,
									"name": "Земфира",
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
							"coverUri": "avatars.yandex.net/get-music-content/5234929/e893dd2c.a.81431-2/%%",
							"derivedColors": {
								"accent": "#999999",
								"average": "#999999",
								"miniPlayer": "#B7B7B7",
								"waveText": "#CCCCCC"
							},
							"disclaimers": [],
							"durationMs": 214040,
							"explicit": false,
							"fade": {
								"inStart": 0.2,
								"inStop": 1.5,
								"outStart": 206.7,
								"outStop": 209.2
							},
							"fileSize": 0,
							"id": 732401,
							"lyricsAvailable": true,
							"lyricsInfo": {
								"hasAvailableSyncLyrics": true,
								"hasAvailableTextLyrics": true
							},
							"major": {
								"id": 135,
								"name": "ZEMFIRA"
							},
							"ogImage": "avatars.yandex.net/get-music-content/5234929/e893dd2c.a.81431-2/%%",
							"previewDurationMs": 30000,
							"r128": {
								"i": -11.30,
								"tp": 0.21
							},
							"realId": "732401",
							"regions": [
								"BELARUS",
								"BELARUS_PREMIUM"
							],
							"rememberPosition": false,
							"storageDir": "",
							"title": "ИСКАЛА",
							"trackSharingFlag": "COVER_ONLY",
							"trackSource": "OWN",
							"type": "music"
						}                
					],
					"total": 166
				},
				"type": "track"
			}
		}
		`,
	}
	TrackFixtureNadezhdaKadyshevaShirokaReka = TrackFixture{
		ID: "33223088",
		Track: &ymusic.Track{
			ID:    33223088.0,
			Title: "Широка река",
			Artists: []ymusic.Artist{
				{
					ID:   164414,
					Name: "Надежда Кадышева",
				},
			},
			Albums: []ymusic.Album{
				{
					ID: 4058886,
				},
			},
		},
		SearchQueryArtist: "надежда кадышева",
		SearchQueryTrack:  "Широка Река",
		SearchResponse: `{
			"invocationInfo": {
				"exec-duration-millis": "151",
				"hostname": "music-stable-back-klg-74.klg.yp-c.yandex.net",
				"req-id": "1704743108236830-13332655568070335715"
			},
			"result": {
				"page": 0,
				"perPage": 20,
				"searchRequestId": "music-stable-back-klg-74.klg.yp-c.yandex.net-1704743108236830-13332655568070335715-1704743108248",
				"text": "надежда кадышева – широка река",
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
												"prefix": "40fc25be.p.164414/",
												"type": "from-artist-photos",
												"uri": "avatars.yandex.net/get-music-content/5502420/40fc25be.p.164414/%%"
											},
											"decomposed": [
												" & ",
												{
													"composer": false,
													"cover": {
														"prefix": "18ebbf20.p.907549/",
														"type": "from-artist-photos",
														"uri": "avatars.yandex.net/get-music-content/2359742/18ebbf20.p.907549/%%"
													},
													"disclaimers": [],
													"genres": [],
													"id": 907549,
													"name": "Золотое кольцо",
													"various": false
												}
											],
											"disclaimers": [],
											"genres": [],
											"id": 164414,
											"name": "Надежда Кадышева",
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
										33223085,
										33223090,
										33223088
									],
									"coverUri": "avatars.yandex.net/get-music-content/139444/b54a3d47.a.4058886-1/%%",
									"disclaimers": [],
									"genre": "rusfolk",
									"id": 4058886,
									"labels": [
										"United Music Group"
									],
									"likesCount": 898,
									"metaType": "music",
									"ogImage": "avatars.yandex.net/get-music-content/139444/b54a3d47.a.4058886-1/%%",
									"recent": false,
									"regions": [
										"GEORGIA",
										"GEORGIA_PREMIUM"
									],
									"releaseDate": "2009-04-02T00:00:00+04:00",
									"storageDir": "b54a3d47.a.4058886-1",
									"title": "Когда-нибудь…",
									"trackCount": 14,
									"trackPosition": {
										"index": 4,
										"volume": 1
									},
									"veryImportant": false,
									"year": 2003
								}
							],                    
							"artists": [
								{
									"composer": false,
									"cover": {
										"prefix": "40fc25be.p.164414/",
										"type": "from-artist-photos",
										"uri": "avatars.yandex.net/get-music-content/5502420/40fc25be.p.164414/%%"
									},
									"decomposed": [
										" & ",
										{
											"composer": false,
											"cover": {
												"prefix": "18ebbf20.p.907549/",
												"type": "from-artist-photos",
												"uri": "avatars.yandex.net/get-music-content/2359742/18ebbf20.p.907549/%%"
											},
											"disclaimers": [],
											"genres": [],
											"id": 907549,
											"name": "Золотое кольцо",
											"various": false
										}
									],
									"disclaimers": [],
									"genres": [],
									"id": 164414,
									"name": "Надежда Кадышева",
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
							"coverUri": "avatars.yandex.net/get-music-content/139444/b54a3d47.a.4058886-1/%%",
							"derivedColors": {
								"accent": "#999999",
								"average": "#8E8E8E",
								"miniPlayer": "#B7B7B7",
								"waveText": "#CCCCCC"
							},
							"disclaimers": [],
							"durationMs": 238720,
							"explicit": false,
							"fade": {
								"inStart": 0.0,
								"inStop": 1.5,
								"outStart": 229.6,
								"outStop": 234.3
							},
							"fileSize": 0,
							"id": 33223088,
							"lyricsAvailable": false,
							"lyricsInfo": {
								"hasAvailableSyncLyrics": true,
								"hasAvailableTextLyrics": true
							},
							"major": {
								"id": 123,
								"name": "IRICOM"
							},
							"ogImage": "avatars.yandex.net/get-music-content/139444/b54a3d47.a.4058886-1/%%",
							"previewDurationMs": 30000,
							"r128": {
								"i": -11.14,
								"tp": 0.78
							},
							"realId": "33223088",
							"regions": [
								"GEORGIA",
								"GEORGIA_PREMIUM"
							],
							"rememberPosition": false,
							"storageDir": "",
							"title": "Широка река",
							"trackSharingFlag": "VIDEO_ALLOWED",
							"trackSource": "OWN",
							"type": "music"
						}
						
					],
					"total": 45
				},
				"type": "track"
			}
		}
		`,
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
	return fmt.Sprintf(
		"%s – %s",
		strings.ToLower(tf.SearchQueryArtist),
		strings.ToLower(tf.SearchQueryTrack),
	)
}
