package soundcloud

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
			input: "https://soundcloud.com/forss/flickermood",
			want:  "forss:flickermood",
		},
		{
			input: "https://www.soundcloud.com/forss/flickermood",
			want:  "forss:flickermood",
		},
		{
			input: "https://soundcloud.com/forss/flickermood?si=abc123",
			want:  "forss:flickermood",
		},
		{
			input: "https://soundcloud.com/forss/sets/soulhack",
			want:  "",
		},
		{
			input: "https://soundcloud.com/forss",
			want:  "",
		},
		{
			input: "https://example.com/forss/flickermood",
			want:  "",
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
			input: "https://soundcloud.com/forss/sets/soulhack",
			want:  "forss:soulhack",
		},
		{
			input: "https://www.soundcloud.com/forss/sets/soulhack",
			want:  "forss:soulhack",
		},
		{
			input: "https://soundcloud.com/forss/sets/soulhack?si=abc123",
			want:  "forss:soulhack",
		},
		{
			input: "https://soundcloud.com/forss/flickermood",
			want:  "",
		},
		{
			input: "https://soundcloud.com/forss/sets/",
			want:  "",
		},
		{
			input: "https://example.com/forss/sets/soulhack",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectAlbumID(tt.input)
			require.Equal(t, tt.want, result)
		})
	}
}
