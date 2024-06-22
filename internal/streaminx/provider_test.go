package streaminx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindProviderByCode(t *testing.T) {
	tests := []struct {
		code string
		want *Provider
	}{
		{
			code: "ap",
			want: Apple,
		},
		{
			code: "sf",
			want: Spotify,
		},
		{
			code: "ya",
			want: Yandex,
		},
		{
			code: "yt",
			want: Youtube,
		},
		{
			code: "unknown",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := FindProviderByCode(tt.code)
			require.Equal(t, tt.want, result)
		})
	}
}
