package soundcloud

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestSoundcloudAdapter_ParseLink(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType release.Type
		expectedID   string
		expectedOK   bool
	}{
		{
			name:         "valid track URL",
			input:        "https://soundcloud.com/forss/flickermood",
			expectedType: release.TypeTrack,
			expectedID:   "forss:flickermood",
			expectedOK:   true,
		},
		{
			name:         "valid track URL with www host",
			input:        "https://www.soundcloud.com/forss/flickermood",
			expectedType: release.TypeTrack,
			expectedID:   "forss:flickermood",
			expectedOK:   true,
		},
		{
			name:         "valid track URL with query",
			input:        "https://soundcloud.com/forss/flickermood?si=abc123",
			expectedType: release.TypeTrack,
			expectedID:   "forss:flickermood",
			expectedOK:   true,
		},
		{
			name:         "valid album URL",
			input:        "https://soundcloud.com/forss/sets/soulhack",
			expectedType: release.TypeAlbum,
			expectedID:   "forss:soulhack",
			expectedOK:   true,
		},
		{
			name:         "valid album URL with www host",
			input:        "https://www.soundcloud.com/forss/sets/soulhack",
			expectedType: release.TypeAlbum,
			expectedID:   "forss:soulhack",
			expectedOK:   true,
		},
		{
			name:         "valid album URL with query",
			input:        "https://soundcloud.com/forss/sets/soulhack?si=abc123",
			expectedType: release.TypeAlbum,
			expectedID:   "forss:soulhack",
			expectedOK:   true,
		},
		{
			name:       "URL without track ID",
			input:      "https://soundcloud.com/forss",
			expectedOK: false,
		},
		{
			name:       "album URL without ID",
			input:      "https://soundcloud.com/forss/sets/",
			expectedOK: false,
		},
		{
			name:       "track URL with invalid host",
			input:      "https://example.com/forss/flickermood",
			expectedOK: false,
		},
		{
			name:       "album URL with invalid host",
			input:      "https://example.com/forss/sets/soulhack",
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

func TestSoundcloudAdapter_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedTrack release.Track
		expectedErr   string
	}{
		{
			name: "found ID",
			id:   "forss:flickermood",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "forss", "flickermood").
					Return(track{
						Title:       "Flickermood",
						Description: "sample track description",
						Permalink:   "flickermood",
						User: user{
							Username:  "Forss",
							Permalink: "forss",
						},
						PublisherMetadata: publisherMetadata{
							Artist: "Real Artist",
							ISRC:   "GBARL9300135",
						},
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:          "forss:flickermood",
				ISRC:        "GBARL9300135",
				Title:       "Flickermood",
				Artist:      "Real Artist",
				URL:         "https://soundcloud.com/forss/flickermood",
				Provider:    release.Soundcloud,
				Creator:     "Forss",
				Description: "sample track description",
			},
		},
		{
			name: "not found ID",
			id:   "forss:flickermood",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "forss", "flickermood").
					Return(nil, errNotFound).
					Once()
			},
			expectedErr: release.ErrNotFound.Error(),
		},
		{
			name:        "invalid composite key",
			id:          "invalid",
			mockClient:  func(_ *clientMock) {},
			expectedErr: "failed to parse track id: failed to parse key parts: invalid key: expected 2 parts, got 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)
			result, err := a.FetchTrack(t.Context(), tt.id)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Zero(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_fetchAlbum(t *testing.T) {
	const (
		firstID    int64 = 3_000_000_001
		secondID   int64 = 3_000_000_002
		missingID  int64 = 3_000_000_003
		foundURN         = "soundcloud:tracks:found"
		missingURN       = "soundcloud:tracks:missing"
	)

	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedAlbum release.Album
		expectedErr   string
	}{
		{
			name: "found ID",
			id:   "forss:soulhack",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "forss", "soulhack").
					Return(album{
						Title:       "Soulhack",
						LabelName:   "Sample Label",
						Description: "sample album description",
						Permalink:   "soulhack",
						User: user{
							Username:  "Forss",
							Permalink: "forss",
						},
						Tracks: []track{
							{
								Title:        "Flickermood",
								Description:  "first track description",
								Permalink:    "flickermood",
								PermalinkURL: "https://soundcloud.com/forss/flickermood",
								PublisherMetadata: publisherMetadata{
									Artist: "Real Artist",
								},
							},
							{
								Title:        "Using Splashes",
								Description:  "second track description",
								Permalink:    "using-splashes",
								PermalinkURL: "https://soundcloud.com/forss/using-splashes",
							},
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:          "forss:soulhack",
				Title:       "Soulhack",
				Artist:      "Real Artist",
				Label:       "Sample Label",
				URL:         "https://soundcloud.com/forss/sets/soulhack",
				Provider:    release.Soundcloud,
				Creator:     "Forss",
				Description: "sample album description",
				TrackIDs:    []string{"forss:flickermood", "forss:using-splashes"},
			},
		},
		{
			name: "enriches numeric track placeholders",
			id:   "forss:soulhack",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "forss", "soulhack").
					Return(album{
						Title:     "Soulhack",
						Permalink: "soulhack",
						User: user{
							Username:  "Forss",
							Permalink: "forss",
						},
						Tracks: []track{
							{
								Title:        "Already full",
								PermalinkURL: "https://soundcloud.com/forss/already-full",
							},
							{ID: firstID},
							{
								Title:     "Full from slugs",
								Permalink: "full-from-slugs",
								User:      user{Permalink: "forss"},
							},
							{ID: secondID},
							{ID: firstID},
							{ID: missingID},
						},
					}, nil).
					Once()
				m.
					On("fetchTracksByNumericIDs", []int64{firstID, secondID, missingID}).
					Return([]track{
						{
							ID:           secondID,
							URN:          "soundcloud:tracks:3000000002",
							Title:        "Second hydrated",
							PermalinkURL: "https://soundcloud.com/forss/second-hydrated",
						},
						{
							ID:           firstID,
							URN:          "soundcloud:tracks:3000000001",
							Title:        "First hydrated",
							PermalinkURL: "https://soundcloud.com/forss/first-hydrated",
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:       "forss:soulhack",
				Title:    "Soulhack",
				Artist:   "Forss",
				URL:      "https://soundcloud.com/forss/sets/soulhack",
				Provider: release.Soundcloud,
				Creator:  "Forss",
				TrackIDs: []string{
					"forss:already-full",
					"forss:first-hydrated",
					"forss:full-from-slugs",
					"forss:second-hydrated",
					"forss:first-hydrated",
				},
			},
		},
		{
			name: "enriches urn-only track placeholders",
			id:   "forss:soulhack",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "forss", "soulhack").
					Return(album{
						Title:     "Soulhack",
						Permalink: "soulhack",
						User: user{
							Username:  "Forss",
							Permalink: "forss",
						},
						Tracks: []track{
							{URN: foundURN},
							{URN: missingURN},
							{},
						},
					}, nil).
					Once()
				m.
					On("fetchTrackByURN", foundURN).
					Return(track{
						URN:          foundURN,
						Title:        "URN hydrated",
						PermalinkURL: "https://soundcloud.com/forss/urn-hydrated",
					}, nil).
					Once()
				m.
					On("fetchTrackByURN", missingURN).
					Return(nil, errNotFound).
					Once()
			},
			expectedAlbum: release.Album{
				ID:       "forss:soulhack",
				Title:    "Soulhack",
				Artist:   "Forss",
				URL:      "https://soundcloud.com/forss/sets/soulhack",
				Provider: release.Soundcloud,
				Creator:  "Forss",
				TrackIDs: []string{"forss:urn-hydrated"},
			},
		},
		{
			name: "not found ID",
			id:   "forss:soulhack",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "forss", "soulhack").
					Return(nil, errNotFound).
					Once()
			},
			expectedErr: release.ErrNotFound.Error(),
		},
		{
			name: "track hydration failure is operational",
			id:   "forss:soulhack",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "forss", "soulhack").
					Return(album{
						Tracks: []track{{ID: 3_000_000_001}},
					}, nil).
					Once()
				m.
					On("fetchTracksByNumericIDs", []int64{3_000_000_001}).
					Return(nil, errors.New("sample hydration error")).
					Once()
			},
			expectedErr: "failed to enrich incomplete soundcloud album tracks: " +
				"failed to enrich numeric track placeholders: sample hydration error",
		},
		{
			name:        "invalid composite key",
			id:          "invalid",
			mockClient:  func(_ *clientMock) {},
			expectedErr: "failed to parse album id: failed to parse key parts: invalid key: expected 2 parts, got 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)
			result, err := a.FetchAlbum(t.Context(), tt.id)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.Zero(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_searchTracks(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		title          string
		mockClient     func(m *clientMock)
		expectedTracks []release.SearchTrack
		expectedErr    error
		errContains    string
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
							Title:        "Flickermood",
							Description:  "first track description",
							PermalinkURL: "https://soundcloud.com/forss/flickermood",
							PublisherMetadata: publisherMetadata{
								ISRC: "GBARL9300135",
							},
							User: user{
								Username: "Forss",
							},
						},
						{
							Title:        "Using Splashes",
							Description:  "second track description",
							PermalinkURL: "https://soundcloud.com/forss/using-splashes",
							User: user{
								Username: "Forss",
							},
						},
					}, nil).
					Once()
			},
			expectedTracks: []release.SearchTrack{
				{
					ID:          "forss:flickermood",
					ISRC:        "GBARL9300135",
					Title:       "Flickermood",
					Artist:      "Forss",
					URL:         "https://soundcloud.com/forss/flickermood",
					Provider:    release.Soundcloud,
					Creator:     "Forss",
					Description: "first track description",
				},
				{
					ID:          "forss:using-splashes",
					Title:       "Using Splashes",
					Artist:      "Forss",
					URL:         "https://soundcloud.com/forss/using-splashes",
					Provider:    release.Soundcloud,
					Creator:     "Forss",
					Description: "second track description",
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
		{
			name:   "invalid candidate URL",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist", "sample track").
					Return([]track{
						{
							Title:        "Invalid",
							PermalinkURL: "https://soundcloud.com/forss/sets/soulhack",
							User: user{
								Username: "Invalid Artist",
							},
						},
					}, nil).
					Once()
			},
			errContains: "failed to parse composite key from track url",
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
			} else if tt.errContains != "" {
				require.Nil(t, result)
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTracks, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_searchAlbums(t *testing.T) {
	tests := []struct {
		name           string
		artist         string
		title          string
		mockClient     func(m *clientMock)
		expectedAlbums []release.SearchAlbum
		expectedErr    error
		errContains    string
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
							Title:        "Soulhack",
							Description:  "first album description",
							PermalinkURL: "https://soundcloud.com/forss/sets/soulhack",
							User: user{
								Username: "Forss",
							},
						},
						{
							Title:        "Journeyman",
							Description:  "second album description",
							PermalinkURL: "https://soundcloud.com/forss/sets/journeyman",
							User: user{
								Username: "Forss",
							},
						},
					}, nil).
					Once()
			},
			expectedAlbums: []release.SearchAlbum{
				{
					ID:          "forss:soulhack",
					Title:       "Soulhack",
					Artist:      "Forss",
					URL:         "https://soundcloud.com/forss/sets/soulhack",
					Provider:    release.Soundcloud,
					Creator:     "Forss",
					Description: "first album description",
				},
				{
					ID:          "forss:journeyman",
					Title:       "Journeyman",
					Artist:      "Forss",
					URL:         "https://soundcloud.com/forss/sets/journeyman",
					Provider:    release.Soundcloud,
					Creator:     "Forss",
					Description: "second album description",
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
		{
			name:   "invalid candidate URL",
			artist: "sample artist",
			title:  "sample album",
			mockClient: func(m *clientMock) {
				m.
					On("searchAlbums", "sample artist", "sample album").
					Return([]album{
						{
							Title:        "Invalid",
							PermalinkURL: "https://soundcloud.com/forss/flickermood",
							User: user{
								Username: "Invalid Artist",
							},
						},
					}, nil).
					Once()
			},
			errContains: "failed to parse composite key from album url",
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
			} else if tt.errContains != "" {
				require.Nil(t, result)
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbums, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestSoundcloudAdapter_fetchTracksByISRC(t *testing.T) {
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

func (m *clientMock) fetchTrack(_ context.Context, userSlug, trackSlug string) (track, error) {
	args := m.Called(userSlug, trackSlug)
	if args.Get(0) == nil {
		return track{}, args.Error(1)
	}
	return args.Get(0).(track), args.Error(1)
}

func (m *clientMock) fetchAlbum(_ context.Context, userSlug, setSlug string) (album, error) {
	args := m.Called(userSlug, setSlug)
	if args.Get(0) == nil {
		return album{}, args.Error(1)
	}
	return args.Get(0).(album), args.Error(1)
}

func (m *clientMock) fetchTracksByNumericIDs(_ context.Context, ids []int64) ([]track, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]track), args.Error(1)
}

func (m *clientMock) fetchTrackByURN(_ context.Context, urn string) (track, error) {
	args := m.Called(urn)
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

func (m *clientMock) searchAlbums(_ context.Context, artist, title string) ([]album, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]album), args.Error(1)
}
