package spotify

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestSpotifyAdapter_ParseLink(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType release.Type
		expectedID   string
		expectedOK   bool
	}{
		{
			name:         "valid track URL",
			input:        "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			expectedType: release.TypeTrack,
			expectedID:   "7uv632EkfwYhXoqf8rhYrg",
			expectedOK:   true,
		},
		{
			name:         "valid track URL with query",
			input:        "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?test=123",
			expectedType: release.TypeTrack,
			expectedID:   "7uv632EkfwYhXoqf8rhYrg",
			expectedOK:   true,
		},
		{
			name:         "valid track URL with intl path",
			input:        "https://open.spotify.com/intl-pt/track/2xmQMKTjiOdkdGVgqDzezo",
			expectedType: release.TypeTrack,
			expectedID:   "2xmQMKTjiOdkdGVgqDzezo",
			expectedOK:   true,
		},
		{
			name:         "valid track URL with intl path and query",
			input:        "https://open.spotify.com/intl-pt/track/2xmQMKTjiOdkdGVgqDzezo?sample=query",
			expectedType: release.TypeTrack,
			expectedID:   "2xmQMKTjiOdkdGVgqDzezo",
			expectedOK:   true,
		},
		{
			name:         "valid track URL with prefix and suffix",
			input:        "prefix https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?test=123 suffix",
			expectedType: release.TypeTrack,
			expectedID:   "7uv632EkfwYhXoqf8rhYrg",
			expectedOK:   true,
		},
		{
			name:         "valid album URL",
			input:        "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg",
			expectedType: release.TypeAlbum,
			expectedID:   "7uv632EkfwYhXoqf8rhYrg",
			expectedOK:   true,
		},
		{
			name:         "valid album URL with intl path",
			input:        "https://open.spotify.com/intl-pt/album/7uv632EkfwYhXoqf8rhYrg",
			expectedType: release.TypeAlbum,
			expectedID:   "7uv632EkfwYhXoqf8rhYrg",
			expectedOK:   true,
		},
		{
			name:         "valid album URL with query",
			input:        "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg?test=123",
			expectedType: release.TypeAlbum,
			expectedID:   "7uv632EkfwYhXoqf8rhYrg",
			expectedOK:   true,
		},
		{
			name:         "valid album URL with prefix and suffix",
			input:        "prefix https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg?test=123 suffix",
			expectedType: release.TypeAlbum,
			expectedID:   "7uv632EkfwYhXoqf8rhYrg",
			expectedOK:   true,
		},
		{
			name:       "empty URL",
			input:      "",
			expectedOK: false,
		},
		{
			name:       "non-Spotify track URL",
			input:      "https://example.com/track/7uv632EkfwYhXoqf8rhYrg",
			expectedOK: false,
		},
		{
			name:       "non-Spotify album URL",
			input:      "https://example.com/album/7uv632EkfwYhXoqf8rhYrg",
			expectedOK: false,
		},
		{
			name:       "track URL without ID",
			input:      "https://open.spotify.com/track/",
			expectedOK: false,
		},
		{
			name:       "album URL without ID",
			input:      "https://open.spotify.com/album/",
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

func TestSpotifyAdapter_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedTrack release.Track
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "sampleID").
					Return(track{
						ID:    "sampleID",
						Name:  "sample name",
						Album: albumInfo{ReleaseDate: "1987-11"},
						Artists: []artist{
							{
								Name: "sample artist",
							},
						},
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://open.spotify.com/track/sampleID",
				Provider: release.Spotify,
				ReleaseDate: release.Date{
					Year: 1987, Month: 11,
				},
			},
		},
		{
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "notFoundID").
					Return(nil, errNotFound).
					Once()
			},
			expectedErr: release.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := &clientMock{}
			tt.mockClient(clientMock)

			a := NewAdapter(clientMock)
			result, err := a.FetchTrack(t.Context(), tt.id)

			if tt.expectedErr != nil {
				require.Zero(t, result)
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

func TestSpotifyAdapter_fetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedAlbum release.Album
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "sampleID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "sampleID").
					Return(album{
						ID:          "sampleID",
						Name:        "sample name",
						Label:       "sample label",
						ReleaseDate: "1987",
						Artists: []artist{
							{
								Name: "sample artist",
							},
						},
						Tracks: albumTracks{
							Items: []track{
								{
									ID:   "firstTrackID",
									Name: "first track",
									Artists: []artist{
										{
											Name: "first artist",
										},
									},
								},
								{
									ID:   "secondTrackID",
									Name: "second track",
								},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:       "sampleID",
				Title:    "sample name",
				Artist:   "sample artist",
				Label:    "sample label",
				URL:      "https://open.spotify.com/album/sampleID",
				Provider: release.Spotify,
				Creator:  "sample label",
				TrackIDs: []string{"firstTrackID", "secondTrackID"},
				ReleaseDate: release.Date{
					Year: 1987,
				},
			},
		},
		{
			name: "not found ID",
			id:   "notFoundID",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "notFoundID").
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

func TestSpotifyAdapter_searchTracks(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		title          string
		mockClient     func(m *clientMock)
		expectedTracks []release.SearchTrack
		expectedErr    error
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
							ID:   "firstTrackID",
							Name: "first track",
							Artists: []artist{
								{
									Name: "first artist",
								},
							},
						},
						{
							ID:   "secondTrackID",
							Name: "second track",
							Artists: []artist{
								{
									Name: "second artist",
								},
							},
						},
					}, nil).
					Once()
			},
			expectedTracks: []release.SearchTrack{
				{
					ID:       "firstTrackID",
					Title:    "first track",
					Artist:   "first artist",
					URL:      "https://open.spotify.com/track/firstTrackID",
					Provider: release.Spotify,
				},
				{
					ID:       "secondTrackID",
					Title:    "second track",
					Artist:   "second artist",
					URL:      "https://open.spotify.com/track/secondTrackID",
					Provider: release.Spotify,
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
			expectedTracks: []release.SearchTrack{},
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
			expectedErr: errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)
			result, err := a.SearchTracks(t.Context(), tt.artist, tt.title)

			if tt.expectedErr != nil {
				require.Nil(t, result)
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTracks, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestSpotifyAdapter_searchAlbums(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		title          string
		mockClient     func(m *clientMock)
		expectedAlbums []release.SearchAlbum
		expectedErr    error
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
							ID:    "firstAlbumID",
							Name:  "first album",
							Label: "first label",
							Artists: []artist{
								{
									Name: "first artist",
								},
							},
						},
						{
							ID:    "secondAlbumID",
							Name:  "second album",
							Label: "second label",
							Artists: []artist{
								{
									Name: "second artist",
								},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbums: []release.SearchAlbum{
				{
					ID:       "firstAlbumID",
					Title:    "first album",
					Artist:   "first artist",
					URL:      "https://open.spotify.com/album/firstAlbumID",
					Provider: release.Spotify,
					Creator:  "first label",
				},
				{
					ID:       "secondAlbumID",
					Title:    "second album",
					Artist:   "second artist",
					URL:      "https://open.spotify.com/album/secondAlbumID",
					Provider: release.Spotify,
					Creator:  "second label",
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
			expectedAlbums: []release.SearchAlbum{},
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
			expectedErr: errNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)
			result, err := a.SearchAlbums(t.Context(), tt.artist, tt.title)

			if tt.expectedErr != nil {
				require.Nil(t, result)
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbums, result)
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
