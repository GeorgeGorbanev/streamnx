package soundcloud

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient_FetchTrack(t *testing.T) {
	tests := []struct {
		name        string
		userSlug    string
		trackSlug   string
		respStatus  int
		respBody    string
		want        *Track
		wantReqPath string
		wantErr     error
	}{
		{
			name:       "when found",
			userSlug:   "forss",
			trackSlug:  "flickermood",
			respStatus: http.StatusOK,
			respBody: `<html><body><script>window.__sc_hydration = [
				{"hydratable":"user","data":{"username":"Forss","permalink":"forss","permalink_url":"https://soundcloud.com/forss"}},
				{"hydratable":"sound","data":{
					"kind":"track",
					"urn":"soundcloud:tracks:293",
					"title":"Flickermood",
					"permalink":"flickermood",
					"permalink_url":"https://soundcloud.com/forss/flickermood",
					"user":{
						"username":"Forss",
						"permalink":"forss",
						"permalink_url":"https://soundcloud.com/forss"
					}
				}}
			];</script></body></html>`,
			wantReqPath: "/forss/flickermood",
			want: &Track{
				Kind:         "track",
				URN:          "soundcloud:tracks:293",
				Title:        "Flickermood",
				Permalink:    "flickermood",
				PermalinkURL: "https://soundcloud.com/forss/flickermood",
				User: User{
					Username:     "Forss",
					Permalink:    "forss",
					PermalinkURL: "https://soundcloud.com/forss",
				},
			},
		},
		{
			name:        "when not found",
			userSlug:    "missing",
			trackSlug:   "track",
			respStatus:  http.StatusNotFound,
			wantReqPath: "/missing/track",
			wantErr:     NotFoundError,
		},
		{
			name:        "when hydration script missing",
			userSlug:    "forss",
			trackSlug:   "flickermood",
			respStatus:  http.StatusOK,
			respBody:    `<html><body>No hydration here</body></html>`,
			wantReqPath: "/forss/flickermood",
			wantErr:     ErrHydrationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respBody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			serverURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewHTTPClient(
				WithAPIClient(srv.Client()),
				WithAPIHost(serverURL.Host),
				WithAPIScheme(serverURL.Scheme),
			)

			result, err := client.FetchTrack(t.Context(), tt.userSlug, tt.trackSlug)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestHTTPClient_FetchAlbum(t *testing.T) {
	tests := []struct {
		name        string
		userSlug    string
		setSlug     string
		respStatus  int
		respBody    string
		want        *Album
		wantReqPath string
		wantErr     error
	}{
		{
			name:       "when found",
			userSlug:   "forss",
			setSlug:    "soulhack",
			respStatus: http.StatusOK,
			respBody: `<html><body><script>window.__sc_hydration = [
				{"hydratable":"playlist","data":{
					"kind":"playlist",
					"urn":"soundcloud:playlists:123",
					"title":"Soulhack",
					"permalink":"soulhack",
					"permalink_url":"https://soundcloud.com/forss/sets/soulhack",
					"set_type":"album",
					"track_count":2,
					"user":{
						"username":"Forss",
						"permalink":"forss",
						"permalink_url":"https://soundcloud.com/forss"
					},
					"tracks":[
						{"kind":"track","title":"Flickermood","permalink":"flickermood","permalink_url":"https://soundcloud.com/forss/flickermood"},
						{"kind":"track","title":"Using Splashes","permalink":"using-splashes","permalink_url":"https://soundcloud.com/forss/using-splashes"}
					]
				}}
			];</script></body></html>`,
			wantReqPath: "/forss/sets/soulhack",
			want: &Album{
				Kind:         "playlist",
				URN:          "soundcloud:playlists:123",
				Title:        "Soulhack",
				Permalink:    "soulhack",
				PermalinkURL: "https://soundcloud.com/forss/sets/soulhack",
				SetType:      "album",
				TrackCount:   2,
				User: User{
					Username:     "Forss",
					Permalink:    "forss",
					PermalinkURL: "https://soundcloud.com/forss",
				},
				Tracks: []Track{
					{
						Kind:         "track",
						Title:        "Flickermood",
						Permalink:    "flickermood",
						PermalinkURL: "https://soundcloud.com/forss/flickermood",
					},
					{
						Kind:         "track",
						Title:        "Using Splashes",
						Permalink:    "using-splashes",
						PermalinkURL: "https://soundcloud.com/forss/using-splashes",
					},
				},
			},
		},
		{
			name:        "when not found",
			userSlug:    "missing",
			setSlug:     "album",
			respStatus:  http.StatusNotFound,
			wantReqPath: "/missing/sets/album",
			wantErr:     NotFoundError,
		},
		{
			name:        "when hydration script missing",
			userSlug:    "forss",
			setSlug:     "soulhack",
			respStatus:  http.StatusOK,
			respBody:    `<html><body>No hydration here</body></html>`,
			wantReqPath: "/forss/sets/soulhack",
			wantErr:     ErrHydrationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respBody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			serverURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewHTTPClient(
				WithAPIClient(srv.Client()),
				WithAPIHost(serverURL.Host),
				WithAPIScheme(serverURL.Scheme),
			)

			result, err := client.FetchAlbum(t.Context(), tt.userSlug, tt.setSlug)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestParseTrackHTML(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    *Track
		wantErr error
	}{
		{
			name: "when sound data found",
			html: `<script>window.__sc_hydration = [
				{"hydratable":"sound","data":{
					"kind":"track",
					"urn":"soundcloud:tracks:293",
					"title":"Flickermood",
					"permalink":"flickermood",
					"permalink_url":"https://soundcloud.com/forss/flickermood",
					"user":{"username":"Forss","permalink":"forss","permalink_url":"https://soundcloud.com/forss"}
				}}
			];</script>`,
			want: &Track{
				Kind:         "track",
				URN:          "soundcloud:tracks:293",
				Title:        "Flickermood",
				Permalink:    "flickermood",
				PermalinkURL: "https://soundcloud.com/forss/flickermood",
				User: User{
					Username:     "Forss",
					Permalink:    "forss",
					PermalinkURL: "https://soundcloud.com/forss",
				},
			},
		},
		{
			name:    "when hydration script missing",
			html:    `<html></html>`,
			wantErr: ErrHydrationNotFound,
		},
		{
			name: "when sound data missing",
			html: `<script>window.__sc_hydration = [
				{"hydratable":"user","data":{"username":"Forss"}}
			];</script>`,
			wantErr: ErrSoundDataNotFound,
		},
		{
			name: "when hydration contains playlist",
			html: `<script>window.__sc_hydration = [
				{"hydratable":"sound","data":{"kind":"playlist","title":"Soulhack"}}
			];</script>`,
			wantErr: NotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTrackHTML([]byte(tt.html))
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestParseAlbumHTML(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    *Album
		wantErr error
	}{
		{
			name: "when playlist data found",
			html: `<script>window.__sc_hydration = [
				{"hydratable":"playlist","data":{
					"kind":"playlist",
					"urn":"soundcloud:playlists:123",
					"title":"Soulhack",
					"permalink":"soulhack",
					"permalink_url":"https://soundcloud.com/forss/sets/soulhack",
					"set_type":"album",
					"track_count":1,
					"user":{"username":"Forss","permalink":"forss","permalink_url":"https://soundcloud.com/forss"},
					"tracks":[
						{"kind":"track","title":"Flickermood","permalink":"flickermood","permalink_url":"https://soundcloud.com/forss/flickermood"}
					]
				}}
			];</script>`,
			want: &Album{
				Kind:         "playlist",
				URN:          "soundcloud:playlists:123",
				Title:        "Soulhack",
				Permalink:    "soulhack",
				PermalinkURL: "https://soundcloud.com/forss/sets/soulhack",
				SetType:      "album",
				TrackCount:   1,
				User: User{
					Username:     "Forss",
					Permalink:    "forss",
					PermalinkURL: "https://soundcloud.com/forss",
				},
				Tracks: []Track{
					{
						Kind:         "track",
						Title:        "Flickermood",
						Permalink:    "flickermood",
						PermalinkURL: "https://soundcloud.com/forss/flickermood",
					},
				},
			},
		},
		{
			name:    "when hydration script missing",
			html:    `<html></html>`,
			wantErr: ErrHydrationNotFound,
		},
		{
			name: "when playlist data missing",
			html: `<script>window.__sc_hydration = [
				{"hydratable":"user","data":{"username":"Forss"}}
			];</script>`,
			wantErr: ErrPlaylistDataNotFound,
		},
		{
			name: "when hydration contains track",
			html: `<script>window.__sc_hydration = [
				{"hydratable":"playlist","data":{"kind":"track","title":"Flickermood"}}
			];</script>`,
			wantErr: NotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAlbumHTML([]byte(tt.html))
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}
