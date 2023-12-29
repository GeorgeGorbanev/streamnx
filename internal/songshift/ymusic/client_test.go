package ymusic_test

import (
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
	ymusic_utils "github.com/GeorgeGorbanev/songshift/tests/utils/ymusic"

	"github.com/stretchr/testify/require"
)

func TestClient_Search(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  *ymusic.SearchResponse
	}{
		{
			name:  "when track found",
			query: ymusic_utils.SampleSearchQuery,
			want:  &ymusic_utils.SampleSearchResponse,
		},
		{
			name:  "when track not found",
			query: "any impossible query",
			want:  &ymusic_utils.EmptySearchResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiServerMock := ymusic_utils.NewAPIServerMock(t, ymusic_utils.SampleSearchQuery)
			defer apiServerMock.Close()

			client := ymusic.NewClient(ymusic.WithAPIURL(apiServerMock.URL))

			result, err := client.Search(tt.query)
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}
