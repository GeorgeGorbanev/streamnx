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

func NewSpotifyAPIServerMock(fm *fixture.FixturesMap) *httptest.Server {
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
