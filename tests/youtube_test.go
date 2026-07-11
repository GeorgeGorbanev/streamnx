package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2"
	"github.com/GeorgeGorbanev/streamnx/v2/tests/fixtures"
)

const (
	youtubeAPIKey = "sample_api_key"

	youtubeTrackID        = "dQw4w9WgXcQ"
	youtubeAlbumID        = "PLtgiPVWmOlIbBLeKL7iH0Bs84cYVF59e8"
	youtubeSearchArtist   = "Rick Astley"
	youtubeSearchTrack    = "Never Gonna Give You Up"
	youtubeSearchAlbum    = "Whenever You Need Somebody"
	youtubeMissingTrackID = "00000000000"
	youtubeMissingAlbumID = "PL00000000000000000000000000000000"
)

func TestYoutubeCatalogFetchTrack(t *testing.T) {
	server := fixtures.NewServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/youtube/v3/videos",
		Query:   map[string]string{"id": youtubeTrackID},
		Status:  http.StatusOK,
		Fixture: "youtube_fetch_track_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			q := r.URL.Query()
			require.Equal(t, youtubeAPIKey, q.Get("key"))
			require.Equal(t, "snippet,contentDetails", q.Get("part"))
			require.Equal(t, youtubeTrackID, q.Get("id"))
		},
	})
	defer server.Close()

	catalog := newYoutubeCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Youtube, youtubeTrackID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Track{
		ID:          youtubeTrackID,
		CoverURL:    "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
		Title:       "Rick Astley - Never Gonna Give You Up (Official Video) (4K Remaster)",
		Provider:    streamnx.Youtube,
		Creator:     "Rick Astley",
		Description: youtubeFetchVideoDescription,
		Duration:    213,
		URL:         "https://www.youtube.com/watch?v=" + youtubeTrackID,
	}, got)
}

func TestYoutubeCatalogFetchTrackNotFound(t *testing.T) {
	server := fixtures.NewServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/youtube/v3/videos",
		Query:   map[string]string{"id": youtubeMissingTrackID},
		Status:  http.StatusOK,
		Fixture: "youtube_fetch_track_not_found_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			q := r.URL.Query()
			require.Equal(t, youtubeAPIKey, q.Get("key"))
			require.Equal(t, "snippet,contentDetails", q.Get("part"))
			require.Equal(t, youtubeMissingTrackID, q.Get("id"))
		},
	})
	defer server.Close()

	catalog := newYoutubeCatalog(t, server.URL)
	got, err := catalog.FetchTrack(t.Context(), streamnx.Youtube, youtubeMissingTrackID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestYoutubeCatalogFetchAlbum(t *testing.T) {
	server := fixtures.NewServer(
		t,
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": youtubeAlbumID},
			Status:  http.StatusOK,
			Fixture: "youtube_fetch_album_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, youtubeAlbumID, q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlistItems",
			Query:   map[string]string{"playlistId": youtubeAlbumID},
			Status:  http.StatusOK,
			Fixture: "youtube_fetch_album_items_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "50", q.Get("maxResults"))
				require.Equal(t, youtubeAlbumID, q.Get("playlistId"))
			},
		},
	)
	defer server.Close()

	catalog := newYoutubeCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Youtube, youtubeAlbumID)

	require.NoError(t, err)
	require.Equal(t, streamnx.Album{
		ID:          youtubeAlbumID,
		CoverURL:    "https://i.ytimg.com/vi/3BFTio5296w/maxresdefault.jpg",
		Title:       "Rick Astley - Whenever You Need Somebody (Deluxe Edition - 2022 Remaster) [Full Album] (Pop)",
		Provider:    streamnx.Youtube,
		Creator:     "Peter Kolya Jr",
		Description: youtubeFetchPlaylistDescription,
		URL:         "https://www.youtube.com/playlist?list=" + youtubeAlbumID,
		TrackIDs:    []string{"firstVideoID", "secondVideoID"},
	}, got)
}

func TestYoutubeCatalogFetchAlbumNotFound(t *testing.T) {
	server := fixtures.NewServer(t, fixtures.Route{
		Method:  http.MethodGet,
		Path:    "/youtube/v3/playlists",
		Query:   map[string]string{"id": youtubeMissingAlbumID},
		Status:  http.StatusOK,
		Fixture: "youtube_fetch_album_not_found_200.json",
		Assert: func(t *testing.T, r *http.Request) {
			q := r.URL.Query()
			require.Equal(t, youtubeAPIKey, q.Get("key"))
			require.Equal(t, "snippet", q.Get("part"))
			require.Equal(t, youtubeMissingAlbumID, q.Get("id"))
		},
	})
	defer server.Close()

	catalog := newYoutubeCatalog(t, server.URL)
	got, err := catalog.FetchAlbum(t.Context(), streamnx.Youtube, youtubeMissingAlbumID)

	require.Zero(t, got)
	require.ErrorIs(t, err, streamnx.ErrNotFound)
}

func TestYoutubeCatalogSearchTracks(t *testing.T) {
	server := fixtures.NewServer(t,
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/search",
			Query:   map[string]string{"type": "video"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_tracks_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, youtubeSearchArtist+" \u2013 "+youtubeSearchTrack, q.Get("q"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "video", q.Get("type"))
				require.Equal(t, "10", q.Get("videoCategoryId"))
				require.Equal(t, "10", q.Get("maxResults"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "dQw4w9WgXcQ"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_01_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "dQw4w9WgXcQ", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "7FwDP17XPlk"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_02_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "7FwDP17XPlk", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "3BFTio5296w"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_03_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "3BFTio5296w", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "QFK09y01me0"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_04_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "QFK09y01me0", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "DLzxrzFCyOs"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_05_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "DLzxrzFCyOs", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "WmttqfVT830"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_06_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "WmttqfVT830", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "LLFhKaqnWwk"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_07_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "LLFhKaqnWwk", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "RL6Jq1hHsco"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_08_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "RL6Jq1hHsco", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "nsCIeklgp1M"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_09_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "nsCIeklgp1M", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/videos",
			Query:   map[string]string{"id": "SbYXkOAoZpI"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_track_candidate_10_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet,contentDetails", q.Get("part"))
				require.Equal(t, "SbYXkOAoZpI", q.Get("id"))
			},
		},
	)
	defer server.Close()

	catalog := newYoutubeCatalog(t, server.URL)
	got, err := catalog.SearchTracks(t.Context(), streamnx.Youtube, streamnx.SearchQuery{
		Artist: youtubeSearchArtist,
		Title:  youtubeSearchTrack,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchTrack{
		{
			ID:          "dQw4w9WgXcQ",
			CoverURL:    "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
			Title:       "Rick Astley - Never Gonna Give You Up (Official Video) (4K Remaster)",
			Provider:    streamnx.Youtube,
			Creator:     "Rick Astley",
			Description: youtubeSearchVideosDescription1,
			URL:         "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		},
		{
			ID:          "7FwDP17XPlk",
			CoverURL:    "https://i.ytimg.com/vi/7FwDP17XPlk/maxresdefault.jpg",
			Title:       "Rick Astley - Never Gonna Give You Up",
			Provider:    streamnx.Youtube,
			Creator:     "Amazing Lyrics",
			Description: youtubeSearchVideosDescription2,
			URL:         "https://www.youtube.com/watch?v=7FwDP17XPlk",
		},
		{
			ID:          "3BFTio5296w",
			CoverURL:    "https://i.ytimg.com/vi/3BFTio5296w/maxresdefault.jpg",
			Title:       "Never Gonna Give You Up (2022 Remaster)",
			Provider:    streamnx.Youtube,
			Creator:     "Rick Astley - Topic",
			Description: youtubeSearchVideosDescription3,
			URL:         "https://www.youtube.com/watch?v=3BFTio5296w",
		},
		{
			ID:          "QFK09y01me0",
			CoverURL:    "https://i.ytimg.com/vi/QFK09y01me0/maxresdefault.jpg",
			Title:       "Rick Astley - Never Gonna Give You Up (Countdown, 1987)",
			Provider:    streamnx.Youtube,
			Creator:     "Countdown",
			Description: youtubeSearchVideosDescription4,
			URL:         "https://www.youtube.com/watch?v=QFK09y01me0",
		},
		{
			ID:          "DLzxrzFCyOs",
			CoverURL:    "https://i.ytimg.com/vi/DLzxrzFCyOs/maxresdefault.jpg",
			Title:       "Rick Astley - Never Gonna Give You Up [HQ]",
			Provider:    streamnx.Youtube,
			Creator:     "AllKindsOfStuff",
			Description: youtubeSearchVideosDescription5,
			URL:         "https://www.youtube.com/watch?v=DLzxrzFCyOs",
		},
		{
			ID:          "WmttqfVT830",
			CoverURL:    "https://i.ytimg.com/vi/WmttqfVT830/maxresdefault.jpg",
			Title:       "Rick Astley - Never Gonna Give You Up (Jeanie Schultheiß) | The Voice of Germany | Blind Audition",
			Provider:    streamnx.Youtube,
			Creator:     "The Voice of Germany - Offiziell",
			Description: youtubeSearchVideosDescription6,
			URL:         "https://www.youtube.com/watch?v=WmttqfVT830",
		},
		{
			ID:          "LLFhKaqnWwk",
			CoverURL:    "https://i.ytimg.com/vi/LLFhKaqnWwk/maxresdefault.jpg",
			Title:       "Rick Astley - Never Gonna Give You Up (Official Animated Video)",
			Provider:    streamnx.Youtube,
			Creator:     "Rick Astley",
			Description: youtubeSearchVideosDescription7,
			URL:         "https://www.youtube.com/watch?v=LLFhKaqnWwk",
		},
		{
			ID:          "RL6Jq1hHsco",
			CoverURL:    "https://i.ytimg.com/vi/RL6Jq1hHsco/maxresdefault.jpg",
			Title:       "InsurAAAnce & Rick Astley - Never Gonna Give You Up 2022 4K/8K 60FPS",
			Provider:    streamnx.Youtube,
			Creator:     "cskillers1",
			Description: youtubeSearchVideosDescription8,
			URL:         "https://www.youtube.com/watch?v=RL6Jq1hHsco",
		},
		{
			ID:          "nsCIeklgp1M",
			CoverURL:    "https://i.ytimg.com/vi/nsCIeklgp1M/maxresdefault.jpg",
			Title:       "Rick Astley - Never Gonna Give You Up | Glastonbury 2023",
			Provider:    streamnx.Youtube,
			Creator:     "BBC Music",
			Description: youtubeSearchVideosDescription9,
			URL:         "https://www.youtube.com/watch?v=nsCIeklgp1M",
		},
		{
			ID:          "SbYXkOAoZpI",
			CoverURL:    "https://i.ytimg.com/vi/SbYXkOAoZpI/sddefault.jpg",
			Title:       "Never Gonna Give You Up (Lyrics) - Rick Astley",
			Provider:    streamnx.Youtube,
			Creator:     "Giovanna Lozano",
			Description: youtubeSearchVideosDescription10,
			URL:         "https://www.youtube.com/watch?v=SbYXkOAoZpI",
		},
	}, got)
}

func TestYoutubeCatalogSearchAlbums(t *testing.T) {
	server := fixtures.NewServer(t,
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/search",
			Query:   map[string]string{"type": "playlist"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_albums_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, youtubeSearchArtist+" \u2013 "+youtubeSearchAlbum, q.Get("q"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "playlist", q.Get("type"))
				require.Equal(t, "10", q.Get("maxResults"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PLtgiPVWmOlIbBLeKL7iH0Bs84cYVF59e8"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_01_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PLtgiPVWmOlIbBLeKL7iH0Bs84cYVF59e8", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PL3aRpGkCm4PEMWuR-3-9oc1ktpNXNpHxk"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_02_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PL3aRpGkCm4PEMWuR-3-9oc1ktpNXNpHxk", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PL_9gWeiShHFHJDgga5SSJCN8eztoilisC"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_03_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PL_9gWeiShHFHJDgga5SSJCN8eztoilisC", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PL5l4oM-gW3urbqHYZfxdYuSFM3F7D9qCU"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_04_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PL5l4oM-gW3urbqHYZfxdYuSFM3F7D9qCU", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PLSYjVXj03FmCkiHQCMCghJSOdOkRNzhAJ"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_05_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PLSYjVXj03FmCkiHQCMCghJSOdOkRNzhAJ", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PLcsWXT-BRUUYx0h3cPPX6Gs3J1hqlypme"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_06_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PLcsWXT-BRUUYx0h3cPPX6Gs3J1hqlypme", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PLMynaxX_I0z_f853CNb7_gDYDN2AN59jH"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_07_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PLMynaxX_I0z_f853CNb7_gDYDN2AN59jH", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PLF-kxtaaTH5UAs61S7du4RmD4lC6RHk4j"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_08_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PLF-kxtaaTH5UAs61S7du4RmD4lC6RHk4j", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PLgaFNC_I_ZkkyIFxmCEXnRZSp0xWwbTCU"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_09_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PLgaFNC_I_ZkkyIFxmCEXnRZSp0xWwbTCU", q.Get("id"))
			},
		},
		fixtures.Route{
			Method:  http.MethodGet,
			Path:    "/youtube/v3/playlists",
			Query:   map[string]string{"id": "PLtgiPVWmOlIYtYfnnutPhJviSrv5Bk8_W"},
			Status:  http.StatusOK,
			Fixture: "youtube_search_album_candidate_10_200.json",
			Assert: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				require.Equal(t, youtubeAPIKey, q.Get("key"))
				require.Equal(t, "snippet", q.Get("part"))
				require.Equal(t, "PLtgiPVWmOlIYtYfnnutPhJviSrv5Bk8_W", q.Get("id"))
			},
		},
	)
	defer server.Close()

	catalog := newYoutubeCatalog(t, server.URL)
	got, err := catalog.SearchAlbums(t.Context(), streamnx.Youtube, streamnx.SearchQuery{
		Artist: youtubeSearchArtist,
		Title:  youtubeSearchAlbum,
	})

	require.NoError(t, err)
	require.Equal(t, []streamnx.SearchAlbum{
		{
			ID:          "PLtgiPVWmOlIbBLeKL7iH0Bs84cYVF59e8",
			CoverURL:    "https://i.ytimg.com/vi/3BFTio5296w/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody (Deluxe Edition - 2022 Remaster) [Full Album] (Pop)",
			Provider:    streamnx.Youtube,
			Creator:     "Peter Kolya Jr",
			Description: youtubeSearchPlaylistsDescription1,
			URL:         "https://www.youtube.com/playlist?list=PLtgiPVWmOlIbBLeKL7iH0Bs84cYVF59e8",
		},
		{
			ID:          "PL3aRpGkCm4PEMWuR-3-9oc1ktpNXNpHxk",
			CoverURL:    "https://i.ytimg.com/vi/lYBUbBu4W08/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody (Album)",
			Provider:    streamnx.Youtube,
			Creator:     "MCADCUVIDEO",
			Description: youtubeSearchPlaylistsDescription2,
			URL:         "https://www.youtube.com/playlist?list=PL3aRpGkCm4PEMWuR-3-9oc1ktpNXNpHxk",
		},
		{
			ID:          "PL_9gWeiShHFHJDgga5SSJCN8eztoilisC",
			CoverURL:    "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody (Full Album)",
			Provider:    streamnx.Youtube,
			Creator:     "Sony Music UK",
			Description: youtubeSearchPlaylistsDescription3,
			URL:         "https://www.youtube.com/playlist?list=PL_9gWeiShHFHJDgga5SSJCN8eztoilisC",
		},
		{
			ID:          "PL5l4oM-gW3urbqHYZfxdYuSFM3F7D9qCU",
			CoverURL:    "https://i.ytimg.com/vi/3BFTio5296w/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody (Deluxe Edition - 2022 Remaster) (Album)",
			Provider:    streamnx.Youtube,
			Creator:     "Denny Li",
			Description: youtubeSearchPlaylistsDescription4,
			URL:         "https://www.youtube.com/playlist?list=PL5l4oM-gW3urbqHYZfxdYuSFM3F7D9qCU",
		},
		{
			ID:          "PLSYjVXj03FmCkiHQCMCghJSOdOkRNzhAJ",
			CoverURL:    "https://i.ytimg.com/vi/B77NF290LDA/maxresdefault.jpg",
			Title:       "Rick Astley (1987) Whenever You Need Somebody",
			Provider:    streamnx.Youtube,
			Creator:     "Zsolt Rajkó",
			Description: youtubeSearchPlaylistsDescription5,
			URL:         "https://www.youtube.com/playlist?list=PLSYjVXj03FmCkiHQCMCghJSOdOkRNzhAJ",
		},
		{
			ID:          "PLcsWXT-BRUUYx0h3cPPX6Gs3J1hqlypme",
			CoverURL:    "https://i.ytimg.com/vi/qwLQkLTOmqo/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody - 432hz",
			Provider:    streamnx.Youtube,
			Creator:     "Hertz Patrol",
			Description: youtubeSearchPlaylistsDescription6,
			URL:         "https://www.youtube.com/playlist?list=PLcsWXT-BRUUYx0h3cPPX6Gs3J1hqlypme",
		},
		{
			ID:          "PLMynaxX_I0z_f853CNb7_gDYDN2AN59jH",
			CoverURL:    "https://i.ytimg.com/vi/lYBUbBu4W08/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody",
			Provider:    streamnx.Youtube,
			Creator:     "Topic Albums Playlists Only",
			Description: youtubeSearchPlaylistsDescription7,
			URL:         "https://www.youtube.com/playlist?list=PLMynaxX_I0z_f853CNb7_gDYDN2AN59jH",
		},
		{
			ID:          "PLF-kxtaaTH5UAs61S7du4RmD4lC6RHk4j",
			CoverURL:    "https://i.ytimg.com/vi/lYBUbBu4W08/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody (1987)",
			Provider:    streamnx.Youtube,
			Creator:     "Jennie Music Archive",
			Description: youtubeSearchPlaylistsDescription8,
			URL:         "https://www.youtube.com/playlist?list=PLF-kxtaaTH5UAs61S7du4RmD4lC6RHk4j",
		},
		{
			ID:          "PLgaFNC_I_ZkkyIFxmCEXnRZSp0xWwbTCU",
			CoverURL:    "https://i.ytimg.com/vi/3BFTio5296w/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody (Full Album, 1987/Album Completo)",
			Provider:    streamnx.Youtube,
			Creator:     "Hits & Best Songs Playlist ",
			Description: youtubeSearchPlaylistsDescription9,
			URL:         "https://www.youtube.com/playlist?list=PLgaFNC_I_ZkkyIFxmCEXnRZSp0xWwbTCU",
		},
		{
			ID:          "PLtgiPVWmOlIYtYfnnutPhJviSrv5Bk8_W",
			CoverURL:    "https://i.ytimg.com/vi/lYBUbBu4W08/maxresdefault.jpg",
			Title:       "Rick Astley - Whenever You Need Somebody (1987) (Album) (Pop)",
			Provider:    streamnx.Youtube,
			Creator:     "Peter Kolya Jr",
			Description: youtubeSearchPlaylistsDescription10,
			URL:         "https://www.youtube.com/playlist?list=PLtgiPVWmOlIYtYfnnutPhJviSrv5Bk8_W",
		},
	}, got)
}

func TestYoutubeCatalogRejectsUnsupportedOperation(t *testing.T) {
	catalog := newYoutubeCatalog(t, "http://127.0.0.1")

	gotType, gotID, err := catalog.Uncloak(t.Context(), streamnx.Youtube, "sample")

	require.Zero(t, gotType)
	require.Zero(t, gotID)
	require.ErrorIs(t, err, streamnx.ErrUnsupportedOperation)
}

func newYoutubeCatalog(t *testing.T, serverURL string) *streamnx.Catalog {
	t.Helper()

	catalog, err := streamnx.NewCatalog(streamnx.WithYoutube(
		streamnx.YoutubeCredentials{APIKey: youtubeAPIKey},
		streamnx.WithYoutubeAPIURL(serverURL),
	))
	require.NoError(t, err)

	return catalog
}

const (
	youtubeFetchVideoDescription = `The official video for “Never Gonna Give You Up” by Rick Astley. 

Never: The Autobiography 📚 OUT NOW! 
Follow this link to get your copy and listen to Rick’s ‘Never’ playlist ❤️ #RickAstleyNever
https://linktr.ee/rickastleynever

“Never Gonna Give You Up” was a global smash on its release in July 1987, topping the charts in 25 countries including Rick’s native UK and the US Billboard Hot 100.  It also won the Brit Award for Best single in 1988. Stock Aitken and Waterman wrote and produced the track which was the lead-off single and lead track from Rick’s debut LP “Whenever You Need Somebody”.  The album was itself a UK number one and would go on to sell over 15 million copies worldwide.

The legendary video was directed by Simon West – who later went on to make Hollywood blockbusters such as Con Air, Lara Croft – Tomb Raider and The Expendables 2.  The video passed the 1bn YouTube views milestone on 28 July 2021.

Subscribe to the official Rick Astley YouTube channel: https://RickAstley.lnk.to/YTSubID

Follow Rick Astley:
Facebook: https://RickAstley.lnk.to/FBFollowID 
Twitter: https://RickAstley.lnk.to/TwitterID 
Instagram: https://RickAstley.lnk.to/InstagramID 
Website: https://RickAstley.lnk.to/storeID 
TikTok: https://RickAstley.lnk.to/TikTokID

Listen to Rick Astley:
Spotify: https://RickAstley.lnk.to/SpotifyID 
Apple Music: https://RickAstley.lnk.to/AppleMusicID 
Amazon Music: https://RickAstley.lnk.to/AmazonMusicID 
Deezer: https://RickAstley.lnk.to/DeezerID 

Lyrics:
We’re no strangers to love
You know the rules and so do I
A full commitment’s what I’m thinking of
You wouldn’t get this from any other guy

I just wanna tell you how I’m feeling
Gotta make you understand

Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

We’ve known each other for so long
Your heart’s been aching but you’re too shy to say it
Inside we both know what’s been going on
We know the game and we’re gonna play it

And if you ask me how I’m feeling
Don’t tell me you’re too blind to see

Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

#RickAstley #NeverGonnaGiveYouUp #WheneverYouNeedSomebody #OfficialMusicVideo`

	youtubeFetchPlaylistDescription = ``

	youtubeSearchVideosDescription1 = `The official video for “Never Gonna Give You Up” by Rick Astley. 

Never: The Autobiography 📚 OUT NOW! 
Follow this link to get your copy and listen to Rick’s ‘Never’ playlist ❤️ #RickAstleyNever
https://linktr.ee/rickastleynever

“Never Gonna Give You Up” was a global smash on its release in July 1987, topping the charts in 25 countries including Rick’s native UK and the US Billboard Hot 100.  It also won the Brit Award for Best single in 1988. Stock Aitken and Waterman wrote and produced the track which was the lead-off single and lead track from Rick’s debut LP “Whenever You Need Somebody”.  The album was itself a UK number one and would go on to sell over 15 million copies worldwide.

The legendary video was directed by Simon West – who later went on to make Hollywood blockbusters such as Con Air, Lara Croft – Tomb Raider and The Expendables 2.  The video passed the 1bn YouTube views milestone on 28 July 2021.

Subscribe to the official Rick Astley YouTube channel: https://RickAstley.lnk.to/YTSubID

Follow Rick Astley:
Facebook: https://RickAstley.lnk.to/FBFollowID 
Twitter: https://RickAstley.lnk.to/TwitterID 
Instagram: https://RickAstley.lnk.to/InstagramID 
Website: https://RickAstley.lnk.to/storeID 
TikTok: https://RickAstley.lnk.to/TikTokID

Listen to Rick Astley:
Spotify: https://RickAstley.lnk.to/SpotifyID 
Apple Music: https://RickAstley.lnk.to/AppleMusicID 
Amazon Music: https://RickAstley.lnk.to/AmazonMusicID 
Deezer: https://RickAstley.lnk.to/DeezerID 

Lyrics:
We’re no strangers to love
You know the rules and so do I
A full commitment’s what I’m thinking of
You wouldn’t get this from any other guy

I just wanna tell you how I’m feeling
Gotta make you understand

Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

We’ve known each other for so long
Your heart’s been aching but you’re too shy to say it
Inside we both know what’s been going on
We know the game and we’re gonna play it

And if you ask me how I’m feeling
Don’t tell me you’re too blind to see

Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

#RickAstley #NeverGonnaGiveYouUp #WheneverYouNeedSomebody #OfficialMusicVideo`

	youtubeSearchVideosDescription2 = `Rick Astley - Never Gonna Give You Up
Get:

Lyrics:

[Intro]
Desert you
Ooh-ooh-ooh-ooh
Hurt you

[Verse 1]
We're no strangers to love
You know the rules and so do I (Do I)
A full commitment's what I'm thinking of
You wouldn't get this from any other guy

[Pre-Chorus]
I just wanna tell you how I'm feeling
Gotta make you understand

[Chorus]
Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

[Verse 2]
We've known each other for so long
Your heart's been aching, but you're too shy to say it (To say it)
Inside, we both know what's been going on (Going on)
We know the game, and we're gonna play it

[Pre-Chorus]
And if you ask me how I'm feeling
Don't tell me you're too blind to see

[Chorus]
Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you
Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

[Bridge]
Ooh (Give you up)
Ooh-ooh (Give you up)
Ooh-ooh
Never gonna give, never gonna give (Give you up)
Ooh-ooh
Never gonna give, never gonna give (Give you up)
[Verse 2]
We've known each other for so long
Your heart's been aching, but you're too shy to say it (To say it)
Inside, we both know what's been going on (Going on)
We know the game, and we're gonna play it

[Pre-Chorus]
I just wanna tell you how I'm feeling
Gotta make you understand

[Chorus]
Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you
Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you
Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you`

	youtubeSearchVideosDescription3 = `Provided to YouTube by BMG Rights Management (UK) Ltd.

Never Gonna Give You Up (2022 Remaster) · Rick Astley

Whenever You Need Somebody

℗ 2022 BMG Rights Management (UK) Limited

Released on: 2022-05-20

Keyboards, Producer: Mike Stock
Guitar, Producer: Matt Aitken
Producer: Pete Waterman
Background  Vocals: Dee Lewis
Background  Vocals: Shirley Lewis
Background  Vocals: Mae McKenna
Background  Vocals: Suzanne Rhatigan
Mixer: Pete Hammond
Sound  Engineer: Mark McGuire
Composer: Mike Stock
Composer: Matt Aitken
Composer: Pete Waterman

Auto-generated by YouTube.`

	youtubeSearchVideosDescription4 = `TV performance by Rick Astley with ‘Never Gonna Give You Up’ on Dutch popshow “Countdown”.  The television show was broadcasted by Veronica from 1976 to 1993 and was "Europe's Number 1 Rock Show" with performances by the greatest artists from all over the world like Paul McCartney, David Bowie, Tina Turner, Janet Jackson, Stevie Wonder, The Rolling Stones, George Harrison, U2, The Police, Eric Clapton, Bruce Springsteen, Lenny Kravitz, Depeche Mode, New Order, Whitney Houston, Duran Duran, R.E.M., Frankie Goes To Hollywood, The Cure, Run DMC, Cyndi Lauper, Iron Maiden, New Kids On The Block, LL Cool J.

 • Subscribe to our channel: https://bit.ly/3kepHLP

• Playlist met Nederlandstalige acts bij Countdown: https://bit.ly/3qJ0NFQ
• Playlist met 70's Countdown clips: https://bit.ly/3bBwY4L
• Playlist met 80's Countdown clips: https://bit.ly/3qVyKE5
• Playlist met 90's Countdown clips: https://bit.ly/3kmXiUf

====
OVER COUNTDOWN
====
Het Nederlandse televisieprogramma "Countdown" was van 1976 t/m 1993 “Europe’s Number 1 Rock Show” en werd uitgezonden bij Veronica. Later werd het programma ook uitgezonden in het buitenland en was het op z'n hoogtepunt in meer dan 20 landen op televisie te zien.

Gedurende deze 17 jaar muziekgeschiedenis zijn er opnames gemaakt met talloze internationale sterren, zoals Paul McCartney, David Bowie, Tina Turner, Janet Jackson, Stevie Wonder, The Rolling Stones, George Harrison, U2, The Police, Eric Clapton, Bruce Springsteen, Lenny Kravitz, Depeche Mode, New Order, Whitney Houston, Duran Duran, R.E.M., Frankie Goes To Hollywood, The Cure, Run DMC, Cyndi Lauper, Iron Maiden, New Kids On The Block, LL Cool J.

Uiteraard kwamen Nederlandse top-acts als Doe Maar, Andre Hazes, Herman Brood, Frank Boeijen, Spargo, Golden Earring, Time Bandits,  De Dijk, De Raggende Manne, Hans de Booij, Gerard Joling, Drukwerk ook langs voor een optreden in Bussum.

Countdown werd gepresenteerd door oa Lex Harding, Erik de Zwart, Adam Curry, Simone Walraven Wessel van Diepen, Jeroen van Inkel, Rob Stenders, Simone Angel en Jasper Faber.

Het Countdown archief is vrijwel helemaal compleet en sinds 2017 eigendom van het Amerikaanse Reelin’ In The Years Productions (David Peck - https://reelinintheyears.com/) en Double 2 Hilversum (Jan Douwe Kroeske - http://www.double2bv.nl). Het archief bestaat uit meer dan 3.000 analoge mastertapes die allemaal gedigitaliseerd en gearchiveerd zijn. Indien u gebruik wilt maken van beeldmateriaal in een documentaire, tv-programma of andere mediaproductie neemt u voor een licentiedeal contact op met Double 2 BV via info@d2bv.nl. Beeldmateriaal wordt enkel verstrekt aan zakelijke projecten en niet aan particulieren.

===
#RickAstley`

	youtubeSearchVideosDescription5 = `Artist: Rick Astley
Title: Never Gonna Give You Up
Difference with original: nothing but the quality of audio and video.

all rights belong to their respective owners and publishers. I do not own any rights and remastering/new quality was solely done from original sources. No copyright infringment is intended and the content owner (PWL/Rick Astley/respective owners) is recognized & acknownledged.

(C) 1987 PWL`

	youtubeSearchVideosDescription6 = `Mit so viel Gefühl haben wir den 80er-Hit von Rick Astley "Never Gonna Give You Up" noch nie gehört! Jeanie Schultheiß verzaubert nicht nur das Publikum zu "Never Gonna Give You Up", sondern auch die Coaches. Ist dieses Cover besser als das Original von Rick Astley?

#RickAstley #NeverGonnaGiveYouUp #Voice

►Ganze Folgen: http://bit.ly/TVOG2018GANZEFOLGEN
►Coach Stories: http://bit.ly/TVOGCoachStories
►Vorab-Blind: http://bit.ly/BlindPreview
►Unplugged-Clips: http://bit.ly/UnpluggedClips

► THE VOICE OF GERMANY 2018: http://bit.ly/TVOG2018
► BLIND AUDITIONS 2018: http://bit.ly/TVOG2018Blinds

► THE VOICE OF GERMANY 2017: http://bit.ly/TVOG2017
► BATTLES 2017: http://bit.ly/TVOG2017Battles  
► BLIND AUDITIONS 2017: http://bit.ly/TVOG2017Blinds 

► SUBSCRIBE TO THE VOICE OF GERMANY:
 https://www.youtube.com/user/VoiceOfGermanyTVOG?sub_confirmation=1

► LOVESONGS: http://bit.ly/TVOGBalladen  
► POP MUSIC: http://bit.ly/TVOGPop  
► ROCK MUSIC: http://bit.ly/ROCKTVOG 

The Voice of Germany: die neue einzigartige Casting Show, in der nur echte Künstler mit großartiger Stimme gesucht werden, ist wieder da. Immer Donnerstag 20.15 Uhr auf #ProSieben und sonntags in SAT.1! Ein Hoch auf diese Coaches: Mark Forster kämpft zusammen mit Yvonne Catterfeld, Michael Patrick Kelly, Smudo & Michi Beck um die besten Stimmen bei The Voice of Germany 2018! #TVOG

Germany’s most successful music TV Show is back! – The Voice of Germany 2018 airs Thursday at 20.15 pm on ProSieben and Sunday on SAT.1! New season, new coaches: Mark Forster is fighting together with Yvonne Catterfeld, Michael Patrick Kelly, Smudo & Michi Beck for the best voice in Germany! 

********************************************
BEST OF THE VOICE OF GERMANY auf einer CD*: http://amzn.to/12BwUwc
The Voice of Germany 2 Wii Game*: http://amzn.to/1vXLbNI
********************************************
*Links sind Affiliate Links

#TVOG #TheVoiceOfGermany2018 #MarkForster #MichaelPatrickKelly #YvonneCatterfeld #Smudo #MichiBeck #TVShow

Impressum: https://www.prosieben.de/service/impressum/`

	youtubeSearchVideosDescription7 = `The official animated video for "Never Gonna Give You Up” by Rick Astley
Taken from the album ‘Whenever You Need Somebody’ – deluxe 2CD and digital deluxe out now. Buy/stream here – https://RickAstley.lnk.to/WYNS2022ID

A 50/50 Media House Production 
Animation: Bruna Pias
Producer: Lene Bausager
©50/50MediaHouse2022

Subscribe to the official Rick Astley YouTube channel: https://RickAstley.lnk.to/YTSubID 

Follow Rick Astley:
Facebook: https://RickAstley.lnk.to/FBFollowID  
Twitter: https://RickAstley.lnk.to/TwitterID  
Instagram: https://RickAstley.lnk.to/InstagramID 
Website: https://RickAstley.lnk.to/storeID 
TikTok: https://RickAstley.lnk.to/TikTokID

Listen to Rick Astley:
Spotify: https://RickAstley.lnk.to/SpotifyID 
Apple Music: https://RickAstley.lnk.to/AppleMusicID 
Amazon Music: https://RickAstley.lnk.to/AmazonMusicID 
Deezer: https://RickAstley.lnk.to/DeezerID 

Lyrics:
We’re no strangers to love
You know the rules and so do I
A full commitment’s what I’m thinking of
You wouldn’t get this from any other guy

I just wanna tell you how I’m feeling
Gotta make you understand

Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

We’ve known each other for so long
Your heart’s been aching but you’re too shy to say it
Inside we both know what’s been going on
We know the game and we’re gonna play it

And if you ask me how I’m feeling
Don’t tell me you’re too blind to see

Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you

#RickAstley #NeverGonnaGiveYouUp #WheneverYouNeedSomebody #OfficialAnimatedVideo`

	youtubeSearchVideosDescription8 = `this is just an edit, all the credits going to -
@csaainsurancegroup
@RickAstleyYT
watch the original here - https://youtu.be/GtL1huin9EE
official shortcut to the original - http://bit.ly/get3a
shortcut to this edit - TBA
also check out the first known edit of 2022 version by @DOCOKOCO - https://youtu.be/tgTUtfb0Ok8

i really wanted to include the ending but couldn't really fit it in, will do this instead, thank you for giving us 2022 version -

AAA.com/Insurance
"Our legendary service is never gonna let you down.  Start saving today on auto and home coverage that would never run around and desert you. Because with 100 years of experience, InsurAAAnce is no stranger to love or exceptional service."
#RickAstley #Rickroll  #RR #InsurAAAnce #CSAAInsuranceGroup`

	youtubeSearchVideosDescription9 = `Rick Astley performs Never Gonna Give You Up at Glastonbury 2023. 

Watch more highlights on BBC iPlayer: https://www.bbc.co.uk/iplayer/episodes/b007r6vx/glastonbury
Listen to sets and highlights on BBC Sounds: https://www.bbc.co.uk/sounds/curation/m001n97q
Follow all on social @BBCGlastonbury.`

	youtubeSearchVideosDescription10 = `Music video and lyrics Never Gonna Give You Up of Rick Astley

We're no strangers to love 
You know the rules and so do I 
A full commitment's what I'm thinking of 
You wouldn't get this from any other guy 
I just wanna tell you how I'm feeling 
Gotta make you understand 

CHORUS 
Never gonna give you up, 
Never gonna let you down 
Never gonna run around and desert you 
Never gonna make you cry, 
Never gonna say goodbye 
Never gonna tell a lie and hurt you 

We've known each other for so long 
Your heart's been aching but you're too shy to say it 
Inside we both know what's been going on 
We know the game and we're gonna play it 
And if you ask me how I'm feeling 
Don't tell me you're too blind to see (CHORUS) 

CHORUSCHORUS 
(Ooh give you up) 
(Ooh give you up) 
(Ooh) never gonna give, never gonna give 
(give you up) 
(Ooh) never gonna give, never gonna give 
(give you up) 

We've known each other for so long 
Your heart's been aching but you're too shy to say it 
Inside we both know what's been going on 
We know the game and we're gonna play it (TO FRONT)`

	youtubeSearchPlaylistsDescription1 = ``

	youtubeSearchPlaylistsDescription2 = `MCADUSAUDIO`

	youtubeSearchPlaylistsDescription3 = `Listen On Spotify: http://smarturl.it/AstleySpotify
Buy On iTunes: http://smarturl.it/AstleyGHiTunes
Amazon: http://smarturl.it/AstleyGHAmazon`

	youtubeSearchPlaylistsDescription4 = ``

	youtubeSearchPlaylistsDescription5 = ``

	youtubeSearchPlaylistsDescription6 = `Artist: Rick Astley
Album: Whenever You Need Somebody

432hz

No copyright infringement intended.
All images and audio belong to Rick Astley.`

	youtubeSearchPlaylistsDescription7 = ``

	youtubeSearchPlaylistsDescription8 = `Tracks 1-4, 6, 10, 12-14: Produced by Mike Stock, Matt Aitken, and Pete Waterman (Stock Aitken Waterman)

Tracks 5, 7, 11: Produced by Phil Harding & Ian Curnow (Harding/Curnow)

Tracks 8 & 9: Produced by Daize Washbourn`

	youtubeSearchPlaylistsDescription9 = `SUBSCRIBE: https://youtube.com/channel/UCCh7iY0jKvjszpQFy9HykAw

RICK ASTLEY GREATEST HITS (GRANDES EXITOS)...

RICK ASTLEY DISCOGRAFIA COMPLETA...


(11 de marzo del 2026)

TRACKLIST: CANCIONES: SONGS:

1. Never gonna give you up

Album Completo, Disco Completo, CD Completo, Full Album, DVD Completo, Ep Completo, Vinilo Completo, Lp completo, Vinilo, Nuevo Disco, New Album, Album De Duetos, Concierto Completo, Full Concert, OST, Original Motion Picture`

	youtubeSearchPlaylistsDescription10 = ``
)
