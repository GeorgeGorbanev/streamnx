package translit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCyrillicToLatin(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "земфира",
			output: "zemfira",
		},
		{
			input:  "надежда кадышева",
			output: "nadezhda kadysheva",
		},
		{
			input:  "игорь стравинский",
			output: "igor stravinskii",
		},
		{
			input:  "ЗЕМФИРА",
			output: "ZEMFIRA",
		},
		{
			input:  "НАДЕЖДА КАДЫШЕВА",
			output: "NADEZhDA KADYShEVA",
		},
		{
			input:  "ИГОРЬ СТРАВИНСКИЙ",
			output: "IGOR STRAVINSKII",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s => %s", tt.input, tt.output), func(t *testing.T) {
			require.Equal(t, tt.output, CyrillicToLatin(tt.input))
		})
	}
}

func TestLatinToCyrillic(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "zemfira",
			output: "земфира",
		},
		{
			input:  "nadezhda kadysheva",
			output: "надежда кадышева",
		},
		{
			input:  "igor stravinskii",
			output: "игор стравинскии",
		},
		{
			input:  "ZEMFIRA",
			output: "ЗЕМФИРА",
		},
		{
			input:  "NADEZhDA KADYShEVA",
			output: "НАДЕЖДА КАДЫШЕВА",
		},
		{
			input:  "IGOR STRAVINSKII",
			output: "ИГОР СТРАВИНСКИИ",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s => %s", tt.input, tt.output), func(t *testing.T) {
			require.Equal(t, tt.output, LatinToCyrillic(tt.input))
		})
	}
}
