package spotify

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	tracksPathRe = regexp.MustCompile(`/v1/tracks/([a-zA-Z0-9]+)`)
	albumsPathRe = regexp.MustCompile(`/v1/albums/([a-zA-Z0-9]+)`)
)

func NewAPIServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, r.Header.Get("Authorization"), "Bearer mock_access_token")

		switch {
		case tracksPathRe.MatchString(r.URL.Path):
			tracksHandler(t)(w, r)
		case r.URL.Path == "/v1/search":
			searchHandler(t)(w, r)
		case albumsPathRe.MatchString(r.URL.Path):
			albumsHandler(t)(w, r)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
}

func tracksHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		splitted := strings.Split(r.URL.Path, "/")
		trackID := splitted[len(splitted)-1]

		switch trackID {
		case TrackFixtureMassiveAttackAngel.Track.ID:
			_, err := w.Write([]byte(TrackFixtureMassiveAttackAngel.GetResponse))
			require.NoError(t, err)
		case TrackFixtureMileyCyrusFlowers.Track.ID:
			_, err := w.Write([]byte(TrackFixtureMileyCyrusFlowers.GetResponse))
			require.NoError(t, err)
		case TrackFixtureZemfiraIskala.Track.ID:
			_, err := w.Write([]byte(TrackFixtureZemfiraIskala.GetResponse))
			require.NoError(t, err)
		case TrackFixtureNadezhdaKadyshevaShirokaReka.Track.ID:
			_, err := w.Write([]byte(TrackFixtureNadezhdaKadyshevaShirokaReka.GetResponse))
			require.NoError(t, err)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func albumsHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		splitted := strings.Split(r.URL.Path, "/")
		albumID := splitted[len(splitted)-1]

		switch albumID {
		case AlbumFixtureRadioheadAmnesiac.Album.ID:
			_, err := w.Write([]byte(AlbumFixtureRadioheadAmnesiac.GetResponse))
			require.NoError(t, err)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func searchHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		searchType := r.URL.Query().Get("type")

		switch searchType {
		case "album":
			switch q {
			case AlbumFixtureRadioheadAmnesiac.SearchQuery:
				_, err := w.Write([]byte(AlbumFixtureRadioheadAmnesiac.SearchResponse))
				require.NoError(t, err)
			default:
				t.Errorf("unexpected search query: %s", q)
			}
		case "track":
			switch q {
			case TrackFixtureMassiveAttackAngel.SearchQuery:
				_, err := w.Write([]byte(TrackFixtureMassiveAttackAngel.SearchResponse))
				require.NoError(t, err)
			default:
				t.Errorf("unexpected search query: %s", q)
			}
		}
	}
}
