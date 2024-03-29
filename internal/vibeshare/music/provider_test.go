package music

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidProvider(t *testing.T) {
	tests := []struct {
		provider Provider
		expected bool
	}{
		{
			provider: Spotify,
			expected: true,
		},
		{
			provider: Yandex,
			expected: true,
		},
		{
			provider: Youtube,
			expected: true,
		},
		{
			provider: Provider("not provider"),
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.provider), func(t *testing.T) {
			result := IsValidProvider(tt.provider)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestProvider_Name(t *testing.T) {
	tests := []struct {
		provider Provider
		expected string
	}{
		{
			provider: Spotify,
			expected: "Spotify",
		},
		{
			provider: Yandex,
			expected: "Yandex",
		},
		{
			provider: Youtube,
			expected: "Youtube",
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.provider), func(t *testing.T) {
			result := tt.provider.Name()
			require.Equal(t, tt.expected, result)
		})
	}
}
