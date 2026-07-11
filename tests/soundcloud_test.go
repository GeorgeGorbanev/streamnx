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
	soundcloudSampleClientID = "QNR5nrdLOvApYERC8AOUr3VjRfHnLjle"

	soundcloudTrackID        = "rick-astley-official:never-gonna-give-you-up"
	soundcloudAlbumID        = "aikostar-music:whenever-you-need-somebody-4"
	soundcloudSearchArtist   = "Rick Astley"
	soundcloudSearchTrack    = "Never Gonna Give You Up"
	soundcloudSearchAlbum    = "Whenever You Need Somebody"
	soundcloudMissingTrackID = "rick-astley-official:missing-never-gonna-give-you-up"
	soundcloudMissingAlbumID = "rick-astley-official:missing-whenever-you-need-somebody"
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
	server := newSoundcloudFixtureServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/aikostar-music/sets/whenever-you-need-somebody-4",
		Status:  http.StatusOK,
		Fixture: "soundcloud_fetch_album_200.html",
	})
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
			require.Equal(t, soundcloudSearchAlbum, r.URL.Query().Get("q"))
			require.Equal(t, soundcloudSampleClientID, r.URL.Query().Get("client_id"))
		},
	})
	defer server.Close()

	catalog := newSoundcloudCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Soundcloud, streamnx.SearchQuery{
		Artist: soundcloudSearchArtist,
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
		{
			ID:       "lindorff:whenever-you-need-somebody",
			CoverURL: "https://i1.sndcdn.com/artworks-e19d0f03-3a3b-4348-abba-c95a083d8665-0-original.jpg",
			Title:    "Whenever you need somebody",
			Artist:   "Lindorff",
			Provider: streamnx.Soundcloud,
			Creator:  "Lindorff",
			URL:      "https://soundcloud.com/lindorff/sets/whenever-you-need-somebody",
		},
		{
			ID:       "nobrainerrecords:malente-zero-cash-ill-be",
			CoverURL: "https://i1.sndcdn.com/artworks-000011838695-fcct6m-original.jpg",
			Title:    "Malente & Zero Cash - I'll Be There (NBR011)",
			Artist:   "No Brainer Records",
			Provider: streamnx.Soundcloud,
			Creator:  "No Brainer Records",
			URL:      "https://soundcloud.com/nobrainerrecords/sets/malente-zero-cash-ill-be",
			Description: "Whenever you need somebody …\r\n\r\n'I'll Be There'. Sounds like a love so" +
				"ng, if it wasn't for the music. Banging techno by ZERO CASH (Televisio" +
				"n Rocks) and MALENTE (No Brainer). Then again: Banging is not that far" +
				" from love, is it?\r\n\r\nMODEK (Keatchen) loves a good party. His remix n" +
				"ails your hands to the sky and your legs on the floor while shaking yo" +
				"ur body. Funny how brutal love can be.\r\n\r\nLIGHT YEAR (Bang Gang) love " +
				"the acid. Their raw and uncompromising take on 'I'll Be There' makes y" +
				"ou wanna be nothing else but their 808.\r\n\r\nKILL FRENZY (Lectroluv) lov" +
				"es bass. He bounces his techno below the surface. I think he might be " +
				"after some Mermaids with what sounds like sexy sub abuse.\r\n\r\nIf you're" +
				" not in love yet MALENTE & ZERO CASH make sure to break your heart wit" +
				"h 'Hardware'. Not too many of you will be able to handle the pure amou" +
				"nt of adrenaline this tune shoots into your body.\r\n\r\n'But don't worry," +
				" I'll be there' \r\n\r\n\r\n\r\nBORIS DLUGOSCH \"Light Year is KILLER!!!!! Full" +
				" support\"\r\nBOBMO \"I like the Light Year remix and Hardware , thanks!\"\r" +
				"\nBRODINSKI \"Love the Light Year rmx! Bigup!\"\r\nBAG RAIDERS \"This is ter" +
				"rific!!\"PUNKS JUMP UP \"Love the record! Hardware is our fav!\"\r\nSINDEN " +
				"\"Yeah bwoy! The OG works for me. Quirky fun dance music\"\r\nTAI \"Lovely " +
				"package. Definitely bunging Modek and Light Year rmx in my next set!\"\r" +
				"\nBENI \"Light Year mix is for me\"\r\nMUMBAI SCIENCE \"Light Year Remix in " +
				"our sets!\"\r\nAC SLATER \"Great release. hard to pick a fav track but Kil" +
				"l Frenzy mix sticks out to me. Hardware is crazy with those detuned ho" +
				"rns.\"\r\nDON RIMINI \"I play the Modek remix. And the Light Year is crazy" +
				" ! I love the break in Hardware. Huge tracks!\"\r\nDON DIABLO \"Light Year" +
				" mix = fire!\"\r\nMASON \"It's mad, disturbing, troubled. Thank you\"\r\n\r\n\r\n" +
				"Zombies For Money, Ado, Russ Chimes, Moonbootica, Feadz, Tony Senghore" +
				", Big Dope P / Moveltraxx, Neoteric, Larry Tee, Mixhell, Botnek, Utah " +
				"Saints, Slap In The Bass, Hickup, The Aston Shuffle, Nick Catchdubs, P" +
				"lump DJs, Jay Robinson, Joyce Muniz, Wolfie, John Roman, Blaze Tripp, " +
				"Burns, Act Yo Age, Lorcan Mak, Jeff Doubleu, Tom Piper, Udachi, Stereo" +
				" MCs, Danny T, Ben Mono, Maelstrom, DJ Gina Turner, Teenage Mutants, T" +
				"om Stephan, Tagteam Terror, Acidkids, Markus Lange, Milt Mortez, Max l" +
				"e Daron, Breakfastclub DJs, Fashen, Smalltown DJs, Voltron, Mat The Al" +
				"ien, Willy Joy, B Rich, aUtOdiDakT, Oh Snap!, Designer Drugs, Ursula 1" +
				"000, Rishi Romero, Fex Fellini, Ajax, Twist It!, Pete Carvell (Bad Lif" +
				"e), 3 Is A Crowd, Max Cherry, BREAKS lda, \r\n\r\n\r\nRambaud Ludovic (Only " +
				"For Djs Mag), Nina (Triple J Radio, Australia), Jérémie Anticlimax (Ts" +
				"ugi Mag France), Andre Langenfeld (Radio Fritz, Berlin),  DJ Pffff (Ra" +
				"zzmatazz BCN), Tommy Yamaha (on 3 radio germany), Dj Dmit.ry (Central " +
				"Station Radost FX, GOX Radi0 Czech), HARPER (Boogie Mafia, Polskie Rad" +
				"io 4), Robert Borzym (Polsike Radio Euro), K. Ramba (TheNewFrenchTouch" +
				"), Richard Heinemann (T M I Radio ARA, Luxemburg), John Buergin (Schwe" +
				"izer Rundfunk DRS Virus), Michal Stolárik (Hochspannung / Slavakia), M" +
				"ister Sushi (ibreaks.co.uk / London)",
		},
		{
			ID:       "oneill-fernandes:you-light-up-my-life",
			CoverURL: "https://i1.sndcdn.com/artworks-vNExEuqkbwmbkdXv-Fz9UjA-original.jpg",
			Title:    "You Light Up My Life",
			Artist:   "O'Neill Fernandes",
			Provider: streamnx.Soundcloud,
			Creator:  "O'Neill Fernandes",
			URL:      "https://soundcloud.com/oneill-fernandes/sets/you-light-up-my-life",
			Description: "This album is ‘You Light Up My Life’ which is my 53rd Album and contin" +
				"ues with more classic dance hits across the genres and generations and" +
				" contains music from 1971 to 2013. So, there is something in it for ev" +
				"eryone…\n\t\t\t\t\t\n‘Together Forever’ is a song recorded by English singer-" +
				"songwriter Rick Astley and released by RCA and BMG as the fourth singl" +
				"e from his debut album, ‘Whenever You Need Somebody’ (1987). \n\n‘Never " +
				"Gonna Give You Up’ is a pop song by English singer Rick Astley, releas" +
				"ed on 27 July 1987. Written and produced by Stock Aitken Waterman, it " +
				"was released as the first single from Astley's debut studio album, ‘Wh" +
				"enever You Need Somebody’.\n\n‘Wake Me Up’ is a song by Swedish DJ and r" +
				"ecord producer Avicii, released as the lead single from his debut stud" +
				"io album ‘True’ on 17 June 2013. \n\n‘Ain't No Sunshine’ is a song by Bi" +
				"ll Withers, from his 1971 debut album ‘Just As I Am’ and produced by B" +
				"ooker T. Jones.\n\n‘The Lazy Song’ is a song by American singer-songwrit" +
				"er Bruno Mars for his debut studio album, ‘Doo-Wops & Hooligans’ and r" +
				"eleased on February 15, 2011.\n\n‘I Need More of You’ is a song written " +
				"by David Bellamy and recorded by American country music duo The Bellam" +
				"y Brothers. It was released in January 1985 as the third single from t" +
				"he album ‘Restless’. \n\n‘Blame It On The Fire In My Heart’ is a song by" +
				" American country music duo The Bellamy Brothers. It was released in J" +
				"anuary 1992 from the album ‘Beggars and Heroes’. \n\n‘Matrimony’ is a so" +
				"ng by Gilbert O’Sullivan and was released in August 1971 from the albu" +
				"m ‘Himself’. \n\n‘Leave a Light On’ is a song by American singer Belinda" +
				" Carlisle, recorded for her third studio album ‘Runaway Horses’ releas" +
				"ed on Sept 25, 1989. \n\n‘Dreams’ is a song by the British American rock" +
				" band Fleetwood Mac, written and sung by Stevie Nicks for the band's e" +
				"leventh studio album, ‘Rumours’ and released on 24 March 1977.\n\n‘Me an" +
				"d You and a Dog Named Boo’ is the March 1971 debut single by Lobo. Wri" +
				"tten by Lobo under his real name Kent LaVoie, it appears on the ‘Intro" +
				"ducing Lobo’ album.\n\n‘Sky High’ is a song by British band Jigsaw. It w" +
				"as released as a single in 1975 and was the main title theme to the fi" +
				"lm ‘The Man from Hong Kong’. The song was a worldwide hit in the latte" +
				"r part of 1975, reaching No. 3 on the Billboard Hot 100. \n\n‘If You Thi" +
				"nk You Know How to Love Me’ is a song by British rock band Smokie. It " +
				"was first released in June 1975 as a single and appeared on the album " +
				"‘Changing All the Time’. \n\n‘Fast Car’ is the debut single by American " +
				"singer-songwriter Tracy Chapman, released on April 6, 1988, as the lea" +
				"d single from her 1988 self-titled debut studio album.\n\n‘Lost in Franc" +
				"e’ is a song recorded by Welsh singer Bonnie Tyler. It was released as" +
				" a single in September 1976 by RCA Records, written by her producers a" +
				"nd songwriters Ronnie Scott and Steve Wolfe. \n\n‘In These Arms’ is a so" +
				"ng by American rock band Bon Jovi, released on May 3, 1993, as the thi" +
				"rd single from the band's fifth studio album, ‘Keep the Faith’ (1992)." +
				" \n\n‘Baby I'm A Want You’ is the fourth album by Bread, released in Jan" +
				" 1972. Its singles included the title cut which reached No. 3 on the B" +
				"illboard Top 100. \n\n‘Magic Woman Touch’ is a song by The Hollies relea" +
				"sed in 1972 from their studio Album ‘A Selection’\n\n‘A Little Peace’ is" +
				" a song recorded by German singer Nicole, with music composed by Ralph" +
				" Siegel and German lyrics written by Bernd Meinunger. It represented G" +
				"ermany in the Eurovision Song Contest 1982.\n\n‘I Have Always Loved You " +
				"is a song by Enrique Iglesias released in 1999 from the album ‘Enrique" +
				"’.\n\n‘You Light Up My Life’ is the first solo album from singer Debby B" +
				"oone released on Jun 15, 1977, and reached No. 1 on the Billboard Hot " +
				"100.\n\nSo, despite all my health issues and challenges, I present to yo" +
				"u a timeless Album filled with nostalgic hits of the past and present." +
				" \n\nMusic Videos for this album is on YouTube…link is below:\nwww.youtub" +
				"e.com/@ONeillFernandes\n\nEnjoy!!!",
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
