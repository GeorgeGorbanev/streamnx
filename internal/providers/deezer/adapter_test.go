package deezer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestDeezerAdapter_ParseLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType release.Type
		wantID   string
		wantOK   bool
	}{
		{
			name:     "valid localized track URL with www",
			input:    "https://www.deezer.com/us/track/441481432",
			wantType: release.TypeTrack,
			wantID:   "441481432",
			wantOK:   true,
		},
		{
			name:     "valid localized track URL",
			input:    "https://deezer.com/us/track/441481432",
			wantType: release.TypeTrack,
			wantID:   "441481432",
			wantOK:   true,
		},
		{
			name:     "valid track URL without locale",
			input:    "https://deezer.com/track/441481432",
			wantType: release.TypeTrack,
			wantID:   "441481432",
			wantOK:   true,
		},
		{
			name:     "valid track URL with query",
			input:    "https://www.deezer.com/track/3129746?host=6024526261",
			wantType: release.TypeTrack,
			wantID:   "3129746",
			wantOK:   true,
		},
		{
			name:     "valid localized album URL with www",
			input:    "https://www.deezer.com/us/album/53564982",
			wantType: release.TypeAlbum,
			wantID:   "53564982",
			wantOK:   true,
		},
		{
			name:     "valid localized album URL",
			input:    "https://deezer.com/us/album/53564982",
			wantType: release.TypeAlbum,
			wantID:   "53564982",
			wantOK:   true,
		},
		{
			name:     "valid album URL without locale",
			input:    "https://deezer.com/album/53564982",
			wantType: release.TypeAlbum,
			wantID:   "53564982",
			wantOK:   true,
		},
		{
			name:     "valid album URL with query",
			input:    "https://www.deezer.com/album/3129746?host=6024526261",
			wantType: release.TypeAlbum,
			wantID:   "3129746",
			wantOK:   true,
		},
		{
			name:     "valid cloak URL",
			input:    "https://link.deezer.com/s/30FbcgrctxIrNQImDnVEZ",
			wantType: release.TypeCloak,
			wantID:   "30FbcgrctxIrNQImDnVEZ",
			wantOK:   true,
		},
		{
			name:     "valid short cloak URL",
			input:    "https://link.deezer.com/s/jHWZaLoRJutY3TiVA",
			wantType: release.TypeCloak,
			wantID:   "jHWZaLoRJutY3TiVA",
			wantOK:   true,
		},
		{
			name:   "not a deezer link",
			input:  "not a deezer link",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter(&clientMock{})

			gotType, gotID, gotOK := a.ParseLink(tt.input)

			require.Equal(t, tt.wantOK, gotOK)
			require.Equal(t, tt.wantType, gotType)
			require.Equal(t, tt.wantID, gotID)
		})
	}
}

func TestDeezerAdapter_fetchTrack(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockClient func(m *clientMock)
		want       release.Track
		wantError  error
	}{
		{
			name: "successful fetch",
			id:   "123456",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "123456").
					Return(track{
						ID:    123456,
						ISRC:  "GBARL9300135",
						Title: "Test Song",
						Artist: artist{
							Name: "Test Artist",
						},
					}, nil).
					Once()
			},
			want: release.Track{
				ID:       "123456",
				ISRC:     "GBARL9300135",
				Title:    "Test Song",
				Artist:   "Test Artist",
				URL:      "https://deezer.com/track/123456",
				Provider: release.Deezer,
			},
		},
		{
			name: "not found error",
			id:   "404",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "404").
					Return(nil, errNotFound).
					Once()
			},
			wantError: release.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			adapter := NewAdapter(cm)

			got, err := adapter.FetchTrack(t.Context(), tt.id)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Zero(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestDeezerAdapter_Uncloak(t *testing.T) {
	tests := []struct {
		name         string
		cloakCode    string
		mockClient   func(m *clientMock)
		wantType     release.Type
		wantID       string
		wantErrorMsg string
	}{
		{
			name:      "successful track resolve",
			cloakCode: "abc123",
			mockClient: func(m *clientMock) {
				m.
					On("followCloak", "abc123").
					Return("https://link.deezer.com/?awf=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1&dest=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1&gwf=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1&iwf=https%3A%2F%2Fwww.deezer.com%2Ftrack%2F3129746%3Fhost%3D6024526261%26utm_campaign%3Dclipboard-generic%26utm_source%3Duser_sharing%26utm_content%3Dtrack-3129746%26deferredFl%3D1%26universal_link%3D1", nil).
					Once()
			},
			wantType: release.TypeTrack,
			wantID:   "3129746",
		},
		{
			name:      "successful album resolve",
			cloakCode: "album123",
			mockClient: func(m *clientMock) {
				m.
					On("followCloak", "album123").
					Return("https://link.com/?dest=https%3A%2F%2Fwww.deezer.com%2Falbum%2F654321", nil).
					Once()
			},
			wantType: release.TypeAlbum,
			wantID:   "654321",
		},
		{
			name:      "follow cloak error",
			cloakCode: "error123",
			mockClient: func(m *clientMock) {
				m.
					On("followCloak", "error123").
					Return("", errors.New("network error")).
					Once()
			},
			wantErrorMsg: "failed to follow cloak link",
		},
		{
			name:      "invalid cloak URL",
			cloakCode: "invalid-url",
			mockClient: func(m *clientMock) {
				m.
					On("followCloak", "invalid-url").
					Return("not a url", nil).
					Once()
			},
			wantErrorMsg: "failed to find cloak 'dest' param in url: not a url",
		},
		{
			name:      "missing destination",
			cloakCode: "missing-dest",
			mockClient: func(m *clientMock) {
				m.
					On("followCloak", "missing-dest").
					Return("https://link.deezer.com?no_dest=true", nil).
					Once()
			},
			wantErrorMsg: "failed to find cloak 'dest' param in url: https://link.deezer.com?no_dest=true",
		},
		{
			name:      "invalid resolved URL",
			cloakCode: "invalid123",
			mockClient: func(m *clientMock) {
				m.
					On("followCloak", "invalid123").
					Return("https://link.com/?dest=https%3A%2F%2Fwww.deezer.com%2Fartist%2F123456", nil).
					Once()
			},
			wantErrorMsg: "release not found: cloak dest is not a track or album (https://www.deezer.com/artist/123456)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			adapter := NewAdapter(cm)

			gotType, gotID, err := adapter.Uncloak(t.Context(), tt.cloakCode)

			if tt.wantErrorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrorMsg)
				require.Empty(t, gotType)
				require.Empty(t, gotID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantType, gotType)
				require.Equal(t, tt.wantID, gotID)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestDeezerAdapter_fetchAlbum(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockClient func(m *clientMock)
		want       release.Album
		wantError  error
	}{
		{
			name: "successful fetch",
			id:   "789",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "789").
					Return(album{
						ID:    789,
						Title: "Test Album",
						Artist: artist{
							Name: "Test Artist",
						},
						Label: "Test Label",
						Tracks: trackData{
							Data: []track{
								{
									ID:    123,
									Title: "First Track",
									Artist: artist{
										Name: "First Artist",
									},
								},
								{
									ID:    456,
									Title: "Second Track",
									Artist: artist{
										Name: "Second Artist",
									},
								},
							},
						},
					}, nil).
					Once()
			},
			want: release.Album{
				ID:       "789",
				Title:    "Test Album",
				Artist:   "Test Artist",
				Label:    "Test Label",
				URL:      "https://deezer.com/album/789",
				Provider: release.Deezer,
				Creator:  "Test Label",
				TrackIDs: []string{"123", "456"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			adapter := NewAdapter(cm)

			got, err := adapter.FetchAlbum(t.Context(), tt.id)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Zero(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestDeezerAdapter_searchTracks(t *testing.T) {
	tests := []struct {
		name       string
		artist     string
		title      string
		mockClient func(m *clientMock)
		want       []release.SearchTrack
		wantError  error
	}{
		{
			name:   "found candidates",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist", "sample track").
					Return([]track{
						{
							ID:    123,
							ISRC:  "GBARL9300135",
							Title: "First Track",
							Artist: artist{
								Name: "First Artist",
							},
						},
						{
							ID:    456,
							Title: "Second Track",
							Artist: artist{
								Name: "Second Artist",
							},
						},
					}, nil).
					Once()
			},
			want: []release.SearchTrack{
				{
					ID:       "123",
					ISRC:     "GBARL9300135",
					Title:    "First Track",
					Artist:   "First Artist",
					URL:      "https://deezer.com/track/123",
					Provider: release.Deezer,
				},
				{
					ID:       "456",
					Title:    "Second Track",
					Artist:   "Second Artist",
					URL:      "https://deezer.com/track/456",
					Provider: release.Deezer,
				},
			},
		},
		{
			name:   "no candidates",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist", "sample track").
					Return([]track{}, nil).
					Once()
			},
			want: []release.SearchTrack{},
		},
		{
			name:   "client error",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist", "sample track").
					Return(nil, errNotFound).
					Once()
			},
			wantError: errNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			adapter := NewAdapter(cm)

			got, err := adapter.SearchTracks(t.Context(), tt.artist, tt.title)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestDeezerAdapter_fetchTracksByISRC(t *testing.T) {
	tests := []struct {
		name       string
		mockClient func(m *clientMock)
		want       []release.Track
	}{
		{
			name: "found",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrackByISRC", "GBARL9300135").
					Return(track{
						ID:          123,
						ISRC:        "GBARL9300135",
						Title:       "Sample Track",
						Artist:      artist{Name: "Sample Artist"},
						Album:       albumInfo{ID: 456, Title: "Sample Album"},
						Duration:    123,
						ReleaseDate: "2024-01-02",
					}, nil).
					Once()
			},
			want: []release.Track{{
				ID:          "123",
				ISRC:        "GBARL9300135",
				Title:       "Sample Track",
				Artist:      "Sample Artist",
				AlbumID:     "456",
				AlbumTitle:  "Sample Album",
				URL:         "https://deezer.com/track/123",
				Duration:    123,
				ReleaseDate: release.Date{Year: 2024, Month: 1, Day: 2},
				Provider:    release.Deezer,
			}},
		},
		{
			name: "not found",
			mockClient: func(m *clientMock) {
				m.On("fetchTrackByISRC", "GBARL9300135").Return(nil, errNotFound).Once()
			},
			want: []release.Track{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			got, err := NewAdapter(cm).FetchTracksByISRC(t.Context(), "GBARL9300135")

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			cm.AssertExpectations(t)
		})
	}
}

func TestDeezerAdapter_searchAlbums(t *testing.T) {
	tests := []struct {
		name       string
		artist     string
		title      string
		mockClient func(m *clientMock)
		want       []release.SearchAlbum
		wantError  error
	}{
		{
			name:   "found candidates",
			artist: "sample artist",
			title:  "sample album",
			mockClient: func(m *clientMock) {
				m.
					On("searchAlbums", "sample artist", "sample album").
					Return([]album{
						{
							ID:    789,
							Title: "First Album",
							Artist: artist{
								Name: "First Artist",
							},
							Label: "First Label",
						},
						{
							ID:    987,
							Title: "Second Album",
							Artist: artist{
								Name: "Second Artist",
							},
							Label: "Second Label",
						},
					}, nil).
					Once()
			},
			want: []release.SearchAlbum{
				{
					ID:       "789",
					Title:    "First Album",
					Artist:   "First Artist",
					URL:      "https://deezer.com/album/789",
					Provider: release.Deezer,
					Creator:  "First Label",
				},
				{
					ID:       "987",
					Title:    "Second Album",
					Artist:   "Second Artist",
					URL:      "https://deezer.com/album/987",
					Provider: release.Deezer,
					Creator:  "Second Label",
				},
			},
		},
		{
			name:   "no candidates",
			artist: "sample artist",
			title:  "sample album",
			mockClient: func(m *clientMock) {
				m.
					On("searchAlbums", "sample artist", "sample album").
					Return([]album{}, nil).
					Once()
			},
			want: []release.SearchAlbum{},
		},
		{
			name:   "client error",
			artist: "sample artist",
			title:  "sample album",
			mockClient: func(m *clientMock) {
				m.
					On("searchAlbums", "sample artist", "sample album").
					Return(nil, errNotFound).
					Once()
			},
			wantError: errNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			adapter := NewAdapter(cm)

			got, err := adapter.SearchAlbums(t.Context(), tt.artist, tt.title)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}

			cm.AssertExpectations(t)
		})
	}
}

type clientMock struct {
	mock.Mock
}

func (m *clientMock) fetchTrack(_ context.Context, id string) (track, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return track{}, args.Error(1)
	}
	return args.Get(0).(track), args.Error(1)
}

func (m *clientMock) searchTracks(_ context.Context, artist, title string) ([]track, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]track), args.Error(1)
}

func (m *clientMock) fetchTrackByISRC(_ context.Context, isrc string) (track, error) {
	args := m.Called(isrc)
	if args.Get(0) == nil {
		return track{}, args.Error(1)
	}
	return args.Get(0).(track), args.Error(1)
}

func (m *clientMock) fetchAlbum(_ context.Context, id string) (album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return album{}, args.Error(1)
	}
	return args.Get(0).(album), args.Error(1)
}

func (m *clientMock) searchAlbums(_ context.Context, artist, title string) ([]album, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]album), args.Error(1)
}

func (m *clientMock) followCloak(_ context.Context, cloakCode string) (string, error) {
	args := m.Called(cloakCode)
	return args.String(0), args.Error(1)
}
