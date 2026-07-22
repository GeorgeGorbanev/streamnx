package bandcamp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_fetchAlbum(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		artistSlug  string
		respStatus  int
		respbody    string
		want        Entity
		wantReqHost string
		wantReqPath string
		wantErr     error
	}{
		{
			name:       "when found",
			id:         "amber",
			artistSlug: "autechre",
			respStatus: http.StatusOK,
			respbody: `<html>
				<head>
					<script type="application/ld+json">
					{
						"name": "Amber",
						"description": "sample album description",
						"byArtist": {
							"name": "Autechre"
						},
						"publisher": {
							"name": "Warp Records"
						},
						"track": {
							"@type": "ItemList",
							"numberOfItems": 2,
							"itemListElement": [
								{
									"@type": "ListItem",
									"position": 1,
									"item": {
										"@type": "MusicRecording",
										"@id": "https://autechre.bandcamp.com/track/foil",
										"name": "Foil",
										"mainEntityOfPage": "https://autechre.bandcamp.com/track/foil"
									}
								},
								{
									"@type": "ListItem",
									"position": 2,
									"item": {
										"@type": "MusicRecording",
										"@id": "https://autechre.bandcamp.com/track/montreal",
										"name": "Montreal",
										"mainEntityOfPage": "https://autechre.bandcamp.com/track/montreal"
									}
								}
							]
						}
					}
					</script>
				</head>
				<body></body>
			</html>`,
			wantReqHost: "autechre.bandcamp.com",
			wantReqPath: "/album/amber",
			want: Entity{
				Name:        "Amber",
				BandName:    "Autechre",
				CreatorName: "Warp Records",
				Description: "sample album description",
				URL:         "http://autechre.bandcamp.com/album/amber",
				TrackURLs: []string{
					"https://autechre.bandcamp.com/track/foil",
					"https://autechre.bandcamp.com/track/montreal",
				},
			},
		},
		{
			name:        "when not found",
			id:          "notfound",
			artistSlug:  "notfound",
			wantReqHost: "notfound.bandcamp.com",
			wantReqPath: "/album/notfound",
			respStatus:  http.StatusNotFound,
			wantErr:     errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqHost, r.Host)
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respbody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			client := NewClient(
				WithAPIClient(mockSubdomainClient(srv.URL)),
				WithAPI("http", "bandcamp.com"),
			)

			result, err := client.fetchAlbum(t.Context(), tt.artistSlug, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_fetchTrack(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		artistSlug  string
		respStatus  int
		respbody    string
		want        Entity
		wantReqHost string
		wantReqPath string
		wantErr     error
	}{
		{
			name:       "when found",
			id:         "nil",
			artistSlug: "autechre",
			respStatus: http.StatusOK,
			respbody: `<html>
				<head>
					<script type="application/ld+json">
					{
						"name": "Nil",
						"description": "sample track description",
						"byArtist": {
							"name": "Autechre"
						},
						"publisher": {
							"name": "Warp Records"
						}
					}
					</script>
				</head>
				<body></body>
			</html>`,
			wantReqHost: "autechre.bandcamp.com",
			wantReqPath: "/track/nil",
			want: Entity{
				Name:        "Nil",
				BandName:    "Autechre",
				CreatorName: "Warp Records",
				Description: "sample track description",
				URL:         "http://autechre.bandcamp.com/track/nil",
			},
		},
		{
			name:        "when not found",
			id:          "notfound",
			artistSlug:  "notfound",
			wantReqHost: "notfound.bandcamp.com",
			wantReqPath: "/track/notfound",
			respStatus:  http.StatusNotFound,
			wantErr:     errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tt.wantReqHost, r.Host)
				require.Equal(t, tt.wantReqPath, r.URL.Path)
				w.WriteHeader(tt.respStatus)
				_, err := w.Write([]byte(tt.respbody))
				require.NoError(t, err)
			}))
			defer srv.Close()

			client := NewClient(
				WithAPIClient(mockSubdomainClient(srv.URL)),
				WithAPI("http", "bandcamp.com"),
			)

			result, err := client.fetchTrack(t.Context(), tt.artistSlug, tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_searchTracks(t *testing.T) {
	tests := []struct {
		name            string
		artist          string
		title           string
		respStatus      int
		respBody        string
		want            []Entity
		wantErr         error
		wantErrContains string
	}{
		{
			name:       "when multiple track candidates found",
			artist:     "Autechre",
			title:      "Nil",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": [
						{
							"id": 123,
							"name": "Nil",
							"band_name": "Autechre",
							"item_url_root": "https://autechre.bandcamp.com",
							"item_url_path": "https://autechre.bandcamp.com/track/nil"
						},
						{
							"name": "Nil (alternate)",
							"band_name": "Autechre",
							"item_url_root": "https://autechre.bandcamp.com",
							"item_url_path": "https://autechre.bandcamp.com/track/nil-alternate"
						}
					]
				}
			}`,
			want: []Entity{
				{
					NumericID: 123,
					Name:      "Nil",
					BandName:  "Autechre",
					URL:       "https://autechre.bandcamp.com/track/nil",
				},
				{
					Name:     "Nil (alternate)",
					BandName: "Autechre",
					URL:      "https://autechre.bandcamp.com/track/nil-alternate",
				},
			},
		},
		{
			name:       "when track not found",
			artist:     "Unknown",
			title:      "Track",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": []
				}
			}`,
			want: []Entity{},
		},
		{
			name:            "when api responds with error",
			artist:          "Error",
			title:           "State",
			respStatus:      http.StatusInternalServerError,
			wantErrContains: "unexpected status code: 500",
		},
		{
			name:            "when api response is malformed",
			artist:          "Malformed",
			title:           "Response",
			respStatus:      http.StatusOK,
			respBody:        `{`,
			wantErrContains: "failed to decode response body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/api/bcsearch_public_api/1/autocomplete_elastic", r.URL.Path)
				require.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Content-Type"))
				var body struct {
					SearchText   string `json:"search_text"`
					SearchFilter string `json:"search_filter"`
					FullPage     bool   `json:"full_page"`
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.Equal(t, tt.artist+" "+tt.title, body.SearchText)
				require.Equal(t, string(trackEntityType), body.SearchFilter)
				require.True(t, body.FullPage)
				w.WriteHeader(tt.respStatus)
				if tt.respBody != "" {
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				}
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewClient(
				WithAPI(srvURL.Scheme, srvURL.Host),
			)

			result, err := client.searchTracks(t.Context(), tt.artist, tt.title)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else if tt.wantErrContains != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestClient_searchAlbums(t *testing.T) {
	tests := []struct {
		name            string
		artist          string
		title           string
		respStatus      int
		respBody        string
		want            []Entity
		wantErr         error
		wantErrContains string
	}{
		{
			name:       "when multiple album candidates found",
			artist:     "Boards of Canada",
			title:      "Music Has The Right To Children",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": [
						{
							"id": 456,
							"name": "Music Has The Right To Children",
							"band_name": "Boards of Canada",
							"item_url_root": "https://boardsofcanada.bandcamp.com",
							"item_url_path": "https://boardsofcanada.bandcamp.com/album/music-has-the-right-to-children"
						},
						{
							"name": "Hi Scores",
							"band_name": "Boards of Canada",
							"item_url_root": "https://boardsofcanada.bandcamp.com",
							"item_url_path": "https://boardsofcanada.bandcamp.com/album/hi-scores"
						}
					]
				}
			}`,
			want: []Entity{
				{
					NumericID: 456,
					Name:      "Music Has The Right To Children",
					BandName:  "Boards of Canada",
					URL:       "https://boardsofcanada.bandcamp.com/album/music-has-the-right-to-children",
				},
				{
					Name:     "Hi Scores",
					BandName: "Boards of Canada",
					URL:      "https://boardsofcanada.bandcamp.com/album/hi-scores",
				},
			},
		},
		{
			name:       "when album not found",
			artist:     "Unknown",
			title:      "Album",
			respStatus: http.StatusOK,
			respBody: `{
				"auto": {
					"results": []
				}
			}`,
			want: []Entity{},
		},
		{
			name:            "when api responds with error",
			artist:          "Boards of Canada",
			title:           "Hi Scores",
			respStatus:      http.StatusBadGateway,
			wantErrContains: "unexpected status code: 502",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/api/bcsearch_public_api/1/autocomplete_elastic", r.URL.Path)
				require.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Content-Type"))
				var body struct {
					SearchText   string `json:"search_text"`
					SearchFilter string `json:"search_filter"`
					FullPage     bool   `json:"full_page"`
				}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.Equal(t, tt.artist+" "+tt.title, body.SearchText)
				require.Equal(t, string(albumEntityType), body.SearchFilter)
				require.True(t, body.FullPage)
				w.WriteHeader(tt.respStatus)
				if tt.respBody != "" {
					_, err := w.Write([]byte(tt.respBody))
					require.NoError(t, err)
				}
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			require.NoError(t, err)

			client := NewClient(
				WithAPI(srvURL.Scheme, srvURL.Host),
			)

			result, err := client.searchAlbums(t.Context(), tt.artist, tt.title)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else if tt.wantErrContains != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func mockSubdomainClient(u string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, _, err := net.SplitHostPort(addr)
				if err == nil && strings.HasSuffix(host, ".bandcamp.com") {
					d := net.Dialer{Timeout: 5 * time.Second}
					return d.DialContext(ctx, network, strings.TrimPrefix(u, "http://"))
				}
				d := net.Dialer{Timeout: 5 * time.Second}
				return d.DialContext(ctx, network, addr)
			},
		},
	}
}

func TestParseEmbeddedPlayerData(t *testing.T) {
	tests := []struct {
		name            string
		html            string
		wantLinkback    string
		wantTrackID     uint64
		wantTrackLink   string
		wantErrContains string
	}{
		{
			name: "double quoted escaped attribute with whitespace",
			html: `<html><body>
				<div class="inline_player"
					data-player-data="{&quot;linkback&quot;:&quot;https://fantastictrax.bandcamp.com/track/torn&quot;,&quot;tracks&quot;:[{&quot;id&quot;:852027615,&quot;title_link&quot;:&quot;https://fantastictrax.bandcamp.com/track/torn&quot;}]}">
				</div>
			</body></html>`,
			wantLinkback:  "https://fantastictrax.bandcamp.com/track/torn",
			wantTrackID:   852027615,
			wantTrackLink: "https://fantastictrax.bandcamp.com/track/torn",
		},
		{
			name:         "single quoted attribute",
			html:         `<div data-player-data='{&quot;linkback&quot;:&quot;https://label.bandcamp.com/album/release&quot;}'></div>`,
			wantLinkback: "https://label.bandcamp.com/album/release",
		},
		{
			name:            "missing attribute",
			html:            `<html><body></body></html>`,
			wantErrContains: "attribute not found",
		},
		{
			name:            "empty attribute",
			html:            `<div data-player-data=""></div>`,
			wantErrContains: "attribute is empty",
		},
		{
			name:            "invalid json",
			html:            `<div data-player-data="{not-json}"></div>`,
			wantErrContains: "failed to decode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient().parseEmbeddedPlayerData([]byte(tt.html))
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantLinkback, got.Linkback)
			if tt.wantTrackID != 0 {
				require.Len(t, got.Tracks, 1)
				require.Equal(t, tt.wantTrackID, got.Tracks[0].ID)
				require.Equal(t, tt.wantTrackLink, got.Tracks[0].TitleLink)
			}
		})
	}
}

func TestEmbeddedPlayerCanonicalURL(t *testing.T) {
	tests := []struct {
		name            string
		data            embeddedPlayerData
		entityType      entityType
		id              uint64
		want            string
		wantErrContains string
	}{
		{
			name:       "track linkback",
			data:       embeddedData("https://fantastictrax.bandcamp.com/track/torn"),
			entityType: trackEntityType,
			id:         852027615,
			want:       "https://fantastictrax.bandcamp.com/track/torn",
		},
		{
			name:       "album linkback",
			data:       embeddedData("https://fantastictrax.bandcamp.com/album/mad-gavs-909-edits"),
			entityType: albumEntityType,
			id:         1884059585,
			want:       "https://fantastictrax.bandcamp.com/album/mad-gavs-909-edits",
		},
		{
			name: "matching track fallback",
			data: embeddedDataWithTracks("https://fantasticvoyage.com/track/torn", []embeddedTrack{
				{id: 1, titleLink: "https://other.bandcamp.com/track/other"},
				{id: 852027615, titleLink: "https://fantastictrax.bandcamp.com/track/torn"},
			}),
			entityType: trackEntityType,
			id:         852027615,
			want:       "https://fantastictrax.bandcamp.com/track/torn",
		},
		{
			name:            "wrong track type",
			data:            embeddedData("https://fantastictrax.bandcamp.com/album/mad-gavs-909-edits"),
			entityType:      trackEntityType,
			id:              852027615,
			wantErrContains: "canonical track url not found",
		},
		{
			name:            "wrong album type",
			data:            embeddedData("https://fantastictrax.bandcamp.com/track/torn"),
			entityType:      albumEntityType,
			id:              1884059585,
			wantErrContains: "canonical album url not found",
		},
		{
			name:            "custom host rejected",
			data:            embeddedData("https://fantasticvoyage.com/track/torn"),
			entityType:      trackEntityType,
			id:              852027615,
			wantErrContains: "canonical track url not found",
		},
		{
			name: "nonmatching track fallback rejected",
			data: embeddedDataWithTracks("", []embeddedTrack{
				{id: 1, titleLink: "https://fantastictrax.bandcamp.com/track/torn"},
			}),
			entityType:      trackEntityType,
			id:              852027615,
			wantErrContains: "canonical track url not found",
		},
		{
			name:            "empty tracks do not panic",
			data:            embeddedPlayerData{},
			entityType:      trackEntityType,
			id:              852027615,
			wantErrContains: "canonical track url not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient().embeddedPlayerCanonicalURL(tt.data, tt.entityType, tt.id)
			if tt.wantErrContains != "" {
				require.Empty(t, got)
				require.ErrorContains(t, err, tt.wantErrContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestClient_resolveSearchResultURL(t *testing.T) {
	tests := []struct {
		name            string
		entityType      entityType
		id              uint64
		status          int
		body            string
		wantPath        string
		want            string
		wantErrContains string
	}{
		{
			name:       "track",
			entityType: trackEntityType,
			id:         852027615,
			status:     http.StatusOK,
			body:       `<html><body><div data-player-data="{&quot;linkback&quot;:&quot;https://fantastictrax.bandcamp.com/track/natalie-imbruglia-torn-mad-gavs-909-edit&quot;}"></div></body></html>`,
			wantPath:   "/EmbeddedPlayer/track=852027615/",
			want:       "https://fantastictrax.bandcamp.com/track/natalie-imbruglia-torn-mad-gavs-909-edit",
		},
		{
			name:       "album",
			entityType: albumEntityType,
			id:         1884059585,
			status:     http.StatusOK,
			body:       `<html><body><div data-player-data="{&quot;linkback&quot;:&quot;https://fantastictrax.bandcamp.com/album/mad-gavs-909-edits&quot;}"></div></body></html>`,
			wantPath:   "/EmbeddedPlayer/album=1884059585/",
			want:       "https://fantastictrax.bandcamp.com/album/mad-gavs-909-edits",
		},
		{
			name:            "unexpected status",
			entityType:      trackEntityType,
			id:              852027615,
			status:          http.StatusBadGateway,
			wantPath:        "/EmbeddedPlayer/track=852027615/",
			wantErrContains: "unexpected embedded player status code: 502",
		},
		{
			name:            "malformed html",
			entityType:      trackEntityType,
			id:              852027615,
			status:          http.StatusOK,
			body:            `<html><body></body></html>`,
			wantPath:        "/EmbeddedPlayer/track=852027615/",
			wantErrContains: "attribute not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var serverHost string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, serverHost, r.Host)
				require.Equal(t, tt.wantPath, r.URL.Path)
				w.WriteHeader(tt.status)
				_, err := w.Write([]byte(tt.body))
				require.NoError(t, err)
			}))
			defer srv.Close()

			srvURL, err := url.Parse(srv.URL)
			require.NoError(t, err)
			serverHost = srvURL.Host
			client := NewClient(WithAPI(srvURL.Scheme, srvURL.Host))

			got, err := client.resolveSearchResultURL(t.Context(), tt.entityType, tt.id)
			if tt.wantErrContains != "" {
				require.Empty(t, got)
				require.ErrorContains(t, err, tt.wantErrContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestClient_resolveSearchResultURLUsesCallerContext(t *testing.T) {
	type contextKey string
	const key contextKey = "test-key"

	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), key, "test-value"))
	client := NewClient(
		WithAPIClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			require.Equal(t, "test-value", r.Context().Value(key))
			cancel()
			return nil, r.Context().Err()
		})}),
	)

	got, err := client.resolveSearchResultURL(ctx, trackEntityType, 852027615)

	require.Empty(t, got)
	require.ErrorIs(t, err, context.Canceled)
}

func TestClient_resolveSearchResultURLRejectsMissingIDWithoutRequest(t *testing.T) {
	client := NewClient(WithAPIClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("unexpected HTTP request")
		return nil, nil
	})}))

	got, err := client.resolveSearchResultURL(t.Context(), trackEntityType, 0)

	require.Empty(t, got)
	require.ErrorContains(t, err, "missing numeric search result id")
}

type embeddedTrack struct {
	id        uint64
	titleLink string
}

func embeddedData(linkback string) embeddedPlayerData {
	return embeddedPlayerData{Linkback: linkback}
}

func embeddedDataWithTracks(linkback string, tracks []embeddedTrack) embeddedPlayerData {
	data := embeddedPlayerData{Linkback: linkback}
	data.Tracks = make([]struct {
		ID        uint64 `json:"id"`
		TitleLink string `json:"title_link"`
	}, len(tracks))
	for i, track := range tracks {
		data.Tracks[i].ID = track.id
		data.Tracks[i].TitleLink = track.titleLink
	}
	return data
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
