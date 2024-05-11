package spotify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentials_authHeader(t *testing.T) {
	сredentials := &Credentials{
		ClientID:     "testID",
		ClientSecret: "testSecret",
	}

	result := сredentials.authHeader()

	require.Equal(t, "Basic dGVzdElEOnRlc3RTZWNyZXQ=", result)
}
