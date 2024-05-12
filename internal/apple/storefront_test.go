package apple

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidStorefront(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid code: us",
			input:    "us",
			expected: true,
		},
		{
			name:     "Valid code: gb",
			input:    "gb",
			expected: true,
		},
		{
			name:     "Valid uppercase code: RU",
			input:    "RU",
			expected: true,
		},
		{
			name:     "Invalid code: xx",
			input:    "xx",
			expected: false,
		},
		{
			name:     "Empty code",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidStorefront(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
