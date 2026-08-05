package spotify

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const sampleBasicAuth = "Basic c2FtcGxlQ2xpZW50SUQ6c2FtcGxlQ2xpZW50U2VjcmV0"

var (
	sampleCredentials = Credentials{
		ClientID:     "sampleClientID",
		ClientSecret: "sampleClientSecret",
	}
	sampleToken = token{
		AccessToken: "mock_access_token",
		TokenType:   "Bearer",
		ExpiresIn:   360,
	}
)

func TestClient_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedTrack track
		expectedErr   string
	}{
		{
			name:       "found id",
			statusCode: http.StatusOK,
			responseBody: `{
				"id": "sampletrackid",
				"artists": [{"name": "Sample Artist"}],
				"name": "Sample Track"
			}`,
			expectedTrack: track{
				ID: "sampletrackid",
				Artists: []artist{
					{
						Name: "Sample Artist",
					},
				},
				Name: "Sample Track",
			},
		},
		{
			name:       "api error",
			statusCode: http.StatusForbidden,
			responseBody: `{
				"error" : {
					"status" : 403,
					"message" : "Spotify is unavailable in this country"
				}
			}`,
			expectedErr: "unexpected api response: 403 Spotify is unavailable in this country",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthServer := newAuthServerMock(t)
			defer mockAuthServer.Close()

			mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer mock_access_token", r.Header.Get("Authorization"))
				require.Equal(t, "/v1/tracks/sampletrackid", r.URL.Path)

				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			}))
			defer mockAPIServer.Close()

			client := NewClient(
				&sampleCredentials,
				WithAuthURL(mockAuthServer.URL),
				WithAPIURL(mockAPIServer.URL),
			)

			track, err := client.fetchTrack(t.Context(), "sampletrackid")
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Zero(t, track)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, track)
			}
		})
	}
}

func TestClient_searchTracks(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedTracks []track
		expectedErr    string
	}{
		{
			name:       "multiple candidates",
			statusCode: http.StatusOK,
			responseBody: `{
				"tracks": {
					"items": [{
						"id": "sampletrackid",
						"artists": [{"name": "Sample Artist"}],
						"name": "Sample Track"
					}, {
						"id": "sampletrackid2",
						"artists": [{"name": "Sample Artist 2"}],
						"name": "Sample Track 2"
					}]
				}
			}`,
			expectedTracks: []track{
				{
					ID: "sampletrackid",
					Artists: []artist{
						{
							Name: "Sample Artist",
						},
					},
					Name: "Sample Track",
				},
				{
					ID: "sampletrackid2",
					Artists: []artist{
						{
							Name: "Sample Artist 2",
						},
					},
					Name: "Sample Track 2",
				},
			},
		},
		{
			name:           "empty candidates",
			statusCode:     http.StatusOK,
			responseBody:   `{"tracks":{"items":[]}}`,
			expectedTracks: []track{},
		},
		{
			name:       "api error",
			statusCode: http.StatusForbidden,
			responseBody: `{
				"error" : {
					"status" : 403,
					"message" : "Spotify is unavailable in this country"
				}
			}`,
			expectedErr: "unexpected api response: 403 Spotify is unavailable in this country",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthServer := newAuthServerMock(t)
			defer mockAuthServer.Close()

			mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer mock_access_token", r.Header.Get("Authorization"))
				require.Equal(t, "/v1/search", r.URL.Path)
				require.Equal(t, "artist:Sample Artist track:Sample Track", r.URL.Query().Get("q"))
				require.Equal(t, "10", r.URL.Query().Get("limit"))

				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			}))
			defer mockAPIServer.Close()

			client := NewClient(
				&sampleCredentials,
				WithAuthURL(mockAuthServer.URL),
				WithAPIURL(mockAPIServer.URL),
			)

			tracks, err := client.searchTracks(t.Context(), "Sample Artist", "Sample Track")
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, tracks)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTracks, tracks)
			}
		})
	}
}

func TestClient_fetchTracksByISRC(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "Bearer mock_access_token", r.Header.Get("Authorization"))
		require.Equal(t, "/v1/search", r.URL.Path)
		require.Equal(t, "isrc:GBARL9300135", r.URL.Query().Get("q"))
		require.Equal(t, "track", r.URL.Query().Get("type"))
		require.Equal(t, "10", r.URL.Query().Get("limit"))
		_, err := w.Write([]byte(`{
			"tracks": {
				"items": [{
					"id": "sampletrackid",
					"external_ids": {"isrc": "GBARL9300135"},
					"name": "Sample Track"
				}]
			}
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)

	tracks, err := client.fetchTracksByISRC(t.Context(), "GBARL9300135")

	require.NoError(t, err)
	require.Equal(t, []track{{
		ID:          "sampletrackid",
		ExternalIDs: externalIDs{ISRC: "GBARL9300135"},
		Name:        "Sample Track",
	}}, tracks)
}

func TestClient_fetchAlbumsByUPC(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "Bearer mock_access_token", r.Header.Get("Authorization"))
		switch r.URL.Path {
		case "/v1/search":
			require.Equal(t, "upc:196006422677", r.URL.Query().Get("q"))
			require.Equal(t, "album", r.URL.Query().Get("type"))
			require.Equal(t, "10", r.URL.Query().Get("limit"))
			_, err := w.Write([]byte(`{
				"albums": {"items": [{"id": "1P0Ox1jln7nn2J9BT6yMhL"}]}
			}`))
			require.NoError(t, err)
		case "/v1/albums/1P0Ox1jln7nn2J9BT6yMhL":
			_, err := w.Write([]byte(`{
				"id": "1P0Ox1jln7nn2J9BT6yMhL",
				"external_ids": {"upc": "196006422677"},
				"name": "კინომუსიკა: გოგი ცაბაძის შემოქმედება ნაწილი XIII"
			}`))
			require.NoError(t, err)
		default:
			require.Fail(t, "unexpected path", r.URL.Path)
		}
	}))
	defer mockAPIServer.Close()

	client := NewClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)

	albums, err := client.fetchAlbumsByUPC(t.Context(), "196006422677")

	require.NoError(t, err)
	require.Equal(t, []album{{
		ID:          "1P0Ox1jln7nn2J9BT6yMhL",
		ExternalIDs: externalIDs{UPC: "196006422677"},
		Name:        "კინომუსიკა: გოგი ცაბაძის შემოქმედება ნაწილი XIII",
	}}, albums)
}

func TestClient_fetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedAlbum album
		expectedErr   string
	}{
		{
			name:       "found ID",
			statusCode: http.StatusOK,
			responseBody: `{
				"id": "samplealbumid",
				"artists": [{"name": "Sample Artist"}],
				"label": "Sample Label",
				"name": "Sample Album",
				"tracks": {
					"items": [{
						"id": "firsttrackid",
						"artists": [{"name": "First Artist"}],
						"name": "First Track"
					}, {
						"id": "secondtrackid",
						"artists": [{"name": "Second Artist"}],
						"name": "Second Track"
					}]
				}
			}`,
			expectedAlbum: album{
				ID:    "samplealbumid",
				Label: "Sample Label",
				Name:  "Sample Album",
				Artists: []artist{
					{
						Name: "Sample Artist",
					},
				},
				Tracks: albumTracks{
					Items: []track{
						{
							ID: "firsttrackid",
							Artists: []artist{
								{
									Name: "First Artist",
								},
							},
							Name: "First Track",
						},
						{
							ID: "secondtrackid",
							Artists: []artist{
								{
									Name: "Second Artist",
								},
							},
							Name: "Second Track",
						},
					},
				},
			},
		},
		{
			name:       "api error",
			statusCode: http.StatusForbidden,
			responseBody: `{
				"error" : {
					"status" : 403,
					"message" : "Spotify is unavailable in this country"
				}
			}`,
			expectedErr: "unexpected api response: 403 Spotify is unavailable in this country",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthServer := newAuthServerMock(t)
			defer mockAuthServer.Close()

			mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer mock_access_token", r.Header.Get("Authorization"))
				require.Equal(t, "/v1/albums/samplealbumid", r.URL.Path)

				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			}))
			defer mockAPIServer.Close()

			client := NewClient(
				&sampleCredentials,
				WithAuthURL(mockAuthServer.URL),
				WithAPIURL(mockAPIServer.URL),
			)

			album, err := client.fetchAlbum(t.Context(), "samplealbumid")
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Zero(t, album)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, album)
			}
		})
	}
}

func TestClient_searchAlbums(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedAlbums []album
		expectedErr    string
	}{
		{
			name:       "multiple candidates",
			statusCode: http.StatusOK,
			responseBody: `{
				"albums": {
					"items": [{
						"id": "samplealbumid",
						"artists": [{"name": "Sample Artist"}],
						"label": "Sample Label",
						"name": "Sample Album"
					}, {
						"id": "samplealbumid2",
						"artists": [{"name": "Sample Artist 2"}],
						"label": "Sample Label 2",
						"name": "Sample Album 2"
					}]
				}
			}`,
			expectedAlbums: []album{
				{
					ID:    "samplealbumid",
					Label: "Sample Label",
					Name:  "Sample Album",
					Artists: []artist{
						{
							Name: "Sample Artist",
						},
					},
				},
				{
					ID:    "samplealbumid2",
					Label: "Sample Label 2",
					Name:  "Sample Album 2",
					Artists: []artist{
						{
							Name: "Sample Artist 2",
						},
					},
				},
			},
		},
		{
			name:           "empty candidates",
			statusCode:     http.StatusOK,
			responseBody:   `{"albums":{"items":[]}}`,
			expectedAlbums: []album{},
		},
		{
			name:       "api error",
			statusCode: http.StatusForbidden,
			responseBody: `{
				"error" : {
					"status" : 403,
					"message" : "Spotify is unavailable in this country"
				}
			}`,
			expectedErr: "unexpected api response: 403 Spotify is unavailable in this country",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthServer := newAuthServerMock(t)
			defer mockAuthServer.Close()

			mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer mock_access_token", r.Header.Get("Authorization"))
				require.Equal(t, "/v1/search", r.URL.Path)
				require.Equal(t, "artist:Sample Artist album:Sample Album", r.URL.Query().Get("q"))
				require.Equal(t, "10", r.URL.Query().Get("limit"))

				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			}))
			defer mockAPIServer.Close()

			client := NewClient(
				&sampleCredentials,
				WithAuthURL(mockAuthServer.URL),
				WithAPIURL(mockAPIServer.URL),
			)

			albums, err := client.searchAlbums(t.Context(), "Sample Artist", "Sample Album")
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, albums)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbums, albums)
			}
		})
	}
}

func TestClient_TokenNotExpired(t *testing.T) {
	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		authorization := r.Header.Get("Authorization")
		require.Equal(t, authorization, "Bearer mock_access_token")
		require.Equal(t, r.URL.Path, "/v1/tracks/sampletrackid")
		_, err := w.Write([]byte(`{
			"id": "sampletrackid",
			"artists": [{"name": "Sample Artist"}],
			"name": "Sample Track"
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewClient(
		&sampleCredentials,
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt:   time.Now(),
		ExpiresIn:   3600,
		AccessToken: "mock_access_token",
	}

	trk, err := client.fetchTrack(t.Context(), "sampletrackid")
	require.NoError(t, err)
	require.Equal(t, track{
		ID: "sampletrackid",
		Artists: []artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, trk)
}

func TestClient_RefreshTokenWhenExpired(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		authorization := r.Header.Get("Authorization")
		require.Equal(t, authorization, "Bearer mock_access_token")
		require.Equal(t, r.URL.Path, "/v1/tracks/sampletrackid")
		_, err := w.Write([]byte(`{
			"id": "sampletrackid",
			"artists": [{"name": "Sample Artist"}],
			"name": "Sample Track"
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt: time.Now().Add(-time.Hour * 24),
		ExpiresIn: 1,
	}

	trk, err := client.fetchTrack(t.Context(), "sampletrackid")
	require.NoError(t, err)
	require.Equal(t, track{
		ID: "sampletrackid",
		Artists: []artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, trk)
}

func TestClient_RefreshTokenWhenUnauthorized(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		authorization := r.Header.Get("Authorization")
		if authorization == "Bearer not_expired_token_to_refresh" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		require.Equal(t, authorization, "Bearer mock_access_token")
		require.Equal(t, r.URL.Path, "/v1/tracks/sampletrackid")
		_, err := w.Write([]byte(`{
			"id": "sampletrackid",
			"artists": [{"name": "Sample Artist"}],
			"name": "Sample Track"
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt:   time.Now(),
		ExpiresIn:   3600,
		AccessToken: "not_expired_token_to_refresh",
	}

	trk, err := client.fetchTrack(t.Context(), "sampletrackid")
	require.NoError(t, err)
	require.Equal(t, track{
		ID: "sampletrackid",
		Artists: []artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, trk)
}

func TestClient_HandlesRequestError(t *testing.T) {
	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mockAPIServer.Close()

	client := NewClient(
		&sampleCredentials,
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt:   time.Now(),
		ExpiresIn:   3600,
		AccessToken: "valid_token",
	}

	track, err := client.fetchTrack(t.Context(), "sampletrackid")
	require.Error(t, err)

	port := mockAPIServer.Listener.Addr().(*net.TCPAddr).Port
	expectedError := fmt.Sprintf(
		`failed to send request: `+
			`Get "%s/v1/tracks/sampletrackid?": `+
			`dial tcp 127.0.0.1:%d: `+
			`connect: connection refused`,
		mockAPIServer.URL, port)
	require.Equal(t, expectedError, err.Error())
	require.Zero(t, track)
}

func newAuthServerMock(t *testing.T) *httptest.Server {
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

func TestClient_ConcurrentTokenRefresh(t *testing.T) {
	var tokenRequests, apiRequests int64

	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&tokenRequests, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":1}`))
	}))
	defer authServer.Close()

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&apiRequests, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"test-track","name":"Test Track"}`))
	}))
	defer apiServer.Close()

	client := NewClient(
		&Credentials{ClientID: "test", ClientSecret: "secret"},
		WithAuthURL(authServer.URL),
		WithAPIURL(apiServer.URL),
	)

	const concurrency = 10
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, err := client.fetchTrack(t.Context(), "test-track")
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
	require.Equal(t, int64(concurrency), atomic.LoadInt64(&apiRequests))
}

func TestClient_ConcurrentTokenRefreshAfter401(t *testing.T) {
	var tokenRequests int64

	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&tokenRequests, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"refreshed-token","token_type":"Bearer","expires_in":3600}`))
	}))
	defer authServer.Close()

	requestCount := make(map[string]int64)
	mu := sync.Mutex{}
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount["total"]++
		auth := r.Header.Get("Authorization")
		count := requestCount[auth]
		requestCount[auth] = count + 1
		mu.Unlock()

		if auth == "Bearer expired-token" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"status":401,"message":"invalid_token"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"test-track","name":"Test Track"}`))
	}))
	defer apiServer.Close()

	client := NewClient(
		&Credentials{ClientID: "test", ClientSecret: "secret"},
		WithAuthURL(authServer.URL),
		WithAPIURL(apiServer.URL),
	)

	client.token = &token{
		AccessToken: "expired-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		fetchedAt:   time.Now(),
	}

	_, err := client.fetchTrack(t.Context(), "test-track")
	require.NoError(t, err)
	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
}

func TestClient_ConcurrentExpiredTokenRefresh(t *testing.T) {
	var tokenRequests int64

	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&tokenRequests, 1)
		time.Sleep(10 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"new-token","token_type":"Bearer","expires_in":3600}`))
	}))
	defer authServer.Close()

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"test-track","name":"Test Track"}`))
	}))
	defer apiServer.Close()

	client := NewClient(
		&Credentials{ClientID: "test", ClientSecret: "secret"},
		WithAuthURL(authServer.URL),
		WithAPIURL(apiServer.URL),
	)

	client.token = &token{
		AccessToken: "expired-token",
		TokenType:   "Bearer",
		ExpiresIn:   1,
		fetchedAt:   time.Now().Add(-2 * time.Second),
	}

	const concurrency = 20
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, err := client.fetchTrack(t.Context(), "test-track")
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
}

func TestClient_TokenDoubleChecking(t *testing.T) {
	var tokenRequests int64

	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt64(&tokenRequests, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"token-` + string(rune('0'+count)) + `","token_type":"Bearer","expires_in":3600}`))
	}))
	defer authServer.Close()

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		require.Contains(t, auth, "token-1")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"test-track","name":"Test Track"}`))
	}))
	defer apiServer.Close()

	client := NewClient(
		&Credentials{ClientID: "test", ClientSecret: "secret"},
		WithAuthURL(authServer.URL),
		WithAPIURL(apiServer.URL),
	)

	const concurrency = 3
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, err := client.fetchTrack(t.Context(), "test-track")
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
}
