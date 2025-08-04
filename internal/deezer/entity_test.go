package deezer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectTrackID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "https://www.deezer.com/us/track/441481432",
			want:  "441481432",
		},
		{
			input: "https://deezer.com/us/track/441481432",
			want:  "441481432",
		},
		{
			input: "https://deezer.com/track/441481432",
			want:  "441481432",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectTrackID(tt.input)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestDetectAlbumID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "https://www.deezer.com/us/album/53564982",
			want:  "53564982",
		},
		{
			input: "https://deezer.com/us/album/53564982",
			want:  "53564982",
		},
		{
			input: "https://deezer.com/album/53564982",
			want:  "53564982",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectAlbumID(tt.input)
			require.Equal(t, tt.want, result)
		})
	}

}

func TestDetectCloakID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "https://link.deezer.com/s/30FbcgrctxIrNQImDnVEZ",
			want:  "30FbcgrctxIrNQImDnVEZ",
		},
		{
			input: "https://link.deezer.com/s/jHWZaLoRJutY3TiVA",
			want:  "jHWZaLoRJutY3TiVA",
		},
		{
			input: "not a deezer link",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectCloakID(tt.input)
			require.Equal(t, tt.want, result)
		})
	}
}
