package isrc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		want      string
		wantValid bool
	}{
		{
			name:      "canonical",
			raw:       "GBARL9300135",
			want:      "GBARL9300135",
			wantValid: true,
		},
		{
			name:      "hyphenated lowercase with whitespace",
			raw:       " gb-arl-93-00135 ",
			want:      "GBARL9300135",
			wantValid: true,
		},
		{
			name: "empty",
		},
		{
			name: "invalid characters",
			raw:  "not-an-isrc",
		},
		{
			name: "misplaced hyphen",
			raw:  "G-BARL9300135",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, valid := Normalize(tt.raw)

			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantValid, valid)
		})
	}
}
