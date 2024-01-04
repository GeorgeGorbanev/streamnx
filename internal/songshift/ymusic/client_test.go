package ymusic_test

import (
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
	ymusic_utils "github.com/GeorgeGorbanev/songshift/tests/utils/ymusic"

	"github.com/stretchr/testify/require"
)

func TestClient_SearchTrack(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  *ymusic.Track
	}{
		{
			name:  "when track found",
			query: ymusic_utils.TrackFixtureMassiveAttackAngel.SearchQuery,
			want:  ymusic_utils.TrackFixtureMassiveAttackAngel.Track,
		},
		{
			name:  "when track not found",
			query: "any impossible query",
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := ymusic_utils.NewAPIServerMock(t)
			defer apiServerMock.Close()

			client := ymusic.NewClient(ymusic.WithAPIURL(apiServerMock.URL))

			result, err := client.SearchTrack(tt.query)
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
