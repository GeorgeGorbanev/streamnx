package utils

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
)

func NewYandexAPIServerMock(fm *fixture.FixturesMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response []byte
		var ok bool

		switch {
		case regexp.MustCompile(`/tracks/\d+`).MatchString(r.URL.Path):
			trackID := strings.Split(r.URL.Path, "/")[2]
			if response, ok = fm.YandexTracks[trackID]; !ok {
				response = fixture.Read("yandex/get_track_not_found.json")
			}
		case regexp.MustCompile(`/albums/\d+`).MatchString(r.URL.Path):
			albumID := strings.Split(r.URL.Path, "/")[2]
			if response, ok = fm.YandexAlbums[albumID]; !ok {
				response = fixture.Read("yandex/get_album_not_found.json")
			}
		case r.URL.Path == "/search":
			var searchMap map[string][]byte
			searchType := r.URL.Query().Get("type")
			query := strings.ToLower(r.URL.Query().Get("text"))

			switch searchType {
			case "album":
				searchMap = fm.YandexSearchAlbums
			case "track":
				searchMap = fm.YandexSearchTracks
			default:
				panic("unexpected search type")
			}

			if response, ok = searchMap[query]; !ok {
				response = fixture.Read("yandex/search_not_found.json")
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
