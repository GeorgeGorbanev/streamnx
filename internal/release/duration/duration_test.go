package duration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsToSeconds(t *testing.T) {
	require.Equal(t, 0, MsToSeconds(0))
	require.Equal(t, 213, MsToSeconds(212740))
	require.Equal(t, 214, MsToSeconds(213573))
}

func TestISO8601ToSeconds(t *testing.T) {
	require.Equal(t, 213, ISO8601ToSeconds("PT3M33S"))
	require.Equal(t, 214, ISO8601ToSeconds("P00H03M34S"))
	require.Equal(t, 3724, ISO8601ToSeconds("PT1H2M3.5S"))
	require.Equal(t, 0, ISO8601ToSeconds(""))
}
