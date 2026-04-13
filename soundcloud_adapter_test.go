package streamnx

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/internal/soundcloud"
)

func TestSoundcloudAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *soundcloud.ClientMock)
		expectedTrack *Entity
		expectedErr   string
	}{
		{
			name: "found ID",
			id:   "forss:flickermood",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("FetchTrack", "forss", "flickermood").
					Return(&soundcloud.Track{
						Title:     "Flickermood",
						Permalink: "flickermood",
						User: soundcloud.User{
							Username:  "Forss",
							Permalink: "forss",
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "forss:flickermood",
				Title:    "Flickermood",
				Artist:   "Forss",
				URL:      "https://soundcloud.com/forss/flickermood",
				Provider: Soundcloud,
				Type:     Track,
			},
		},
		{
			name: "not found ID",
			id:   "forss:flickermood",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("FetchTrack", "forss", "flickermood").
					Return(nil, soundcloud.ErrNotFound).
					Once()
			},
			expectedErr: EntityNotFoundError.Error(),
		},
		{
			name:        "invalid composite key",
			id:          "invalid",
			mockClient:  func(_ *soundcloud.ClientMock) {},
			expectedErr: "failed to unmarshal track id: invalid composite key: invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &soundcloud.ClientMock{}
			tt.mockClient(clientMock)

			a := newSoundcloudAdapter(clientMock)
			result, err := a.FetchTrack(t.Context(), tt.id)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *soundcloud.ClientMock)
		expectedTrack *Entity
		expectedErr   string
	}{
		{
			name:       "found query",
			artist:     "Forss",
			searchName: "Flickermood",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("SearchTrack", "Forss", "Flickermood").
					Return(&soundcloud.Track{
						Title:     "Flickermood",
						Permalink: "flickermood",
						User: soundcloud.User{
							Username:  "Forss",
							Permalink: "forss",
						},
					}, nil).
					Once()
			},
			expectedTrack: &Entity{
				ID:       "forss:flickermood",
				Title:    "Flickermood",
				Artist:   "Forss",
				URL:      "https://soundcloud.com/forss/flickermood",
				Provider: Soundcloud,
				Type:     Track,
			},
		},
		{
			name:       "not found query",
			artist:     "missing",
			searchName: "track",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("SearchTrack", "missing", "track").
					Return(nil, soundcloud.ErrNotFound).
					Once()
			},
			expectedErr: EntityNotFoundError.Error(),
		},
		{
			name:       "malformed result URL",
			artist:     "Forss",
			searchName: "Flickermood",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("SearchTrack", "Forss", "Flickermood").
					Return(&soundcloud.Track{
						Title:     "Flickermood",
						Permalink: "flickermood",
						User: soundcloud.User{
							Username: "Forss",
						},
					}, nil).
					Once()
			},
			expectedErr: "failed to parse composite key from track url: invalid track url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &soundcloud.ClientMock{}
			tt.mockClient(clientMock)

			a := newSoundcloudAdapter(clientMock)
			result, err := a.SearchTrack(t.Context(), tt.artist, tt.searchName)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *soundcloud.ClientMock)
		expectedAlbum *Entity
		expectedErr   string
	}{
		{
			name: "found ID",
			id:   "forss:soulhack",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("FetchAlbum", "forss", "soulhack").
					Return(&soundcloud.Album{
						Title:     "Soulhack",
						Permalink: "soulhack",
						User: soundcloud.User{
							Username:  "Forss",
							Permalink: "forss",
						},
					}, nil).
					Once()
			},
			expectedAlbum: &Entity{
				ID:       "forss:soulhack",
				Title:    "Soulhack",
				Artist:   "Forss",
				URL:      "https://soundcloud.com/forss/sets/soulhack",
				Provider: Soundcloud,
				Type:     Album,
			},
		},
		{
			name: "not found ID",
			id:   "forss:soulhack",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("FetchAlbum", "forss", "soulhack").
					Return(nil, soundcloud.ErrNotFound).
					Once()
			},
			expectedErr: EntityNotFoundError.Error(),
		},
		{
			name:        "invalid composite key",
			id:          "invalid",
			mockClient:  func(_ *soundcloud.ClientMock) {},
			expectedErr: "failed to unmarshal album id: invalid composite key: invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &soundcloud.ClientMock{}
			tt.mockClient(clientMock)

			a := newSoundcloudAdapter(clientMock)
			result, err := a.FetchAlbum(t.Context(), tt.id)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artist        string
		searchName    string
		mockClient    func(m *soundcloud.ClientMock)
		expectedAlbum *Entity
		expectedErr   string
	}{
		{
			name:       "found query",
			artist:     "Forss",
			searchName: "Soulhack",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("SearchAlbum", "Forss", "Soulhack").
					Return(&soundcloud.Album{
						Title:     "Soulhack",
						Permalink: "soulhack",
						User: soundcloud.User{
							Username:  "Forss",
							Permalink: "forss",
						},
					}, nil).
					Once()
			},
			expectedAlbum: &Entity{
				ID:       "forss:soulhack",
				Title:    "Soulhack",
				Artist:   "Forss",
				URL:      "https://soundcloud.com/forss/sets/soulhack",
				Provider: Soundcloud,
				Type:     Album,
			},
		},
		{
			name:       "not found query",
			artist:     "missing",
			searchName: "album",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("SearchAlbum", "missing", "album").
					Return(nil, soundcloud.ErrNotFound).
					Once()
			},
			expectedErr: EntityNotFoundError.Error(),
		},
		{
			name:       "malformed result URL",
			artist:     "Forss",
			searchName: "Soulhack",
			mockClient: func(m *soundcloud.ClientMock) {
				m.
					On("SearchAlbum", "Forss", "Soulhack").
					Return(&soundcloud.Album{
						Title:     "Soulhack",
						Permalink: "soulhack",
						User: soundcloud.User{
							Username: "Forss",
						},
					}, nil).
					Once()
			},
			expectedErr: "failed to parse composite key from album url: invalid album url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &soundcloud.ClientMock{}
			tt.mockClient(clientMock)

			a := newSoundcloudAdapter(clientMock)
			result, err := a.SearchAlbum(t.Context(), tt.artist, tt.searchName)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_FetchCloak(t *testing.T) {
	a := newSoundcloudAdapter(&soundcloud.ClientMock{})

	result, err := a.FetchCloak(t.Context(), "cloak-code")

	require.ErrorIs(t, err, UnsupportedEntityTypeError)
	require.Nil(t, result)
}
