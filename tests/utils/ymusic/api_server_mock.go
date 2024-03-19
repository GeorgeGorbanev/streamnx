package ymusic

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	albumPathRe  = regexp.MustCompile(`/albums/\d+`)
	tracksPathRe = regexp.MustCompile(`/tracks/\d+`)
)

func NewAPIServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		switch {
		case tracksPathRe.MatchString(r.URL.Path):
			tracksHandler(t)(w, r)
		case albumPathRe.MatchString(r.URL.Path):
			albumsHandler(t)(w, r)
		case r.URL.Path == "/search":
			searchHandler(t)(w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
}

func tracksHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackID := strings.Split(r.URL.Path, "/")[2]

		if trackID == TrackFixtureMassiveAttackAngel.ID {
			_, err := w.Write([]byte(TrackFixtureMassiveAttackAngel.GetResponse))
			require.NoError(t, err)
		} else {
			_, err := w.Write([]byte(TrackFixtureNotFound.GetResponse))
			require.NoError(t, err)
		}
	}
}

func albumsHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		albumID := strings.Split(r.URL.Path, "/")[2]

		if albumID == AlbumFixtureRadioheadAmnesiac.ID {
			_, err := w.Write([]byte(AlbumFixtureRadioheadAmnesiac.GetResponse))
			require.NoError(t, err)
		} else {
			_, err := w.Write([]byte(AlbumFixtureNotFound.GetResponse))
			require.NoError(t, err)
		}
	}
}

func searchHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchType := r.URL.Query().Get("type")
		query := strings.ToLower(r.URL.Query().Get("text"))

		var err error

		switch searchType {
		case "album":
			switch query {
			case AlbumFixtureRadioheadAmnesiac.SearchQuery():
				_, err = w.Write([]byte(AlbumFixtureRadioheadAmnesiac.SearchResponse))
			default:
				_, err = w.Write([]byte(TrackFixtureNotFound.SearchResponse))
			}
		case "track":
			switch query {
			case TrackFixtureMassiveAttackAngel.SearchQuery():
				_, err = w.Write([]byte(TrackFixtureMassiveAttackAngel.SearchResponse))
			case TrackFixtureDJAmor20Flowers.SearchQuery():
				_, err = w.Write([]byte(TrackFixtureDJAmor20Flowers.SearchResponse))
			case TrackFixtureZemfiraIskala.SearchQuery():
				_, err = w.Write([]byte(TrackFixtureZemfiraIskala.SearchResponse))
			case TrackFixtureNadezhdaKadyshevaShirokaReka.SearchQuery():
				_, err = w.Write([]byte(TrackFixtureNadezhdaKadyshevaShirokaReka.SearchResponse))
			default:
				_, err = w.Write([]byte(TrackFixtureNotFound.SearchResponse))
			}
		default:
			t.Errorf("unexpected search type: %s", searchType)
		}

		require.NoError(t, err)
	}
}
