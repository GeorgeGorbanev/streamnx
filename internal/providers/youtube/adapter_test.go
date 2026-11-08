package youtube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/providers/youtube/client"
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
					On("FetchVideo", "sampleID").
					Return(client.Video{
						Metadata: client.Metadata{
							ID:          "sampleID",
							Title:       "sample artist – sample track (official client.Video)",
							Creator:     "sample channel",
							Description: "sample raw description",
						},
						Duration: 213,
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:          "sampleID",
				Title:       "sample artist – sample track (official client.Video)",
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
					On("FetchVideo", "sampleID").
					Return(client.Video{
						Metadata: client.Metadata{
							ID:          "sampleID",
							Title:       "track name (remastered)",
							Description: "raw provider generated description",
							Creator:     "sample artist topic channel",
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
					On("FetchVideo", "notFoundID").
					Return(nil, client.ErrNotFound).
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
	tests := []struct {
		name           string
		mockClient     func(*clientMock)
		expectedTracks []release.SearchTrack
	}{
		{
			name: "hydrates every available search result in order",
			mockClient: func(m *clientMock) {
				m.
					On("SearchVideos", "sample artist – sample track").
					Return([]string{"firstID", "secondID"}, nil).
					Once()
				m.
					On("FetchVideo", "firstID").
					Return(client.Video{
						Metadata: client.Metadata{
							ID:          "firstID",
							Title:       "first raw title [official]",
							Creator:     "first channel",
							Description: "first raw description",
						},
					}, nil).
					Once()
				m.
					On("FetchVideo", "secondID").
					Return(client.Video{
						Metadata: client.Metadata{
							ID:          "secondID",
							Title:       "second raw title (live)",
							Creator:     "second channel",
							Description: "second raw description",
						},
					}, nil).
					Once()
			},
			expectedTracks: []release.SearchTrack{
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
			},
		},
		{
			name: "skips unavailable videos",
			mockClient: func(m *clientMock) {
				m.
					On("SearchVideos", "sample artist – sample track").
					Return([]string{"unavailableID", "availableID"}, nil).
					Once()
				m.
					On("FetchVideo", "unavailableID").
					Return(nil, client.ErrNotFound).
					Once()
				m.
					On("FetchVideo", "availableID").
					Return(client.Video{
						Metadata: client.Metadata{
							ID:      "availableID",
							Title:   "available track",
							Creator: "available channel",
						},
					}, nil).
					Once()
			},
			expectedTracks: []release.SearchTrack{{
				ID:       "availableID",
				Title:    "available track",
				Provider: release.Youtube,
				Creator:  "available channel",
				URL:      "https://www.youtube.com/watch?v=availableID",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			result, err := NewAdapter(cm).SearchTracks(t.Context(), "sample artist", "sample track")

			require.NoError(t, err)
			require.Equal(t, tt.expectedTracks, result)
			cm.AssertExpectations(t)
		})
	}
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
					On("FetchPlaylistWithItems", "sampleID").
					Return(client.Metadata{
						ID:          "sampleID",
						Title:       "sample artist – sample album [full playlist]",
						Creator:     "playlist channel",
						Description: "playlist raw description",
					}, []string{"firstVideoID", "secondVideoID"}, nil).
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
					On("FetchPlaylistWithItems", "sampleID").
					Return(client.Metadata{
						ID:          "sampleID",
						Title:       "Album - sample name",
						Creator:     "YouTube",
						Description: "playlist raw description",
					}, []string{"sampleTrackID"}, nil).
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
			name: "empty playlist item IDs are skipped",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("FetchPlaylistWithItems", "sampleID").
					Return(client.Metadata{
						ID:          "sampleID",
						Title:       "Album - sample name",
						Creator:     "YouTube",
						Description: "playlist raw description",
					}, []string{"", "sampleTrackID"}, nil).
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
					On("FetchPlaylistWithItems", "notFoundID").
					Return(nil, nil, client.ErrNotFound).
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
	tests := []struct {
		name           string
		mockClient     func(*clientMock)
		expectedAlbums []release.SearchAlbum
	}{
		{
			name: "hydrates every available search result in order",
			mockClient: func(m *clientMock) {
				m.
					On("SearchPlaylists", "sample artist – sample album").
					Return([]string{"firstID", "secondID"}, nil).
					Once()
				m.
					On("FetchPlaylist", "firstID").
					Return(client.Metadata{
						ID:          "firstID",
						Title:       "first raw playlist [full album]",
						Creator:     "first playlist channel",
						Description: "first playlist raw description",
					}, nil).
					Once()
				m.
					On("FetchPlaylist", "secondID").
					Return(client.Metadata{
						ID:          "secondID",
						Title:       "Album - second raw playlist",
						Creator:     "YouTube",
						Description: "second playlist raw description",
					}, nil).
					Once()
			},
			expectedAlbums: []release.SearchAlbum{
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
			},
		},
		{
			name: "skips unavailable playlists",
			mockClient: func(m *clientMock) {
				m.
					On("SearchPlaylists", "sample artist – sample album").
					Return([]string{"unavailableID", "availableID"}, nil).
					Once()
				m.
					On("FetchPlaylist", "unavailableID").
					Return(nil, client.ErrNotFound).
					Once()
				m.
					On("FetchPlaylist", "availableID").
					Return(client.Metadata{
						ID:      "availableID",
						Title:   "available album",
						Creator: "available channel",
					}, nil).
					Once()
			},
			expectedAlbums: []release.SearchAlbum{{
				ID:       "availableID",
				Title:    "available album",
				Provider: release.Youtube,
				Creator:  "available channel",
				URL:      "https://www.youtube.com/playlist?list=availableID",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			result, err := NewAdapter(cm).SearchAlbums(t.Context(), "sample artist", "sample album")

			require.NoError(t, err)
			require.Equal(t, tt.expectedAlbums, result)
			cm.AssertExpectations(t)
		})
	}
}

func TestYoutubeAdapter_fetchTracksByISRC(t *testing.T) {
	tracks, err := (&Adapter{}).FetchTracksByISRC(t.Context(), "GBARL9300135")

	require.Nil(t, tracks)
	require.ErrorIs(t, err, release.ErrUnsupportedOperation)

	albums, err := (&Adapter{}).FetchAlbumsByUPC(t.Context(), "196006422677")

	require.Nil(t, albums)
	require.ErrorIs(t, err, release.ErrUnsupportedOperation)
}

type clientMock struct {
	mock.Mock
}

func (m *clientMock) FetchVideo(_ context.Context, id string) (client.Video, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return client.Video{}, args.Error(1)
	}
	return args.Get(0).(client.Video), args.Error(1)
}

func (m *clientMock) SearchVideos(_ context.Context, term string) ([]string, error) {
	args := m.Called(term)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *clientMock) FetchPlaylist(_ context.Context, id string) (client.Metadata, error) {
	args := m.Called(id)
	metadata, _ := args.Get(0).(client.Metadata)
	return metadata, args.Error(1)
}

func (m *clientMock) FetchPlaylistWithItems(_ context.Context, id string) (client.Metadata, []string, error) {
	args := m.Called(id)
	metadata, _ := args.Get(0).(client.Metadata)
	items, _ := args.Get(1).([]string)
	return metadata, items, args.Error(2)
}

func (m *clientMock) SearchPlaylists(_ context.Context, term string) ([]string, error) {
	args := m.Called(term)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
