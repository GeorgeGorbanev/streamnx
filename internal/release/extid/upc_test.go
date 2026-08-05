package extid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeUPC(t *testing.T) {
	upcTests := []struct {
		name      string
		rawUPC    string
		wantUPC   string
		wantValid bool
	}{
		{
			name:      "upc a",
			rawUPC:    "196006422677",
			wantUPC:   "196006422677",
			wantValid: true,
		},
		{
			name:      "upc a with whitespace",
			rawUPC:    " 196006422677\n",
			wantUPC:   "196006422677",
			wantValid: true,
		},
		{
			name:      "ean 13 equivalent to upc a",
			rawUPC:    "0859381157694",
			wantUPC:   "859381157694",
			wantValid: true,
		},
		{
			name:      "ean 13 without upc a equivalent",
			rawUPC:    "4006381333931",
			wantUPC:   "4006381333931",
			wantValid: true,
		},
		{name: "empty"},
		{name: "wrong length", rawUPC: "19600642267"},
		{name: "invalid characters", rawUPC: "19600642267X"},
		{name: "invalid check digit", rawUPC: "196006422678"},
	}

	for _, upcTest := range upcTests {
		t.Run(upcTest.name, func(t *testing.T) {
			upcGot, upcValid := NormalizeUPC(upcTest.rawUPC)

			require.Equal(t, upcTest.wantUPC, upcGot)
			require.Equal(t, upcTest.wantValid, upcValid)
		})
	}
}
