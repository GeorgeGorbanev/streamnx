package streamnx

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/deezer"
)

func TestDeezerAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockClient func(m *deezer.ClientMock)
		want       *Entity
		wantError  error
	}{
		{
			name: "successful fetch",
			id:   "123456",
			mockClient: func(m *deezer.ClientMock) {
				m.
					On("FetchTrack", "123456").
					Return(&deezer.Track{
						ID:    123456,
						Title: "Test Song",
						Artist: deezer.Artist{
							Name: "Test Artist",
						},
					}, nil).
					Once()
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
			name: "not found error",
			id:   "404",
			mockClient: func(m *deezer.ClientMock) {
				m.
					On("FetchTrack", "404").
					Return(nil, deezer.NotFoundError).
					Once()
			},
			wantError: EntityNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &deezer.ClientMock{}
			tt.mockClient(clientMock)

			adapter := newDeezerAdapter(clientMock)

			got, err := adapter.FetchTrack(context.Background(), tt.id)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestDeezerAdapter_FetchCloak(t *testing.T) {
	tests := []struct {
		name       string
		cloakCode  string
		mockClient func(m *deezer.ClientMock)
		want       *Entity
		wantErr    string
	}{
		{
			name:      "successful track resolution",
			cloakCode: "abc123",
			mockClient: func(m *deezer.ClientMock) {
				m.
					On("FollowCloak", "abc123").
					Return("https://link.deezer.com/?dest=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F123456", nil).
					Once()
				m.
					On("FetchTrack", "123456").
					Return(&deezer.Track{
						ID:    123456,
						Title: "Test Song",
						Artist: deezer.Artist{
							Name: "Test Artist",
						},
					}, nil).
					Once()
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
			name:      "follow cloak error",
			cloakCode: "error123",
			mockClient: func(m *deezer.ClientMock) {
				m.
					On("FollowCloak", "error123").
					Return("", errors.New("network error")).
					Once()
			},
			wantErr: "failed to follow cloak link",
		},
		{
			name:      "invalid resolved URL",
			cloakCode: "invalid123",
			mockClient: func(m *deezer.ClientMock) {
				m.
					On("FollowCloak", "invalid123").
					Return("https://link.deezer.com/?dest=https%3A%2F%2Fwww.deezer.com%2Fartist%2F123456", nil).
					Once()
			},
			wantErr: "entity not found: cloak dest is not a track or album (https://www.deezer.com/artist/123456)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &deezer.ClientMock{}
			tt.mockClient(clientMock)

			adapter := newDeezerAdapter(clientMock)

			got, err := adapter.FetchCloak(context.Background(), tt.cloakCode)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestDeezerAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockClient func(m *deezer.ClientMock)
		want       *Entity
		wantError  error
	}{
		{
			name: "successful fetch",
			id:   "789",
			mockClient: func(m *deezer.ClientMock) {
				m.
					On("FetchAlbum", "789").
					Return(&deezer.Album{
						ID:    789,
						Title: "Test Album",
						Artist: deezer.Artist{
							Name: "Test Artist",
						},
					}, nil).
					Once()
			},
			want: &Entity{
				ID:       "789",
				Title:    "Test Album",
				Artist:   "Test Artist",
				URL:      "https://deezer.com/album/789",
				Provider: Deezer,
				Type:     Album,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &deezer.ClientMock{}
			tt.mockClient(clientMock)

			adapter := newDeezerAdapter(clientMock)

			got, err := adapter.FetchAlbum(context.Background(), tt.id)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			clientMock.AssertExpectations(t)
		})
	}
}
