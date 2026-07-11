package fixtures

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed responses/*
var responses embed.FS

type Route struct {
	Method  string
	Path    string
	Query   map[string]string
	Status  int
	Fixture string
	Headers map[string]string
	Assert  func(*testing.T, *http.Request)
}

func NewServer(t *testing.T, routes ...Route) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, ok := findRoute(routes, r)
		require.Truef(t, ok, "unexpected request: %s %s", r.Method, r.URL.String())

		if route.Assert != nil {
			route.Assert(t, r)
		}

		for key, value := range route.Headers {
			w.Header().Set(key, value)
		}

		body := readFixture(t, route.Fixture)
		if strings.HasSuffix(route.Fixture, ".json") {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(route.Status)
		_, err := w.Write(body)
		require.NoError(t, err)
	}))
}

func findRoute(routes []Route, r *http.Request) (Route, bool) {
	for _, route := range routes {
		if !routeMatches(route, r) {
			continue
		}
		return route, true
	}
	return Route{}, false
}

func routeMatches(route Route, r *http.Request) bool {
	if route.Method != r.Method || route.Path != r.URL.Path {
		return false
	}
	for key, value := range route.Query {
		if r.URL.Query().Get(key) != value {
			return false
		}
	}
	return true
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()

	body, err := responses.ReadFile(path.Join("responses", name))
	require.NoError(t, err)
	return body
}
