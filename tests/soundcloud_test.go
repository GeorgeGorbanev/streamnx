package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2"
	"github.com/GeorgeGorbanev/streamnx/v2/tests/fixtures"
)

const (
	soundcloudSampleClientID = "QNR5nrdLOvApYERC8AOUr3VjRfHnLjle"

	soundcloudTrackID           = "rick-astley-official:never-gonna-give-you-up"
	soundcloudAlbumID           = "aikostar-music:whenever-you-need-somebody-4"
	soundcloudSearchArtist      = "Rick Astley"
	soundcloudSearchAlbumArtist = "Aiko Star"
	soundcloudSearchTrack       = "Never Gonna Give You Up"
	soundcloudSearchAlbum       = "Whenever You Need Somebody"
	soundcloudMissingTrackID    = "rick-astley-official:missing-never-gonna-give-you-up"
	soundcloudMissingAlbumID    = "rick-astley-official:missing-whenever-you-need-somebody"
)

func TestSoundcloudCatalogFetchTrack(t *testing.T) {
	server := newSoundcloudFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/rick-astley-official/never-gonna-give-you-up",
		Status:  http.StatusOK,
		Fixture: "soundcloud_fetch_track_200.html",
	})
	defer server.Close()

	catalog := newSoundcloudCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Soundcloud, soundcloudTrackID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:         soundcloudTrackID,
		CoverURL:   "https://i1.sndcdn.com/artworks-keikGz8pxIhJ-0-original.jpg",
		Title:      "Never Gonna Give You Up",
		Artist:     "Rick Astley",
		AlbumTitle: "Whenever You Need Somebody",
		Duration:   214,
		ReleaseDate: streamnx.ReleaseDate{
			Year: 2026, Month: 1, Day: 28,
		},
		Provider: streamnx.Soundcloud,
		Creator:  "Rick Astley",
		URL:      "https://soundcloud.com/rick-astley-official/never-gonna-give-you-up",
	}, got)
}

func TestSoundcloudCatalogFetchTrackNotFound(t *testing.T) {
	server := newSoundcloudFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/rick-astley-official/missing-never-gonna-give-you-up",
		Status:  http.StatusOK,
		Fixture: "soundcloud_fetch_track_not_found_200.html",
	})
	defer server.Close()

	catalog := newSoundcloudCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Soundcloud, soundcloudMissingTrackID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestSoundcloudCatalogFetchAlbum(t *testing.T) {
	server := newSoundcloudFixtureServer(t,
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/aikostar-music/sets/whenever-you-need-somebody-4",
			Status:  http.StatusOK,
			Fixture: "soundcloud_fetch_album_200.html",
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/tracks",
			Status:  http.StatusOK,
			Fixture: "soundcloud_fetch_album_tracks_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				require.Equal(t, url.Values{
					"client_id": {soundcloudSampleClientID},
					"ids":       {"2206838003"},
				}, r.URL.Query())
			},
		},
	)
	defer server.Close()

	catalog := newSoundcloudCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Soundcloud, soundcloudAlbumID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:       soundcloudAlbumID,
		CoverURL: "https://i1.sndcdn.com/artworks-e7f343c1-fdef-41b5-ab34-211a36b3ff39-0-original.jpg",
		Title:    "Whenever You Need Somebody",
		Artist:   "Aiko Star",
		ReleaseDate: streamnx.ReleaseDate{
			Year: 2025, Month: 11, Day: 14,
		},
		Provider: streamnx.Soundcloud,
		Creator:  "Aiko Star",
		URL:      "https://soundcloud.com/aikostar-music/sets/whenever-you-need-somebody-4",
		TrackIDs: []string{
			"aikostar-music:wild-boy",
			"aikostar-music:forget-me",
			"aikostar-music:edge-of-destruction",
			"aikostar-music:more-than-life",
			"aikostar-music:paper-cuts",
			"aikostar-music:all-night-long",
		},
	}, got)
}

func TestSoundcloudCatalogFetchAlbumNotFound(t *testing.T) {
	server := newSoundcloudFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/rick-astley-official/sets/missing-whenever-you-need-somebody",
		Status:  http.StatusOK,
		Fixture: "soundcloud_fetch_album_not_found_200.html",
	})
	defer server.Close()

	catalog := newSoundcloudCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Soundcloud, soundcloudMissingAlbumID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestSoundcloudCatalogSearchTracks(t *testing.T) {
	server := newSoundcloudFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/search/tracks",
		Status:  http.StatusOK,
		Fixture: "soundcloud_search_tracks_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, soundcloudSearchArtist+" "+soundcloudSearchTrack, r.URL.Query().Get("q"))
			require.Equal(t, soundcloudSampleClientID, r.URL.Query().Get("client_id"))
		},
	})
	defer server.Close()

	catalog := newSoundcloudCatalog(t, server.URL)
	got, err := catalog.SearchTracks(t.Context(), streamnx.Soundcloud, streamnx.SearchQuery{
		Artist: soundcloudSearchArtist,
		Title:  soundcloudSearchTrack,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchTrack{
		{
			ID:         "rick-astley-official:never-gonna-give-you-up",
			AlbumID:    "",
			AlbumTitle: "Whenever You Need Somebody",
			CoverURL:   "https://i1.sndcdn.com/artworks-keikGz8pxIhJ-0-original.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley",
			Provider:   streamnx.Soundcloud,
			Creator:    "Rick Astley",
			URL:        "https://soundcloud.com/rick-astley-official/never-gonna-give-you-up",
		},
		{
			ID:         "rick-astley-official:never-gonna-give-you-up-2022",
			AlbumID:    "",
			AlbumTitle: "Chart Toppers",
			CoverURL:   "https://i1.sndcdn.com/artworks-YE40JIgqUheZ-0-original.jpg",
			Title:      "Never Gonna Give You Up (Remastered 2022)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Soundcloud,
			Creator:    "Rick Astley",
			URL:        "https://soundcloud.com/rick-astley-official/never-gonna-give-you-up-2022",
		},
		{
			ID:         "timmohendriks:never-gonna-give-you-up-remix",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://i1.sndcdn.com/artworks-000244572658-tgywx3-original.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Timmo Hendriks Remix)",
			Artist:     "Timmo Hendriks",
			Provider:   streamnx.Soundcloud,
			Creator:    "Timmo Hendriks",
			URL:        "https://soundcloud.com/timmohendriks/never-gonna-give-you-up-remix",
			Description: "Free download: http://thehusk.ca/gate.asp?t=404988519\n\n►Timmo Hendriks" +
				":\n@timmohendriks\nwww.facebook.com/TimmoHendriksofficial\ntwitter.com/Ti" +
				"mmo_Hendriks\nwww.instagram.com/timmohendriks/\nwww.youtube.com/c/TimmoH" +
				"endriks\nSnapchat: Timmo_Hendriks\n\nArtwork by: www.facebook.com/AkralDe" +
				"sign/",
		},
		{
			ID:         "primeshocklive:rick-astley-never-gonna-give-you-up-primeshock-bootleg",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://i1.sndcdn.com/artworks-nhUhMzYnIzWgL9l7-Bu64NQ-original.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Primeshock Bootleg)",
			Artist:     "Primeshock",
			Provider:   streamnx.Soundcloud,
			Creator:    "Primeshock",
			URL: "https://soundcloud.com/primeshocklive/" +
				"rick-astley-never-gonna-give-you-up-primeshock-bootleg",
			Description: "After rickrolling some livestreams here and there, it's time to give a" +
				"way our bootleg of 'Never Gonna Give You Up' for free! Enjoy the Prime" +
				"shock-injected bootleg of this alltime classic by Rick Astley! 🤩⚡\n\n♫ " +
				"Free download: https://primeshocklive.com/music/nggyubootleg\n\nEnergy, " +
				"power, entertainment; keywords that describe Primeshock on their way t" +
				"o the top. Devouring the Hard Dance scene with a flood of energy, thei" +
				"r shockwave is expanding at the speed of sound!\n\nPrimeshock delivers c" +
				"razy powerful sets that will pump up crowds and blow up stages. In add" +
				"ition they represent the best Hardstyle in their monthly Powermode Pod" +
				"cast, broadcasted at Q-dance Radio.\n\nPrimeshock is ready to take over " +
				"the world of Hardstyle with their energy, positive vibes and powerful " +
				"music. Switch into Powermode! ⚡\n\n▼ Follow Primeshock:\nwww.primeshockli" +
				"ve.com\nwww.primeshocklive.com/instagram\nwww.primeshocklive.com/spotify" +
				"\nwww.primeshocklive.com/facebook\nwww.primeshocklive.com/soundcloud\nwww" +
				".primeshocklive.com/mixcloud\nwww.primeshocklive.com/youtube\nwww.primes" +
				"hocklive.com/twitter\n\n♫ Follow the Powermode Hardstyle Playlist on Spo" +
				"tify!\nwww.primeshocklive.com/powermode-spotify\n\n⌁ Join our 'Primeshock" +
				" | Plugged In' server on Discord!\nhttps://discord.gg/8UM2zrX",
		},
		{
			ID:         "djericfaria:eric-faria-oni-remix-rick-astley-never-gonna-give-you-up-out-soon",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://i1.sndcdn.com/artworks-000196755418-x4u4ma-original.jpg",
			Title:      "Eric Faria & Oni Remix - Rick Astley - Never Gonna Give You Up ------FREE DOWNLOAD",
			Artist:     "djericfaria",
			Provider:   streamnx.Soundcloud,
			Creator:    "Eric Faria",
			URL: "https://soundcloud.com/djericfaria/" +
				"eric-faria-oni-remix-rick-astley-never-gonna-give-you-up-out-soon",
			Description: "BPM STUDIOS : www.facebook.com/bpmstudiospt/\nOFFICIAL FACEBOOCK : www." +
				"facebook.com/EricFariaOficial\nOFFICIAL FACEBOOCK PAGE : www.facebook.c" +
				"om/EricFariaOfficialPage\n\n<a href=\"https://theartistunion.com/tracks/5" +
				"21ebe\" rel=\"nofollow\">Download for free on The Artist Union</a>",
		},
		{
			ID:         "lesbisousmusic2:rick-astley-never-gonna-give-3",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://i1.sndcdn.com/artworks-5bj7WTkBvCZhtuNO-Dz8vpg-original.png",
			Title:      "Rick Astley - Never Gonna Give You Up ( LES BISOUS REMIX ) SAMPLE FOR SOUDCLOUD",
			Artist:     "LES BISOUS",
			Provider:   streamnx.Soundcloud,
			Creator:    "LES BISOUS",
			URL:        "https://soundcloud.com/lesbisousmusic2/rick-astley-never-gonna-give-3",
		},
		{
			ID:         "daymaan:rick-astley-never-gonna-give",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://i1.sndcdn.com/artworks-P14eXXDwFuGt47fw-B1A2XA-original.png",
			Title:      "Rick Astley - Never Gonna Give You Up (Daymaan Remix)",
			Artist:     "Daymaan",
			Provider:   streamnx.Soundcloud,
			Creator:    "Daymaan",
			URL:        "https://soundcloud.com/daymaan/rick-astley-never-gonna-give",
			Description: "Enjoy my new remix of « Never Gonna Give You Up » by Rick Astley 🩵 \n\n" +
				"FREE DOWNLOAD : https://hypeddit.com/daymaan/rickastleynevergonnagivey" +
				"ouupdaymaanremix\n\nMix & Mastered by @neus 📀",
		},
		{
			ID:          "capsounds:rick-roll-flip",
			AlbumID:     "",
			AlbumTitle:  "",
			CoverURL:    "https://i1.sndcdn.com/artworks-qAlHEW9eiQy229sC-uPjpyA-original.jpg",
			Title:       "Never Gonna Give You Up - Rick Astley (CAP Bootleg)",
			Artist:      "CAP",
			Provider:    streamnx.Soundcloud,
			Creator:     "CAP",
			URL:         "https://soundcloud.com/capsounds/rick-roll-flip",
			Description: "Rick rolled",
		},
		{
			ID:         "damicoandvalax:never-gonna-give-you-up-damicoandvalax-bootleg-edit-rick-astley",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://i1.sndcdn.com/artworks-QoTvqdk3kIJh1JMZ-HyyKrw-original.jpg",
			Title:      "Never Gonna Give You Up (D'Amico & Valax Bootleg Edit) - Rick Astley",
			Artist:     "D'Amico & Valax",
			Provider:   streamnx.Soundcloud,
			Creator:    "D'Amico & Valax",
			URL: "https://soundcloud.com/damicoandvalax/" +
				"never-gonna-give-you-up-damicoandvalax-bootleg-edit-rick-astley",
			Description: "Questa traccia è monetizzabile.\nClick \"Buy/Acquista\" for FREEDOWNLOAD." +
				"\n\nFollow us//\nInstagram: www.instagram.com/damicoandvalax/\nFacebook: w" +
				"ww.facebook.com/damicoandvalax",
		},
		{
			ID:         "ruanelias:ruan-c-elias-never-gonna-give",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://i1.sndcdn.com/artworks-000134411874-8yipeh-original.jpg",
			Title:      "Branime Studios - Never Gonna Give You Up ( Rick Astley Cover )",
			Artist:     "Ruan Elias",
			Provider:   streamnx.Soundcloud,
			Creator:    "Ruan Elias",
			URL:        "https://soundcloud.com/ruanelias/ruan-c-elias-never-gonna-give",
			Description: "Check out my Spotify page for more songs! \nLink: https://open.spotify." +
				"com/artist/5ASrs5UBOQ8FAKX2cl4D7n",
		},
	}, got)
}

func TestSoundcloudCatalogSearchAlbums(t *testing.T) {
	server := newSoundcloudFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/search/albums",
		Status:  http.StatusOK,
		Fixture: "soundcloud_search_albums_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, soundcloudSearchAlbumArtist+" "+soundcloudSearchAlbum, r.URL.Query().Get("q"))
			require.Equal(t, soundcloudSampleClientID, r.URL.Query().Get("client_id"))
		},
	})
	defer server.Close()

	catalog := newSoundcloudCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Soundcloud, streamnx.SearchQuery{
		Artist: soundcloudSearchAlbumArtist,
		Title:  soundcloudSearchAlbum,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchAlbum{
		{
			ID:       "aikostar-music:whenever-you-need-somebody-4",
			CoverURL: "https://i1.sndcdn.com/artworks-e7f343c1-fdef-41b5-ab34-211a36b3ff39-0-original.jpg",
			Title:    "Whenever You Need Somebody",
			Artist:   "Aiko Star",
			Provider: streamnx.Soundcloud,
			Creator:  "Aiko Star",
			URL:      "https://soundcloud.com/aikostar-music/sets/whenever-you-need-somebody-4",
		},
	}, got)
}

func TestSoundcloudCatalogRejectsUnsupportedOperation(t *testing.T) {
	catalog := newSoundcloudCatalog(t, "http://127.0.0.1")

	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.Soundcloud, "sample")

	require.Zero(t, gotType)
	require.Zero(t, gotID)
	require.ErrorIs(t, err, streamnx.ErrUnsupportedOperation)
}

func newSoundcloudFixtureServer(t *testing.T, apiRoutes ...fixtures.Route) *httptest.Server {
	t.Helper()

	rootRoute := fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Status:  http.StatusOK,
		Fixture: "soundcloud_web_root_200.html",
	}
	return fixtures.NewServer(t, append([]fixtures.Route{rootRoute}, apiRoutes...)...)
}

func newSoundcloudCatalog(t *testing.T, serverURL string) *streamnx.Catalog {
	t.Helper()

	catalog, err := streamnx.NewCatalog(streamnx.WithSoundcloud(
		streamnx.WithSoundcloudAPIURL(serverURL),
		streamnx.WithSoundcloudWebURL(serverURL),
	))
	require.NoError(t, err)

	return catalog
}
