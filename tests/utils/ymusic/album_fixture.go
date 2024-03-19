package ymusic

import (
	"fmt"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/ymusic"
)

type AlbumFixture struct {
	ID                string
	Album             *ymusic.Album
	GetResponse       string
	SearchQueryAlbum  string
	SearchQueryArtist string
	SearchResponse    string
}

var (
	AlbumFixtureRadioheadAmnesiac = AlbumFixture{
		ID:                "3389008",
		SearchQueryArtist: "radiohead",
		SearchQueryAlbum:  "amnesiac",
		Album: &ymusic.Album{
			ID:    3389008,
			Title: "Amnesiac",
			Artists: []ymusic.Artist{
				{
					ID:   36825,
					Name: "Radiohead",
				},
			},
		},
		GetResponse: `
			{
				"invocationInfo": {
					"hostname": "music-stable-back-vla-49.vla.yp-c.yandex.net",
					"req-id": "1710827171729552-10782237058029685308",
					"exec-duration-millis": "0"
				},
				"result": {
					"id": 3389008,
					"title": "Amnesiac",
					"metaType": "music",
					"year": 2001,
					"releaseDate": "2001-03-12T00:00:00+03:00",
					"coverUri": "avatars.yandex.net/get-music-content/4399834/13c8d84b.a.3389008-2/%%",
					"ogImage": "avatars.yandex.net/get-music-content/4399834/13c8d84b.a.3389008-2/%%",
					"genre": "indie",
					"metaTagId": "58e3c0796ebb3476087d190c",
					"trackCount": 11,
					"likesCount": 7321,
					"recent": false,
					"veryImportant": false,
					"artists": [
						{
							"id": 36825,
							"name": "Radiohead",
							"various": false,
							"composer": false,
							"cover": {
								"type": "from-artist-photos",
								"uri": "avatars.yandex.net/get-music-content/5314916/561eaa85.p.36825/%%",
								"prefix": "561eaa85.p.36825/"
							},
							"genres": [],
							"disclaimers": []
						}
					],
					"labels": [
						{
							"id": 13,
							"name": "XL"
						}
					],
					"available": true,
					"availableForPremiumUsers": true,
					"availableForOptions": [
						"bookmate"
					],
					"availableForMobile": true,
					"availablePartially": false,
					"bests": [
						333447,
						333456,
						333449
					],
					"disclaimers": []
				}
			}
		`,
		SearchResponse: `
			{
				"invocationInfo": {
					"hostname": "music-stable-back-vla-16.vla.yp-c.yandex.net",
					"req-id": "1710781983719419-6063053006080689686",
					"exec-duration-millis": "86"
				},
				"result": {
					"type": "album",
					"page": 0,
					"perPage": 10,
					"text": "Radiohead – Amnesiac",
					"searchRequestId": "music-stable-back-vla-16.vla.yp-c.yandex.net-1710781983719419-6063053006080689686-1710781983729",
					"albums": {
						"total": 5,
						"perPage": 10,
						"order": 0,
						"results": [
							{
								"id": 3389008,
								"title": "Amnesiac",
								"metaType": "music",
								"year": 2001,
								"releaseDate": "2001-03-12T00:00:00+03:00",
								"coverUri": "avatars.yandex.net/get-music-content/4399834/13c8d84b.a.3389008-2/%%",
								"ogImage": "avatars.yandex.net/get-music-content/4399834/13c8d84b.a.3389008-2/%%",
								"genre": "indie",
								"trackCount": 11,
								"likesCount": 7319,
								"recent": false,
								"veryImportant": false,
								"artists": [
									{
										"id": 36825,
										"name": "Radiohead",
										"various": false,
										"composer": false,
										"cover": {
											"type": "from-artist-photos",
											"uri": "avatars.yandex.net/get-music-content/5314916/561eaa85.p.36825/%%",
											"prefix": "561eaa85.p.36825/"
										},
										"genres": [],
										"disclaimers": []
									}
								],
								"available": true,
								"availableForPremiumUsers": true,
								"availableForOptions": [
									"bookmate"
								],
								"availableForMobile": true,
								"availablePartially": false,
								"bests": [
									333447,
									333456,
									333449
								],
								"disclaimers": [],
								"labels": [
									"XL"
								],
								"storageDir": "13c8d84b.a.3389008-2",
								"regions": [
									"RUSSIA",
									"RUSSIA_PREMIUM"
								],
								"availableRegions": [
									"kz",
									"tm",
									"kg",
									"by",
									"om",
									"eg",
									"tj",
									"tn",
									"am",
									"ae",
									"uz",
									"bh",
									"ge",
									"il",
									"qa",
									"az",
									"md",
									"ua",
									"kw",
									"ru"
								]
							},
							{
								"id": 19028505,
								"title": "KID A MNESIA",
								"metaType": "music",
								"year": 2021,
								"releaseDate": "2021-11-05T00:00:00+03:00",
								"coverUri": "avatars.yandex.net/get-music-content/4401814/85ff9ba4.a.19028505-1/%%",
								"ogImage": "avatars.yandex.net/get-music-content/4401814/85ff9ba4.a.19028505-1/%%",
								"genre": "indie",
								"trackCount": 34,
								"likesCount": 7329,
								"recent": false,
								"veryImportant": false,
								"artists": [
									{
										"id": 36825,
										"name": "Radiohead",
										"various": false,
										"composer": false,
										"cover": {
											"type": "from-artist-photos",
											"uri": "avatars.yandex.net/get-music-content/5314916/561eaa85.p.36825/%%",
											"prefix": "561eaa85.p.36825/"
										},
										"genres": [],
										"disclaimers": []
									}
								],
								"available": true,
								"availableForPremiumUsers": true,
								"availableForOptions": [
									"bookmate"
								],
								"availableForMobile": true,
								"availablePartially": false,
								"bests": [
									333456,
									333470,
									333480
								],
								"disclaimers": [],
								"labels": [
									"XL"
								],
								"storageDir": "85ff9ba4.a.19028505-1",
								"regions": [
									"RUSSIA",
									"RUSSIA_PREMIUM"
								],
								"availableRegions": [
									"kz",
									"tm",
									"kg",
									"by",
									"om",
									"eg",
									"tj",
									"tn",
									"am",
									"ae",
									"uz",
									"bh",
									"ge",
									"il",
									"qa",
									"az",
									"md",
									"ua",
									"kw",
									"ru"
								]
							},
							{
								"id": 22623579,
								"title": "Tribute to Radiohead",
								"metaType": "music",
								"year": 2010,
								"releaseDate": "2010-01-01T00:00:00+03:00",
								"coverUri": "avatars.yandex.net/get-music-content/6202531/69587b78.a.22623579-1/%%",
								"ogImage": "avatars.yandex.net/get-music-content/6202531/69587b78.a.22623579-1/%%",
								"genre": "jazz",
								"trackCount": 5,
								"likesCount": 3,
								"recent": false,
								"veryImportant": false,
								"artists": [
									{
										"id": 16778281,
										"name": "Amnesiac Quartet",
										"various": false,
										"composer": false,
										"cover": {
											"type": "from-album-cover",
											"uri": "avatars.yandex.net/get-music-content/6021799/9570690c.a.22771300-1/%%",
											"prefix": "9570690c.a.22771300-1"
										},
										"genres": [],
										"disclaimers": []
									}
								],
								"available": true,
								"availableForPremiumUsers": true,
								"availableForOptions": [
									"bookmate"
								],
								"availableForMobile": true,
								"availablePartially": false,
								"bests": [],
								"disclaimers": [],
								"labels": [
									"Sébastien Paindestre"
								],
								"storageDir": "69587b78.a.22623579-1",
								"regions": [
									"RUSSIA",
									"RUSSIA_PREMIUM"
								],
								"availableRegions": [
									"kz",
									"tm",
									"sa",
									"kg",
									"by",
									"om",
									"eg",
									"tj",
									"tn",
									"am",
									"ae",
									"uz",
									"bh",
									"ge",
									"il",
									"qa",
									"az",
									"md",
									"ua",
									"kw",
									"ru"
								]
							},
							{
								"id": 22623572,
								"title": "Tribute to Radiohead, Vol. 2",
								"metaType": "music",
								"year": 2013,
								"releaseDate": "2013-01-10T00:00:00+04:00",
								"coverUri": "avatars.yandex.net/get-music-content/6447985/d3aacc63.a.22623572-1/%%",
								"ogImage": "avatars.yandex.net/get-music-content/6447985/d3aacc63.a.22623572-1/%%",
								"genre": "conjazz",
								"trackCount": 7,
								"likesCount": 2,
								"recent": false,
								"veryImportant": false,
								"artists": [
									{
										"id": 16778281,
										"name": "Amnesiac Quartet",
										"various": false,
										"composer": false,
										"cover": {
											"type": "from-album-cover",
											"uri": "avatars.yandex.net/get-music-content/6021799/9570690c.a.22771300-1/%%",
											"prefix": "9570690c.a.22771300-1"
										},
										"genres": [],
										"disclaimers": []
									}
								],
								"available": true,
								"availableForPremiumUsers": true,
								"availableForOptions": [
									"bookmate"
								],
								"availableForMobile": true,
								"availablePartially": false,
								"bests": [],
								"disclaimers": [],
								"labels": [
									"Sébastien Paindestre"
								],
								"storageDir": "d3aacc63.a.22623572-1",
								"regions": [
									"RUSSIA",
									"RUSSIA_PREMIUM"
								],
								"availableRegions": [
									"kz",
									"tm",
									"sa",
									"kg",
									"by",
									"om",
									"eg",
									"tj",
									"tn",
									"am",
									"ae",
									"uz",
									"bh",
									"ge",
									"il",
									"qa",
									"az",
									"md",
									"ua",
									"kw",
									"ru"
								]
							},
							{
								"id": 22771300,
								"title": "Tribute to Radiohead, Vol. 3",
								"metaType": "music",
								"year": 2022,
								"releaseDate": "2022-07-01T00:00:00+03:00",
								"coverUri": "avatars.yandex.net/get-music-content/6021799/9570690c.a.22771300-1/%%",
								"ogImage": "avatars.yandex.net/get-music-content/6021799/9570690c.a.22771300-1/%%",
								"genre": "jazz",
								"trackCount": 17,
								"likesCount": 2,
								"recent": false,
								"veryImportant": false,
								"artists": [
									{
										"id": 16778281,
										"name": "Amnesiac Quartet",
										"various": false,
										"composer": false,
										"cover": {
											"type": "from-album-cover",
											"uri": "avatars.yandex.net/get-music-content/6021799/9570690c.a.22771300-1/%%",
											"prefix": "9570690c.a.22771300-1"
										},
										"genres": [],
										"disclaimers": []
									}
								],
								"available": true,
								"availableForPremiumUsers": true,
								"availableForOptions": [
									"bookmate"
								],
								"availableForMobile": true,
								"availablePartially": false,
								"bests": [],
								"disclaimers": [],
								"labels": [
									"Sébastien Paindestre"
								],
								"storageDir": "9570690c.a.22771300-1",
								"regions": [
									"RUSSIA",
									"RUSSIA_PREMIUM"
								],
								"availableRegions": [
									"kz",
									"tm",
									"sa",
									"kg",
									"by",
									"om",
									"eg",
									"tj",
									"tn",
									"am",
									"ae",
									"uz",
									"bh",
									"ge",
									"il",
									"qa",
									"az",
									"md",
									"ua",
									"kw",
									"ru"
								]
							}
						]
					}
				}
			}
		`,
	}
	AlbumFixtureNotFound = AlbumFixture{
		ID: "0",
		GetResponse: `{
			"invocationInfo": {
				"hostname": "music-stable-back-vla-24.vla.yp-c.yandex.net",
				"req-id": "1710868512438509-4151477972976085889",
				"exec-duration-millis": "1"
			},
			"result": {
				"id": 0,
				"error": "not-found"
			}
		}`,
	}
)

func (af AlbumFixture) SearchQuery() string {
	return fmt.Sprintf(
		"%s – %s",
		strings.ToLower(af.SearchQueryArtist),
		strings.ToLower(af.SearchQueryAlbum),
	)
}
