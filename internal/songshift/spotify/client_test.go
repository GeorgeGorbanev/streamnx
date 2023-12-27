package spotify_test

import (
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	spotify_utils "github.com/GeorgeGorbanev/songshift/tests/utils/spotify"

	"github.com/stretchr/testify/require"
)

func TestFetchToken(t *testing.T) {
	mockServer := spotify_utils.NewAuthServerMock(t)
	defer mockServer.Close()

	client := spotify.NewClient(&spotify_utils.SampleCredentials, spotify.WithAuthURL(mockServer.URL))
	token, err := client.FetchToken()
	require.NoError(t, err)
	require.NotNil(t, token)
	require.Equal(t, &spotify_utils.SampleToken, token)
}

func TestGetTrack(t *testing.T) {
	mockAuthServer := spotify_utils.NewAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := spotify_utils.NewAPIServerMock(t)
	defer mockAPIServer.Close()

	client := spotify.NewClient(
		&spotify_utils.SampleCredentials,
		spotify.WithAuthURL(mockAuthServer.URL),
		spotify.WithAPIURL(mockAPIServer.URL),
	)

	track, err := client.GetTrack(spotify_utils.SampleTrack.ID)
	require.NoError(t, err)
	require.Equal(t, &spotify_utils.SampleTrack, track)
}
