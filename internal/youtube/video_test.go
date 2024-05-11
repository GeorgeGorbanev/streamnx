package youtube

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVideo_URL(t *testing.T) {
	video := Video{ID: "dQw4w9WgXcQ"}
	result := video.URL()
	require.Equal(t, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", result)
}
