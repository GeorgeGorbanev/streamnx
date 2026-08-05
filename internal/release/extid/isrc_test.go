package extid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeISRC(t *testing.T) {
	isrcTests := []struct {
		name      string
		rawISRC   string
		wantISRC  string
		wantValid bool
	}{
		{
			name:      "canonical",
			rawISRC:   "GBARL9300135",
			wantISRC:  "GBARL9300135",
			wantValid: true,
		},
		{
			name:      "hyphenated lowercase with whitespace",
			rawISRC:   " gb-arl-93-00135 ",
			wantISRC:  "GBARL9300135",
			wantValid: true,
		},
		{
			name: "empty",
		},
		{
			name:    "invalid characters",
			rawISRC: "not-an-isrc",
		},
		{
			name:    "misplaced hyphen",
			rawISRC: "G-BARL9300135",
		},
	}

	for _, isrcTest := range isrcTests {
		t.Run(isrcTest.name, func(t *testing.T) {
			isrcGot, isrcValid := NormalizeISRC(isrcTest.rawISRC)

			require.Equal(t, isrcTest.wantISRC, isrcGot)
			require.Equal(t, isrcTest.wantValid, isrcValid)
		})
	}
}
