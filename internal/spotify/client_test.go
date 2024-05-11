package spotify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestHTTPClient_GetTrack(t *testing.T) {
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

	track, err := client.GetTrack("sampletrackid")
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

	track, err := client.SearchTrack("Sample Artist", "Sample Track")
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

func TestHTTPClient_GetAlbum(t *testing.T) {
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

	album, err := client.GetAlbum("samplealbumid")
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

	album, err := client.SearchAlbum("Sample Artist", "Sample Album")
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

	track, err := client.GetTrack("sampletrackid")
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

	track, err := client.GetTrack("sampletrackid")
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

	track, err := client.GetTrack("sampletrackid")
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
