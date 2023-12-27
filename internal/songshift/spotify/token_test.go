package spotify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToken_authHeader(t *testing.T) {
	token := Token{AccessToken: "sampleAccessToken"}
	result := token.authHeader()
	require.Equal(t, "Bearer sampleAccessToken", result)
}
