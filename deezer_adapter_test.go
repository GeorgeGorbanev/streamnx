package streamnx

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/deezer"
)

type deezerClientMock struct {
	fetchTrackResult  *deezer.Track
	fetchTrackError   error
	searchTrackResult *deezer.Track
	searchTrackError  error
	fetchAlbumResult  *deezer.Album
	fetchAlbumError   error
	searchAlbumResult *deezer.Album
	searchAlbumError  error
	followCloakResult string
	followCloakError  error
}

func (m *deezerClientMock) FetchTrack(ctx context.Context, id string) (*deezer.Track, error) {
	return m.fetchTrackResult, m.fetchTrackError
}

func (m *deezerClientMock) SearchTrack(ctx context.Context, artistName, trackName string) (*deezer.Track, error) {
	return m.searchTrackResult, m.searchTrackError
}

func (m *deezerClientMock) FetchAlbum(ctx context.Context, id string) (*deezer.Album, error) {
	return m.fetchAlbumResult, m.fetchAlbumError
}

func (m *deezerClientMock) SearchAlbum(ctx context.Context, artistName, albumName string) (*deezer.Album, error) {
	return m.searchAlbumResult, m.searchAlbumError
}

func (m *deezerClientMock) FollowCloak(ctx context.Context, cloakCode string) (string, error) {
	return m.followCloakResult, m.followCloakError
}

func TestDeezerAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		track       *deezer.Track
		clientError error
		want        *Entity
		wantError   error
	}{
		{
			name: "successful fetch",
			id:   "123456",
			track: &deezer.Track{
				ID:    123456,
				Title: "Test Song",
				Artist: deezer.Artist{
					Name: "Test Artist",
				},
			},
			want: &Entity{
				ID:       "123456",
				Title:    "Test Song",
				Artist:   "Test Artist",
				URL:      "https://deezer.com/track/123456",
				Provider: Deezer,
				Type:     Track,
			},
		},
		{
			name:        "not found error",
			id:          "404",
			clientError: deezer.NotFoundError,
			wantError:   EntityNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &deezerClientMock{
				fetchTrackResult: tt.track,
				fetchTrackError:  tt.clientError,
			}
			adapter := newDeezerAdapter(client)

			got, err := adapter.FetchTrack(context.Background(), tt.id)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeezerAdapter_FetchCloak(t *testing.T) {
	tests := []struct {
		name        string
		cloakCode   string
		resolvedURL string
		followError error
		fetchResult *deezer.Track
		fetchError  error
		want        *Entity
		wantErr     string
	}{
		{
			name:        "successful track resolution",
			cloakCode:   "abc123",
			resolvedURL: "https://link.deezer.com/?dest=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F123456",
			fetchResult: &deezer.Track{
				ID:    123456,
				Title: "Test Song",
				Artist: deezer.Artist{
					Name: "Test Artist",
				},
			},
			want: &Entity{
				ID:       "123456",
				Title:    "Test Song",
				Artist:   "Test Artist",
				URL:      "https://deezer.com/track/123456",
				Provider: Deezer,
				Type:     Track,
			},
		},
		{
			name:        "follow cloak error",
			cloakCode:   "error123",
			followError: errors.New("network error"),
			wantErr:     "failed to follow cloak link",
		},
		{
			name:        "invalid resolved URL",
			cloakCode:   "invalid123",
			resolvedURL: "https://link.deezer.com/?dest=https%3A%2F%2Fwww.deezer.com%2Fartist%2F123456",
			wantErr:     "entity not found: cloak dest is not a track or album (https://www.deezer.com/artist/123456)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &deezerClientMock{
				followCloakResult: tt.resolvedURL,
				followCloakError:  tt.followError,
				fetchTrackResult:  tt.fetchResult,
				fetchTrackError:   tt.fetchError,
			}
			adapter := newDeezerAdapter(client)

			got, err := adapter.FetchCloak(context.Background(), tt.cloakCode)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeezerAdapter_FetchAlbum(t *testing.T) {
	client := &deezerClientMock{
		fetchAlbumResult: &deezer.Album{
			ID:    789,
			Title: "Test Album",
			Artist: deezer.Artist{
				Name: "Test Artist",
			},
		},
	}
	adapter := newDeezerAdapter(client)

	got, err := adapter.FetchAlbum(context.Background(), "789")

	require.NoError(t, err)
	require.Equal(t, &Entity{
		ID:       "789",
		Title:    "Test Album",
		Artist:   "Test Artist",
		URL:      "https://deezer.com/album/789",
		Provider: Deezer,
		Type:     Album,
	}, got)
}
