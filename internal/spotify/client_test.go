package spotify

import (
	"context"
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
	sampleBasicAuth = "Basic c2FtcGxlQ2xpZW50SUQ6c2FtcGxlQ2xpZW50U2VjcmV0"
)

func TestHTTPClient_FetchTrack(t *testing.T) {
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

	client := NewHTTPClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)

	track, err := client.FetchTrack(t.Context(), "sampletrackid")
	require.NoError(t, err)
	require.Equal(t, &Track{
		ID: "sampletrackid",
		Artists: []Artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, track)
}

func TestHTTPClient_SearchTrack(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		authorization := r.Header.Get("Authorization")
		require.Equal(t, authorization, "Bearer mock_access_token")
		require.Equal(t, r.URL.Path, "/v1/search")
		require.Equal(t, r.URL.Query().Get("q"), "artist:Sample Artist track:Sample Track")
		_, err := w.Write([]byte(`{
			"tracks": {
				"items": [{		
					"id": "sampletrackid",	
					"artists": [{"name": "Sample Artist"}],	
					"name": "Sample Track"	
				}]	
			}
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewHTTPClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)

	track, err := client.SearchTrack(t.Context(), "Sample Artist", "Sample Track")
	require.NoError(t, err)
	require.Equal(t, &Track{
		ID: "sampletrackid",
		Artists: []Artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, track)
}

func TestHTTPClient_FetchAlbum(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		authorization := r.Header.Get("Authorization")
		require.Equal(t, authorization, "Bearer mock_access_token")
		require.Equal(t, r.URL.Path, "/v1/albums/samplealbumid")
		_, err := w.Write([]byte(`{
			"id": "samplealbumid",
			"artists": [{"name": "Sample Artist"}],
			"name": "Sample Album"
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewHTTPClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)

	album, err := client.FetchAlbum(t.Context(), "samplealbumid")
	require.NoError(t, err)
	require.Equal(t, &Album{
		ID:   "samplealbumid",
		Name: "Sample Album",
		Artists: []Artist{
			{
				Name: "Sample Artist",
			},
		},
	}, album)
}

func TestHTTPClient_SearchAlbum(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		authorization := r.Header.Get("Authorization")
		require.Equal(t, authorization, "Bearer mock_access_token")
		require.Equal(t, r.URL.Path, "/v1/search")
		require.Equal(t, r.URL.Query().Get("q"), "artist:Sample Artist album:Sample Album")
		_, err := w.Write([]byte(`{
			"albums": {
				"items": [{		
					"id": "samplealbumid",
					"artists": [{"name": "Sample Artist"}],
					"name": "Sample Album"
				}]	
			}
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewHTTPClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	album, err := client.SearchAlbum(ctx, "Sample Artist", "Sample Album")
	require.NoError(t, err)
	require.Equal(t, &Album{
		ID:   "samplealbumid",
		Name: "Sample Album",
		Artists: []Artist{
			{
				Name: "Sample Artist",
			},
		},
	}, album)
}

func TestHTTPClient_TokenNotExpired(t *testing.T) {
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

	client := NewHTTPClient(
		&sampleCredentials,
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt:   time.Now(),
		ExpiresIn:   3600,
		AccessToken: "mock_access_token",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	track, err := client.FetchTrack(ctx, "sampletrackid")
	require.NoError(t, err)
	require.Equal(t, &Track{
		ID: "sampletrackid",
		Artists: []Artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, track)
}

func TestHTTPClient_RefreshTokenWhenExpired(t *testing.T) {
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

	client := NewHTTPClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt: time.Now().Add(-time.Hour * 24),
		ExpiresIn: 1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	track, err := client.FetchTrack(ctx, "sampletrackid")
	require.NoError(t, err)
	require.Equal(t, &Track{
		ID: "sampletrackid",
		Artists: []Artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, track)
}

func TestHTTPClient_RefreshTokenWhenUnauthorized(t *testing.T) {
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

	client := NewHTTPClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt:   time.Now(),
		ExpiresIn:   3600,
		AccessToken: "not_expired_token_to_refresh",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	track, err := client.FetchTrack(ctx, "sampletrackid")
	require.NoError(t, err)
	require.Equal(t, &Track{
		ID: "sampletrackid",
		Artists: []Artist{
			{
				Name: "Sample Artist",
			},
		},
		Name: "Sample Track",
	}, track)
}

func TestHTTPClient_APIError(t *testing.T) {
	mockAuthServer := newAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		authorization := r.Header.Get("Authorization")
		require.Equal(t, authorization, "Bearer mock_access_token")
		require.Equal(t, r.URL.Path, "/v1/tracks/sampletrackid")
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(`{
			"error" : {
				"status" : 403,
				"message" : "Spotify is unavailable in this country"
			}
		}`))
		require.NoError(t, err)
	}))
	defer mockAPIServer.Close()

	client := NewHTTPClient(
		&sampleCredentials,
		WithAuthURL(mockAuthServer.URL),
		WithAPIURL(mockAPIServer.URL),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	track, err := client.FetchTrack(ctx, "sampletrackid")
	require.Errorf(t, err,
		"failed to send request: unexpected API response: 403 Spotify is unavailable in this country")
	require.Nil(t, track)
}

func TestHTTPClient_HandlesRequestError(t *testing.T) {
	mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mockAPIServer.Close()

	client := NewHTTPClient(
		&sampleCredentials,
		WithAPIURL(mockAPIServer.URL),
	)
	client.token = &token{
		fetchedAt:   time.Now(),
		ExpiresIn:   3600,
		AccessToken: "valid_token",
	}

	track, err := client.FetchTrack(t.Context(), "sampletrackid")
	require.Error(t, err)

	port := mockAPIServer.Listener.Addr().(*net.TCPAddr).Port
	expectedError := fmt.Sprintf(
		`failed to send request: `+
			`Get "%s/v1/tracks/sampletrackid?": `+
			`dial tcp 127.0.0.1:%d: `+
			`connect: connection refused`,
		mockAPIServer.URL, port)
	require.Equal(t, expectedError, err.Error())
	require.Nil(t, track)
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

func TestHTTPClient_ConcurrentTokenRefresh(t *testing.T) {
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

	client := NewHTTPClient(
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
			_, err := client.FetchTrack(context.Background(), "test-track")
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
	require.Equal(t, int64(concurrency), atomic.LoadInt64(&apiRequests))
}

func TestHTTPClient_ConcurrentTokenRefreshAfter401(t *testing.T) {
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

	client := NewHTTPClient(
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

	_, err := client.FetchTrack(context.Background(), "test-track")
	require.NoError(t, err)
	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
}

func TestHTTPClient_ConcurrentExpiredTokenRefresh(t *testing.T) {
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

	client := NewHTTPClient(
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
			_, err := client.FetchTrack(context.Background(), "test-track")
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
}

func TestHTTPClient_TokenDoubleChecking(t *testing.T) {
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

	client := NewHTTPClient(
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
			_, err := client.FetchTrack(context.Background(), "test-track")
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	require.Equal(t, int64(1), atomic.LoadInt64(&tokenRequests))
}
