package youtube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestYoutubeAdapter_ParseLink(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType release.Type
		expectedID   string
		expectedOK   bool
	}{
		{
			name:         "short track URL",
			input:        "https://youtu.be/dQw4w9WgXcQ",
			expectedType: release.TypeTrack,
			expectedID:   "dQw4w9WgXcQ",
			expectedOK:   true,
		},
		{
			name:         "long track URL",
			input:        "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expectedType: release.TypeTrack,
			expectedID:   "dQw4w9WgXcQ",
			expectedOK:   true,
		},
		{
			name:         "track URL without scheme",
			input:        "youtube.com/watch?v=dQw4w9WgXcQ",
			expectedType: release.TypeTrack,
			expectedID:   "dQw4w9WgXcQ",
			expectedOK:   true,
		},
		{
			name:         "track URL with extra parameters",
			input:        "https://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=youtu.be",
			expectedType: release.TypeTrack,
			expectedID:   "dQw4w9WgXcQ",
			expectedOK:   true,
		},
		{
			name:       "youtube music track URL belongs to another provider",
			input:      "https://music.youtube.com/watch?v=5PgdZDXg0z0&si=LkthPMI6H_I04dhP",
			expectedOK: false,
		},
		{
			name:         "standard album URL",
			input:        "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expectedType: release.TypeAlbum,
			expectedID:   "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expectedOK:   true,
		},
		{
			name:         "shortened album URL",
			input:        "https://youtu.be/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expectedType: release.TypeAlbum,
			expectedID:   "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expectedOK:   true,
		},
		{
			name:         "album URL with extra parameters",
			input:        "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj&feature=share",
			expectedType: release.TypeAlbum,
			expectedID:   "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expectedOK:   true,
		},
		{
			name:       "youtube music album URL belongs to another provider",
			input:      "https://music.youtube.com/playlist?list=OLAK5uy_n4xauusTJSj6Mtt4cIuq4KZziSfjABYWU",
			expectedOK: false,
		},
		{
			name:       "track URL without ID",
			input:      "https://www.youtube.com/watch?v=",
			expectedOK: false,
		},
		{
			name:       "non-youtube track URL",
			input:      "https://www.example.com/watch?v=dQw4w9WgXcQ",
			expectedOK: false,
		},
		{
			name:       "short track ID",
			input:      "https://www.youtube.com/watch?v=notFound",
			expectedOK: false,
		},
		{
			name:       "non-youtube album URL",
			input:      "https://www.example.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expectedOK: false,
		},
		{
			name:       "empty string",
			input:      "",
			expectedOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter(&clientMock{})

			linkType, id, ok := a.ParseLink(tt.input)

			require.Equal(t, tt.expectedOK, ok)
			require.Equal(t, tt.expectedType, linkType)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestYoutubeAdapter_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedTrack release.Track
		expectedErr   error
	}{
		{
			name: "found ID returns raw video fields",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchVideo", "sampleID").
					Return(video{
						ID: "sampleID",
						Snippet: snippet{
							Title:        "sample artist – sample track (official video)",
							ChannelTitle: "sample channel",
							Description:  "sample raw description",
						},
						ContentDetails: contentDetails{
							Duration: "PT3M33S",
						},
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:          "sampleID",
				Title:       "sample artist – sample track (official video)",
				Artist:      "",
				URL:         "https://www.youtube.com/watch?v=sampleID",
				Provider:    release.Youtube,
				Creator:     "sample channel",
				Description: "sample raw description",
				Duration:    213,
			},
		},
		{
			name: "found ID preserves raw topic-style channel data",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchVideo", "sampleID").
					Return(video{
						ID: "sampleID",
						Snippet: snippet{
							Title:        "track name (remastered)",
							Description:  "raw provider generated description",
							ChannelTitle: "sample artist topic channel",
						},
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:          "sampleID",
				Title:       "track name (remastered)",
				Artist:      "",
				URL:         "https://www.youtube.com/watch?v=sampleID",
				Provider:    release.Youtube,
				Creator:     "sample artist topic channel",
				Description: "raw provider generated description",
			},
		},
		{
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchVideo", "notFoundID").
					Return(nil, errNotFound).
					Once()
			},
			expectedErr: release.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)

			result, err := a.FetchTrack(t.Context(), tt.id)

			if tt.expectedErr != nil {
				require.Zero(t, result)
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestYoutubeAdapter_searchTracks(t *testing.T) {
	cm := &clientMock{}
	cm.
		On("searchVideos", "sample artist – sample track").
		Return([]videoSearchResult{
			{ID: videoSearchID{VideoID: "firstID"}},
			{ID: videoSearchID{VideoID: "secondID"}},
		}, nil).
		Once()
	cm.
		On("fetchVideo", "firstID").
		Return(video{
			ID: "firstID",
			Snippet: snippet{
				Title:        "first raw title [official]",
				ChannelTitle: "first channel",
				Description:  "first raw description",
			},
		}, nil).
		Once()
	cm.
		On("fetchVideo", "secondID").
		Return(video{
			ID: "secondID",
			Snippet: snippet{
				Title:        "second raw title (live)",
				ChannelTitle: "second channel",
				Description:  "second raw description",
			},
		}, nil).
		Once()

	a := NewAdapter(cm)

	result, err := a.SearchTracks(t.Context(), "sample artist", "sample track")

	require.NoError(t, err)
	require.Equal(t, []release.SearchTrack{
		{
			ID:          "firstID",
			Title:       "first raw title [official]",
			Artist:      "",
			URL:         "https://www.youtube.com/watch?v=firstID",
			Provider:    release.Youtube,
			Creator:     "first channel",
			Description: "first raw description",
		},
		{
			ID:          "secondID",
			Title:       "second raw title (live)",
			Artist:      "",
			URL:         "https://www.youtube.com/watch?v=secondID",
			Provider:    release.Youtube,
			Creator:     "second channel",
			Description: "second raw description",
		},
	}, result)
	cm.AssertExpectations(t)
}

func TestYoutubeAdapter_fetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedAlbum release.Album
		expectedErr   error
	}{
		{
			name: "found ID returns raw playlist fields",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchPlaylist", "sampleID").
					Return(playlist{
						ID: "sampleID",
						Snippet: snippet{
							Title:        "sample artist – sample album [full playlist]",
							ChannelTitle: "playlist channel",
							Description:  "playlist raw description",
						},
					}, nil).
					Once()
				m.
					On("fetchPlaylistItems", "sampleID").
					Return([]playlistItem{
						{
							Snippet: snippet{
								ResourceID: resourceID{VideoID: "firstVideoID"},
							},
						},
						{
							Snippet: snippet{
								ResourceID: resourceID{VideoID: "secondVideoID"},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:          "sampleID",
				Title:       "sample artist – sample album [full playlist]",
				Artist:      "",
				URL:         "https://www.youtube.com/playlist?list=sampleID",
				Provider:    release.Youtube,
				Creator:     "playlist channel",
				Description: "playlist raw description",
				TrackIDs:    []string{"firstVideoID", "secondVideoID"},
			},
		},
		{
			name: "generic playlist owner keeps raw creator and exposes track ids",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchPlaylist", "sampleID").
					Return(playlist{
						ID: "sampleID",
						Snippet: snippet{
							Title:        "Album - sample name",
							ChannelTitle: "YouTube",
							Description:  "playlist raw description",
						},
					}, nil).
					Once()
				m.
					On("fetchPlaylistItems", "sampleID").
					Return([]playlistItem{
						{
							ID: "sampleTrackID",
							Snippet: snippet{
								Title:                  "first item raw title",
								VideoOwnerChannelTitle: "sample artist - Topic",
								Description:            "first item raw description",
								ResourceID:             resourceID{VideoID: "sampleTrackID"},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:          "sampleID",
				Title:       "Album - sample name",
				Artist:      "",
				URL:         "https://www.youtube.com/playlist?list=sampleID",
				Provider:    release.Youtube,
				Creator:     "YouTube",
				Description: "playlist raw description",
				TrackIDs:    []string{"sampleTrackID"},
			},
		},
		{
			name: "generic playlist owner ignores item channel title",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchPlaylist", "sampleID").
					Return(playlist{
						ID: "sampleID",
						Snippet: snippet{
							Title:        "Album - sample name",
							ChannelTitle: "YouTube",
							Description:  "playlist raw description",
						},
					}, nil).
					Once()
				m.
					On("fetchPlaylistItems", "sampleID").
					Return([]playlistItem{
						{
							ID: "sampleTrackID",
							Snippet: snippet{
								Title:        "first item raw title",
								ChannelTitle: "first item channel",
								Description:  "first item raw description",
								ResourceID:   resourceID{VideoID: "sampleTrackID"},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:          "sampleID",
				Title:       "Album - sample name",
				Artist:      "",
				URL:         "https://www.youtube.com/playlist?list=sampleID",
				Provider:    release.Youtube,
				Creator:     "YouTube",
				Description: "playlist raw description",
				TrackIDs:    []string{"sampleTrackID"},
			},
		},
		{
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchPlaylist", "notFoundID").
					Return(nil, errNotFound).
					Once()
			},
			expectedErr: release.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)

			result, err := a.FetchAlbum(t.Context(), tt.id)

			if tt.expectedErr != nil {
				require.Zero(t, result)
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestYoutubeAdapter_searchAlbums(t *testing.T) {
	cm := &clientMock{}
	cm.
		On("searchPlaylists", "sample artist – sample album").
		Return([]playlistSearchResult{
			{ID: playlistSearchID{PlaylistID: "firstID"}},
			{ID: playlistSearchID{PlaylistID: "secondID"}},
		}, nil).
		Once()
	cm.
		On("fetchPlaylist", "firstID").
		Return(playlist{
			ID: "firstID",
			Snippet: snippet{
				Title:        "first raw playlist [full album]",
				ChannelTitle: "first playlist channel",
				Description:  "first playlist raw description",
			},
		}, nil).
		Once()
	cm.
		On("fetchPlaylist", "secondID").
		Return(playlist{
			ID: "secondID",
			Snippet: snippet{
				Title:        "Album - second raw playlist",
				ChannelTitle: "YouTube",
				Description:  "second playlist raw description",
			},
		}, nil).
		Once()

	a := NewAdapter(cm)

	result, err := a.SearchAlbums(t.Context(), "sample artist", "sample album")

	require.NoError(t, err)
	require.Equal(t, []release.SearchAlbum{
		{
			ID:          "firstID",
			Title:       "first raw playlist [full album]",
			Artist:      "",
			URL:         "https://www.youtube.com/playlist?list=firstID",
			Provider:    release.Youtube,
			Creator:     "first playlist channel",
			Description: "first playlist raw description",
		},
		{
			ID:          "secondID",
			Title:       "Album - second raw playlist",
			Artist:      "",
			URL:         "https://www.youtube.com/playlist?list=secondID",
			Provider:    release.Youtube,
			Creator:     "YouTube",
			Description: "second playlist raw description",
		},
	}, result)
	cm.AssertExpectations(t)
}

type clientMock struct {
	mock.Mock
}

func (m *clientMock) fetchVideo(_ context.Context, id string) (video, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return video{}, args.Error(1)
	}
	return args.Get(0).(video), args.Error(1)
}

func (m *clientMock) searchVideos(_ context.Context, term string) ([]videoSearchResult, error) {
	args := m.Called(term)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]videoSearchResult), args.Error(1)
}

func (m *clientMock) fetchPlaylist(_ context.Context, id string) (playlist, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return playlist{}, args.Error(1)
	}
	return args.Get(0).(playlist), args.Error(1)
}

func (m *clientMock) searchPlaylists(_ context.Context, term string) ([]playlistSearchResult, error) {
	args := m.Called(term)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]playlistSearchResult), args.Error(1)
}

func (m *clientMock) fetchPlaylistItems(_ context.Context, id string) ([]playlistItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]playlistItem), args.Error(1)
}
