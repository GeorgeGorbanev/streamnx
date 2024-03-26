package utils

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

type FixturesMap struct {
	SpotifyTracks       map[string][]byte
	SpotifyAlbums       map[string][]byte
	SpotifySearchTracks map[string][]byte
	SpotifySearchAlbums map[string][]byte
	YandexTracks        map[string][]byte
	YandexAlbums        map[string][]byte
	YandexSearchTracks  map[string][]byte
	YandexSearchAlbums  map[string][]byte
}

var (
	SampleCredentials = spotify.Credentials{
		ClientID:     "sampleClientID",
		ClientSecret: "sampleClientSecret",
	}
	SampleToken = spotify.Token{
		AccessToken: "mock_access_token",
		TokenType:   "Bearer",
		ExpiresIn:   360,
	}
	SampleBasicAuth = "Basic c2FtcGxlQ2xpZW50SUQ6c2FtcGxlQ2xpZW50U2VjcmV0"
)

func NewSpotifyAuthServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/token", r.URL.Path)
		require.Equal(t, r.Header.Get("Authorization"), SampleBasicAuth)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		err := json.NewEncoder(w).Encode(map[string]any{
			"access_token": SampleToken.AccessToken,
			"token_type":   SampleToken.TokenType,
			"expires_in":   SampleToken.ExpiresIn,
		})
		require.NoError(t, err)
	}))
}

func NewYandexAPIServerMock(t *testing.T, fm FixturesMap) *httptest.Server {
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

func NewSpotifyAPIServerMock(t *testing.T, fm FixturesMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Header.Get("Authorization"), "Bearer mock_access_token")

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
			var searchMap map[string][]byte
			query := r.URL.Query().Get("q")
			searchType := r.URL.Query().Get("type")

			switch searchType {
			case "album":
				searchMap = fm.SpotifySearchAlbums
			case "track":
				searchMap = fm.SpotifySearchTracks
			default:
				t.Errorf("unexpected search type: %s", searchType)
			}

			if response, ok = searchMap[query]; !ok {
				t.Errorf("unexpected search query: %s", query)
			}
		case regexp.MustCompile(`/v1/albums/([a-zA-Z0-9]+)`).MatchString(r.URL.Path):
			splitted := strings.Split(r.URL.Path, "/")
			albumID := splitted[len(splitted)-1]

			if response, ok = fm.SpotifyAlbums[albumID]; !ok {
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
