package spotify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthHeader(t *testing.T) {
	client := &Client{
		clientID:     "testID",
		clientSecret: "testSecret",
	}

	result := client.authHeader()

	require.Equal(t, "Basic dGVzdElEOnRlc3RTZWNyZXQ=", result)
}

func TestFetchToken(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/token", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		err := json.NewEncoder(w).Encode(map[string]any{
			"access_token": "test_access_token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
		require.NoError(t, err)
	}))

	defer mockServer.Close()

	client := NewClient("", mockServer.URL, "client_id", "client_secret")
	token, err := client.FetchToken()
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Equal(t, "test_access_token", token.AccessToken)
	require.Equal(t, "Bearer", token.TokenType)
	require.Equal(t, 3600, token.ExpiresIn)
}
