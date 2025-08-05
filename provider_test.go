package streamnx

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
			code: "dz",
			want: Deezer,
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

func TestProvider_DetectCloakEntityID(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		want     string
	}{
		{
			name: "has cloakIDParser",
			provider: Provider{
				cloakIDParser: func(url string) string {
					return "sample result"
				},
			},
			want: "sample result",
		},
		{
			name:     "has no cloakIDParser",
			provider: Provider{},
			want:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.provider.DetectCloakEntityID("sample url")
			require.Equal(t, tt.want, result)
		})
	}
}
