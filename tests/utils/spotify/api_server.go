package spotify

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/stretchr/testify/require"
)

var (
	SampleTrack = spotify.Track{
		Album: spotify.Album{
			AlbumType: "album",
			Artists: []spotify.Artist{
				{
					ExternalURLs: map[string]string{"spotify": "https://spotify.com/artist1"},
					Href:         "https://api.spotify.com/v1/artists/artist1",
					ID:           "artist1",
					Name:         "Massive Attack",
					Type:         "artist",
					URI:          "spotify:artist:artist1",
				},
			},
			AvailableMarkets: []string{"US", "GB"},
			ExternalURLs:     map[string]string{"spotify": "https://spotify.com/album1"},
			Href:             "https://api.spotify.com/v1/albums/album1",
			ID:               "album1",
			Images: []spotify.Image{
				{Height: 640, URL: "https://image.url/1", Width: 640},
				{Height: 300, URL: "https://image.url/2", Width: 300},
			},
			Name:                 "Album One",
			ReleaseDate:          "2023-01-01",
			ReleaseDatePrecision: "day",
			Type:                 "album",
			URI:                  "spotify:album:album1",
		},
		Artists: []spotify.Artist{
			{
				ExternalURLs: map[string]string{"spotify": "https://spotify.com/artist1"},
				Href:         "https://api.spotify.com/v1/artists/artist1",
				ID:           "artist1",
				Name:         "Massive Attack",
				Type:         "artist",
				URI:          "spotify:artist:artist1",
			},
		},
		AvailableMarkets: []string{"US", "GB"},
		DiscNumber:       1,
		DurationMs:       240000,
		Explicit:         false,
		ExternalIDs:      map[string]string{"isrc": "US1234567890"},
		ExternalURLs:     map[string]string{"spotify": "https://spotify.com/track1"},
		Href:             "https://api.spotify.com/v1/tracks/track1",
		ID:               "track1",
		IsPlayable:       true,
		LinkedFrom:       nil,
		Name:             "Angel",
		Popularity:       80,
		PreviewURL:       "https://preview.url/track1",
		TrackNumber:      5,
		Type:             "track",
		URI:              "spotify:track:track1",
	}
	SampleTrackJSON = `{
        "album": {
            "album_type": "album",
            "artists": [
                {
                    "external_urls": {"spotify": "https://spotify.com/artist1"},
                    "href": "https://api.spotify.com/v1/artists/artist1",
                    "id": "artist1",
                    "name": "Massive Attack",
                    "type": "artist",
                    "uri": "spotify:artist:artist1"
                }
            ],
            "available_markets": ["US", "GB"],
            "external_urls": {"spotify": "https://spotify.com/album1"},
            "href": "https://api.spotify.com/v1/albums/album1",
            "id": "album1",
            "images": [
                {"height": 640, "url": "https://image.url/1", "width": 640},
                {"height": 300, "url": "https://image.url/2", "width": 300}
            ],
            "name": "Album One",
            "release_date": "2023-01-01",
            "release_date_precision": "day",
            "type": "album",
            "uri": "spotify:album:album1"
        },
        "artists": [
            {
                "external_urls": {"spotify": "https://spotify.com/artist1"},
                "href": "https://api.spotify.com/v1/artists/artist1",
                "id": "artist1",
                "name": "Massive Attack",
                "type": "artist",
                "uri": "spotify:artist:artist1"
            }
        ],
        "available_markets": ["US", "GB"],
        "disc_number": 1,
        "duration_ms": 240000,
        "explicit": false,
        "external_ids": {"isrc": "US1234567890"},
        "external_urls": {"spotify": "https://spotify.com/track1"},
        "href": "https://api.spotify.com/v1/tracks/track1",
        "id": "track1",
        "is_playable": true,
        "linked_from": null,
        "name": "Angel",
        "popularity": 80,
        "preview_url": "https://preview.url/track1",
        "track_number": 5,
        "type": "track",
        "uri": "spotify:track:track1"
    }`

	SampleYMusicNotFoundTrack = spotify.Track{
		Album: spotify.Album{
			AlbumType: "album",
			Artists: []spotify.Artist{
				{
					ExternalURLs: map[string]string{"spotify": "https://spotify.com/artist1"},
					Href:         "https://api.spotify.com/v1/artists/artist1",
					ID:           "artist1",
					Name:         "Not Y Music",
					Type:         "artist",
					URI:          "spotify:artist:artist1",
				},
			},
			AvailableMarkets: []string{"US", "GB"},
			ExternalURLs:     map[string]string{"spotify": "https://spotify.com/album1"},
			Href:             "https://api.spotify.com/v1/albums/album1",
			ID:               "album1",
			Images: []spotify.Image{
				{Height: 640, URL: "https://image.url/1", Width: 640},
				{Height: 300, URL: "https://image.url/2", Width: 300},
			},
			Name:                 "Album One",
			ReleaseDate:          "2023-01-01",
			ReleaseDatePrecision: "day",
			Type:                 "album",
			URI:                  "spotify:album:album1",
		},
		Artists: []spotify.Artist{
			{
				ExternalURLs: map[string]string{"spotify": "https://spotify.com/artist1"},
				Href:         "https://api.spotify.com/v1/artists/artist1",
				ID:           "artist1",
				Name:         "Not Y Music",
				Type:         "artist",
				URI:          "spotify:artist:artist1",
			},
		},
		AvailableMarkets: []string{"US", "GB"},
		DiscNumber:       1,
		DurationMs:       240000,
		Explicit:         false,
		ExternalIDs:      map[string]string{"isrc": "US1234567890"},
		ExternalURLs:     map[string]string{"spotify": "https://spotify.com/track1"},
		Href:             "https://api.spotify.com/v1/tracks/track1",
		ID:               "ymN07f0und",
		IsPlayable:       true,
		LinkedFrom:       nil,
		Name:             "Not y music",
		Popularity:       80,
		PreviewURL:       "https://preview.url/track1",
		TrackNumber:      5,
		Type:             "track",
		URI:              "spotify:track:track1",
	}
	SampleYMusicNotFoundTrackJSON = `{
        "album": {
            "album_type": "album",
            "artists": [
                {
                    "external_urls": {"spotify": "https://spotify.com/artist1"},
                    "href": "https://api.spotify.com/v1/artists/artist1",
                    "id": "artist1",
                    "name": "Not Y Music",
                    "type": "artist",
                    "uri": "spotify:artist:artist1"
                }
            ],
            "available_markets": ["US", "GB"],
            "external_urls": {"spotify": "https://spotify.com/album1"},
            "href": "https://api.spotify.com/v1/albums/album1",
            "id": "album1",
            "images": [
                {"height": 640, "url": "https://image.url/1", "width": 640},
                {"height": 300, "url": "https://image.url/2", "width": 300}
            ],
            "name": "Album One",
            "release_date": "2023-01-01",
            "release_date_precision": "day",
            "type": "album",
            "uri": "spotify:album:album1"
        },
        "artists": [
            {
                "external_urls": {"spotify": "https://spotify.com/artist1"},
                "href": "https://api.spotify.com/v1/artists/artist1",
                "id": "artist1",
                "name": "Not Y Music",
                "type": "artist",
                "uri": "spotify:artist:artist1"
            }
        ],
        "available_markets": ["US", "GB"],
        "disc_number": 1,
        "duration_ms": 240000,
        "explicit": false,
        "external_ids": {"isrc": "US1234567890"},
        "external_urls": {"spotify": "https://spotify.com/track1"},
        "href": "https://api.spotify.com/v1/tracks/track1",
        "id": "ymN07f0und",
        "is_playable": true,
        "linked_from": null,
        "name": "Not y music",
        "popularity": 80,
        "preview_url": "https://preview.url/track1",
        "track_number": 5,
        "type": "track",
        "uri": "spotify:track:track1"
    }`

	SampleTrackTitle         = "The Raveonettes – Come On Santa"
	SearchSampleQuery        = "artist:The Raveonettes track:Come On Santa"
	SearchSampleResponseJSON = `{
	  "tracks" : {
		"href" : "https://api.spotify.com/v1/search?query=artist%3AThe+Raveonettes+track%3ACome+On+Santa&type=track&offset=0&limit=20",
		"items" : [ {
		  "album" : {
			"album_type" : "single",
			"artists" : [ {
			  "external_urls" : {
				"spotify" : "https://open.spotify.com/artist/3LTXHU3DhiYzGIgF2PP8Q8"
			  },
			  "href" : "https://api.spotify.com/v1/artists/3LTXHU3DhiYzGIgF2PP8Q8",
			  "id" : "3LTXHU3DhiYzGIgF2PP8Q8",
			  "name" : "The Raveonettes",
			  "type" : "artist",
			  "uri" : "spotify:artist:3LTXHU3DhiYzGIgF2PP8Q8"
			} ],
			"available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/album/6ncXiBXuxWYybb9DsUVCIc"
			},
			"href" : "https://api.spotify.com/v1/albums/6ncXiBXuxWYybb9DsUVCIc",
			"id" : "6ncXiBXuxWYybb9DsUVCIc",
			"images" : [ {
			  "height" : 640,
			  "url" : "https://i.scdn.co/image/ab67616d0000b2730288c044e8eb80392f439803",
			  "width" : 640
			}, {
			  "height" : 300,
			  "url" : "https://i.scdn.co/image/ab67616d00001e020288c044e8eb80392f439803",
			  "width" : 300
			}, {
			  "height" : 64,
			  "url" : "https://i.scdn.co/image/ab67616d000048510288c044e8eb80392f439803",
			  "width" : 64
			} ],
			"name" : "Wishing You A Rave Christmas",
			"release_date" : "2008-11-25",
			"release_date_precision" : "day",
			"total_tracks" : 4,
			"type" : "album",
			"uri" : "spotify:album:6ncXiBXuxWYybb9DsUVCIc"
		  },
		  "artists" : [ {
			"external_urls" : {
			  "spotify" : "https://open.spotify.com/artist/3LTXHU3DhiYzGIgF2PP8Q8"
			},
			"href" : "https://api.spotify.com/v1/artists/3LTXHU3DhiYzGIgF2PP8Q8",
			"id" : "3LTXHU3DhiYzGIgF2PP8Q8",
			"name" : "The Raveonettes",
			"type" : "artist",
			"uri" : "spotify:artist:3LTXHU3DhiYzGIgF2PP8Q8"
		  } ],
		  "available_markets" : [ "AR", "AU", "AT", "BE", "BO", "BR", "BG", "CA", "CL", "CO", "CR", "CY", "CZ", "DK", "DO", "DE", "EC", "EE", "SV", "FI", "FR", "GR", "GT", "HN", "HK", "HU", "IS", "IE", "IT", "LV", "LT", "LU", "MY", "MT", "MX", "NL", "NZ", "NI", "NO", "PA", "PY", "PE", "PH", "PL", "PT", "SG", "SK", "ES", "SE", "CH", "TW", "TR", "UY", "US", "GB", "AD", "LI", "MC", "ID", "JP", "TH", "VN", "RO", "IL", "ZA", "SA", "AE", "BH", "QA", "OM", "KW", "EG", "MA", "DZ", "TN", "LB", "JO", "PS", "IN", "BY", "KZ", "MD", "UA", "AL", "BA", "HR", "ME", "MK", "RS", "SI", "KR", "BD", "PK", "LK", "GH", "KE", "NG", "TZ", "UG", "AG", "AM", "BS", "BB", "BZ", "BT", "BW", "BF", "CV", "CW", "DM", "FJ", "GM", "GE", "GD", "GW", "GY", "HT", "JM", "KI", "LS", "LR", "MW", "MV", "ML", "MH", "FM", "NA", "NR", "NE", "PW", "PG", "WS", "SM", "ST", "SN", "SC", "SL", "SB", "KN", "LC", "VC", "SR", "TL", "TO", "TT", "TV", "VU", "AZ", "BN", "BI", "KH", "CM", "TD", "KM", "GQ", "SZ", "GA", "GN", "KG", "LA", "MO", "MR", "MN", "NP", "RW", "TG", "UZ", "ZW", "BJ", "MG", "MU", "MZ", "AO", "CI", "DJ", "ZM", "CD", "CG", "IQ", "LY", "TJ", "VE", "ET", "XK" ],
		  "disc_number" : 1,
		  "duration_ms" : 221648,
		  "explicit" : false,
		  "external_ids" : {
			"isrc" : "USA560861700"
		  },
		  "external_urls" : {
			"spotify" : "https://open.spotify.com/track/1wsw6jqC9KuJfqSshK3437"
		  },
		  "href" : "https://api.spotify.com/v1/tracks/1wsw6jqC9KuJfqSshK3437",
		  "id" : "1wsw6jqC9KuJfqSshK3437",
		  "is_local" : false,
		  "name" : "Come On Santa",
		  "popularity" : 41,
		  "preview_url" : "https://p.scdn.co/mp3-preview/27e15bc702b938b3517581220d53578d92e9f708?cid=e7b45730e3774355ae6a15e7e4d188da",
		  "track_number" : 2,
		  "type" : "track",
		  "uri" : "spotify:track:1wsw6jqC9KuJfqSshK3437"
		} ],
		"limit" : 20,
		"next" : null,
		"offset" : 0,
		"previous" : null,
		"total" : 1
	  }
	}
	`
)

var tracksPathRe = regexp.MustCompile(`/v1/tracks/([a-zA-Z0-9]+)`)

func NewAPIServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, r.Header.Get("Authorization"), "Bearer mock_access_token")

		switch {
		case tracksPathRe.MatchString(r.URL.Path):
			tracksHandler(t)(w, r)
		case r.URL.Path == "/v1/search":
			searchHandler(t)(w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
}

func tracksHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		splitted := strings.Split(r.URL.Path, "/")
		trackID := splitted[len(splitted)-1]

		switch trackID {
		case SampleTrack.ID:
			_, err := w.Write([]byte(SampleTrackJSON))
			require.NoError(t, err)
		case SampleYMusicNotFoundTrack.ID:
			_, err := w.Write([]byte(SampleYMusicNotFoundTrackJSON))
			require.NoError(t, err)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func searchHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "track", r.URL.Query().Get("type"))

		q := r.URL.Query().Get("q")
		switch q {
		case SearchSampleQuery:
			_, err := w.Write([]byte(SearchSampleResponseJSON))
			require.NoError(t, err)
		}
	}
}
