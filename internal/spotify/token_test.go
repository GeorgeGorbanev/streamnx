package spotify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToken_authHeader(t *testing.T) {
	token := token{AccessToken: "sampleAccessToken"}
	result := token.authHeader()
	require.Equal(t, "Bearer sampleAccessToken", result)
}

func TestToken_isExpired(t *testing.T) {
	tests := []struct {
		name  string
		token token
		want  bool
	}{
		{
			name: "when token is expired",
			token: token{
				ExpiresIn: 3600,
				fetchedAt: time.Now().Add(-3601 * time.Second),
			},
			want: true,
		},
		{
			name: "when token is not expired",
			token: token{
				ExpiresIn: 3600,
				fetchedAt: time.Now().Add(-3599 * time.Second),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			result := tt.token.isExpired()
			require.Equal(t1, tt.want, result)
		})
	}
}
