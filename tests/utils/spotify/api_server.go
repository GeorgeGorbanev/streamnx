package spotify

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/stretchr/testify/require"
)

var (
	SampleTrack = spotify.Track{
		Album: spotify.Album{
			AlbumType: "album",
			Artists: []spotify.Artist{
				{
					ExternalURLs: map[string]string{"spotify": "https://spotify.com/artist1"},
					Href:         "https://api.spotify.com/v1/artists/artist1",
					ID:           "artist1",
					Name:         "Artist One",
					Type:         "artist",
					URI:          "spotify:artist:artist1",
				},
			},
			AvailableMarkets: []string{"US", "GB"},
			ExternalURLs:     map[string]string{"spotify": "https://spotify.com/album1"},
			Href:             "https://api.spotify.com/v1/albums/album1",
			ID:               "album1",
			Images: []spotify.Image{
				{Height: 640, URL: "https://image.url/1", Width: 640},
				{Height: 300, URL: "https://image.url/2", Width: 300},
			},
			Name:                 "Album One",
			ReleaseDate:          "2023-01-01",
			ReleaseDatePrecision: "day",
			Type:                 "album",
			URI:                  "spotify:album:album1",
		},
		Artists: []spotify.Artist{
			{
				ExternalURLs: map[string]string{"spotify": "https://spotify.com/artist1"},
				Href:         "https://api.spotify.com/v1/artists/artist1",
				ID:           "artist1",
				Name:         "Artist One",
				Type:         "artist",
				URI:          "spotify:artist:artist1",
			},
		},
		AvailableMarkets: []string{"US", "GB"},
		DiscNumber:       1,
		DurationMs:       240000,
		Explicit:         false,
		ExternalIDs:      map[string]string{"isrc": "US1234567890"},
		ExternalURLs:     map[string]string{"spotify": "https://spotify.com/track1"},
		Href:             "https://api.spotify.com/v1/tracks/track1",
		ID:               "track1",
		IsPlayable:       true,
		LinkedFrom:       nil,
		Name:             "Track One",
		Popularity:       80,
		PreviewURL:       "https://preview.url/track1",
		TrackNumber:      5,
		Type:             "track",
		URI:              "spotify:track:track1",
	}
	SampleTrackJSON = `{
        "album": {
            "album_type": "album",
            "artists": [
                {
                    "external_urls": {"spotify": "https://spotify.com/artist1"},
                    "href": "https://api.spotify.com/v1/artists/artist1",
                    "id": "artist1",
                    "name": "Artist One",
                    "type": "artist",
                    "uri": "spotify:artist:artist1"
                }
            ],
            "available_markets": ["US", "GB"],
            "external_urls": {"spotify": "https://spotify.com/album1"},
            "href": "https://api.spotify.com/v1/albums/album1",
            "id": "album1",
            "images": [
                {"height": 640, "url": "https://image.url/1", "width": 640},
                {"height": 300, "url": "https://image.url/2", "width": 300}
            ],
            "name": "Album One",
            "release_date": "2023-01-01",
            "release_date_precision": "day",
            "type": "album",
            "uri": "spotify:album:album1"
        },
        "artists": [
            {
                "external_urls": {"spotify": "https://spotify.com/artist1"},
                "href": "https://api.spotify.com/v1/artists/artist1",
                "id": "artist1",
                "name": "Artist One",
                "type": "artist",
                "uri": "spotify:artist:artist1"
            }
        ],
        "available_markets": ["US", "GB"],
        "disc_number": 1,
        "duration_ms": 240000,
        "explicit": false,
        "external_ids": {"isrc": "US1234567890"},
        "external_urls": {"spotify": "https://spotify.com/track1"},
        "href": "https://api.spotify.com/v1/tracks/track1",
        "id": "track1",
        "is_playable": true,
        "linked_from": null,
        "name": "Track One",
        "popularity": 80,
        "preview_url": "https://preview.url/track1",
        "track_number": 5,
        "type": "track",
        "uri": "spotify:track:track1"
    }`
)

func NewAPIServerMock(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Contains(t, r.URL.Path, "/v1/tracks/")
		require.Equal(t, r.Header.Get("Authorization"), "Bearer mock_access_token")

		if r.URL.Path == fmt.Sprintf("/v1/tracks/%s", SampleTrack.ID) {
			_, err := w.Write([]byte(SampleTrackJSON))
			require.NoError(t, err)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}))
}
