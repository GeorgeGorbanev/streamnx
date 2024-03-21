package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/stretchr/testify/require"
)

type fixturesMap struct {
	spotifyTracks       map[string][]byte
	spotifyAlbums       map[string][]byte
	spotifySearchTracks map[string][]byte
	spotifySearchAlbums map[string][]byte

	yandexTracks       map[string][]byte
	yandexAlbums       map[string][]byte
	yandexSearchTracks map[string][]byte
	yandexSearchAlbums map[string][]byte
}

var (
	sampleCredentials = spotify.Credentials{
		ClientID:     "sampleClientID",
		ClientSecret: "sampleClientSecret",
	}
	sampleToken = spotify.Token{
		AccessToken: "mock_access_token",
		TokenType:   "Bearer",
		ExpiresIn:   360,
	}
	sampleBasicAuth = "Basic c2FtcGxlQ2xpZW50SUQ6c2FtcGxlQ2xpZW50U2VjcmV0"
)

func newSpotifyAuthServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/token", r.URL.Path)
		require.Equal(t, r.Header.Get("Authorization"), sampleBasicAuth)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		err := json.NewEncoder(w).Encode(map[string]any{
			"access_token": sampleToken.AccessToken,
			"token_type":   sampleToken.TokenType,
			"expires_in":   sampleToken.ExpiresIn,
		})
		require.NoError(t, err)
	}))
}

func newYandexAPIServerMock(t *testing.T, fm fixturesMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response []byte
		var ok bool

		switch {
		case regexp.MustCompile(`/tracks/\d+`).MatchString(r.URL.Path):
			trackID := strings.Split(r.URL.Path, "/")[2]
			if response, ok = fm.yandexTracks[trackID]; !ok {
				response = fixture.Read("yandex/get_track_not_found.json")
			}
		case regexp.MustCompile(`/albums/\d+`).MatchString(r.URL.Path):
			albumID := strings.Split(r.URL.Path, "/")[2]
			if response, ok = fm.yandexAlbums[albumID]; !ok {
				response = fixture.Read("yandex/get_album_not_found.json")
			}
		case r.URL.Path == "/search":
			var searchMap map[string][]byte
			searchType := r.URL.Query().Get("type")
			query := strings.ToLower(r.URL.Query().Get("text"))

			switch searchType {
			case "album":
				searchMap = fm.yandexSearchAlbums
			case "track":
				searchMap = fm.yandexSearchTracks
			default:
				t.Errorf("unexpected search type: %s", searchType)
			}

			if response, ok = searchMap[query]; !ok {
				response = fixture.Read("yandex/search_not_found.json")
			}
		default:
			t.Errorf("unexpected request: %s", r.URL.Path)
		}

		_, err := w.Write(response)
		require.NoError(t, err)
	}))
}

func newSpotifyAPIServerMock(t *testing.T, fm fixturesMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Header.Get("Authorization"), "Bearer mock_access_token")

		var response []byte
		var ok bool

		switch {
		case regexp.MustCompile(`/v1/tracks/([a-zA-Z0-9]+)`).MatchString(r.URL.Path):
			splitted := strings.Split(r.URL.Path, "/")
			trackID := splitted[len(splitted)-1]

			if response, ok = fm.spotifyTracks[trackID]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case r.URL.Path == "/v1/search":
			var searchMap map[string][]byte
			query := r.URL.Query().Get("q")
			searchType := r.URL.Query().Get("type")

			switch searchType {
			case "album":
				searchMap = fm.spotifySearchAlbums
			case "track":
				searchMap = fm.spotifySearchTracks
			default:
				t.Errorf("unexpected search type: %s", searchType)
			}

			if response, ok = searchMap[query]; !ok {
				t.Errorf("unexpected search query: %s", query)
			}
		case regexp.MustCompile(`/v1/albums/([a-zA-Z0-9]+)`).MatchString(r.URL.Path):
			splitted := strings.Split(r.URL.Path, "/")
			albumID := splitted[len(splitted)-1]

			if response, ok = fm.spotifyAlbums[albumID]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			t.Errorf("unexpected request: %s", r.URL.Path)
		}

		_, err := w.Write(response)
		require.NoError(t, err)
	}))
}
