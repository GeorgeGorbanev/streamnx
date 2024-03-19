package spotify_test

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	spotify_utils "github.com/GeorgeGorbanev/vibeshare/tests/utils/spotify"

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

	track, err := client.GetTrack(spotify_utils.TrackFixtureMassiveAttackAngel.Track.ID)
	require.NoError(t, err)
	require.Equal(t, spotify_utils.TrackFixtureMassiveAttackAngel.Track, track)
}

func TestSearchTrack(t *testing.T) {
	mockAuthServer := spotify_utils.NewAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := spotify_utils.NewAPIServerMock(t)
	defer mockAPIServer.Close()

	client := spotify.NewClient(
		&spotify_utils.SampleCredentials,
		spotify.WithAuthURL(mockAuthServer.URL),
		spotify.WithAPIURL(mockAPIServer.URL),
	)

	track, err := client.SearchTrack(
		spotify_utils.TrackFixtureMassiveAttackAngel.Track.Artists[0].Name,
		spotify_utils.TrackFixtureMassiveAttackAngel.Track.Name,
	)
	require.NoError(t, err)
	require.Equal(t, spotify_utils.TrackFixtureMassiveAttackAngel.Track, track)
}

func TestGetAlbum(t *testing.T) {
	mockAuthServer := spotify_utils.NewAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := spotify_utils.NewAPIServerMock(t)
	defer mockAPIServer.Close()

	client := spotify.NewClient(
		&spotify_utils.SampleCredentials,
		spotify.WithAuthURL(mockAuthServer.URL),
		spotify.WithAPIURL(mockAPIServer.URL),
	)

	album, err := client.GetAlbum(spotify_utils.AlbumFixtureRadioheadAmnesiac.Album.ID)
	require.NoError(t, err)
	require.Equal(t, spotify_utils.AlbumFixtureRadioheadAmnesiac.Album, album)
}

func TestSearchAlbum(t *testing.T) {
	mockAuthServer := spotify_utils.NewAuthServerMock(t)
	defer mockAuthServer.Close()

	mockAPIServer := spotify_utils.NewAPIServerMock(t)
	defer mockAPIServer.Close()

	client := spotify.NewClient(
		&spotify_utils.SampleCredentials,
		spotify.WithAuthURL(mockAuthServer.URL),
		spotify.WithAPIURL(mockAPIServer.URL),
	)

	album, err := client.SearchAlbum(
		spotify_utils.AlbumFixtureRadioheadAmnesiac.Album.Artists[0].Name,
		spotify_utils.AlbumFixtureRadioheadAmnesiac.Album.Name,
	)
	require.NoError(t, err)
	require.Equal(t, spotify_utils.AlbumFixtureRadioheadAmnesiac.Album, album)
}
