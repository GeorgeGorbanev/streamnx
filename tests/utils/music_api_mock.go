package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
)

type FixturesMap struct {
	SpotifyTracks       map[string][]byte
	SpotifyAlbums       map[string][]byte
	SpotifySearchTracks map[string][]byte
	SpotifySearchAlbums map[string][]byte

	YandexTracks       map[string][]byte
	YandexAlbums       map[string][]byte
	YandexSearchTracks map[string][]byte
	YandexSearchAlbums map[string][]byte

	YoutubeTracks       map[string][]byte
	YoutubeAlbums       map[string][]byte
	YoutubeSearchTracks map[string][]byte
	YoutubeSearchAlbums map[string][]byte
}

var (
	SpotifyCredentials = spotify.Credentials{
		ClientID:     "sampleClientID",
		ClientSecret: "sampleClientSecret",
	}
	SpotifyToken = spotify.Token{
		AccessToken: "mock_access_token",
		TokenType:   "Bearer",
		ExpiresIn:   360,
	}
	SpotifyBasicAuth = "Basic c2FtcGxlQ2xpZW50SUQ6c2FtcGxlQ2xpZW50U2VjcmV0"
	YoutubeAPIKey    = "sampleAPIKey"
)

func NewSpotifyAuthServerMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/token" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get("Authorization") != SpotifyBasicAuth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err := json.NewEncoder(w).Encode(map[string]any{
			"access_token": SpotifyToken.AccessToken,
			"token_type":   SpotifyToken.TokenType,
			"expires_in":   SpotifyToken.ExpiresIn,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func NewSpotifyAPIServerMock(fm *FixturesMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mock_access_token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var response []byte
		var ok bool

		switch {
		case regexp.MustCompile(`/v1/tracks/([a-zA-Z0-9]+)`).MatchString(r.URL.Path):
			splitted := strings.Split(r.URL.Path, "/")
			trackID := splitted[len(splitted)-1]

			if response, ok = fm.SpotifyTracks[trackID]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case r.URL.Path == "/v1/search":
			query := r.URL.Query().Get("q")
			searchType := r.URL.Query().Get("type")

			switch searchType {
			case "album":
				if response, ok = fm.SpotifySearchAlbums[query]; !ok {
					response = fixture.Read("spotify/search_album_not_found.json")
				}
			case "track":
				if response, ok = fm.SpotifySearchTracks[query]; !ok {
					response = fixture.Read("spotify/search_track_not_found.json")
				}
			default:
				panic("unexpected search type")
			}
		case regexp.MustCompile(`/v1/albums/([a-zA-Z0-9]+)`).MatchString(r.URL.Path):
			splitted := strings.Split(r.URL.Path, "/")
			albumID := splitted[len(splitted)-1]

			if response, ok = fm.SpotifyAlbums[albumID]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
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

func NewYandexAPIServerMock(fm *FixturesMap) *httptest.Server {
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

func NewYoutubeAPIServerMock(fm *FixturesMap) *httptest.Server {
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
		default:
			panic("unexpected request")
		}

		_, err := w.Write(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}

func (fm *FixturesMap) Merge(mergeFm *FixturesMap) {
	fm.SpotifyAlbums = mergeFm.SpotifyAlbums
	fm.SpotifyTracks = mergeFm.SpotifyTracks
	fm.SpotifySearchAlbums = mergeFm.SpotifySearchAlbums
	fm.SpotifySearchTracks = mergeFm.SpotifySearchTracks

	fm.YandexAlbums = mergeFm.YandexAlbums
	fm.YandexTracks = mergeFm.YandexTracks
	fm.YandexSearchAlbums = mergeFm.YandexSearchAlbums
	fm.YandexSearchTracks = mergeFm.YandexSearchTracks

	fm.YoutubeAlbums = mergeFm.YoutubeAlbums
	fm.YoutubeTracks = mergeFm.YoutubeTracks
	fm.YoutubeSearchAlbums = mergeFm.YoutubeSearchAlbums
	fm.YoutubeSearchTracks = mergeFm.YoutubeSearchTracks
}

func (fm *FixturesMap) Reset() {
	fm.SpotifyAlbums = map[string][]byte{}
	fm.SpotifyTracks = map[string][]byte{}
	fm.SpotifySearchAlbums = map[string][]byte{}
	fm.SpotifySearchTracks = map[string][]byte{}

	fm.YandexAlbums = map[string][]byte{}
	fm.YandexTracks = map[string][]byte{}
	fm.YandexSearchAlbums = map[string][]byte{}
	fm.YandexSearchTracks = map[string][]byte{}

	fm.YoutubeAlbums = map[string][]byte{}
	fm.YoutubeTracks = map[string][]byte{}
	fm.YoutubeSearchAlbums = map[string][]byte{}
	fm.YoutubeSearchTracks = map[string][]byte{}
}
