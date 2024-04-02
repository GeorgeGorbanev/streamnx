package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
)

var (
	AppleToken = "sampleAppleToken"
)

func NewAppleWebPlayerServerMock() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, err := w.Write([]byte(`
				<!DOCTYPE html>
				<html>
					<head><script src="/assets/index-Samp1eBund13.js"></script></head>
				</html>
			`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		case "/assets/index-Samp1eBund13.js":
			_, err := w.Write([]byte(`
				tokenVar = "` + AppleToken + `" 
				headers.Authorization = ` + "`Bearer ${tokenVar}`",
			))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func NewAppleAPIServerMock(fm *fixture.FixturesMap) *httptest.Server {
	authHeader := fmt.Sprintf("Bearer %s", AppleToken)
	trackRe := regexp.MustCompile(`^/v1/catalog/us/songs/(\d+)$`)
	albumRe := regexp.MustCompile(`^/v1/catalog/us/albums/(\d+)$`)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != authHeader || r.Header.Get("Origin") != "https://music.apple.com" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var (
			response []byte
			ok       bool
			status   = http.StatusOK
		)

		switch {
		case trackRe.MatchString(r.URL.Path):
			trackID := trackRe.FindSubmatch([]byte(r.URL.Path))[1]

			if response, ok = fm.AppleTracks[string(trackID)]; !ok {
				status = http.StatusNotFound
			}
		case r.URL.Path == "/v1/catalog/us/search/suggestions":
			term := r.URL.Query().Get("term")
			found := false

			response, found = fm.AppleSearchAlbums[term]
			if !found {
				response, found = fm.AppleSearchTracks[term]
			}
			if !found {
				response = fixture.Read("apple/not_found.json")
			}
		case albumRe.MatchString(r.URL.Path):
			albumID := albumRe.FindSubmatch([]byte(r.URL.Path))[1]

			if response, ok = fm.AppleAlbums[string(albumID)]; !ok {
				status = http.StatusNotFound
			}
		default:
			panic("unexpected request")
		}

		w.WriteHeader(status)
		_, err := w.Write(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
}
