package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2"
	"github.com/GeorgeGorbanev/streamnx/v2/tests/fixtures"
)

const (
	bandcampTrackID        = "archiveformeleins:never-gonna-give-you-up"
	bandcampAlbumID        = "matthewrobertadamski2:whenever-you-need-somebody"
	bandcampSearchArtist   = "Rick Astley"
	bandcampSearchTrack    = "Never Gonna Give You Up"
	bandcampSearchAlbum    = "Whenever You Need Somebody"
	bandcampMissingTrackID = "archiveformeleins:missing-never-gonna-give-you-up"
	bandcampMissingAlbumID = "matthewrobertadamski2:missing-whenever-you-need-somebody"
)

func TestBandcampCatalogFetchTrack(t *testing.T) {
	server := newBandcampFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/track/never-gonna-give-you-up",
		Status:  http.StatusOK,
		Fixture: "bandcamp_fetch_track_200.html",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, "archiveformeleins.bandcamp.com", r.Host)
		},
	})
	defer server.Close()

	catalog := newBandcampCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Bandcamp, bandcampTrackID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:         bandcampTrackID,
		CoverURL:   "https://f4.bcbits.com/img/a3819069079_10.jpg",
		Title:      "Never Gonna Give You Up",
		Artist:     "Rick Astley",
		AlbumID:    "archiveformeleins:formel-eins-space-hits",
		AlbumTitle: "Formel Eins∶ Space Hits",
		Duration:   214,
		ReleaseDate: streamnx.ReleaseDate{
			Year: 1987, Month: 1, Day: 1,
		},
		Provider: streamnx.Bandcamp,
		Creator:  "Archive Formel Eins",
		URL:      "https://archiveformeleins.bandcamp.com/track/never-gonna-give-you-up",
	}, got)
}

func TestBandcampCatalogFetchTrackNotFound(t *testing.T) {
	server := newBandcampFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/track/missing-never-gonna-give-you-up",
		Status:  http.StatusNotFound,
		Fixture: "bandcamp_fetch_track_404.html",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, "archiveformeleins.bandcamp.com", r.Host)
		},
	})
	defer server.Close()

	catalog := newBandcampCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Bandcamp, bandcampMissingTrackID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestBandcampCatalogFetchAlbum(t *testing.T) {
	server := newBandcampFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/album/whenever-you-need-somebody",
		Status:  http.StatusOK,
		Fixture: "bandcamp_fetch_album_200.html",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, "matthewrobertadamski2.bandcamp.com", r.Host)
		},
	})
	defer server.Close()

	catalog := newBandcampCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Bandcamp, bandcampAlbumID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:       bandcampAlbumID,
		CoverURL: "https://f4.bcbits.com/img/a3419300230_10.jpg",
		Title:    "Whenever You Need Somebody",
		Artist:   "Rick Astley",
		Label:    "Matthew Robert Adamski",
		ReleaseDate: streamnx.ReleaseDate{
			Year: 2021, Month: 3, Day: 24,
		},
		Provider: streamnx.Bandcamp,
		Creator:  "Matthew Robert Adamski",
		URL:      "https://matthewrobertadamski2.bandcamp.com/album/whenever-you-need-somebody",
		TrackIDs: []string{
			"matthewrobertadamski2:never-gonna-give-you-up",
			"matthewrobertadamski2:whenever-you-need-somebody",
			"matthewrobertadamski2:together-forever-2",
			"matthewrobertadamski2:when-you-gonna",
			"matthewrobertadamski2:never-need-to-give-you-up-2",
			"matthewrobertadamski2:she-wants-to-dance-with-me",
			"matthewrobertadamski2:it-would-take-a-strong-strong-man",
			"matthewrobertadamski2:never-need-to-give-you-up",
			"matthewrobertadamski2:take-me-to-your-heart-2",
			"matthewrobertadamski2:when-i-fall-in-love",
		},
	}, got)
}

func TestBandcampCatalogFetchAlbumNotFound(t *testing.T) {
	server := newBandcampFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/album/missing-whenever-you-need-somebody",
		Status:  http.StatusNotFound,
		Fixture: "bandcamp_fetch_album_404.html",
		Assert: func(t *testing.T, r *http.Request) {
			require.Equal(t, "matthewrobertadamski2.bandcamp.com", r.Host)
		},
	})
	defer server.Close()

	catalog := newBandcampCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Bandcamp, bandcampMissingAlbumID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestBandcampCatalogSearchTracks(t *testing.T) {
	server := newBandcampFixtureServer(t, fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/api/bcsearch_public_api/1/autocomplete_elastic",
		Status:  http.StatusOK,
		Fixture: "bandcamp_search_tracks_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertBandcampSearchRequest(t, r, bandcampSearchArtist+" "+bandcampSearchTrack, "t")
		},
	})
	defer server.Close()

	catalog := newBandcampCatalog(t, server.URL)
	got, err := catalog.SearchTracks(t.Context(), streamnx.Bandcamp, streamnx.SearchQuery{
		Artist: bandcampSearchArtist,
		Title:  bandcampSearchTrack,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchTrack{
		{
			ID:         "nivedan:rick-astley-never-gonna-give-you-up-nivedan-remix",
			AlbumID:    "",
			AlbumTitle: "Rick Astley - Never Gonna Give You Up (Nivedan Remix)",
			CoverURL:   "https://f4.bcbits.com/img/2531631648_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Nivedan Remix)",
			Artist:     "Nivedan",
			Provider:   streamnx.Bandcamp,
			Creator:    "Nivedan",
			URL:        "https://nivedan.bandcamp.com/track/rick-astley-never-gonna-give-you-up-nivedan-remix",
		},
		{
			ID:         "nivedan:rick-astley-never-gonna-give-you-up-nivedan-remix-dj-intro",
			AlbumID:    "",
			AlbumTitle: "Rick Astley - Never Gonna Give You Up (Nivedan Remix)",
			CoverURL:   "https://f4.bcbits.com/img/2531631648_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Nivedan Remix) - DJ Intro",
			Artist:     "Nivedan",
			Provider:   streamnx.Bandcamp,
			Creator:    "Nivedan",
			URL:        "https://nivedan.bandcamp.com/track/rick-astley-never-gonna-give-you-up-nivedan-remix-dj-intro",
		},
		{
			ID:         "funkastik:rick-astley-never-gonna-give-you-up-funkastik-remix",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://f4.bcbits.com/img/4081158824_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Funkastik remix)",
			Artist:     "Funkastik",
			Provider:   streamnx.Bandcamp,
			Creator:    "Funkastik",
			URL:        "https://funkastik.bandcamp.com/track/rick-astley-never-gonna-give-you-up-funkastik-remix",
		},
		{
			ID:         "fororchestra:rick-astley-never-gonna-give-you-up",
			AlbumID:    "",
			AlbumTitle: "Volume 5",
			CoverURL:   "https://f4.bcbits.com/img/0607888457_3.jpg",
			Title:      "Rick Astley 'Never Gonna Give You Up'",
			Artist:     "Walt Ribeiro",
			Provider:   streamnx.Bandcamp,
			Creator:    "Walt Ribeiro",
			URL:        "https://fororchestra.bandcamp.com/track/rick-astley-never-gonna-give-you-up",
		},
		{
			ID:         "deepsound:rick-astley-never-gonna-give-you-up-deepsound-remix",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://f4.bcbits.com/img/3840263489_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (DEEPSOUND REMIX)",
			Artist:     "Deepsound",
			Provider:   streamnx.Bandcamp,
			Creator:    "Deepsound",
			URL:        "https://deepsound.bandcamp.com/track/rick-astley-never-gonna-give-you-up-deepsound-remix",
		},
		{
			ID:         "kikesundance:rick-astley-never-gonna-give-you-up-kike-regroove",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://f4.bcbits.com/img/0559451941_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up - Kike Regroove",
			Artist:     "Kike Sundance",
			Provider:   streamnx.Bandcamp,
			Creator:    "Kike Sundance",
			URL:        "https://kikesundance.bandcamp.com/track/rick-astley-never-gonna-give-you-up-kike-regroove",
		},
		{
			ID:         "elzexd:elzexd-rick-astley-never-gonna-give-you-up",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://f4.bcbits.com/img/2431503587_3.jpg",
			Title:      "elzeXD™ - Rick Astley - Never Gonna Give You Up",
			Artist:     "elzeXD",
			Provider:   streamnx.Bandcamp,
			Creator:    "elzeXD",
			URL:        "https://elzexd.bandcamp.com/track/elzexd-rick-astley-never-gonna-give-you-up",
		},
		{
			ID:         "djkikegarage:rick-astley-never-gonna-give-you-up-kike-regroove",
			AlbumID:    "",
			AlbumTitle: "",
			CoverURL:   "https://f4.bcbits.com/img/1682149501_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up - Kike ReGroove",
			Artist:     "Dj Kike Garage",
			Provider:   streamnx.Bandcamp,
			Creator:    "Dj Kike Garage",
			URL:        "https://djkikegarage.bandcamp.com/track/rick-astley-never-gonna-give-you-up-kike-regroove",
		},
		{
			ID:         "fororchestra:rick-astley-never-gonna-give-you-up-2",
			AlbumID:    "",
			AlbumTitle: "Every Song!",
			CoverURL:   "https://f4.bcbits.com/img/2188361452_3.jpg",
			Title:      "Rick Astley 'Never Gonna Give You Up'",
			Artist:     "Walt Ribeiro",
			Provider:   streamnx.Bandcamp,
			Creator:    "Walt Ribeiro",
			URL:        "https://fororchestra.bandcamp.com/track/rick-astley-never-gonna-give-you-up-2",
		},
		{
			ID:         "clubremix:rick-astley-never-gonna-give-you-up-daniel-adame-remix",
			AlbumID:    "",
			AlbumTitle: "MASTERMIX 001",
			CoverURL:   "https://f4.bcbits.com/img/2457294734_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Daniel Adame Remix)",
			Artist:     "ClubRemix",
			Provider:   streamnx.Bandcamp,
			Creator:    "ClubRemix",
			URL:        "https://clubremix.bandcamp.com/track/rick-astley-never-gonna-give-you-up-daniel-adame-remix",
		},
		{
			ID:         "zeelatunes:rick-astley-never-gonna-give-you-up-chrono-trigger-style",
			AlbumID:    "",
			AlbumTitle: "2023 Releases",
			CoverURL:   "https://f4.bcbits.com/img/0577116200_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up [Chrono Trigger-style]",
			Artist:     "ZeelaTunes",
			Provider:   streamnx.Bandcamp,
			Creator:    "ZeelaTunes",
			URL:        "https://zeelatunes.bandcamp.com/track/rick-astley-never-gonna-give-you-up-chrono-trigger-style",
		},
		{
			ID:         "jompatott:rick-astley-never-gonna-give-you-up-jompatott-reggae-blend",
			AlbumID:    "",
			AlbumTitle: "Reggae Blends 3",
			CoverURL:   "https://f4.bcbits.com/img/1702451392_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Jompatott Reggae Blend)",
			Artist:     "Jompatott",
			Provider:   streamnx.Bandcamp,
			Creator:    "Jompatott",
			URL:        "https://jompatott.bandcamp.com/track/rick-astley-never-gonna-give-you-up-jompatott-reggae-blend",
		},
		{
			ID:         "archiveformeleins:never-gonna-give-you-up",
			AlbumID:    "",
			AlbumTitle: "Formel Eins∶ Space Hits",
			CoverURL:   "https://f4.bcbits.com/img/3819069079_3.jpg",
			Title:      "Never Gonna Give You Up",
			Artist:     "Rick Astley",
			Provider:   streamnx.Bandcamp,
			Creator:    "Rick Astley",
			URL:        "https://archiveformeleins.bandcamp.com/track/never-gonna-give-you-up",
		},
		{
			ID:         "zerolaggaming:rick-astley-never-gonna-give-you-up-chrono-trigger-soundfont",
			AlbumID:    "",
			AlbumTitle: "Videogame Soundfont Songs",
			CoverURL:   "https://f4.bcbits.com/img/2534073844_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Chrono Trigger Soundfont)",
			Artist:     "ZerolagGaming",
			Provider:   streamnx.Bandcamp,
			Creator:    "ZerolagGaming",
			URL:        "https://zerolaggaming.bandcamp.com/track/rick-astley-never-gonna-give-you-up-chrono-trigger-soundfont",
		},
		{
			ID:         "danieladame:rick-astley-never-gonna-give-you-up-daniel-adame-remix",
			AlbumID:    "",
			AlbumTitle: "Master Mix 001",
			CoverURL:   "https://f4.bcbits.com/img/3776071639_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Daniel Adame Remix)",
			Artist:     "Daniel Adame",
			Provider:   streamnx.Bandcamp,
			Creator:    "Daniel Adame",
			URL:        "https://danieladame.bandcamp.com/track/rick-astley-never-gonna-give-you-up-daniel-adame-remix",
		},
		{
			ID:         "djheart1:rick-astley-never-gonna-give-you-up",
			AlbumID:    "",
			AlbumTitle: "Band Dance 1989 - 2000 Vol:3",
			CoverURL:   "https://f4.bcbits.com/img/2455404821_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up",
			Artist:     "DJHeart",
			Provider:   streamnx.Bandcamp,
			Creator:    "DJHeart",
			URL:        "https://djheart1.bandcamp.com/track/rick-astley-never-gonna-give-you-up",
		},
		{
			ID:         "ssproductions1:rick-astley-never-gonna-give-you-up-ss-edit-revibe",
			AlbumID:    "",
			AlbumTitle: "Classic Revibes",
			CoverURL:   "https://f4.bcbits.com/img/3854808208_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up {SS Edit Revibe}",
			Artist:     "SS PRODUCTIONS",
			Provider:   streamnx.Bandcamp,
			Creator:    "SS PRODUCTIONS",
			URL:        "https://ssproductions1.bandcamp.com/track/rick-astley-never-gonna-give-you-up-ss-edit-revibe",
		},
		{
			ID:         "djskytrini:rick-astley-never-gonna-give-you-up-les-bisous-remix",
			AlbumID:    "",
			AlbumTitle: "Mashups/Euroremixes #3",
			CoverURL:   "https://f4.bcbits.com/img/1771115437_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (LES BISOUS REMIX)]",
			Artist:     "sky trini",
			Provider:   streamnx.Bandcamp,
			Creator:    "sky trini",
			URL:        "https://djskytrini.bandcamp.com/track/rick-astley-never-gonna-give-you-up-les-bisous-remix",
		},
		{
			ID:         "studiomasters:rick-astley-never-gonna-give-you-up-redrum-114-bpm",
			AlbumID:    "",
			AlbumTitle: "Rick Astley (SM Remix)",
			CoverURL:   "https://f4.bcbits.com/img/2438376551_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Redrum) 114 bpm",
			Artist:     "Studio Masters",
			Provider:   streamnx.Bandcamp,
			Creator:    "Studio Masters",
			URL:        "https://studiomasters.bandcamp.com/track/rick-astley-never-gonna-give-you-up-redrum-114-bpm",
		},
		{
			ID:         "studiomasters:rick-astley-never-gonna-give-you-up-remix-120-bpm",
			AlbumID:    "",
			AlbumTitle: "Rick Astley (SM Remix)",
			CoverURL:   "https://f4.bcbits.com/img/0501364509_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Remix) 120 bpm",
			Artist:     "Studio Masters",
			Provider:   streamnx.Bandcamp,
			Creator:    "Studio Masters",
			URL:        "https://studiomasters.bandcamp.com/track/rick-astley-never-gonna-give-you-up-remix-120-bpm",
		},
		{
			ID:         "garework:rick-astley-never-gonna-give-you-up-g-a-rework-2026",
			AlbumID:    "",
			AlbumTitle: "G A RHEWORK 2026",
			CoverURL:   "https://f4.bcbits.com/img/1574623383_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (G A REWORK 2026)",
			Artist:     "Garework",
			Provider:   streamnx.Bandcamp,
			Creator:    "Garework",
			URL:        "https://garework.bandcamp.com/track/rick-astley-never-gonna-give-you-up-g-a-rework-2026",
		},
		{
			ID:         "djtools:rick-astley-never-gonna-give-you-up-cake-mix",
			AlbumID:    "",
			AlbumTitle: "80's Remixed & Extended",
			CoverURL:   "https://f4.bcbits.com/img/0133422820_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Cake Mix)",
			Artist:     "DJ Tools",
			Provider:   streamnx.Bandcamp,
			Creator:    "DJ Tools",
			URL:        "https://djtools.bandcamp.com/track/rick-astley-never-gonna-give-you-up-cake-mix",
		},
		{
			ID:         "djselphi:rick-astley-never-gonna-give-you-up-dj-selphi-bachata-remix-ft-camilo-bass",
			AlbumID:    "",
			AlbumTitle: "Bachata Remixes 2019 Pt 2",
			CoverURL:   "https://f4.bcbits.com/img/0879307196_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (DJ Selphi bachata remix ft Camilo Bass)",
			Artist:     "DJ Selphi",
			Provider:   streamnx.Bandcamp,
			Creator:    "DJ Selphi",
			URL:        "https://djselphi.bandcamp.com/track/rick-astley-never-gonna-give-you-up-dj-selphi-bachata-remix-ft-camilo-bass",
		},
		{
			ID:         "gglab:rick-astley-never-gonna-give-you-up-aka-rickroll",
			AlbumID:    "",
			AlbumTitle: "The 8bitified hits",
			CoverURL:   "https://f4.bcbits.com/img/1001684710_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (aka RickRoll)",
			Artist:     "vladkorotnev 8bit",
			Provider:   streamnx.Bandcamp,
			Creator:    "vladkorotnev 8bit",
			URL:        "https://gglab.bandcamp.com/track/rick-astley-never-gonna-give-you-up-aka-rickroll",
		},
		{
			ID:         "deha:rick-astley-never-gonna-give-you-up",
			AlbumID:    "",
			AlbumTitle: "Pandechristmas!",
			CoverURL:   "https://f4.bcbits.com/img/0892388933_3.jpg",
			Title:      "Rick Astley - Never gonna give you up",
			Artist:     "Déhà & The Social Distancing Choir",
			Provider:   streamnx.Bandcamp,
			Creator:    "Déhà & The Social Distancing Choir",
			URL:        "https://deha.bandcamp.com/track/rick-astley-never-gonna-give-you-up",
		},
		{
			ID:         "mcrpmusic:rick-astley-never-gonna-give-you-up-noiseflower-remix",
			AlbumID:    "",
			AlbumTitle: "MCRP V6 - The Bad Song Edition",
			CoverURL:   "https://f4.bcbits.com/img/1867174918_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (noiseflower remix)",
			Artist:     "MCRP",
			Provider:   streamnx.Bandcamp,
			Creator:    "MCRP",
			URL:        "https://mcrpmusic.bandcamp.com/track/rick-astley-never-gonna-give-you-up-noiseflower-remix",
		},
		{
			ID:         "vincentbastille:rick-astley-never-gonna-give-you-up-vincent-bastille-house-remix",
			AlbumID:    "",
			AlbumTitle: "Vincent Bastille - Remixes 2026",
			CoverURL:   "https://f4.bcbits.com/img/3115815517_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Vincent Bastille House Remix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Bandcamp,
			Creator:    "Rick Astley",
			URL:        "https://vincentbastille.bandcamp.com/track/rick-astley-never-gonna-give-you-up-vincent-bastille-house-remix",
		},
		{
			ID:         "polarmaxbit:rick-astley-never-gonna-give-you-up-mix-polarmaxbit-2",
			AlbumID:    "",
			AlbumTitle: "POLARMAXBIT Remixes y Mashups for djs N\u200b\u200b\u200b-\u200b\u200b\u200b3",
			CoverURL:   "https://f4.bcbits.com/img/2676347592_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up MIX POLARMAXBIT",
			Artist:     "Polarmaxbit",
			Provider:   streamnx.Bandcamp,
			Creator:    "Polarmaxbit",
			URL:        "https://polarmaxbit.bandcamp.com/track/rick-astley-never-gonna-give-you-up-mix-polarmaxbit-2",
		},
		{
			ID:         "polarmaxbit:rick-astley-never-gonna-give-you-up-mix-polarmaxbit",
			AlbumID:    "",
			AlbumTitle: "POLARMAXBIT Remixes y Mashups for djs N\u200b-\u200b2",
			CoverURL:   "https://f4.bcbits.com/img/3285644934_3.jpg",
			Title:      "Rick Astley- Never Gonna Give You Up MIX POLARMAXBIT",
			Artist:     "Polarmaxbit",
			Provider:   streamnx.Bandcamp,
			Creator:    "Polarmaxbit",
			URL:        "https://polarmaxbit.bandcamp.com/track/rick-astley-never-gonna-give-you-up-mix-polarmaxbit",
		},
		{
			ID:         "djrapha:rick-astley-never-gonna-give-you-up-remix",
			AlbumID:    "",
			AlbumTitle: "80's Remix Party - Vol. 4",
			CoverURL:   "https://f4.bcbits.com/img/3804535553_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Remix)",
			Artist:     "Listen to the Radio",
			Provider:   streamnx.Bandcamp,
			Creator:    "Listen to the Radio",
			URL:        "https://djrapha.bandcamp.com/track/rick-astley-never-gonna-give-you-up-remix",
		},
		{
			ID:         "deejaydisc:rick-astley-never-gonna-give-you-up-les-bisous-remix",
			AlbumID:    "",
			AlbumTitle: "Remix For Djs (Vol. 17)",
			CoverURL:   "https://f4.bcbits.com/img/0919062034_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Les Bisous Remix)",
			Artist:     "DJ DISC !",
			Provider:   streamnx.Bandcamp,
			Creator:    "DJ DISC !",
			URL:        "https://deejaydisc.bandcamp.com/track/rick-astley-never-gonna-give-you-up-les-bisous-remix",
		},
		{
			ID:         "windowslogic:never-gonna-give-you-up-porn-game-mix",
			AlbumID:    "",
			AlbumTitle: "Rick Astley - Never Gonna Give You Up Remix",
			CoverURL:   "https://f4.bcbits.com/img/2763727142_3.jpg",
			Title:      "Never Gonna Give You Up (Porn Game Mix)",
			Artist:     "WindowsLogic Productions",
			Provider:   streamnx.Bandcamp,
			Creator:    "WindowsLogic Productions",
			URL:        "https://windowslogic.bandcamp.com/track/never-gonna-give-you-up-porn-game-mix",
		},
		{
			ID:         "cryincookieedits:rick-astley-never-gonna-give-you-up-cryincookie-short-edit",
			AlbumID:    "",
			AlbumTitle: "FreeBee Short Edit Pack #03",
			CoverURL:   "https://f4.bcbits.com/img/1792010672_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Cryincookie Short Edit)",
			Artist:     "Cryin Cookie Edits",
			Provider:   streamnx.Bandcamp,
			Creator:    "Cryin Cookie Edits",
			URL:        "https://cryincookieedits.bandcamp.com/track/rick-astley-never-gonna-give-you-up-cryincookie-short-edit",
		},
		{
			ID:         "djtools:rick-astley-never-gonna-give-you-up-dario-caminita-revibe",
			AlbumID:    "",
			AlbumTitle: "Dario Caminita - Classic Revibes Collection (Volume 02)",
			CoverURL:   "https://f4.bcbits.com/img/0610161221_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Dario Caminita Revibe)",
			Artist:     "DJ Tools",
			Provider:   streamnx.Bandcamp,
			Creator:    "DJ Tools",
			URL:        "https://djtools.bandcamp.com/track/rick-astley-never-gonna-give-you-up-dario-caminita-revibe",
		},
		{
			ID:         "theartsmixof2:rick-astley-never-gonna-give-you-up-the-art-mix-rcx-3",
			AlbumID:    "",
			AlbumTitle: "THE! ART$ EDITION VOL.110|",
			CoverURL:   "https://f4.bcbits.com/img/3935152916_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up [The! Art$ Mix RCX 3]",
			Artist:     "Art$.Mix",
			Provider:   streamnx.Bandcamp,
			Creator:    "Art$.Mix",
			URL:        "https://theartsmixof2.bandcamp.com/track/rick-astley-never-gonna-give-you-up-the-art-mix-rcx-3",
		},
		{
			ID:         "drborkez:rick-astley-never-gonna-give-you-up-borkez-music-bootleg-2",
			AlbumID:    "",
			AlbumTitle: "80s English Pop Vol.3 (79 Tracks)",
			CoverURL:   "https://f4.bcbits.com/img/3713870347_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Borkez Music Bootleg)",
			Artist:     "Dr. Borkez",
			Provider:   streamnx.Bandcamp,
			Creator:    "Dr. Borkez",
			URL:        "https://drborkez.bandcamp.com/track/rick-astley-never-gonna-give-you-up-borkez-music-bootleg-2",
		},
		{
			ID:         "windowslogic:never-gonna-give-you-up-70s-strip-club-mix",
			AlbumID:    "",
			AlbumTitle: "Rick Astley - Never Gonna Give You Up Remix",
			CoverURL:   "https://f4.bcbits.com/img/2763727142_3.jpg",
			Title:      "Never Gonna Give You Up (70s Strip Club Mix)",
			Artist:     "WindowsLogic Productions",
			Provider:   streamnx.Bandcamp,
			Creator:    "WindowsLogic Productions",
			URL:        "https://windowslogic.bandcamp.com/track/never-gonna-give-you-up-70s-strip-club-mix",
		},
		{
			ID:         "12inchrecord:never-gonna-give-you-up-12-cake-mix",
			AlbumID:    "",
			AlbumTitle: "Twelve Inch Eighties∶ You Spin Me Round",
			CoverURL:   "https://f4.bcbits.com/img/2200607829_3.jpg",
			Title:      "Never Gonna Give You Up (12″ Cake mix)",
			Artist:     "Rick Astley",
			Provider:   streamnx.Bandcamp,
			Creator:    "Rick Astley",
			URL:        "https://12inchrecord.bandcamp.com/track/never-gonna-give-you-up-12-cake-mix",
		},
		{
			ID:         "joshuaabbottrhythm:rick-astley-never-gonna-give-you-up-80s-bounce-banger-dirty-128-track-46",
			AlbumID:    "",
			AlbumTitle: "Remix Planet 1504",
			CoverURL:   "https://f4.bcbits.com/img/2562883311_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (80s Bounce Banger) Dirty 128 - Track 46",
			Artist:     "Joshua Abbott Rhythm",
			Provider:   streamnx.Bandcamp,
			Creator:    "Joshua Abbott Rhythm",
			URL:        "https://joshuaabbottrhythm.bandcamp.com/track/rick-astley-never-gonna-give-you-up-80s-bounce-banger-dirty-128-track-46",
		},
		{
			ID:         "theartsmixof2:rick-astley-never-gonna-give-you-up-the-art-mix-rcx-3-intro",
			AlbumID:    "",
			AlbumTitle: "THE! ART$ EDITION VOL.110|",
			CoverURL:   "https://f4.bcbits.com/img/3935152916_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up [The! Art$ Mix RCX 3] INTRO",
			Artist:     "Art$.Mix",
			Provider:   streamnx.Bandcamp,
			Creator:    "Art$.Mix",
			URL:        "https://theartsmixof2.bandcamp.com/track/rick-astley-never-gonna-give-you-up-the-art-mix-rcx-3-intro",
		},
		{
			ID:         "kathleenallenclub:rick-astley-never-gonna-give-you-up-acapella-intro-outro-breakz-clean-113-track-14",
			AlbumID:    "",
			AlbumTitle: "Beatfreakz 0705",
			CoverURL:   "https://f4.bcbits.com/img/4150926818_3.jpg",
			Title:      "- Rick Astley Never Gonna Give You Up Acapella Intro Outro Breakz Clean 113 - Track 14",
			Artist:     "Kathleen Allen Club",
			Provider:   streamnx.Bandcamp,
			Creator:    "Kathleen Allen Club",
			URL:        "https://kathleenallenclub.bandcamp.com/track/rick-astley-never-gonna-give-you-up-acapella-intro-outro-breakz-clean-113-track-14",
		},
		{
			ID:         "djrapha:rick-astley-never-gonna-give-you-up-extended-version",
			AlbumID:    "",
			AlbumTitle: "Back to the Beat: 80's Extended Vol. 3",
			CoverURL:   "https://f4.bcbits.com/img/1207811559_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Extended Version)",
			Artist:     "Listen to the Radio",
			Provider:   streamnx.Bandcamp,
			Creator:    "Listen to the Radio",
			URL:        "https://djrapha.bandcamp.com/track/rick-astley-never-gonna-give-you-up-extended-version",
		},
		{
			ID:         "jonathanlovevibes:deejay-oninz-x-rick-astley-never-gonna-give-you-up-80s-redrum-clean-128-track-9",
			AlbumID:    "",
			AlbumTitle: "Beatfreakz 2503",
			CoverURL:   "https://f4.bcbits.com/img/4185235240_3.jpg",
			Title:      "Deejay Oninz X Rick Astley - Never Gonna Give You Up (80's Redrum) (Clean) 128 - Track 9",
			Artist:     "Jonathan Love Vibes",
			Provider:   streamnx.Bandcamp,
			Creator:    "Jonathan Love Vibes",
			URL:        "https://jonathanlovevibes.bandcamp.com/track/deejay-oninz-x-rick-astley-never-gonna-give-you-up-80s-redrum-clean-128-track-9",
		},
		{
			ID:         "delicamarr:rick-astley-never-gonna-give-you-up-delicamarr-latin-house-edit",
			AlbumID:    "",
			AlbumTitle: "델리카마의 \"주크\" 노래방 - Delic'amarr's \"Jukaraoke\" Vol.2",
			CoverURL:   "https://f4.bcbits.com/img/3196051375_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (Delic'amarr Latin House Edit)",
			Artist:     "Delic'amarr",
			Provider:   streamnx.Bandcamp,
			Creator:    "Delic'amarr",
			URL:        "https://delicamarr.bandcamp.com/track/rick-astley-never-gonna-give-you-up-delicamarr-latin-house-edit",
		},
		{
			ID:         "wendyortizproject:deejay-oninz-x-rick-astley-never-gonna-give-you-up-80s-redrum-clean-128-track-18",
			AlbumID:    "",
			AlbumTitle: "Beatfreakz 0102",
			CoverURL:   "https://f4.bcbits.com/img/1736221638_3.jpg",
			Title:      "Deejay Oninz X Rick Astley - Never Gonna Give You Up (80's Redrum) (Clean) 128 - Track 18",
			Artist:     "Wendy Ortiz Project",
			Provider:   streamnx.Bandcamp,
			Creator:    "Wendy Ortiz Project",
			URL:        "https://wendyortizproject.bandcamp.com/track/deejay-oninz-x-rick-astley-never-gonna-give-you-up-80s-redrum-clean-128-track-18",
		},
		{
			ID:         "djmarcand:rick-astley-never-gonna-give-you-up-dj-marcand-alternative-version",
			AlbumID:    "",
			AlbumTitle: "Alternative Versions of Popular Songs Vol-7",
			CoverURL:   "https://f4.bcbits.com/img/4036607259_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (DJ Marcand Alternative Version)",
			Artist:     "DJ Marcand",
			Provider:   streamnx.Bandcamp,
			Creator:    "DJ Marcand",
			URL:        "https://djmarcand.bandcamp.com/track/rick-astley-never-gonna-give-you-up-dj-marcand-alternative-version",
		},
		{
			ID:         "musictotheworld:rick-astley-never-gonna-give-you-up",
			AlbumID:    "",
			AlbumTitle: "ThA BosS Present Classic Dance Craze",
			CoverURL:   "https://f4.bcbits.com/img/0205520552_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up",
			Artist:     "ThA BosS [Music To The World]",
			Provider:   streamnx.Bandcamp,
			Creator:    "ThA BosS [Music To The World]",
			URL:        "https://musictotheworld.bandcamp.com/track/rick-astley-never-gonna-give-you-up",
		},
		{
			ID:         "wendyortizproject:rick-astley-never-gonna-give-you-up-damico-valax-bootleg-edit-clean-126-track-81",
			AlbumID:    "",
			AlbumTitle: "Crack 4 DJs 2602",
			CoverURL:   "https://f4.bcbits.com/img/0740982072_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up (D'amico & Valax Bootleg Edit) (Clean) 126 - Track 81",
			Artist:     "Wendy Ortiz Project",
			Provider:   streamnx.Bandcamp,
			Creator:    "Wendy Ortiz Project",
			URL:        "https://wendyortizproject.bandcamp.com/track/rick-astley-never-gonna-give-you-up-damico-valax-bootleg-edit-clean-126-track-81",
		},
		{
			ID:         "louisgriffinclub:dj-allan-x-rick-astley-never-gonna-give-you-up-80s-dance-redrum-v-2-114-track-11",
			AlbumID:    "",
			AlbumTitle: "Beatfreakz 3105",
			CoverURL:   "https://f4.bcbits.com/img/3162135423_3.jpg",
			Title:      "Dj Allan X Rick Astley - Never Gonna Give You Up (80's Dance Redrum V.2) 114 - Track 11",
			Artist:     "Louis Griffin Club",
			Provider:   streamnx.Bandcamp,
			Creator:    "Louis Griffin Club",
			URL:        "https://louisgriffinclub.bandcamp.com/track/dj-allan-x-rick-astley-never-gonna-give-you-up-80s-dance-redrum-v-2-114-track-11",
		},
		{
			ID:         "edits4you:rick-astley-never-gonna-give-you-up-dj-ugeezy-edit-pn",
			AlbumID:    "",
			AlbumTitle: "80's Vol 2 ( DJ INTROS )",
			CoverURL:   "https://f4.bcbits.com/img/2295777862_3.jpg",
			Title:      "Rick Astley - Never Gonna Give You Up ( DJ UGEEZY EDIT )_PN",
			Artist:     "edits4you",
			Provider:   streamnx.Bandcamp,
			Creator:    "edits4you",
			URL:        "https://edits4you.bandcamp.com/track/rick-astley-never-gonna-give-you-up-dj-ugeezy-edit-pn",
		},
	}, got)
}

func TestBandcampCatalogSearchAlbums(t *testing.T) {
	server := newBandcampFixtureServer(t, fixtures.Route{
		Method:  http.MethodPost,
		Path:    "/api/bcsearch_public_api/1/autocomplete_elastic",
		Status:  http.StatusOK,
		Fixture: "bandcamp_search_albums_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			assertBandcampSearchRequest(t, r, bandcampSearchArtist+" "+bandcampSearchAlbum, "a")
		},
	})
	defer server.Close()

	catalog := newBandcampCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Bandcamp, streamnx.SearchQuery{
		Artist: bandcampSearchArtist,
		Title:  bandcampSearchAlbum,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchAlbum{
		{
			ID:       "matthewrobertadamski2:whenever-you-need-somebody",
			CoverURL: "https://f4.bcbits.com/img/3419300230_3.jpg",
			Title:    "Whenever You Need Somebody",
			Artist:   "Rick Astley",
			Provider: streamnx.Bandcamp,
			Creator:  "Rick Astley",
			URL:      "https://matthewrobertadamski2.bandcamp.com/album/whenever-you-need-somebody",
		},
	}, got)
}

func TestBandcampCatalogRejectsUnsupportedOperation(t *testing.T) {
	catalog := newBandcampCatalog(t, "http://127.0.0.1")

	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.Bandcamp, "sample")

	require.Zero(t, gotType)
	require.Zero(t, gotID)
	require.ErrorIs(t, err, streamnx.ErrUnsupportedOperation)
}

func newBandcampFixtureServer(t *testing.T, apiRoutes ...fixtures.Route) *httptest.Server {
	t.Helper()

	return fixtures.NewServer(t, apiRoutes...)
}

func newBandcampCatalog(t *testing.T, serverURL string) *streamnx.Catalog {
	t.Helper()

	catalog, err := streamnx.NewCatalog(streamnx.WithBandcamp(
		streamnx.WithBandcampAPI("https", "bandcamp.com"),
		streamnx.WithBandcampHTTPClient(newBandcampFixtureHTTPClient(t, serverURL)),
	))
	require.NoError(t, err)

	return catalog
}

func newBandcampFixtureHTTPClient(t *testing.T, serverURL string) *http.Client {
	t.Helper()

	target, err := url.Parse(serverURL)
	require.NoError(t, err)

	return &http.Client{
		Transport: bandcampFixtureTransport{
			target: target,
		},
	}
}

type bandcampFixtureTransport struct {
	target *url.URL
}

func (t bandcampFixtureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	clone := r.Clone(r.Context())
	clone.URL.Scheme = t.target.Scheme
	clone.URL.Host = t.target.Host
	clone.Host = r.URL.Host

	return http.DefaultTransport.RoundTrip(clone)
}

func assertBandcampSearchRequest(t *testing.T, r *http.Request, searchText, searchFilter string) {
	t.Helper()

	require.Equal(t, "bandcamp.com", r.Host)
	require.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Content-Type"))

	var body struct {
		SearchText   string `json:"search_text"`
		SearchFilter string `json:"search_filter"`
		FullPage     bool   `json:"full_page"`
	}
	require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
	require.Equal(t, searchText, body.SearchText)
	require.Equal(t, searchFilter, body.SearchFilter)
	require.True(t, body.FullPage)
}
