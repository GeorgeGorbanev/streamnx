package ymusic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchResult_AnyTracksFound(t *testing.T) {
	tests := []struct {
		name string
		sr   SearchResult
		want bool
	}{
		{
			name: "With Tracks",
			sr: SearchResult{
				Tracks: TracksSection{
					Results: []Track{{}, {}},
				},
			},
			want: true,
		},
		{
			name: "No Tracks",
			sr: SearchResult{
				Tracks: TracksSection{
					Results: []Track{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sr.AnyTracksFound()
			require.Equal(t, tt.want, got)
		})
	}
}
