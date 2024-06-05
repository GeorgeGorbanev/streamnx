package utils

import (
	"net/http"
	"net/http/httptest"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
)

var (
	YoutubeAPIKey = "sampleAPIKey"
)

func NewYoutubeAPIServerMock(fm *fixture.FixturesMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != YoutubeAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var response []byte
		var ok bool

		switch r.URL.Path {
		case "/youtube/v3/videos":
			trackID := r.URL.Query().Get("id")
			if response, ok = fm.YoutubeTracks[trackID]; !ok {
				response = fixture.Read("youtube/not_found.json")
			}
		case "/youtube/v3/search":
			query := r.URL.Query().Get("q")
			searchType := r.URL.Query().Get("type")

			switch searchType {
			case "video":
				if response, ok = fm.YoutubeSearchTracks[query]; !ok {
					response = fixture.Read("youtube/not_found.json")
				}
			case "playlist":
				if response, ok = fm.YoutubeSearchAlbums[query]; !ok {
					response = fixture.Read("youtube/not_found.json")
				}
			default:
				panic("unexpected search type")
			}
		case "/youtube/v3/playlists":
			albumID := r.URL.Query().Get("id")

			if response, ok = fm.YoutubeAlbums[albumID]; !ok {
				response = fixture.Read("youtube/not_found.json")
			}
		case "/youtube/v3/playlistItems":
			albumID := r.URL.Query().Get("playlistId")

			if response, ok = fm.YoutubePlaylistItems[albumID]; !ok {
				response = fixture.Read("youtube/not_found.json")
			}
		default:
			panic("unexpected request")
		}

		_, err := w.Write(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}
