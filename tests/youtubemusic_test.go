package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2"
	"github.com/GeorgeGorbanev/streamnx/v2/tests/fixtures"
)

const (
	youtubeMusicTrackID        = "lYBUbBu4W08"
	youtubeMusicAlbumID        = "MPREb_dcYZhAh5urI"
	youtubeMusicArtist         = "Rick Astley"
	youtubeMusicTitle          = "Never Gonna Give You Up"
	youtubeMusicAlbum          = "Whenever You Need Somebody"
	youtubeMusicMissingTrackID = "00000000000"
)

func TestYoutubeMusicCatalogFetchTrack(t *testing.T) {
	server := newYoutubeMusicFixtureServer(t,
		fixtures.Route{
			Method:  http.MethodPost,
			Path:    "/player",
			Query:   map[string]string{"alt": "json"},
			Status:  http.StatusOK,
			Fixture: "youtubemusic_fetch_track_200.json",
			Assert:  assertYoutubeMusicBody("video_id", youtubeMusicTrackID),
		},
		fixtures.Route{
			Method:  http.MethodPost,
			Path:    "/next",
			Query:   map[string]string{"alt": "json"},
			Status:  http.StatusOK,
			Fixture: "youtubemusic_fetch_track_next_200.json",
			Assert:  assertYoutubeMusicBody("videoId", youtubeMusicTrackID),
		},
	)
	defer server.Close()

	got, err := newYoutubeMusicCatalog(t, server.URL).FetchTrack(
		t.Context(),
		streamnx.YoutubeMusic,
		youtubeMusicTrackID,
	)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:          youtubeMusicTrackID,
		Title:       youtubeMusicTitle,
		Artist:      youtubeMusicArtist,
		AlbumID:     youtubeMusicAlbumID,
		AlbumTitle:  youtubeMusicAlbum,
		URL:         "https://music.youtube.com/watch?v=" + youtubeMusicTrackID,
		CoverURL:    "https://yt3.googleusercontent.com/eC9DfRcYSk4FE-fvDCJSu_4xsKdVMKxwmFTYFZwP8OqB7R4TKxAjKoR-Kp1lXeRi2WddPFYulSte4eW-=w544-h544-l90-rj",
		Duration:    214,
		ReleaseDate: streamnx.ReleaseDate{Year: 1987},
		Provider:    streamnx.YoutubeMusic,
		Creator:     "Rick Astley - Topic",
		Description: "Rick Astley",
	}, got)
}

func TestYoutubeMusicCatalogFetchTrackNotFound(t *testing.T) {
	server := newYoutubeMusicFixtureServer(t, fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/player",
		Query:   map[string]string{"alt": "json"},
		Status:  http.StatusOK,
		Fixture: "youtubemusic_fetch_track_not_found_200.json",
	})
	defer server.Close()

	got, err := newYoutubeMusicCatalog(t, server.URL).FetchTrack(
		t.Context(),
		streamnx.YoutubeMusic,
		youtubeMusicMissingTrackID,
	)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestYoutubeMusicCatalogFetchTrackScrapingBlocked(t *testing.T) {
	server := newYoutubeMusicFixtureServer(t, fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/player",
		Query:   map[string]string{"alt": "json"},
		Status:  http.StatusOK,
		Fixture: "youtubemusic_fetch_track_login_required_200.json",
	})
	defer server.Close()

	got, err := newYoutubeMusicCatalog(t, server.URL).FetchTrack(
		t.Context(),
		streamnx.YoutubeMusic,
		"challenged",
	)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrScrapingBlocked)
	require.NotErrorIs(t, err, streamnx.ErrNotFound)
}

func TestYoutubeMusicCatalogFetchAlbum(t *testing.T) {
	server := newYoutubeMusicFixtureServer(t, fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/browse",
		Query:   map[string]string{"alt": "json"},
		Status:  http.StatusOK,
		Fixture: "youtubemusic_fetch_album_200.json",
		Assert:  assertYoutubeMusicBody("browseId", youtubeMusicAlbumID),
	})
	defer server.Close()

	got, err := newYoutubeMusicCatalog(t, server.URL).FetchAlbum(
		t.Context(),
		streamnx.YoutubeMusic,
		youtubeMusicAlbumID,
	)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:          youtubeMusicAlbumID,
		Title:       youtubeMusicAlbum,
		Artist:      youtubeMusicArtist,
		URL:         "https://music.youtube.com/browse/" + youtubeMusicAlbumID,
		CoverURL:    "https://yt3.googleusercontent.com/eC9DfRcYSk4FE-fvDCJSu_4xsKdVMKxwmFTYFZwP8OqB7R4TKxAjKoR-Kp1lXeRi2WddPFYulSte4eW-=w544-h544-l90-rj",
		ReleaseDate: streamnx.ReleaseDate{Year: 1987},
		Provider:    streamnx.YoutubeMusic,
		Description: `Whenever You Need Somebody is the debut studio album by English singer Rick Astley, released on 16 November 1987 by RCA Records. It is his highest-selling album and has sold 15.2 million copies worldwide. The album is listed as the 136th best-selling album in Spain and was the seventh best-selling album of 1987 in the United Kingdom. A remastered version, containing rare remixes and extended versions, was released on 12 April 2010.

From Wikipedia (https://en.wikipedia.org/wiki/Wheneve...) under Creative Commons Attribution CC-BY-SA 3.0 (https://creativecommons.org/licenses/...)`,
		TrackIDs: []string{
			"dQw4w9WgXcQ",
			"BeyEGebJ1l4",
			"yPYZpwSpKmA",
			"VBssFzSx6sI",
			"V8AvyCpCVJI",
			"X6uzlX51kbo",
			"czDxG-7SFsc",
			"tbmKG-vJank",
			"RXNFzrZTQy4",
			"rF1NHU_0NQE",
		},
	}, got)
}

func TestYoutubeMusicCatalogSearchTracks(t *testing.T) {
	server := newYoutubeMusicFixtureServer(t, fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/search",
		Query:   map[string]string{"alt": "json"},
		Status:  http.StatusOK,
		Fixture: "youtubemusic_search_tracks_200.json",
		Assert: assertYoutubeMusicSearch(
			t,
			"EgWKAQIIAWoMEA4QChADEAQQCRAF",
			youtubeMusicArtist+" – "+youtubeMusicTitle,
		),
	})
	defer server.Close()

	got, err := newYoutubeMusicCatalog(t, server.URL).SearchTracks(
		t.Context(),
		streamnx.YoutubeMusic,
		streamnx.SearchQuery{Artist: youtubeMusicArtist, Title: youtubeMusicTitle},
	)

	require.NoError(t, err)
	require.Equal(t, []string{
		"lYBUbBu4W08",
		"3BFTio5296w",
		"i_Q88T1HI_w",
		"QpbhSlce_ek",
		"sf4F2_r07M0",
		"n3jMY98KIEw",
		"raBobo3GZYA",
		"hdsehurnJt0",
		"pDWrOMKsnyY",
		"0k7BvzQRrOI",
		"HOhUtRlhrT0",
		"_uK77EalGq8",
		"UV5r1TMIDus",
		"131wf0e6ACk",
		"gzssOzSCvbM",
		"5u00ts7tDYQ",
		"-XakJZMWjmM",
		"FMvRxFkVubw",
		"32i5Tw8abmM",
		"0Om-vyXny3Y",
	}, searchTrackIDs(got))
	require.Equal(t, streamnx.SearchTrack{
		ID:         youtubeMusicTrackID,
		Title:      youtubeMusicTitle,
		Artist:     youtubeMusicArtist,
		AlbumID:    youtubeMusicAlbumID,
		AlbumTitle: youtubeMusicAlbum,
		URL:        "https://music.youtube.com/watch?v=" + youtubeMusicTrackID,
		CoverURL:   "https://yt3.googleusercontent.com/eC9DfRcYSk4FE-fvDCJSu_4xsKdVMKxwmFTYFZwP8OqB7R4TKxAjKoR-Kp1lXeRi2WddPFYulSte4eW-=w120-h120-l90-rj",
		Provider:   streamnx.YoutubeMusic,
	}, got[0])
	require.Equal(t, streamnx.SearchTrack{
		ID:         "0Om-vyXny3Y",
		Title:      "I'm Never Gonna Give You Up",
		Artist:     "Frank Stallone",
		AlbumID:    "MPREb_vo3gSsl5j9I",
		AlbumTitle: "Staying Alive (Original Motion Picture Soundtrack)",
		URL:        "https://music.youtube.com/watch?v=0Om-vyXny3Y",
		CoverURL:   "https://yt3.googleusercontent.com/_uMmjd35dKSmPC2YsdUPgq1DHBM8_fIaK2VQlT2t6hO76dULNn7To4X2Mr6lurXG0SClysnQPl5dnME9=w120-h120-l90-rj",
		Provider:   streamnx.YoutubeMusic,
	}, got[len(got)-1])
}

func TestYoutubeMusicCatalogSearchAlbums(t *testing.T) {
	server := newYoutubeMusicFixtureServer(t, fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/search",
		Query:   map[string]string{"alt": "json"},
		Status:  http.StatusOK,
		Fixture: "youtubemusic_search_albums_200.json",
		Assert: assertYoutubeMusicSearch(
			t,
			"EgWKAQIYAWoMEA4QChADEAQQCRAF",
			youtubeMusicArtist+" – "+youtubeMusicAlbum,
		),
	})
	defer server.Close()

	got, err := newYoutubeMusicCatalog(t, server.URL).SearchAlbums(
		t.Context(),
		streamnx.YoutubeMusic,
		streamnx.SearchQuery{Artist: youtubeMusicArtist, Title: youtubeMusicAlbum},
	)

	require.NoError(t, err)
	require.Equal(t, []string{
		"MPREb_dcYZhAh5urI",
		"MPREb_08pJEqlUwXy",
		"MPREb_LH7sUdCAbjP",
		"MPREb_9QrpwHdnH04",
		"MPREb_49zZLkGVsmC",
		"MPREb_7EcYNH9M1sG",
		"MPREb_BOEC90a78us",
		"MPREb_apkMMtu39hY",
		"MPREb_FHLw3Ybu2RR",
		"MPREb_KsdTwKv9zSE",
		"MPREb_fzOaMeZri49",
		"MPREb_ovBQgQI2Vwa",
		"MPREb_TiRZ1W121Y6",
		"MPREb_Y70NAu9CJu7",
		"MPREb_IbzQnyHsX6k",
		"MPREb_ful5YBihbYD",
		"MPREb_l269x3ITY2r",
		"MPREb_ZOKioXS9OCN",
		"MPREb_XF1sas6yjRc",
		"MPREb_DlXOfGV5gPv",
	}, searchAlbumIDs(got))
	require.Equal(t, streamnx.SearchAlbum{
		ID:       youtubeMusicAlbumID,
		Title:    youtubeMusicAlbum,
		Artist:   youtubeMusicArtist,
		URL:      "https://music.youtube.com/browse/" + youtubeMusicAlbumID,
		CoverURL: "https://yt3.googleusercontent.com/eC9DfRcYSk4FE-fvDCJSu_4xsKdVMKxwmFTYFZwP8OqB7R4TKxAjKoR-Kp1lXeRi2WddPFYulSte4eW-=w544-h544-l90-rj",
		Provider: streamnx.YoutubeMusic,
	}, got[0])
	require.Equal(t, streamnx.SearchAlbum{
		ID:       "MPREb_DlXOfGV5gPv",
		Title:    "80s Karaoke Hits, Vol. 7",
		Artist:   "A* Karaoke Jukebox",
		URL:      "https://music.youtube.com/browse/MPREb_DlXOfGV5gPv",
		CoverURL: "https://yt3.googleusercontent.com/YHJeYASGXDgt80HwoW4YrWjx7mt32csGrCKSm1e8l66FC5134_CHuT-NeAAbvLRHIQTEi1J6EZvDdhzwmg=w544-h544-l90-rj",
		Provider: streamnx.YoutubeMusic,
	}, got[len(got)-1])
}

func TestYoutubeMusicCatalogRejectsUnsupportedOperation(t *testing.T) {
	catalog := newYoutubeMusicCatalog(t, "http://127.0.0.1")

	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.YoutubeMusic, "sample")

	require.Empty(t, gotType)
	require.Empty(t, gotID)
	require.ErrorIs(t, err, streamnx.ErrUnsupportedOperation)
}

func newYoutubeMusicFixtureServer(t *testing.T, routes ...fixtures.Route) *httptest.Server {
	t.Helper()
	return fixtures.NewServer(t, routes...)
}

func newYoutubeMusicCatalog(t *testing.T, apiURL string) *streamnx.Catalog {
	t.Helper()
	catalog, err := streamnx.NewCatalog(streamnx.WithYoutubeMusic(
		streamnx.WithYoutubeMusicAPIURL(apiURL),
	))
	require.NoError(t, err)
	return catalog
}

func assertYoutubeMusicBody(key, value string) func(*testing.T, *http.Request) {
	return func(t *testing.T, r *http.Request) {
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, value, body[key])
		assertYoutubeMusicContext(t, body)
	}
}

func assertYoutubeMusicSearch(t *testing.T, params, query string) func(*testing.T, *http.Request) {
	t.Helper()
	return func(t *testing.T, r *http.Request) {
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, params, body["params"])
		require.Equal(t, query, body["query"])
		assertYoutubeMusicContext(t, body)
	}
}

func assertYoutubeMusicContext(t *testing.T, body map[string]any) {
	t.Helper()
	contextValue, ok := body["context"].(map[string]any)
	require.True(t, ok)
	clientValue, ok := contextValue["client"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "WEB_REMIX", clientValue["clientName"])
}

func searchTrackIDs(tracks []streamnx.SearchTrack) []string {
	ids := make([]string, len(tracks))
	for i, track := range tracks {
		ids[i] = track.ID
	}
	return ids
}

func searchAlbumIDs(albums []streamnx.SearchAlbum) []string {
	ids := make([]string, len(albums))
	for i, album := range albums {
		ids[i] = album.ID
	}
	return ids
}
