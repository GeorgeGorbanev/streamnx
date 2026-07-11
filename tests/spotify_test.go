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
	spotifyClientID        = "sampleClientID"
	spotifyClientSecret    = "sampleClientSecret"
	spotifyAuthorization   = "Bearer " + spotifyAccessToken
	spotifyAccessToken     = "sample_access_token"
	spotifyBasicCredential = "Basic c2FtcGxlQ2xpZW50SUQ6c2FtcGxlQ2xpZW50U2VjcmV0"

	spotifyTrackID        = "4cOdK2wGLETKBW3PvgPWqT"
	spotifyAlbumID        = "5Z9iiGl2FcIfa3BMiv6OIw"
	spotifySearchArtist   = "Rick Astley"
	spotifySearchTrack    = "Never Gonna Give You Up"
	spotifySearchAlbum    = "Whenever You Need Somebody"
	spotifyMissingTrackID = "0000000000000000000000"
	spotifyMissingAlbumID = "0000000000000000000000"
)

func TestSpotifyCatalogFetchTrack(t *testing.T) {
	server := newSpotifyFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/tracks/" + spotifyTrackID,
		Status:  http.StatusOK,
		Fixture: "spotify_fetch_track_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, spotifyAuthorization, r.Header.Get("Authorization"))
		},
	})
	defer server.Close()

	catalog := newSpotifyCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Spotify, spotifyTrackID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:         spotifyTrackID,
		CoverURL:   "https://i.scdn.co/image/ab67616d0000b273baf89eb11ec7c657805d2da0",
		Title:      "Never Gonna Give You Up",
		Artist:     "Rick Astley",
		AlbumID:    spotifyAlbumID,
		AlbumTitle: "Whenever You Need Somebody",
		Duration:   214,
		ReleaseDate: streamnx.ReleaseDate{
			Year: 1987, Month: 11, Day: 12,
		},
		Provider: streamnx.Spotify,
		URL:      "https://open.spotify.com/track/" + spotifyTrackID,
	}, got)
}

func TestSpotifyCatalogFetchTrackNotFound(t *testing.T) {
	server := newSpotifyFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/tracks/" + spotifyMissingTrackID,
		Status:  http.StatusNotFound,
		Fixture: "spotify_fetch_track_404.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, spotifyAuthorization, r.Header.Get("Authorization"))
		},
	})
	defer server.Close()

	catalog := newSpotifyCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Spotify, spotifyMissingTrackID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestSpotifyCatalogFetchAlbum(t *testing.T) {
	server := newSpotifyFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/albums/" + spotifyAlbumID,
		Status:  http.StatusOK,
		Fixture: "spotify_fetch_album_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, spotifyAuthorization, r.Header.Get("Authorization"))
		},
	})
	defer server.Close()

	catalog := newSpotifyCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Spotify, spotifyAlbumID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:       spotifyAlbumID,
		CoverURL: "https://i.scdn.co/image/ab67616d0000b273baf89eb11ec7c657805d2da0",
		Title:    "Whenever You Need Somebody",
		Artist:   "Rick Astley",
		Label:    "BMG Rights Management (UK) Ltd",
		ReleaseDate: streamnx.ReleaseDate{
			Year: 1987, Month: 11, Day: 12,
		},
		Provider: streamnx.Spotify,
		Creator:  "BMG Rights Management (UK) Ltd",
		URL:      "https://open.spotify.com/album/" + spotifyAlbumID,
		TrackIDs: []string{
			"4PTG3Z6ehGkBFwjybzWkR8",
			"3X4r1d2VDYjgPYjAEjbEM7",
			"00isIFJWVpXIQ8HkGICSQp",
			"29K9hk3LXfbg12O2tH355c",
			"7DljnkkhQ9fiByyuOKZVgZ",
			"7EyrWqVQ3bAYGWPOFG6HMZ",
			"2vT7vbzj0IpWSn7u1893Cr",
			"7KzUl0WI3PPorUqwxg0tlD",
			"2H8A6GdlmOcG3iahhOmZNZ",
			"7bQ3PzSOamBuf77yAfbSUc",
		},
	}, got)
}

func TestSpotifyCatalogFetchAlbumNotFound(t *testing.T) {
	server := newSpotifyFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/albums/" + spotifyMissingAlbumID,
		Status:  http.StatusNotFound,
		Fixture: "spotify_fetch_album_404.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, spotifyAuthorization, r.Header.Get("Authorization"))
		},
	})
	defer server.Close()

	catalog := newSpotifyCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Spotify, spotifyMissingAlbumID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestSpotifyCatalogSearchTracks(t *testing.T) {
	server := newSpotifyFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/search",
		Status:  http.StatusOK,
		Fixture: "spotify_search_tracks_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, spotifyAuthorization, r.Header.Get("Authorization"))
			require.Equal(t, "artist:Rick Astley track:Never Gonna Give You Up", r.URL.Query().Get("q"))
			require.Equal(t, "track", r.URL.Query().Get("type"))
			require.Equal(t, "10", r.URL.Query().Get("limit"))
		},
	})
	defer server.Close()

	catalog := newSpotifyCatalog(t, server.URL)
	got, err := catalog.SearchTracks(t.Context(), streamnx.Spotify, streamnx.SearchQuery{
		Artist: spotifySearchArtist,
		Title:  spotifySearchTrack,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchTrack{
		{
			ID:         "4PTG3Z6ehGkBFwjybzWkR8",
			AlbumID:    "6eUW0wxWtzkFdaEFsTJto6",
			AlbumTitle: "Whenever You Need Somebody",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b27315ebbedaacef61af244262a8",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/4PTG3Z6ehGkBFwjybzWkR8",
		},
		{
			ID:         "4L7qMw8HI3vM57hHRMyb4Y",
			AlbumID:    "6CWdaSQN5rsdDkOrhFcZ0E",
			AlbumTitle: "The Best of Me",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b273ad292ed19701d9f5c8cf71b6",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/4L7qMw8HI3vM57hHRMyb4Y",
		},
		{
			ID:         "5AErEepoomKZK54EWm5i9a",
			AlbumID:    "4C4LvbYS0pxXLW5sGD9EK5",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b2733a67639779ccabd632e1a80e",
			Title:      "Never Gonna Give You Up - Cake Mix",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/5AErEepoomKZK54EWm5i9a",
		},
		{
			ID:         "2MPdnhqZpLjD97bFQxIZZY",
			AlbumID:    "58VYKotsIdQJVR6MrCZQSY",
			AlbumTitle: "80s Party Anthems",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b2735c37a644a590fc82490f2028",
			Title:      "Never Gonna Give You Up - 2022 Remaster",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/2MPdnhqZpLjD97bFQxIZZY",
		},
		{
			ID:         "0yNttAVwMr39qyODHNIkrY",
			AlbumID:    "4C4LvbYS0pxXLW5sGD9EK5",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b2733a67639779ccabd632e1a80e",
			Title:      "Never Gonna Give You Up - 2022 Remaster",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/0yNttAVwMr39qyODHNIkrY",
		},
		{
			ID:         "2jkJe0SyxAa9rpaRcbZcf1",
			AlbumID:    "4C4LvbYS0pxXLW5sGD9EK5",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b2733a67639779ccabd632e1a80e",
			Title:      "Never Gonna Give You Up - Pianoforte",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/2jkJe0SyxAa9rpaRcbZcf1",
		},
		{
			ID:         "17hHET630wWQ2NLEQeaRGC",
			AlbumID:    "4C4LvbYS0pxXLW5sGD9EK5",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b2733a67639779ccabd632e1a80e",
			Title:      "Never Gonna Give You Up - Phil Harding 12\" Mix",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/17hHET630wWQ2NLEQeaRGC",
		},
		{
			ID:         "3NdPVhbbubJGjrg9LWlc0N",
			AlbumID:    "6CWdaSQN5rsdDkOrhFcZ0E",
			AlbumTitle: "The Best of Me",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b273ad292ed19701d9f5c8cf71b6",
			Title:      "Never Gonna Give You Up - Pianoforte",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/3NdPVhbbubJGjrg9LWlc0N",
		},
		{
			ID:         "2XoYH6OudHoyQUKTlSNq2P",
			AlbumID:    "4C4LvbYS0pxXLW5sGD9EK5",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b2733a67639779ccabd632e1a80e",
			Title:      "Never Gonna Give You Up - Instrumental",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/2XoYH6OudHoyQUKTlSNq2P",
		},
		{
			ID:         "0hvbmkrDotGgNSbjdSEVUY",
			AlbumID:    "4C4LvbYS0pxXLW5sGD9EK5",
			AlbumTitle: "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			CoverURL:   "https://i.scdn.co/image/ab67616d0000b2733a67639779ccabd632e1a80e",
			Title:      "Never Gonna Give You Up - Escape from Newton Mix",
			Artist:     "Rick Astley",
			Provider:   streamnx.Spotify,
			URL:        "https://open.spotify.com/track/0hvbmkrDotGgNSbjdSEVUY",
		},
	}, got)
}

func TestSpotifyCatalogSearchAlbums(t *testing.T) {
	server := newSpotifyFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/v1/search",
		Status:  http.StatusOK,
		Fixture: "spotify_search_albums_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, spotifyAuthorization, r.Header.Get("Authorization"))
			require.Equal(t, "artist:Rick Astley album:Whenever You Need Somebody", r.URL.Query().Get("q"))
			require.Equal(t, "album", r.URL.Query().Get("type"))
			require.Equal(t, "10", r.URL.Query().Get("limit"))
		},
	})
	defer server.Close()

	catalog := newSpotifyCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Spotify, streamnx.SearchQuery{
		Artist: spotifySearchArtist,
		Title:  spotifySearchAlbum,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchAlbum{
		{
			ID:       "6eUW0wxWtzkFdaEFsTJto6",
			CoverURL: "https://i.scdn.co/image/ab67616d0000b27315ebbedaacef61af244262a8",
			Title:    "Whenever You Need Somebody",
			Artist:   "Rick Astley",
			Provider: streamnx.Spotify,
			URL:      "https://open.spotify.com/album/6eUW0wxWtzkFdaEFsTJto6",
		},
		{
			ID:       "4C4LvbYS0pxXLW5sGD9EK5",
			CoverURL: "https://i.scdn.co/image/ab67616d0000b2733a67639779ccabd632e1a80e",
			Title:    "Whenever You Need Somebody (Deluxe Edition - 2022 Remaster)",
			Artist:   "Rick Astley",
			Provider: streamnx.Spotify,
			URL:      "https://open.spotify.com/album/4C4LvbYS0pxXLW5sGD9EK5",
		},
	}, got)
}

func TestSpotifyCatalogRejectsUnsupportedOperation(t *testing.T) {
	catalog := newSpotifyCatalog(t, "http://127.0.0.1")

	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.Spotify, "sample")

	require.Zero(t, gotType)
	require.Zero(t, gotID)
	require.ErrorIs(t, err, streamnx.ErrUnsupportedOperation)
}

func newSpotifyFixtureServer(t *testing.T, apiRoute fixtures.Route) *httptest.Server {
	t.Helper()

	authRoute := fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/api/token",
		Status:  http.StatusOK,
		Fixture: "spotify_auth_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, spotifyBasicCredential, r.Header.Get("Authorization"))
			require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		},
	}

	return fixtures.NewServer(t, authRoute, apiRoute)
}

func newSpotifyCatalog(t *testing.T, serverURL string) *streamnx.Catalog {
	t.Helper()

	catalog, err := streamnx.NewCatalog(streamnx.WithSpotify(
		streamnx.SpotifyCredentials{
			ClientID:     spotifyClientID,
			ClientSecret: spotifyClientSecret,
		},
		streamnx.WithSpotifyAuthURL(serverURL),
		streamnx.WithSpotifyAPIURL(serverURL),
	))
	require.NoError(t, err)

	return catalog
}
