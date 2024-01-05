package ymusic_test

import (
	"strings"
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
	ymusic_utils "github.com/GeorgeGorbanev/songshift/tests/utils/ymusic"

	"github.com/stretchr/testify/require"
)

func TestClient_SearchTrack(t *testing.T) {
	tests := []struct {
		name        string
		queryArtist string
		queryTrack  string
		want        *ymusic.Track
	}{
		{
			name:        "when track found",
			queryArtist: strings.ToLower(ymusic_utils.TrackFixtureMassiveAttackAngel.SearchQueryArtist),
			queryTrack:  strings.ToLower(ymusic_utils.TrackFixtureMassiveAttackAngel.SearchQueryTrack),
			want:        ymusic_utils.TrackFixtureMassiveAttackAngel.Track,
		},
		{
			name:        "when track not found",
			queryArtist: "any impossible artist",
			queryTrack:  "any impossible track",
			want:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := ymusic_utils.NewAPIServerMock(t)
			defer apiServerMock.Close()

			client := ymusic.NewClient(ymusic.WithAPIURL(apiServerMock.URL))

			result, err := client.SearchTrack(tt.queryArtist, tt.queryTrack)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestClient_GetTrack(t *testing.T) {
	tests := []struct {
		name    string
		trackID string
		want    *ymusic.Track
	}{
		{
			name:    "when track found",
			trackID: ymusic_utils.TrackFixtureMassiveAttackAngel.ID,
			want:    ymusic_utils.TrackFixtureMassiveAttackAngel.TrackWithIDString(),
		},
		{
			name:    "when track not found",
			trackID: "0",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := ymusic_utils.NewAPIServerMock(t)
			defer apiServerMock.Close()

			client := ymusic.NewClient(ymusic.WithAPIURL(apiServerMock.URL))

			result, err := client.GetTrack(tt.trackID)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}
