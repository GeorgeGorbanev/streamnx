package tests

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2"
	"github.com/GeorgeGorbanev/streamnx/v2/tests/fixtures"
)

func TestYoutubeUnofficialCatalog(t *testing.T) {
	const (
		prefix      = "youtube_unofficial_"
		longAlbumID = "UUYO_jab_esuFRV4b17AJtAw"
	)
	trackCandidates := strings.Fields(`
		dQw4w9WgXcQ 3BFTio5296w 5SZYz7lZRRI LLFhKaqnWwk 6PLatPMoxGw wuggc60TVzU
		XGxIE1hr0w4 RrESvSRNpeo QpbhSlce_ek nsCIeklgp1M
	`)
	albumCandidates := strings.Fields(`
		PLtgiPVWmOlIbBLeKL7iH0Bs84cYVF59e8 PL_9gWeiShHFHJDgga5SSJCN8eztoilisC PLiVYJJzaD4wv3D5Jw-odTdp4Bv0jRLHHD PL5l4oM-gW3urbqHYZfxdYuSFM3F7D9qCU PLSYjVXj03FmCkiHQCMCghJSOdOkRNzhAJ PLFAcddgaFN8zVPj0vI8gdzrcNmNQMcAhZ
		PL5l4oM-gW3uoIWCbP2rlhrZSbej0mxi7u PLum8kNIteiJwvfpZ2gS_Xqy3elRDja3Iq PLDT8XUkhfwIu0V5LByXwzOygnB4L6kNXN PLtgiPVWmOlIYtYfnnutPhJviSrv5Bk8_W
	`)

	page := func(path, key, value, fixture string) fixtures.Route {
		return fixtures.Route{
			Method:  http.MethodGet,
			Path:    path,
			Query:   map[string]string{key: value, "hl": "en"},
			Status:  http.StatusOK,
			Fixture: prefix + fixture + ".html",
		}
	}
	search := func(filter, query, fixture string) fixtures.Route {
		route := page("/results", "sp", filter, fixture)
		route.Query["search_query"] = query
		return route
	}
	continuation := func(token, fixture string) fixtures.Route {
		return fixtures.Route{
			Method:   http.MethodPost,
			Path:     "/youtubei/v1/browse",
			JSONBody: map[string]string{"continuation": token},
			Status:   http.StatusOK,
			Fixture:  prefix + fixture + ".json",
			Assert: func(t *testing.T, r *http.Request) {
				var request map[string]json.RawMessage
				require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
				require.JSONEq(t, `{"client":{"clientName":"WEB","clientVersion":"2.20261105.05.01","hl":"en","gl":"GE"}}`, string(request["context"]))
			},
		}
	}
	var longAlbumRequests atomic.Int32
	longAlbumPage := page("/playlist", "list", longAlbumID, "fetch_long_album_200")
	longAlbumPage.Assert = func(t *testing.T, _ *http.Request) {
		t.Helper()
		longAlbumRequests.Add(1)
	}
	routes := []fixtures.Route{
		page("/watch", "v", youtubeTrackID, "fetch_track_200"),
		page("/watch", "v", youtubeMissingTrackID, "fetch_track_missing_200"),
		page("/playlist", "list", youtubeAlbumID, "fetch_album_200"),
		longAlbumPage,
		page("/playlist", "list", youtubeMissingAlbumID, "fetch_album_missing_200"),
		search("EgIQAfABAQ==", youtubeSearchArtist+" – "+youtubeSearchTrack, "search_tracks_200"),
		search("EgIQAw==", youtubeSearchArtist+" – "+youtubeSearchAlbum, "search_albums_200"),
		continuation("4qmFsgJtEhpWTFVVWU9famFiX2VzdUZSVjRiMTdBSnRBdxo0Q0FGNkhsQlVPa05IVVdsRlJHUkdUVlJLUWsxVlRUSk9WVXBFVDBSVk5GSnFTUSUzRCUzRJoCGFVVWU9famFiX2VzdUZSVjRiMTdBSnRBdw%3D%3D", "fetch_long_album_continuation_01_200"),
		continuation("4qmFsgJrEhpWTFVVWU9famFiX2VzdUZSVjRiMTdBSnRBdxoyQ0FKNkgxQlVPa05OWjBKSmFFRjNVa1JrUjAxNldUUk5SVXBGVFVSSk5WRnFUVEElM0SaAhhVVVlPX2phYl9lc3VGUlY0YjE3QUp0QXc%3D", "fetch_long_album_continuation_02_200"),
	}
	for _, id := range trackCandidates[1:] {
		routes = append(routes, page("/watch", "v", id, "search_track_candidate_"+id+"_200"))
	}
	for _, id := range albumCandidates[1:] {
		routes = append(routes, page("/playlist", "list", id, "search_album_candidate_"+id+"_200"))
	}
	server := fixtures.NewServer(t, routes...)
	defer server.Close()
	catalog, err := streamnx.NewCatalog(streamnx.WithYoutubeUnofficial(
		streamnx.WithYoutubeUnofficialAPIURL(server.URL),
		streamnx.WithYoutubeUnofficialHTTPClient(server.Client()),
	))
	require.NoError(t, err)

	track, err := catalog.FetchTrack(t.Context(), streamnx.Youtube, youtubeTrackID)
	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:          youtubeTrackID,
		Title:       "Rick Astley - Never Gonna Give You Up (Official Video) (4K Remaster)",
		Creator:     "Rick Astley",
		Provider:    streamnx.Youtube,
		CoverURL:    "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
		URL:         "https://www.youtube.com/watch?v=" + youtubeTrackID,
		Duration:    213,
		Description: track.Description,
	}, track)
	require.Contains(t, track.Description, "The official video for “Never Gonna Give You Up” by Rick Astley.")

	album, err := catalog.FetchAlbum(t.Context(), streamnx.Youtube, youtubeAlbumID)
	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:          youtubeAlbumID,
		Title:       "Rick Astley - Whenever You Need Somebody (Deluxe Edition - 2022 Remaster) [Full Album] (Pop)",
		Creator:     "Peter Kolya Jr",
		Provider:    streamnx.Youtube,
		Description: "",
		CoverURL:    "https://i.ytimg.com/vi/3BFTio5296w/hqdefault.jpg?sqp=-oaymwEXCOADEI4CSFryq4qpAwkIARUAAIhCGAE=&rs=AOn4CLCuMf1nxchCCWwrFV4ZES_h_xrN5A&days_since_epoch=20765",
		URL:         "https://www.youtube.com/playlist?list=" + youtubeAlbumID,
		TrackIDs: strings.Fields(`
		3BFTio5296w JzH2j4ifKQ8 Mh_nFI56oWc W275ib34fpk 68tVxlqYxmc un09kY9Sqmg
		P4GhwTo6NkI r3GR4xDQ0xI oN42i5s4tlU IRdvVILjjkA zAhPlNANuoc FMyJyWUzfno
		6cByHdbxM0Q 3blMEszLPMU Ush4jEfcSVA zMxZJQuloOw D59GnzlsSt8 -XakJZMWjmM
		w5e0FKkiYmk qwuwCdq2-qg 4JFafbrQ6as p_en_Lw56jk eisF4qMYybs n3jMY98KIEw
		NZZ32f96XO0 mq6XiV1hKtw FMvRxFkVubw U4VXDTZ5VKA 5DEXtyAD43I _QSx57P9gwE
		RZ-e-s6U848 M1r0i7m12g4 nZKaxKFK97Q 0k7BvzQRrOI BTpGEzC9TnU -bsvbimLuhg
	`),
	}, album)

	longAlbum, err := catalog.FetchAlbum(t.Context(), streamnx.Youtube, longAlbumID)
	require.NoError(t, err)
	require.Equal(t, int32(1), longAlbumRequests.Load())
	require.Equal(t, "Uploads from 3Blue1Brown", longAlbum.Title)
	require.Len(t, longAlbum.TrackIDs, 247)
	require.Equal(t, []string{"V_uq3v6xK-4", "TaopSi3Ucfc", "BIbHGfdjc_g", "S9JGmA5_unY", "bBC-nXj3Ng4", "F_0yfvm0UoU"}, []string{
		longAlbum.TrackIDs[0], longAlbum.TrackIDs[99], longAlbum.TrackIDs[100],
		longAlbum.TrackIDs[199], longAlbum.TrackIDs[200], longAlbum.TrackIDs[246],
	})

	tracks, err := catalog.SearchTracks(t.Context(), streamnx.Youtube,
		streamnx.SearchQuery{Artist: youtubeSearchArtist, Title: youtubeSearchTrack})
	require.NoError(t, err)
	trackIDs := make([]string, 0, len(tracks))
	for _, result := range tracks {
		trackIDs = append(trackIDs, result.ID)
	}
	require.Equal(t, trackCandidates, trackIDs)
	require.Equal(t, track.Description, tracks[0].Description)

	albums, err := catalog.SearchAlbums(t.Context(), streamnx.Youtube,
		streamnx.SearchQuery{Artist: youtubeSearchArtist, Title: youtubeSearchAlbum})
	require.NoError(t, err)
	albumIDs := make([]string, 0, len(albums))
	for _, result := range albums {
		albumIDs = append(albumIDs, result.ID)
	}
	require.Equal(t, albumCandidates, albumIDs)
	require.Equal(t, album.Creator, albums[0].Creator)

	_, err = catalog.FetchTrack(t.Context(), streamnx.Youtube, youtubeMissingTrackID)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
	_, err = catalog.FetchAlbum(t.Context(), streamnx.Youtube, youtubeMissingAlbumID)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}
