package spotify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrack_URL(t *testing.T) {
	track := Track{ID: "sample_id"}
	result := track.URL()
	require.Equal(t, "https://open.spotify.com/track/sample_id", result)
}
