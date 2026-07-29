package yandex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYandexTrackKeyRoundTrip(t *testing.T) {
	key, err := newTrackKey("14599266", "609676")
	require.NoError(t, err)

	id, err := keyScheme.Dump(key)
	require.NoError(t, err)
	require.Equal(t, "14599266:609676", id)

	albumID, trackID, err := parseTrackKeyParts(id)
	require.NoError(t, err)
	require.Equal(t, "14599266", albumID)
	require.Equal(t, "609676", trackID)
}

func TestYandexTrackKeyRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "missing album id", id: ":609676"},
		{name: "missing track id", id: "14599266:"},
		{name: "non-numeric album id", id: "album:609676"},
		{name: "non-numeric track id", id: "14599266:track"},
		{name: "extra part", id: "14599266:609676:extra"},
		{name: "plain legacy track id", id: "609676"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseTrackKeyParts(tt.id)
			require.Error(t, err)
		})
	}
}
