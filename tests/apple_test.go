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
	appleAccessToken     = "sample_developer_token"
	appleAuthorization   = "Bearer " + appleAccessToken
	appleWebPlayerJSPath = "/assets/index~6da982354d.js"

	appleTrackID        = "us-1559885421"
	appleTrackCatalogID = "1559885421"
	appleAlbumID        = "us-1559885420"
	appleAlbumCatalogID = "1559885420"
	appleMissingID      = "us-0"
	appleMissingRawID   = "0"
	appleSearchArtist   = "Rick Astley"
	appleSearchTrack    = "Never Gonna Give You Up"
	appleSearchAlbum    = "Whenever You Need Somebody"
)

func TestAppleCatalogFetchTrack(t *testing.T) {
	var serverURL string
	server := newAppleFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/catalog/us/songs/" + appleTrackCatalogID,
		Status:  http.StatusOK,
		Fixture: "apple_fetch_track_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertAppleAPIRequest(t, r, serverURL)
		},
	})
	serverURL = server.URL
	defer server.Close()

	catalog := newAppleCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Apple, appleTrackID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:         appleTrackID,
		CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music124/v4/ce/6d/5b/ce6d5b48-8c36-b990-3b9c-81862fadb459/0859381157694.jpg/1200x1200bb.jpg",
		Title:      "Never Gonna Give You Up",
		Artist:     "Rick Astley",
		AlbumID:    "us-1559885420",
		AlbumTitle: "Whenever You Need Somebody",
		Duration:   214,
		ReleaseDate: streamnx.ReleaseDate{
			Year: 1987, Month: 7, Day: 27,
		},
		Provider: streamnx.Apple,
		URL:      "https://music.apple.com/us/album/never-gonna-give-you-up/1559885420?i=1559885421",
	}, got)
}

func TestAppleCatalogFetchTrackNotFound(t *testing.T) {
	var serverURL string
	server := newAppleFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/catalog/us/songs/" + appleMissingRawID,
		Status:  http.StatusNotFound,
		Fixture: "apple_fetch_track_404.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertAppleAPIRequest(t, r, serverURL)
		},
	})
	serverURL = server.URL
	defer server.Close()

	catalog := newAppleCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Apple, appleMissingID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestAppleCatalogFetchAlbum(t *testing.T) {
	var serverURL string
	server := newAppleFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/catalog/us/albums/" + appleAlbumCatalogID,
		Status:  http.StatusOK,
		Fixture: "apple_fetch_album_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertAppleAPIRequest(t, r, serverURL)
		},
	})
	serverURL = server.URL
	defer server.Close()

	catalog := newAppleCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Apple, appleAlbumID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:       appleAlbumID,
		CoverURL: "https://is1-ssl.mzstatic.com/image/thumb/Music124/v4/ce/6d/5b/ce6d5b48-8c36-b990-3b9c-81862fadb459/0859381157694.jpg/1200x1200bb.jpg",
		Title:    "Whenever You Need Somebody",
		Artist:   "Rick Astley",
		Label:    "BMG Rights Management (UK) Ltd.",
		ReleaseDate: streamnx.ReleaseDate{
			Year: 1987, Month: 11, Day: 12,
		},
		Provider: streamnx.Apple,
		URL:      "https://music.apple.com/us/album/whenever-you-need-somebody/1559885420",
		TrackIDs: []string{
			"us-1559885421",
			"us-1559885422",
			"us-1559885423",
			"us-1559885424",
			"us-1559885425",
			"us-1559885426",
			"us-1559885427",
			"us-1559885428",
			"us-1559885429",
			"us-1559885430",
		},
	}, got)
}

func TestAppleCatalogFetchAlbumNotFound(t *testing.T) {
	var serverURL string
	server := newAppleFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/catalog/us/albums/" + appleMissingRawID,
		Status:  http.StatusNotFound,
		Fixture: "apple_fetch_album_404.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertAppleAPIRequest(t, r, serverURL)
		},
	})
	serverURL = server.URL
	defer server.Close()

	catalog := newAppleCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Apple, appleMissingID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestAppleCatalogSearchTracks(t *testing.T) {
	var serverURL string
	server := newAppleFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/catalog/us/search",
		Status:  http.StatusOK,
		Fixture: "apple_search_tracks_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertAppleAPIRequest(t, r, serverURL)
			assertAppleSearchQuery(t, r, appleSearchArtist+" "+appleSearchTrack)
		},
	})
	serverURL = server.URL
	defer server.Close()

	catalog := newAppleCatalog(t, server.URL)
	got, err := catalog.SearchTracks(t.Context(), streamnx.Apple, streamnx.SearchQuery{
		Artist: appleSearchArtist,
		Title:  appleSearchTrack,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchTrack{
		{
			ID:         "us-1559885421",
			AlbumID:    "us-1559885420",
			AlbumTitle: "Whenever You Need Somebody",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music124/v4/ce/6d/5b/ce6d5b48-8c36-b990-3b9c-81862fadb459/0859381157694.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up/1559885420?i=1559885421",
		},
		{
			ID:         "us-1612648434",
			AlbumID:    "us-1612648318",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music126/v4/b3/21/d3/b321d3e4-edfe-124b-d0cd-a64ad1df3290/4050538793840.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up (Cake Mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up-cake-mix/1612648318?i=1612648434",
		},
		{
			ID:         "us-1600742155",
			AlbumID:    "us-1600742151",
			AlbumTitle: "Supermix (DJ Mix)",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music126/v4/35/6b/88/356b88a4-fab2-87cb-9238-0109c40ebfe5/DF_Supermix_4000.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up (Mixed)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up-mixed/1600742151?i=1600742155",
		},
		{
			ID:         "us-1612648427",
			AlbumID:    "us-1612648318",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music126/v4/b3/21/d3/b321d3e4-edfe-124b-d0cd-a64ad1df3290/4050538793840.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up (Phil Harding 12\" Mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up-phil-harding-12-mix/1612648318?i=1612648427",
		},
		{
			ID:         "us-1612648440",
			AlbumID:    "us-1612648318",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music126/v4/b3/21/d3/b321d3e4-edfe-124b-d0cd-a64ad1df3290/4050538793840.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up (Escape to New York Mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up-escape-to-new-york-mix/1612648318?i=1612648440",
		},
		{
			ID:         "us-1503050130",
			AlbumID:    "us-1503050108",
			AlbumTitle: "Never Gonna Give You Up - Single",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music113/v4/0a/cb/04/0acb041e-0d6c-7e9e-329d-60f140e6eae5/195081167275.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "RadioClub",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up/1503050108?i=1503050130",
		},
		{
			ID:         "us-1446924456",
			AlbumID:    "us-1446924085",
			AlbumTitle: "Anthology",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music122/v4/99/be/51/99be51d9-08c1-2a51-b298-813cb644e8c8/16UMGIM70815.rgb.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Musical Youth",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up/1446924085?i=1446924456",
		},
		{
			ID:         "us-1559804346",
			AlbumID:    "us-1559804099",
			AlbumTitle: "Never Gonna Give You Up - Single",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music114/v4/41/50/39/4150396d-4bff-43dc-338f-806aefcf8374/196006358976.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Home Free",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up/1559804099?i=1559804346",
		},
		{
			ID:         "us-1612648444",
			AlbumID:    "us-1612648318",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://is1-ssl.mzstatic.com/image/thumb/Music126/v4/b3/21/d3/b321d3e4-edfe-124b-d0cd-a64ad1df3290/4050538793840.jpg/1200x1200bb.jpg",
			Title:      "Never Gonna Give You Up (Instrumental)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Apple,
			URL:        "https://music.apple.com/us/album/never-gonna-give-you-up-instrumental/1612648318?i=1612648444",
		},
	}, got)
}

func TestAppleCatalogSearchAlbums(t *testing.T) {
	var serverURL string
	server := newAppleFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/catalog/us/search",
		Status:  http.StatusOK,
		Fixture: "apple_search_albums_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertAppleAPIRequest(t, r, serverURL)
			assertAppleSearchQuery(t, r, appleSearchArtist+" "+appleSearchAlbum)
		},
	})
	serverURL = server.URL
	defer server.Close()

	catalog := newAppleCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Apple, streamnx.SearchQuery{
		Artist: appleSearchArtist,
		Title:  appleSearchAlbum,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchAlbum{
		{
			ID:       "us-1612648318",
			CoverURL: "https://is1-ssl.mzstatic.com/image/thumb/Music126/v4/b3/21/d3/b321d3e4-edfe-124b-d0cd-a64ad1df3290/4050538793840.jpg/1200x1200bb.jpg",
			Title:    "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			Artist:   "Rick Astley",
			Provider: streamnx.Apple,
			URL:      "https://music.apple.com/us/album/whenever-you-need-somebody-deluxe-edition-2022-remaster/1612648318",
		},
		{
			ID:       "us-1773292758",
			CoverURL: "https://is1-ssl.mzstatic.com/image/thumb/Music221/v4/db/9e/2a/db9e2ae0-cb9f-f2a9-2774-de399dff2580/4099964133639.jpg/1200x1200bb.jpg",
			Title:    "The Best Of Me: Never Edition",
			Artist:   "Rick Astley",
			Provider: streamnx.Apple,
			URL:      "https://music.apple.com/us/album/the-best-of-me-never-edition/1773292758",
		},
		{
			ID:       "us-1485596039",
			CoverURL: "https://is1-ssl.mzstatic.com/image/thumb/Music123/v4/17/3c/cb/173ccb70-6d48-24cc-3b5a-ff41560a0fda/4050538545074.jpg/1200x1200bb.jpg",
			Title:    "The Best of Me",
			Artist:   "Rick Astley",
			Provider: streamnx.Apple,
			URL:      "https://music.apple.com/us/album/the-best-of-me/1485596039",
		},
		{
			ID:       "us-784225734",
			CoverURL: "https://is1-ssl.mzstatic.com/image/thumb/Music6/v4/8d/91/87/8d9187d0-e7c2-16ee-ce0f-87ef8e5fea95/888003796942.jpg/1200x1200bb.jpg",
			Title:    "Whenever You Need Somebody (In the Style of Rick Astley) [Karaoke Version] - Single",
			Artist:   "Ameritz - Karaoke",
			Provider: streamnx.Apple,
			URL:      "https://music.apple.com/us/album/whenever-you-need-somebody-in-the-style-of-rick/784225734",
		},
		{
			ID:       "us-1559523357",
			CoverURL: "https://is1-ssl.mzstatic.com/image/thumb/Music114/v4/69/5b/e3/695be316-9ddf-7262-177c-e37edd599602/0888880777768.jpg/1200x1200bb.jpg",
			Title:    "3 Originals",
			Artist:   "Rick Astley",
			Provider: streamnx.Apple,
			URL:      "https://music.apple.com/us/album/3-originals/1559523357",
		},
	}, got)
}

func TestAppleCatalogRejectsUnsupportedOperation(t *testing.T) {
	catalog := newAppleCatalog(t, "http://127.0.0.1")

	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.Apple, "sample")

	require.Zero(t, gotType)
	require.Zero(t, gotID)
	require.ErrorIs(t, err, streamnx.ErrUnsupportedOperation)
}

func newAppleFixtureServer(t *testing.T, apiRoute fixtures.Route) *httptest.Server {
	t.Helper()

	webPlayerRoute := fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Status:  http.StatusOK,
		Fixture: "apple_web_player_200.html",
	}
	webPlayerJSRoute := fixtures.Route{
		Method:  http.MethodGet,
		Path:    appleWebPlayerJSPath,
		Status:  http.StatusOK,
		Fixture: "apple_web_player_js_200.js",
	}

	return fixtures.NewServer(t, webPlayerRoute, webPlayerJSRoute, apiRoute)
}

func newAppleCatalog(t *testing.T, serverURL string) *streamnx.Catalog {
	t.Helper()

	catalog, err := streamnx.NewCatalog(streamnx.WithApple(
		streamnx.WithAppleWebPlayerURL(serverURL),
		streamnx.WithAppleAPIURL(serverURL),
	))
	require.NoError(t, err)

	return catalog
}

func assertAppleAPIRequest(t *testing.T, r *http.Request, serverURL string) {
	t.Helper()

	require.Equal(t, appleAuthorization, r.Header.Get("Authorization"))
	require.Equal(t, serverURL, r.Header.Get("Origin"))
}

func assertAppleSearchQuery(t *testing.T, r *http.Request, term string) {
	t.Helper()

	q := r.URL.Query()
	require.Equal(t, term, q.Get("term"))
	require.Equal(t, "21", q.Get("limit"))
	require.Equal(t, "web", q.Get("platform"))
	require.Equal(t, "map", q.Get("format[resources]"))
	require.Equal(t, "albums", q.Get("relate[songs]"))
}
