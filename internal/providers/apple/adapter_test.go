package apple

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestAppleAdapter_ParseLink(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType release.Type
		expectedID   string
		expectedOK   bool
	}{
		{
			name:         "valid URL with track ID",
			input:        "https://music.apple.com/us/album/song-name/1234567890?i=987654321",
			expectedType: release.TypeTrack,
			expectedID:   "us-987654321",
			expectedOK:   true,
		},
		{
			name:         "valid URL with track ID and th storefront",
			input:        "https://music.apple.com/th/album/song-name/1234567890?i=987654321",
			expectedType: release.TypeTrack,
			expectedID:   "th-987654321",
			expectedOK:   true,
		},
		{
			name:         "valid URL with track ID and invalid iso3611 storefront",
			input:        "https://music.apple.com/invalidstorefront/album/song-name/1234567890?i=987654321",
			expectedType: "",
			expectedID:   "",
			expectedOK:   false,
		},
		{
			name:         "valid URL with track ID and unsupported storefront",
			input:        "https://music.apple.com/zz/album/song-name/1234567890?i=987654321",
			expectedType: "",
			expectedID:   "",
			expectedOK:   false,
		},
		{
			name:         "valid URL without album",
			input:        "https://music.apple.com/us/song/angel/724466660",
			expectedType: release.TypeTrack,
			expectedID:   "us-724466660",
			expectedOK:   true,
		},
		{
			name:         "valid URL without album and invalid storefront",
			input:        "https://music.apple.com/invalidstorefront/song/angel/724466660",
			expectedType: "",
			expectedID:   "",
			expectedOK:   false,
		},
		{
			name:         "valid URL with album ID",
			input:        "https://music.apple.com/us/album/album-name/123456789",
			expectedType: release.TypeAlbum,
			expectedID:   "us-123456789",
			expectedOK:   true,
		},
		{
			name:         "valid URL with album ID and gb locale",
			input:        "https://music.apple.com/gb/album/another-album/987654321",
			expectedType: release.TypeAlbum,
			expectedID:   "gb-987654321",
			expectedOK:   true,
		},
		{
			name:         "valid URL with album ID and invalid iso3611 storefront",
			input:        "https://music.apple.com/invalidstorefront/album/another-album/987654321",
			expectedType: "",
			expectedID:   "",
			expectedOK:   false,
		},
		{
			name:         "URL without album ID",
			input:        "https://music.apple.com/us/album/album-name",
			expectedType: "",
			expectedID:   "",
			expectedOK:   false,
		},

		{
			name:         "invalid host URL",
			input:        "https://music.orange.com/us/album/song-name/1234567890?i=987654321",
			expectedType: "",
			expectedID:   "",
			expectedOK:   false,
		},
		{
			name:         "empty string",
			input:        "",
			expectedType: "",
			expectedID:   "",
			expectedOK:   false,
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

func TestAppleAdapter_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedTrack release.Track
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "ru-123",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "123", "ru").
					Return(entity{
						ID: "ru-123",
						Attributes: entityAttributes{
							ArtistName: "sample artist",
							Name:       "sample name",
							URL:        "https://music.apple.com/ru/album/song-name/1234567890?i=123",
							EditorialNotes: editorialNotes{
								Standard: "sample track editorial notes",
							},
						},
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:          "ru-123",
				Title:       "sample name",
				Artist:      "sample artist",
				URL:         "https://music.apple.com/ru/album/song-name/1234567890?i=123",
				Provider:    release.Apple,
				Description: "sample track editorial notes",
			},
		},
		{
			name: "not found ID",
			id:   "ru-123",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "123", "ru").
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

func TestAppleAdapter_fetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		storefront    string
		mockClient    func(m *clientMock)
		expectedAlbum release.Album
		expectedErr   error
	}{
		{
			name:       "found ID",
			id:         "ru-456",
			storefront: "sampleStorefront",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "456", "ru").
					Return(entity{
						ID: "ru-456",
						Attributes: entityAttributes{
							ArtistName:  "sample artist",
							Name:        "sample name",
							RecordLabel: "sample label",
							URL:         "https://music.apple.com/ru/album/name/456",
							EditorialNotes: editorialNotes{
								Standard: "sample album editorial notes",
							},
						},
						Relationships: entityRelationships{
							Tracks: tracksRelationship{
								Data: []entity{
									{
										Attributes: entityAttributes{
											ArtistName: "first track artist",
											Name:       "first track",
											URL:        "https://music.apple.com/ru/album/name/456?i=111",
										},
									},
									{
										Attributes: entityAttributes{
											ArtistName: "second track artist",
											Name:       "second track",
											URL:        "https://music.apple.com/ru/song/second-track/222",
											EditorialNotes: editorialNotes{
												Short: "second track notes",
											},
										},
									},
								},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:          "ru-456",
				Title:       "sample name",
				Artist:      "sample artist",
				Label:       "sample label",
				URL:         "https://music.apple.com/ru/album/name/456",
				Provider:    release.Apple,
				Description: "sample album editorial notes",
				TrackIDs:    []string{"ru-111", "ru-222"},
			},
		},
		{
			name:       "not found ID",
			id:         "ru-456",
			storefront: "notFoundStorefront",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "456", "ru").
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

func TestAppleAdapter_searchTracks(t *testing.T) {
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
					Return([]entity{
						{
							Attributes: entityAttributes{
								ArtistName: "first artist",
								Name:       "first track",
								URL:        "https://music.apple.com/us/album/first-album/111?i=222",
								EditorialNotes: editorialNotes{
									Short: "first track short notes",
								},
							},
						},
						{
							Attributes: entityAttributes{
								ArtistName: "second artist",
								Name:       "second track",
								URL:        "https://music.apple.com/gb/song/second-track/333",
								EditorialNotes: editorialNotes{
									Standard: "second track editorial notes",
								},
							},
						},
					}, nil).
					Once()
			},
			expectedTracks: []release.SearchTrack{
				{
					ID:          "us-222",
					Title:       "first track",
					Artist:      "first artist",
					URL:         "https://music.apple.com/us/album/first-album/111?i=222",
					Provider:    release.Apple,
					Description: "first track short notes",
				},
				{
					ID:          "gb-333",
					Title:       "second track",
					Artist:      "second artist",
					URL:         "https://music.apple.com/gb/song/second-track/333",
					Provider:    release.Apple,
					Description: "second track editorial notes",
				},
			},
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

func TestAppleAdapter_searchAlbums(t *testing.T) {
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
					Return([]entity{
						{
							Attributes: entityAttributes{
								ArtistName: "first artist",
								Name:       "first album",
								URL:        "https://music.apple.com/us/album/first-album/444",
								EditorialNotes: editorialNotes{
									Short: "first album short notes",
								},
							},
						},
						{
							Attributes: entityAttributes{
								ArtistName: "second artist",
								Name:       "second album",
								URL:        "https://music.apple.com/gb/album/second-album/555",
								EditorialNotes: editorialNotes{
									Standard: "second album editorial notes",
								},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbums: []release.SearchAlbum{
				{
					ID:          "us-444",
					Title:       "first album",
					Artist:      "first artist",
					URL:         "https://music.apple.com/us/album/first-album/444",
					Provider:    release.Apple,
					Description: "first album short notes",
				},
				{
					ID:          "gb-555",
					Title:       "second album",
					Artist:      "second artist",
					URL:         "https://music.apple.com/gb/album/second-album/555",
					Provider:    release.Apple,
					Description: "second album editorial notes",
				},
			},
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

func (m *clientMock) fetchTrack(_ context.Context, id, storefront string) (entity, error) {
	args := m.Called(id, storefront)
	if args.Get(0) == nil {
		return entity{}, args.Error(1)
	}
	return args.Get(0).(entity), args.Error(1)
}

func (m *clientMock) searchTracks(_ context.Context, artist, title string) ([]entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity), args.Error(1)
}

func (m *clientMock) fetchAlbum(_ context.Context, id, storefront string) (entity, error) {
	args := m.Called(id, storefront)
	if args.Get(0) == nil {
		return entity{}, args.Error(1)
	}
	return args.Get(0).(entity), args.Error(1)
}

func (m *clientMock) searchAlbums(_ context.Context, artist, title string) ([]entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity), args.Error(1)
}
