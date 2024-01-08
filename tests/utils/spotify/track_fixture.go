package spotify

import (
	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
)

type TrackFixture struct {
	Track          *spotify.Track
	GetResponse    string
	SearchQuery    string
	SearchResponse string
}

var (
	TrackFixtureMassiveAttackAngel = TrackFixture{
		Track: &spotify.Track{
			ID: "7uv632EkfwYhXoqf8rhYrg",
			Artists: []spotify.Artist{
				{
					Name: "Massive Attack",
				},
			},
			Name: "Angel",
		},
		SearchQuery: "artist:Massive Attack track:Angel",
		SearchResponse: `{
		  "tracks" : {
			"href" : "https://api.spotify.com/v1/search?query=artist%3AMassive+Attack+track%3AAngel&type=track&offset=0&limit=1",
			"items" : [ {
			  "album" : {
				"album_type" : "album",
				"artists" : [ {
				  "external_urls" : {
					"spotify" : "https://open.spotify.com/artist/6FXMGgJwohJLUSr5nVlf9X"
				  },
				  "href" : "https://api.spotify.com/v1/artists/6FXMGgJwohJLUSr5nVlf9X",
				  "id" : "6FXMGgJwohJLUSr5nVlf9X",
				  "name" : "Massive Attack",
				  "type" : "artist",
				  "uri" : "spotify:artist:6FXMGgJwohJLUSr5nVlf9X"
				} ],
				"available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
				"external_urls" : {
				  "spotify" : "https://open.spotify.com/album/49MNmJhZQewjt06rpwp6QR"
				},
				"href" : "https://api.spotify.com/v1/albums/49MNmJhZQewjt06rpwp6QR",
				"id" : "49MNmJhZQewjt06rpwp6QR",
				"images" : [ {
				  "height" : 640,
				  "url" : "https://i.scdn.co/image/ab67616d0000b2732fcb0a3c7a66e516b11cd26e",
				  "width" : 640
				}, {
				  "height" : 300,
				  "url" : "https://i.scdn.co/image/ab67616d00001e022fcb0a3c7a66e516b11cd26e",
				  "width" : 300
				}, {
				  "height" : 64,
				  "url" : "https://i.scdn.co/image/ab67616d000048512fcb0a3c7a66e516b11cd26e",
				  "width" : 64
				} ],
				"name" : "Mezzanine",
				"release_date" : "1998-01-01",
				"release_date_precision" : "day",
				"total_tracks" : 11,
				"type" : "album",
				"uri" : "spotify:album:49MNmJhZQewjt06rpwp6QR"
			  },
			  "artists" : [ {
				"external_urls" : {
				  "spotify" : "https://open.spotify.com/artist/6FXMGgJwohJLUSr5nVlf9X"
				},
				"href" : "https://api.spotify.com/v1/artists/6FXMGgJwohJLUSr5nVlf9X",
				"id" : "6FXMGgJwohJLUSr5nVlf9X",
				"name" : "Massive Attack",
				"type" : "artist",
				"uri" : "spotify:artist:6FXMGgJwohJLUSr5nVlf9X"
			  }],
			  "available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
			  "disc_number" : 1,
			  "duration_ms" : 379533,
			  "explicit" : false,
			  "external_ids" : {
				"isrc" : "GBAAA9800327"
			  },
			  "external_urls" : {
				"spotify" : "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg"
			  },
			  "href" : "https://api.spotify.com/v1/tracks/7uv632EkfwYhXoqf8rhYrg",
			  "id" : "7uv632EkfwYhXoqf8rhYrg",
			  "is_local" : false,
			  "name" : "Angel",
			  "popularity" : 65,
			  "preview_url" : null,
			  "track_number" : 1,
			  "type" : "track",
			  "uri" : "spotify:track:7uv632EkfwYhXoqf8rhYrg"
			} ],
			"limit" : 1,
			"next" : "https://api.spotify.com/v1/search?query=artist%3AMassive+Attack+track%3AAngel&type=track&offset=1&limit=1",
			"offset" : 0,
			"previous" : null,
			"total" : 25
		  }
		}`,
		GetResponse: `{
		  "album" : {
			"album_type" : "album",
			"artists" : [ {
			  "external_urls" : {
				"spotify" : "https://open.spotify.com/artist/6FXMGgJwohJLUSr5nVlf9X"
			  },
			  "href" : "https://api.spotify.com/v1/artists/6FXMGgJwohJLUSr5nVlf9X",
			  "id" : "6FXMGgJwohJLUSr5nVlf9X",
			  "name" : "Massive Attack",
			  "type" : "artist",
			  "uri" : "spotify:artist:6FXMGgJwohJLUSr5nVlf9X"
			} ],
			"available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/album/49MNmJhZQewjt06rpwp6QR"
			},
			"href" : "https://api.spotify.com/v1/albums/49MNmJhZQewjt06rpwp6QR",
			"id" : "49MNmJhZQewjt06rpwp6QR",
			"images" : [ {
			  "height" : 640,
			  "url" : "https://i.scdn.co/image/ab67616d0000b2732fcb0a3c7a66e516b11cd26e",
			  "width" : 640
			}, {
			  "height" : 300,
			  "url" : "https://i.scdn.co/image/ab67616d00001e022fcb0a3c7a66e516b11cd26e",
			  "width" : 300
			}, {
			  "height" : 64,
			  "url" : "https://i.scdn.co/image/ab67616d000048512fcb0a3c7a66e516b11cd26e",
			  "width" : 64
			} ],
			"name" : "Mezzanine",
			"release_date" : "1998-01-01",
			"release_date_precision" : "day",
			"total_tracks" : 11,
			"type" : "album",
			"uri" : "spotify:album:49MNmJhZQewjt06rpwp6QR"
		  },
		  "artists" : [ {
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/artist/6FXMGgJwohJLUSr5nVlf9X"
			},
			"href" : "https://api.spotify.com/v1/artists/6FXMGgJwohJLUSr5nVlf9X",
			"id" : "6FXMGgJwohJLUSr5nVlf9X",
			"name" : "Massive Attack",
			"type" : "artist",
			"uri" : "spotify:artist:6FXMGgJwohJLUSr5nVlf9X"
		  }],
		  "available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
		  "disc_number" : 1,
		  "duration_ms" : 379533,
		  "explicit" : false,
		  "external_ids" : {
			"isrc" : "GBAAA9800327"
		  },
		  "external_urls" : {
			"spotify" : "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg"
		  },
		  "href" : "https://api.spotify.com/v1/tracks/7uv632EkfwYhXoqf8rhYrg",
		  "id" : "7uv632EkfwYhXoqf8rhYrg",
		  "is_local" : false,
		  "name" : "Angel",
		  "popularity" : 65,
		  "preview_url" : null,
		  "track_number" : 1,
		  "type" : "track",
		  "uri" : "spotify:track:7uv632EkfwYhXoqf8rhYrg"
		}`,
	}
	TrackFixtureMileyCyrusFlowers = TrackFixture{
		Track: &spotify.Track{
			ID:   "7DSAEUvxU8FajXtRloy8M0",
			Name: "Flowers",
			Artists: []spotify.Artist{
				{
					Name: "Miley Cyrus",
				},
			},
		},
		GetResponse: `
			{
			  "album" : {
				"album_type" : "album",
				"artists" : [ {
				  "external_urls" : {
					"spotify" : "https://open.spotify.com/artist/5YGY8feqx7naU7z4HrwZM6"
				  },
				  "href" : "https://api.spotify.com/v1/artists/5YGY8feqx7naU7z4HrwZM6",
				  "id" : "5YGY8feqx7naU7z4HrwZM6",
				  "name" : "Miley Cyrus",
				  "type" : "artist",
				  "uri" : "spotify:artist:5YGY8feqx7naU7z4HrwZM6"
				} ],
				"available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
				"external_urls" : {
				  "spotify" : "https://open.spotify.com/album/5DvJgsMLbaR1HmAI6VhfcQ"
				},
				"href" : "https://api.spotify.com/v1/albums/5DvJgsMLbaR1HmAI6VhfcQ",
				"id" : "5DvJgsMLbaR1HmAI6VhfcQ",
				"images" : [ {
				  "height" : 640,
				  "url" : "https://i.scdn.co/image/ab67616d0000b273cd222052a2594be29a6616b5",
				  "width" : 640
				}, {
				  "height" : 300,
				  "url" : "https://i.scdn.co/image/ab67616d00001e02cd222052a2594be29a6616b5",
				  "width" : 300
				}, {
				  "height" : 64,
				  "url" : "https://i.scdn.co/image/ab67616d00004851cd222052a2594be29a6616b5",
				  "width" : 64
				} ],
				"name" : "Endless Summer Vacation",
				"release_date" : "2023-08-18",
				"release_date_precision" : "day",
				"total_tracks" : 14,
				"type" : "album",
				"uri" : "spotify:album:5DvJgsMLbaR1HmAI6VhfcQ"
			  },
			  "artists" : [ {
				"external_urls" : {
				  "spotify" : "https://open.spotify.com/artist/5YGY8feqx7naU7z4HrwZM6"
				},
				"href" : "https://api.spotify.com/v1/artists/5YGY8feqx7naU7z4HrwZM6",
				"id" : "5YGY8feqx7naU7z4HrwZM6",
				"name" : "Miley Cyrus",
				"type" : "artist",
				"uri" : "spotify:artist:5YGY8feqx7naU7z4HrwZM6"
			  } ],
			  "available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
			  "disc_number" : 1,
			  "duration_ms" : 200600,
			  "explicit" : false,
			  "external_ids" : {
				"isrc" : "USSM12209777"
			  },
			  "external_urls" : {
				"spotify" : "https://open.spotify.com/track/7DSAEUvxU8FajXtRloy8M0"
			  },
			  "href" : "https://api.spotify.com/v1/tracks/7DSAEUvxU8FajXtRloy8M0",
			  "id" : "7DSAEUvxU8FajXtRloy8M0",
			  "is_local" : false,
			  "name" : "Flowers",
			  "popularity" : 92,
			  "preview_url" : "https://p.scdn.co/mp3-preview/5184d19d1b7fcc3e7c067e38af45a7cc80851440?cid=e7b45730e3774355ae6a15e7e4d188da",
			  "track_number" : 1,
			  "type" : "track",
			  "uri" : "spotify:track:7DSAEUvxU8FajXtRloy8M0"
			}
		`,
	}
	TrackFixtureZemfiraIskala = TrackFixture{
		Track: &spotify.Track{
			ID:   "3NHSz1GyC5IeK1soZSjIIX",
			Name: "ИСКАЛА",
			Artists: []spotify.Artist{
				{
					Name: "Zemfira",
				},
			},
		},
		GetResponse: `{
		  "album" : {
			"album_type" : "album",
			"artists" : [ {
			  "external_urls" : {
				"spotify" : "https://open.spotify.com/artist/6oO3QiWdVj5FZQwbdRtsRh"
			  },
			  "href" : "https://api.spotify.com/v1/artists/6oO3QiWdVj5FZQwbdRtsRh",
			  "id" : "6oO3QiWdVj5FZQwbdRtsRh",
			  "name" : "Zemfira",
			  "type" : "artist",
			  "uri" : "spotify:artist:6oO3QiWdVj5FZQwbdRtsRh"
			} ],
			"available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/album/7wq0Fe0wGDMlIbmFs8BERa"
			},
			"href" : "https://api.spotify.com/v1/albums/7wq0Fe0wGDMlIbmFs8BERa",
			"id" : "7wq0Fe0wGDMlIbmFs8BERa",
			"images" : [ {
			  "height" : 640,
			  "url" : "https://i.scdn.co/image/ab67616d0000b273c044eb2bcc1c20f76916fcf7",
			  "width" : 640
			}, {
			  "height" : 300,
			  "url" : "https://i.scdn.co/image/ab67616d00001e02c044eb2bcc1c20f76916fcf7",
			  "width" : 300
			}, {
			  "height" : 64,
			  "url" : "https://i.scdn.co/image/ab67616d00004851c044eb2bcc1c20f76916fcf7",
			  "width" : 64
			} ],
			"name" : "ПРОСТИ МЕНЯ МОЯ ЛЮБОВЬ",
			"release_date" : "2000-03-28",
			"release_date_precision" : "day",
			"total_tracks" : 13,
			"type" : "album",
			"uri" : "spotify:album:7wq0Fe0wGDMlIbmFs8BERa"
		  },
		  "artists" : [ {
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/artist/6oO3QiWdVj5FZQwbdRtsRh"
			},
			"href" : "https://api.spotify.com/v1/artists/6oO3QiWdVj5FZQwbdRtsRh",
			"id" : "6oO3QiWdVj5FZQwbdRtsRh",
			"name" : "Zemfira",
			"type" : "artist",
			"uri" : "spotify:artist:6oO3QiWdVj5FZQwbdRtsRh"
		  } ],
		  "available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
		  "disc_number" : 1,
		  "duration_ms" : 214000,
		  "explicit" : false,
		  "external_ids" : {
			"isrc" : "FR6V82120156"
		  },
		  "external_urls" : {
			"spotify" : "https://open.spotify.com/track/3NHSz1GyC5IeK1soZSjIIX"
		  },
		  "href" : "https://api.spotify.com/v1/tracks/3NHSz1GyC5IeK1soZSjIIX",
		  "id" : "3NHSz1GyC5IeK1soZSjIIX",
		  "is_local" : false,
		  "name" : "ИСКАЛА",
		  "popularity" : 58,
		  "preview_url" : "https://p.scdn.co/mp3-preview/c7c85a6c8b719f848fd48acb180c9ad41989f962?cid=e7b45730e3774355ae6a15e7e4d188da",
		  "track_number" : 11,
		  "type" : "track",
		  "uri" : "spotify:track:3NHSz1GyC5IeK1soZSjIIX"
		}
		`,
	}
	TrackFixtureNadezhdaKadyshevaShirokaReka = TrackFixture{
		Track: &spotify.Track{
			ID:   "2sP5VgY8PWb6c9DhgZEpSv",
			Name: "Широка река",
			Artists: []spotify.Artist{
				{
					Name: "Nadezhda Kadysheva",
				},
			},
		},
		GetResponse: `{
		  "album" : {
			"album_type" : "single",
			"artists" : [ {
			  "external_urls" : {
				"spotify" : "https://open.spotify.com/artist/6MnbF4D0Fv2zUVnf08Df5t"
			  },
			  "href" : "https://api.spotify.com/v1/artists/6MnbF4D0Fv2zUVnf08Df5t",
			  "id" : "6MnbF4D0Fv2zUVnf08Df5t",
			  "name" : "Nadezhda Kadysheva",
			  "type" : "artist",
			  "uri" : "spotify:artist:6MnbF4D0Fv2zUVnf08Df5t"
			}, {
			  "external_urls" : {
				"spotify" : "https://open.spotify.com/artist/2gpELJDuhOQi1DPcI0VFvw"
			  },
			  "href" : "https://api.spotify.com/v1/artists/2gpELJDuhOQi1DPcI0VFvw",
			  "id" : "2gpELJDuhOQi1DPcI0VFvw",
			  "name" : "Al Bano",
			  "type" : "artist",
			  "uri" : "spotify:artist:2gpELJDuhOQi1DPcI0VFvw"
			} ],
			"available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/album/4r1nB3kCpczB7c5MzNjCBX"
			},
			"href" : "https://api.spotify.com/v1/albums/4r1nB3kCpczB7c5MzNjCBX",
			"id" : "4r1nB3kCpczB7c5MzNjCBX",
			"images" : [ {
			  "height" : 640,
			  "url" : "https://i.scdn.co/image/ab67616d0000b273b7e56d5c61cffdcf0e7f2092",
			  "width" : 640
			}, {
			  "height" : 300,
			  "url" : "https://i.scdn.co/image/ab67616d00001e02b7e56d5c61cffdcf0e7f2092",
			  "width" : 300
			}, {
			  "height" : 64,
			  "url" : "https://i.scdn.co/image/ab67616d00004851b7e56d5c61cffdcf0e7f2092",
			  "width" : 64
			} ],
			"name" : "Felicita",
			"release_date" : "2014-07-03",
			"release_date_precision" : "day",
			"total_tracks" : 2,
			"type" : "album",
			"uri" : "spotify:album:4r1nB3kCpczB7c5MzNjCBX"
		  },
		  "artists" : [ {
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/artist/6MnbF4D0Fv2zUVnf08Df5t"
			},
			"href" : "https://api.spotify.com/v1/artists/6MnbF4D0Fv2zUVnf08Df5t",
			"id" : "6MnbF4D0Fv2zUVnf08Df5t",
			"name" : "Nadezhda Kadysheva",
			"type" : "artist",
			"uri" : "spotify:artist:6MnbF4D0Fv2zUVnf08Df5t"
		  }, {
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/artist/2gpELJDuhOQi1DPcI0VFvw"
			},
			"href" : "https://api.spotify.com/v1/artists/2gpELJDuhOQi1DPcI0VFvw",
			"id" : "2gpELJDuhOQi1DPcI0VFvw",
			"name" : "Al Bano",
			"type" : "artist",
			"uri" : "spotify:artist:2gpELJDuhOQi1DPcI0VFvw"
		  } ],
		  "available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
		  "disc_number" : 1,
		  "duration_ms" : 239984,
		  "explicit" : false,
		  "external_ids" : {
			"isrc" : "RUB731406450"
		  },
		  "external_urls" : {
			"spotify" : "https://open.spotify.com/track/2sP5VgY8PWb6c9DhgZEpSv"
		  },
		  "href" : "https://api.spotify.com/v1/tracks/2sP5VgY8PWb6c9DhgZEpSv",
		  "id" : "2sP5VgY8PWb6c9DhgZEpSv",
		  "is_local" : false,
		  "name" : "Широка река",
		  "popularity" : 28,
		  "preview_url" : "https://p.scdn.co/mp3-preview/bea92b770e0905299bbbb28304bbff3419de3112?cid=e7b45730e3774355ae6a15e7e4d188da",
		  "track_number" : 1,
		  "type" : "track",
		  "uri" : "spotify:track:2sP5VgY8PWb6c9DhgZEpSv"
		}`,
	}
)
