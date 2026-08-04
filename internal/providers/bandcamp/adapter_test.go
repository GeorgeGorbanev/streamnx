package bandcamp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestBandcampAdapter_ParseLink(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType release.Type
		expectedID   string
		expectedOK   bool
	}{
		{
			name:         "valid track URL",
			input:        "https://Autechre.bandcamp.com/track/nil",
			expectedType: release.TypeTrack,
			expectedID:   "autechre:nil",
			expectedOK:   true,
		},
		{
			name:         "valid album URL",
			input:        "https://autechre.bandcamp.com/album/amber",
			expectedType: release.TypeAlbum,
			expectedID:   "autechre:amber",
			expectedOK:   true,
		},
		{
			name:       "invalid URL",
			input:      "https://bandcamp.com/login",
			expectedOK: false,
		},
		{
			name:       "track URL without ID",
			input:      "https://autechre.bandcamp.com/track/",
			expectedOK: false,
		},
		{
			name:       "album URL without ID",
			input:      "https://autechre.bandcamp.com/album/",
			expectedOK: false,
		},
		{
			name:       "track URL with invalid host",
			input:      "https://autechre.danbcamp.com/track/nil",
			expectedOK: false,
		},
		{
			name:       "custom domain is not recognized without search context",
			input:      "https://fantasticvoyage.com/track/natalie-imbruglia-torn-mad-gavs-909-edit",
			expectedOK: false,
		},
		{
			name:       "album URL with invalid host",
			input:      "https://autechre.danbcamp.com/album/amber",
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

func TestBandcampAdapter_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedTrack release.Track
		expectedErr   error
		errContains   string
	}{
		{
			name: "found ID",
			id:   "autechre:nil",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "autechre", "nil").
					Return(Entity{
						Name:        "Nil",
						ISRC:        "GBARL9300135",
						BandName:    "Autechre",
						CreatorName: "Warp Records",
						Description: "sample track description",
						URL:         "https://autechre.com/track/nil",
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:          "autechre:nil",
				ISRC:        "GBARL9300135",
				Title:       "Nil",
				Artist:      "Autechre",
				URL:         "https://autechre.com/track/nil",
				Provider:    release.Bandcamp,
				Creator:     "Warp Records",
				Description: "sample track description",
			},
		},
		{
			name: "not found ID",
			id:   "notfound:notfound",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "notfound", "notfound").
					Return(nil, errNotFound).
					Once()
			},
			expectedErr: release.ErrNotFound,
		},
		{
			name:        "invalid composite ID",
			id:          "autechrenil",
			mockClient:  func(m *clientMock) {},
			errContains: "failed to parse track id",
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
			} else if tt.errContains != "" {
				require.Zero(t, result)
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTrack, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestBandcampAdapter_fetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedAlbum release.Album
		expectedErr   error
		errContains   string
	}{
		{
			name: "found ID",
			id:   "autechre:amber",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "autechre", "amber").
					Return(Entity{
						Name:        "Amber",
						BandName:    "Autechre",
						CreatorName: "Warp Records",
						Description: "sample album description",
						URL:         "https://autechre.com/album/amber",
						TrackURLs: []string{
							"https://autechre.bandcamp.com/track/foil",
							"https://autechre.bandcamp.com/track/montreal",
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:          "autechre:amber",
				Title:       "Amber",
				Artist:      "Autechre",
				Label:       "Warp Records",
				URL:         "https://autechre.com/album/amber",
				Provider:    release.Bandcamp,
				Creator:     "Warp Records",
				Description: "sample album description",
				TrackIDs:    []string{"autechre:foil", "autechre:montreal"},
			},
		},
		{
			name: "not found ID",
			id:   "notfound:notfound",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "notfound", "notfound").
					Return(nil, errNotFound).
					Once()
			},
			expectedErr: release.ErrNotFound,
		},
		{
			name:        "invalid composite ID",
			id:          "autechreamber",
			mockClient:  func(m *clientMock) {},
			errContains: "failed to parse album id",
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
			} else if tt.errContains != "" {
				require.Zero(t, result)
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAlbum, result)
			}

			cm.AssertExpectations(t)
		})
	}
}

func TestBandcampAdapter_searchTracks(t *testing.T) {
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
					Return([]Entity{
						{
							Name:        "Nil",
							BandName:    "Autechre",
							Description: "first track description",
							URL:         "https://autechre.bandcamp.com/track/nil",
						},
						{
							NumericID:  852027615,
							Name:       "Natalie Imbruglia - Torn (mad gavs 909 Edit)",
							AlbumTitle: "mad gavs' 909 Edits",
							BandName:   "mad gavs",
							CoverURL:   "https://f4.bcbits.com/img/2537644019_3.jpg",
							URL:        "https://fantasticvoyage.com/track/natalie-imbruglia-torn-mad-gavs-909-edit",
						},
						{
							Name:        "Foil",
							BandName:    "Autechre",
							Description: "second track description",
							URL:         "https://autechre.bandcamp.com/track/foil",
						},
					}, nil).
					Once()
				m.
					On("resolveSearchResultURL", trackEntityType, uint64(852027615)).
					Return("https://fantastictrax.bandcamp.com/track/natalie-imbruglia-torn-mad-gavs-909-edit", nil).
					Once()
			},
			expectedTracks: []release.SearchTrack{
				{
					ID:          "autechre:nil",
					Title:       "Nil",
					Artist:      "Autechre",
					URL:         "https://autechre.bandcamp.com/track/nil",
					Provider:    release.Bandcamp,
					Creator:     "Autechre",
					Description: "first track description",
				},
				{
					ID:         "fantastictrax:natalie-imbruglia-torn-mad-gavs-909-edit",
					Title:      "Natalie Imbruglia - Torn (mad gavs 909 Edit)",
					Artist:     "mad gavs",
					AlbumTitle: "mad gavs' 909 Edits",
					URL:        "https://fantastictrax.bandcamp.com/track/natalie-imbruglia-torn-mad-gavs-909-edit",
					CoverURL:   "https://f4.bcbits.com/img/2537644019_3.jpg",
					Provider:   release.Bandcamp,
					Creator:    "mad gavs",
				},
				{
					ID:          "autechre:foil",
					Title:       "Foil",
					Artist:      "Autechre",
					URL:         "https://autechre.bandcamp.com/track/foil",
					Provider:    release.Bandcamp,
					Creator:     "Autechre",
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
					Return([]Entity{}, nil).
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
			errContains: "failed to search track on bandcamp",
		},
		{
			name:   "unresolvable candidate fails the whole search",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist", "sample track").
					Return([]Entity{
						{Name: "First", BandName: "First Artist", URL: "https://first.bandcamp.com/track/first"},
						{NumericID: 852027615, Name: "Invalid", BandName: "Invalid Artist", URL: "https://custom.example/track/invalid"},
						{Name: "Last", BandName: "Last Artist", URL: "https://last.bandcamp.com/track/last"},
					}, nil).
					Once()
				m.
					On("resolveSearchResultURL", trackEntityType, uint64(852027615)).
					Return("", errNotFound).
					Once()
			},
			expectedErr: errNotFound,
			errContains: "failed to adapt track search result at index 1 (\"https://custom.example/track/invalid\")",
		},
		{
			name:   "unresolvable candidate without numeric id fails",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist", "sample track").
					Return([]Entity{
						{
							Name:     "Invalid",
							BandName: "Invalid Artist",
							URL:      "https://bandcamp.com/login",
						},
					}, nil).
					Once()
				m.
					On("resolveSearchResultURL", trackEntityType, uint64(0)).
					Return("", errNotFound).
					Once()
			},
			expectedErr: errNotFound,
			errContains: "failed to adapt track search result at index 0 (\"https://bandcamp.com/login\")",
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
				if tt.errContains != "" {
					require.ErrorContains(t, err, tt.errContains)
				}
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

func TestBandcampAdapter_searchAlbums(t *testing.T) {
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
					Return([]Entity{
						{
							Name:        "Amber",
							BandName:    "Autechre",
							Description: "first album description",
							URL:         "https://autechre.bandcamp.com/album/amber",
						},
						{
							Name:        "Tri Repetae",
							BandName:    "Autechre",
							Description: "second album description",
							URL:         "https://autechre.bandcamp.com/album/tri-repetae",
						},
						{
							NumericID:   1884059585,
							Name:        "mad gavs' 909 Edits",
							BandName:    "mad gavs",
							Description: "custom-domain album",
							URL:         "https://fantasticvoyage.com/album/mad-gavs-909-edits",
						},
					}, nil).
					Once()
				m.
					On("resolveSearchResultURL", albumEntityType, uint64(1884059585)).
					Return("https://fantastictrax.bandcamp.com/album/mad-gavs-909-edits", nil).
					Once()
			},
			expectedAlbums: []release.SearchAlbum{
				{
					ID:          "autechre:amber",
					Title:       "Amber",
					Artist:      "Autechre",
					URL:         "https://autechre.bandcamp.com/album/amber",
					Provider:    release.Bandcamp,
					Creator:     "Autechre",
					Description: "first album description",
				},
				{
					ID:          "autechre:tri-repetae",
					Title:       "Tri Repetae",
					Artist:      "Autechre",
					URL:         "https://autechre.bandcamp.com/album/tri-repetae",
					Provider:    release.Bandcamp,
					Creator:     "Autechre",
					Description: "second album description",
				},
				{
					ID:          "fantastictrax:mad-gavs-909-edits",
					Title:       "mad gavs' 909 Edits",
					Artist:      "mad gavs",
					URL:         "https://fantastictrax.bandcamp.com/album/mad-gavs-909-edits",
					Provider:    release.Bandcamp,
					Creator:     "mad gavs",
					Description: "custom-domain album",
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
					Return([]Entity{}, nil).
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
			errContains: "failed to search album on bandcamp",
		},
		{
			name:   "unresolvable candidate fails the whole search",
			artist: "sample artist",
			title:  "sample album",
			mockClient: func(m *clientMock) {
				m.
					On("searchAlbums", "sample artist", "sample album").
					Return([]Entity{
						{
							Name:     "Invalid",
							BandName: "Invalid Artist",
							URL:      "https://bandcamp.com/login",
						},
					}, nil).
					Once()
				m.
					On("resolveSearchResultURL", albumEntityType, uint64(0)).
					Return("", errNotFound).
					Once()
			},
			expectedErr: errNotFound,
			errContains: "failed to adapt album search result at index 0 (\"https://bandcamp.com/login\")",
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
				if tt.errContains != "" {
					require.ErrorContains(t, err, tt.errContains)
				}
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

func TestBandcampAdapter_fetchTracksByISRC(t *testing.T) {
	tracks, err := (&Adapter{}).FetchTracksByISRC(t.Context(), "GBARL9300135")

	require.Nil(t, tracks)
	require.ErrorIs(t, err, release.ErrUnsupportedOperation)
}

type clientMock struct {
	mock.Mock
}

func (m *clientMock) fetchTrack(_ context.Context, artistSlug, id string) (Entity, error) {
	args := m.Called(artistSlug, id)
	if args.Get(0) == nil {
		return Entity{}, args.Error(1)
	}
	return args.Get(0).(Entity), args.Error(1)
}

func (m *clientMock) searchTracks(_ context.Context, artist, title string) ([]Entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *clientMock) fetchAlbum(_ context.Context, artistSlug, id string) (Entity, error) {
	args := m.Called(artistSlug, id)
	if args.Get(0) == nil {
		return Entity{}, args.Error(1)
	}
	return args.Get(0).(Entity), args.Error(1)
}

func (m *clientMock) searchAlbums(_ context.Context, artist, title string) ([]Entity, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *clientMock) resolveSearchResultURL(_ context.Context, et entityType, id uint64) (string, error) {
	args := m.Called(et, id)
	return args.String(0), args.Error(1)
}
