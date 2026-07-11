package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2"
	"github.com/GeorgeGorbanev/streamnx/v2/tests/fixtures"
)

const (
	yandexTrackID        = "609676"
	yandexAlbumID        = "14599266"
	yandexSearchArtist   = "Rick Astley"
	yandexSearchTrack    = "Never Gonna Give You Up"
	yandexSearchAlbum    = "Whenever You Need Somebody"
	yandexMissingTrackID = "987654321"
	yandexMissingAlbumID = "-1231231231"
)

func TestYandexCatalogFetchTrack(t *testing.T) {
	server := newYandexFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/tracks/" + yandexTrackID,
		Status:  http.StatusOK,
		Fixture: "yandex_fetch_track_200.json",
	})
	defer server.Close()

	catalog := newYandexCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Yandex, yandexTrackID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:         yandexTrackID,
		CoverURL:   "https://avatars.yandex.net/get-music-content/19999910/c1ed9300.a.33628686-3/1000x1000",
		Title:      "Never Gonna Give You Up",
		Artist:     "Rick Astley",
		AlbumID:    "33628686",
		AlbumTitle: "The Best Of Me: Never Edition",
		Duration:   213,
		ReleaseDate: streamnx.ReleaseDate{
			Year: 2019, Month: 12, Day: 29,
		},
		Provider: streamnx.Yandex,
		URL:      "https://music.yandex.com/album/33628686/track/" + yandexTrackID,
	}, got)
}

func TestYandexCatalogFetchTrackNotFound(t *testing.T) {
	server := newYandexFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/tracks/" + yandexMissingTrackID,
		Status:  http.StatusOK,
		Fixture: "yandex_fetch_track_not_found_200.json",
	})
	defer server.Close()

	catalog := newYandexCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Yandex, yandexMissingTrackID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestYandexCatalogFetchAlbum(t *testing.T) {
	server := newYandexFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/albums/" + yandexAlbumID + "/with-tracks",
		Status:  http.StatusOK,
		Fixture: "yandex_fetch_album_200.json",
	})
	defer server.Close()

	catalog := newYandexCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Yandex, yandexAlbumID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:       yandexAlbumID,
		CoverURL: "https://avatars.yandex.net/get-music-content/20031403/274495a5.a.14599266-4/1000x1000",
		Title:    "Whenever You Need Somebody",
		Artist:   "Rick Astley",
		Label:    "BMG Rights Management (UK)",
		ReleaseDate: streamnx.ReleaseDate{
			Year: 1987, Month: 11, Day: 12,
		},
		Provider: streamnx.Yandex,
		URL:      "https://music.yandex.com/album/" + yandexAlbumID,
		TrackIDs: []string{
			"609676",
			"648090",
			"609700",
			"647753",
			"648043",
			"648044",
			"648045",
			"648046",
			"648047",
			"633410",
		},
	}, got)
}

func TestYandexCatalogFetchAlbumNotFound(t *testing.T) {
	server := newYandexFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/albums/" + yandexMissingAlbumID + "/with-tracks",
		Status:  http.StatusNotFound,
		Fixture: "yandex_fetch_album_not_found_404.json",
	})
	defer server.Close()

	catalog := newYandexCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Yandex, yandexMissingAlbumID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestYandexCatalogSearchTracks(t *testing.T) {
	server := newYandexFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/search",
		Status:  http.StatusOK,
		Fixture: "yandex_search_tracks_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, yandexSearchArtist+" \u2013 "+yandexSearchTrack, r.URL.Query().Get("text"))
			require.Equal(t, "track", r.URL.Query().Get("type"))
			require.Equal(t, "0", r.URL.Query().Get("page"))
		},
	})
	defer server.Close()

	catalog := newYandexCatalog(t, server.URL)
	got, err := catalog.SearchTracks(t.Context(), streamnx.Yandex, streamnx.SearchQuery{
		Artist: yandexSearchArtist,
		Title:  yandexSearchTrack,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchTrack{
		{
			ID:         "609676",
			AlbumID:    "14599232",
			AlbumTitle: "3 Originals",
			CoverURL:   "https://avatars.yandex.net/get-music-content/18132539/5ef3a5a4.a.14599232-4/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/14599232/track/609676",
		},
		{
			ID:         "6228762",
			AlbumID:    "681909",
			AlbumTitle: "World's Greatest 80's Disco - The Only 80's Disco Album You'll Ever Need",
			CoverURL:   "https://avatars.yandex.net/get-music-content/42108/7f22dbdf.a.681909-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley Tribute Band",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/681909/track/6228762",
		},
		{
			ID:         "24917594",
			AlbumID:    "2924562",
			AlbumTitle: "Favourite 80's Karaoke Vol. 3",
			CoverURL:   "https://avatars.yandex.net/get-music-content/28589/9ad97f0e.a.2924562-1/1000x1000",
			Title:      "Never Gonna Give You Up Originally Performed By Rick Astley",
			Artist:     "Sunfly Karaoke",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/2924562/track/24917594",
		},
		{
			ID:         "54424835",
			AlbumID:    "9486773",
			AlbumTitle: "Back To 80 Retro Party",
			CoverURL:   "https://avatars.yandex.net/get-music-content/2358262/60e86876.a.9486773-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Disco Fever",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/9486773/track/54424835",
		},
		{
			ID:         "60143894",
			AlbumID:    "9266700",
			AlbumTitle: "The Ultimate Tribute To Rick Astley",
			CoverURL:   "https://avatars.yandex.net/get-music-content/1880735/bd4a82af.a.9266700-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Tutt",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/9266700/track/60143894",
		},
		{
			ID:         "4543246",
			AlbumID:    "512721",
			AlbumTitle: "School Reunion: The 80's",
			CoverURL:   "https://avatars.yandex.net/get-music-content/38044/296c9826.a.512721-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "It's a Cover Up",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/512721/track/4543246",
		},
		{
			ID:         "1107906",
			AlbumID:    "300175",
			AlbumTitle: "Music From Queer As Folk Series 2",
			CoverURL:   "https://avatars.yandex.net/get-music-content/41288/b3516c41.a.300175-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Jukebox Junction",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/300175/track/1107906",
		},
		{
			ID:         "19006424",
			AlbumID:    "2125142",
			AlbumTitle: "Zoom Karaoke Party, Vol. 46",
			CoverURL:   "https://avatars.yandex.net/get-music-content/34131/2e3cc405.a.2125142-1/1000x1000",
			Title:      "Never Gonna Give You up",
			Artist:     "Zoom Karaoke",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/2125142/track/19006424",
		},
		{
			ID:         "9649285",
			AlbumID:    "1026526",
			AlbumTitle: "100% 80's",
			CoverURL:   "https://avatars.yandex.net/get-music-content/32236/e55ce820.a.1026526-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Sing Karaoke Sing",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/1026526/track/9649285",
		},
		{
			ID:         "8334780",
			AlbumID:    "869139",
			AlbumTitle: "Karaoke Party: Vol. 5",
			CoverURL:   "https://avatars.yandex.net/get-music-content/42108/515bb214.a.869139-1/1000x1000",
			Title:      "Never Gonna Give You Up in the Style of Rick Astley",
			Artist:     "Sunfly Karaoke",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/869139/track/8334780",
		},
		{
			ID:         "3135519",
			AlbumID:    "332926",
			AlbumTitle: "1987 Karaoke Classics Volume 1",
			CoverURL:   "https://avatars.yandex.net/get-music-content/41288/96303662.a.332926-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "1980s Karaoke Band",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/332926/track/3135519",
		},
		{
			ID:         "147561609",
			AlbumID:    "40374075",
			AlbumTitle: "Metal songs No. 1",
			CoverURL:   "https://avatars.yandex.net/get-music-content/17730916/9965682b.a.40374075-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Crazy Metal",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/40374075/track/147561609",
		},
		{
			ID:         "17878763",
			AlbumID:    "1977346",
			AlbumTitle: "Pop AllStars - Party Classics, Vol. 3",
			CoverURL:   "https://avatars.yandex.net/get-music-content/41288/20f86819.a.1977346-1/1000x1000",
			Title:      "Never Gonna Give You up Originally Performed by Rick Astley",
			Artist:     "New Tribute Kings",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/1977346/track/17878763",
		},
		{
			ID:         "1360485",
			AlbumID:    "141764",
			AlbumTitle: "Karaoke Downloads - Disco Vol.9",
			CoverURL:   "https://avatars.yandex.net/get-music-content/34131/e00d51e1.a.141764-1/1000x1000",
			Title:      "Never Gonna Give You Up (In The Style Of Rick Astley)",
			Artist:     "Karaoke",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/141764/track/1360485",
		},
		{
			ID:         "49865450",
			AlbumID:    "6858321",
			AlbumTitle: "Never Gonna Give You Up",
			CoverURL:   "https://avatars.yandex.net/get-music-content/117546/4f90a024.a.6858321-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Chart Anthems",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/6858321/track/49865450",
		},
		{
			ID:         "7873543",
			AlbumID:    "276843",
			AlbumTitle: "Never Gonna Give You Up",
			CoverURL:   "https://avatars.yandex.net/get-music-content/34131/257d7737.a.276843-10/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Amazing Karaoke",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/276843/track/7873543",
		},
		{
			ID:         "31955383",
			AlbumID:    "3885034",
			AlbumTitle: "80's at the Piano Vol. 2",
			CoverURL:   "https://avatars.yandex.net/get-music-content/33216/c7328d70.a.3885034-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Lang Project",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/3885034/track/31955383",
		},
		{
			ID:         "5467705",
			AlbumID:    "602178",
			AlbumTitle: "Karaoke (In the Style of Rick Astley)",
			CoverURL:   "https://avatars.yandex.net/get-music-content/42108/e7b25e45.a.602178-1/1000x1000",
			Title:      "Never Gonna Give You Up (In the Style of Rick Astley)",
			Artist:     "Ameritz Karaoke Entertainment",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/602178/track/5467705",
		},
		{
			ID:         "50074066",
			AlbumID:    "6894920",
			AlbumTitle: "Never Gonna Give You Up",
			CoverURL:   "https://avatars.yandex.net/get-music-content/114728/fbe1f938.a.6894920-1/1000x1000",
			Title:      "Never Gonna Give You Up",
			Artist:     "Pop Anthems",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/6894920/track/50074066",
		},
		{
			ID:         "1148741",
			AlbumID:    "125923",
			AlbumTitle: "Karaoke Guys From The 80's And 90's Part 1",
			CoverURL:   "https://avatars.yandex.net/get-music-content/49707/7e593cba.a.125923-1/1000x1000",
			Title:      "Never Gonna Give You Up As Made Famous By: Rick Astley",
			Artist:     "Karaoke",
			Provider:   streamnx.Yandex,
			URL:        "https://music.yandex.com/album/125923/track/1148741",
		},
	}, got)
}

func TestYandexCatalogSearchAlbums(t *testing.T) {
	server := newYandexFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/search",
		Status:  http.StatusOK,
		Fixture: "yandex_search_albums_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, yandexSearchArtist+" \u2013 "+yandexSearchAlbum, r.URL.Query().Get("text"))
			require.Equal(t, "album", r.URL.Query().Get("type"))
			require.Equal(t, "0", r.URL.Query().Get("page"))
		},
	})
	defer server.Close()

	catalog := newYandexCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Yandex, streamnx.SearchQuery{
		Artist: yandexSearchArtist,
		Title:  yandexSearchAlbum,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchAlbum{
		{
			ID:       "14599266",
			CoverURL: "https://avatars.yandex.net/get-music-content/20031403/274495a5.a.14599266-4/1000x1000",
			Title:    "Whenever You Need Somebody",
			Artist:   "Rick Astley",
			Provider: streamnx.Yandex,
			URL:      "https://music.yandex.com/album/14599266",
		},
		{
			ID:       "1083828",
			CoverURL: "https://avatars.yandex.net/get-music-content/49876/af55045c.a.1083828-1/1000x1000",
			Title:    "Whenever You Need Somebody (In the Style of Rick Astley) - Single",
			Artist:   "Ameritz Digital Karaoke",
			Provider: streamnx.Yandex,
			URL:      "https://music.yandex.com/album/1083828",
		},
	}, got)
}

func newYandexFixtureServer(t *testing.T, route fixtures.Route) *httptest.Server {
	t.Helper()

	return fixtures.NewServer(t, route)
}

func newYandexCatalog(t *testing.T, serverURL string) *streamnx.Catalog {
	t.Helper()

	catalog, err := streamnx.NewCatalog(streamnx.WithYandex(
		streamnx.WithYandexAPIURL(serverURL),
	))
	require.NoError(t, err)

	return catalog
}
