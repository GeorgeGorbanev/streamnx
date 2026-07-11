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
	deezerTrackID        = "3786023472"
	deezerAlbumID        = "901415162"
	deezerSearchArtist   = "Rick Astley"
	deezerSearchTrack    = "Never Gonna Give You Up"
	deezerSearchAlbum    = "Whenever You Need Somebody"
	deezerMissingTrackID = "0"
	deezerMissingAlbumID = "0"
	deezerTrackCloakID   = "14408104"
	deezerTrackCloakCode = "33x8xurn8M73S6diT7Mv1"
	deezerAlbumCloakCode = "33x8z5JqZdnmoIrEBsJDC"
)

func TestDeezerCatalogFetchTrack(t *testing.T) {
	server := newDeezerFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/track/" + deezerTrackID,
		Status:  http.StatusOK,
		Fixture: "deezer_fetch_track_200.json",
	})
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Deezer, deezerTrackID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:         deezerTrackID,
		CoverURL:   "https://cdn-images.dzcdn.net/images/cover/07de029f6cbc9ce63d9b6064f68b7455/1000x1000-000000-80-0-0.jpg",
		Title:      "Never Gonna Give You Up",
		Artist:     "Rick Astley",
		AlbumID:    deezerAlbumID,
		AlbumTitle: "Whenever You Need Somebody",
		Duration:   213,
		ReleaseDate: streamnx.ReleaseDate{
			Year: 2026, Month: 1, Day: 28,
		},
		Provider: streamnx.Deezer,
		URL:      "https://deezer.com/track/" + deezerTrackID,
	}, got)
}

func TestDeezerCatalogFetchTrackNotFound(t *testing.T) {
	server := newDeezerFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/track/" + deezerMissingTrackID,
		Status:  http.StatusOK,
		Fixture: "deezer_fetch_track_not_found_200.json",
	})
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Deezer, deezerMissingTrackID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestDeezerCatalogFetchAlbum(t *testing.T) {
	server := newDeezerFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/album/" + deezerAlbumID,
		Status:  http.StatusOK,
		Fixture: "deezer_fetch_album_200.json",
	})
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Deezer, deezerAlbumID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:       deezerAlbumID,
		CoverURL: "https://cdn-images.dzcdn.net/images/cover/07de029f6cbc9ce63d9b6064f68b7455/1000x1000-000000-80-0-0.jpg",
		Title:    "Whenever You Need Somebody",
		Artist:   "Rick Astley",
		Label:    "BMG Rights Management (UK) Ltd.",
		ReleaseDate: streamnx.ReleaseDate{
			Year: 1987, Month: 11, Day: 12,
		},
		Provider: streamnx.Deezer,
		Creator:  "BMG Rights Management (UK) Ltd.",
		URL:      "https://deezer.com/album/" + deezerAlbumID,
		TrackIDs: []string{
			"3786023472",
			"3786023492",
			"3786023512",
			"3786023532",
			"3786023552",
			"3786023572",
			"3786023592",
			"3786023602",
			"3786023612",
			"3786023622",
		},
	}, got)
}

func TestDeezerCatalogFetchAlbumNotFound(t *testing.T) {
	server := newDeezerFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/album/" + deezerMissingAlbumID,
		Status:  http.StatusOK,
		Fixture: "deezer_fetch_album_not_found_200.json",
	})
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Deezer, deezerMissingAlbumID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestDeezerCatalogFetchTrackCloak(t *testing.T) {
	server := newDeezerFixtureServer(t,
		fixtures.Route{
			Method:  http.MethodHead,
			Path:    "/s/" + deezerTrackCloakCode,
			Status:  http.StatusMovedPermanently,
			Fixture: "deezer_cloak_track_301.json",
			Headers: map[string]string{
				"Location": "https://link.deezer.com/?dest=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F" +
					deezerTrackCloakID +
					"%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-" +
					deezerTrackCloakID +
					"%26deferredFl%3D1%26universal_link%3D1",
			},
		},
	)
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.Deezer, deezerTrackCloakCode)

	require.NoError(t, err)
	require.Equal(t, streamnx.ReleaseTypeTrack, gotType)
	require.Equal(t, deezerTrackCloakID, gotID)
}

func TestDeezerCatalogFetchAlbumCloak(t *testing.T) {
	server := newDeezerFixtureServer(t,
		fixtures.Route{
			Method:  http.MethodHead,
			Path:    "/s/" + deezerAlbumCloakCode,
			Status:  http.StatusMovedPermanently,
			Fixture: "deezer_cloak_album_301.json",
			Headers: map[string]string{
				"Location": "https://link.deezer.com/?dest=https%3A%2F%2Fwww.deezer.com%2Falbum%2F" +
					deezerAlbumID +
					"%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dalbum-" +
					deezerAlbumID +
					"%26deferredFl%3D1",
			},
		},
	)
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.Deezer, deezerAlbumCloakCode)

	require.NoError(t, err)
	require.Equal(t, streamnx.ReleaseTypeAlbum, gotType)
	require.Equal(t, deezerAlbumID, gotID)
}

func TestDeezerCatalogSearchTracks(t *testing.T) {
	server := newDeezerFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/search",
		Status:  http.StatusOK,
		Fixture: "deezer_search_tracks_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, `artist:"Rick Astley" track:"Never Gonna Give You Up"`, r.URL.Query().Get("q"))
		},
	})
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	got, err := catalog.SearchTracks(t.Context(), streamnx.Deezer, streamnx.SearchQuery{
		Artist: deezerSearchArtist,
		Title:  deezerSearchTrack,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchTrack{
		{
			ID:         "3786363472",
			AlbumID:    "901480252",
			AlbumTitle: "The Best of Me",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/1fcfca61ca4e05027612a1af865b2e03/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/3786363472",
		},
		{
			ID:         "10794700",
			AlbumID:    "986064",
			AlbumTitle: "World's Greatest 80's Disco - The Only 80's Disco Album You'll Ever Need",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/0a0e66cf7d26af9d80716073930cf9e5/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley Tribute Band",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/10794700",
		},
		{
			ID:         "3786024802",
			AlbumID:    "901415272",
			AlbumTitle: "PWL Extended: Big Hits & Surprises (Vols. 1 & 2)",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/2e1f422c0d1883d42d121f396c6bdb14/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up (Cake Mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/3786024802",
		},
		{
			ID:         "3786363752",
			AlbumID:    "901480252",
			AlbumTitle: "The Best of Me",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/1fcfca61ca4e05027612a1af865b2e03/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up (Pianoforte)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/3786363752",
		},
		{
			ID:         "3809420952",
			AlbumID:    "907102572",
			AlbumTitle: "The Hit Factory Ultimate Collection",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/0276d581675201173472fdbc86ce60f8/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up (Escape to NY Mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/3809420952",
		},
		{
			ID:         "3816551312",
			AlbumID:    "909591702",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition / Remastered 2022)",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/f93d54aa0c8411c092d93c65546a8f41/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up (Phil Harding 12\" Mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/3816551312",
		},
		{
			ID:         "3816551402",
			AlbumID:    "909591702",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition / Remastered 2022)",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/f93d54aa0c8411c092d93c65546a8f41/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up (Escape from Newton Mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/3816551402",
		},
		{
			ID:         "3816551472",
			AlbumID:    "909591702",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition / Remastered 2022)",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/f93d54aa0c8411c092d93c65546a8f41/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up (Instrumental)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/3816551472",
		},
		{
			ID:         "87155985",
			AlbumID:    "8795315",
			AlbumTitle: "Where Are You on the Ruby, Pt. 2",
			CoverURL:   "https://cdn-images.dzcdn.net/images/cover/3ccc729302c9c96ae49dfd8a8238f6f6/1000x1000-000000-80-0-0.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Hollywood",
			Provider:   streamnx.Deezer,
			URL:        "https://deezer.com/track/87155985",
		},
	}, got)
}

func TestDeezerCatalogSearchAlbums(t *testing.T) {
	server := newDeezerFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/search/album",
		Status:  http.StatusOK,
		Fixture: "deezer_search_albums_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, `artist:"Rick Astley" album:"Whenever You Need Somebody"`, r.URL.Query().Get("q"))
		},
	})
	defer server.Close()

	catalog := newDeezerCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Deezer, streamnx.SearchQuery{
		Artist: deezerSearchArtist,
		Title:  deezerSearchAlbum,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchAlbum{
		{
			ID:       "901415162",
			CoverURL: "https://cdn-images.dzcdn.net/images/cover/07de029f6cbc9ce63d9b6064f68b7455/1000x1000-000000-80-0-0.jpg",
			Title:    "Whenever You Need Somebody",
			Artist:   "Rick Astley",
			Provider: streamnx.Deezer,
			URL:      "https://deezer.com/album/901415162",
		},
		{
			ID:       "909591702",
			CoverURL: "https://cdn-images.dzcdn.net/images/cover/f93d54aa0c8411c092d93c65546a8f41/1000x1000-000000-80-0-0.jpg",
			Title:    "Whenever You Need Somebody (Deluxe Edition / Remastered 2022)",
			Artist:   "Rick Astley",
			Provider: streamnx.Deezer,
			URL:      "https://deezer.com/album/909591702",
		},
	}, got)
}

func newDeezerFixtureServer(t *testing.T, apiRoutes ...fixtures.Route) *httptest.Server {
	t.Helper()

	return fixtures.NewServer(t, apiRoutes...)
}

func newDeezerCatalog(t *testing.T, serverURL string) *streamnx.Catalog {
	t.Helper()

	catalog, err := streamnx.NewCatalog(streamnx.WithDeezer(
		streamnx.WithDeezerAPIURL(serverURL),
		streamnx.WithDeezerCloakBaseURL(serverURL),
	))
	require.NoError(t, err)

	return catalog
}
