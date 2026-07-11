package yandex

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestYandexAdapter_ParseLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType release.Type
		wantID   string
		wantOK   bool
	}{
		{
			name:     "valid track URL .com",
			input:    "https://music.yandex.com/album/3192570/track/1197793",
			wantType: release.TypeTrack,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:     "valid track URL .ru",
			input:    "https://music.yandex.ru/album/3192570/track/1197793",
			wantType: release.TypeTrack,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:     "valid track URL .by",
			input:    "https://music.yandex.by/album/3192570/track/1197793",
			wantType: release.TypeTrack,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:     "valid track URL .kz",
			input:    "https://music.yandex.kz/album/3192570/track/1197793",
			wantType: release.TypeTrack,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:     "valid track URL .uz",
			input:    "https://music.yandex.uz/album/3192570/track/1197793",
			wantType: release.TypeTrack,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:   "invalid track URL missing ID",
			input:  "https://music.yandex.ru/album/3192570/track/",
			wantOK: false,
		},
		{
			name:   "invalid track URL non-numeric ID",
			input:  "https://music.yandex.ru/album/3192570/track/abc",
			wantOK: false,
		},
		{
			name:   "invalid track URL host",
			input:  "https://example.com/album/3192570/track/1197793",
			wantOK: false,
		},
		{
			name:     "valid album URL .by",
			input:    "https://music.yandex.by/album/1197793",
			wantType: release.TypeAlbum,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:     "valid album URL .kz",
			input:    "https://music.yandex.kz/album/1197793",
			wantType: release.TypeAlbum,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:     "valid album URL .uz",
			input:    "https://music.yandex.uz/album/1197793",
			wantType: release.TypeAlbum,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:     "valid album URL .ru",
			input:    "https://music.yandex.ru/album/1197793",
			wantType: release.TypeAlbum,
			wantID:   "1197793",
			wantOK:   true,
		},
		{
			name:   "invalid album URL missing ID",
			input:  "https://music.yandex.ru/album/",
			wantOK: false,
		},
		{
			name:   "invalid album URL non-numeric ID",
			input:  "https://music.yandex.ru/album/letters",
			wantOK: false,
		},
		{
			name:   "invalid album URL host",
			input:  "https://example.com/album/3192570",
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

func TestYandexAdapter_fetchTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedTrack release.Track
		expectedErr   error
	}{
		{
			name: "found ID",
			id:   "42",
			mockClient: func(m *clientMock) {
				m.
					On("fetchTrack", "42").
					Return(track{
						ID:    "42",
						Title: "sample name",
						Artists: []artist{
							{Name: "sample artist"},
						},
						Albums: []albumRef{
							{ID: 41},
						},
					}, nil).
					Once()
			},
			expectedTrack: release.Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				AlbumID:  "41",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: release.Yandex,
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

func TestYandexAdapter_fetchTrackAllowsMissingArtist(t *testing.T) {
	cm := &clientMock{}
	cm.
		On("fetchTrack", "42").
		Return(track{
			ID:     "42",
			Title:  "sample name",
			Albums: []albumRef{{ID: 41}},
		}, nil).
		Once()
	a := NewAdapter(cm)

	result, err := a.FetchTrack(t.Context(), "42")

	require.NoError(t, err)
	require.Equal(t, release.Track{
		ID:       "42",
		Title:    "sample name",
		Artist:   "",
		AlbumID:  "41",
		URL:      "https://music.yandex.com/album/41/track/42",
		Provider: release.Yandex,
	}, result)
	cm.AssertExpectations(t)
}

func TestYandexAdapter_fetchTrackRejectsMissingAlbum(t *testing.T) {
	cm := &clientMock{}
	cm.
		On("fetchTrack", "42").
		Return(track{
			ID:    "42",
			Title: "sample name",
			Artists: []artist{
				{Name: "sample artist"},
			},
		}, nil).
		Once()
	a := NewAdapter(cm)

	result, err := a.FetchTrack(t.Context(), "42")

	require.Zero(t, result)
	require.ErrorContains(t, err, "missing album")
	cm.AssertExpectations(t)
}

func TestYandexAdapter_fetchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockClient    func(m *clientMock)
		expectedAlbum release.Album
		expectedErr   error
	}{
		{
			name: "found id",
			id:   "42",
			mockClient: func(m *clientMock) {
				m.
					On("fetchAlbum", "42").
					Return(album{
						ID:     42,
						Title:  "sample name",
						Labels: []label{{Name: "sample label"}},
						Artists: []artist{
							{Name: "sample artist"},
						},
						Volumes: [][]track{
							{
								{
									ID:    "100",
									Title: "first track",
									Artists: []artist{
										{Name: "sample artist"},
									},
									Albums: []albumRef{
										{ID: 42},
									},
								},
								{
									ID:    "101",
									Title: "second track",
									Artists: []artist{
										{Name: "sample artist"},
									},
									Albums: []albumRef{
										{ID: 42},
									},
								},
							},
						},
					}, nil).
					Once()
			},
			expectedAlbum: release.Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				Label:    "sample label",
				URL:      "https://music.yandex.com/album/42",
				Provider: release.Yandex,
				TrackIDs: []string{"100", "101"},
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

func TestYandexAdapter_fetchAlbumAllowsMissingArtist(t *testing.T) {
	cm := &clientMock{}
	cm.
		On("fetchAlbum", "42").
		Return(album{
			ID:    42,
			Title: "sample name",
		}, nil).
		Once()
	a := NewAdapter(cm)

	result, err := a.FetchAlbum(t.Context(), "42")

	require.NoError(t, err)
	require.Equal(t, release.Album{
		ID:       "42",
		Title:    "sample name",
		Artist:   "",
		URL:      "https://music.yandex.com/album/42",
		Provider: release.Yandex,
		TrackIDs: []string{},
	}, result)
	cm.AssertExpectations(t)
}

func TestYandexAdapter_searchTracks(t *testing.T) {
	errUnexpected := errors.New("unexpected error")

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
					On("searchTracks", "sample artist \u2013 sample track").
					Return([]searchTrack{
						{
							ID:    123,
							Title: "First Track",
							Artists: []artist{
								{Name: "First Artist"},
							},
							Albums: []albumRef{
								{ID: 456},
							},
						},
						{
							ID:    789,
							Title: "Second Track",
							Artists: []artist{
								{Name: "Second Artist"},
							},
							Albums: []albumRef{
								{ID: 987},
							},
						},
					}, nil).
					Once()
			},
			want: []release.SearchTrack{
				{
					ID:       "123",
					Title:    "First Track",
					Artist:   "First Artist",
					AlbumID:  "456",
					URL:      "https://music.yandex.com/album/456/track/123",
					Provider: release.Yandex,
				},
				{
					ID:       "789",
					Title:    "Second Track",
					Artist:   "Second Artist",
					AlbumID:  "987",
					URL:      "https://music.yandex.com/album/987/track/789",
					Provider: release.Yandex,
				},
			},
		},
		{
			name:   "no candidates",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist \u2013 sample track").
					Return([]searchTrack{}, nil).
					Once()
			},
			want: []release.SearchTrack{},
		},
		{
			name:   "not found returns no candidates",
			artist: "sample artist",
			title:  "sample track",
			mockClient: func(m *clientMock) {
				m.
					On("searchTracks", "sample artist \u2013 sample track").
					Return(nil, errNotFound).
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
					On("searchTracks", "sample artist \u2013 sample track").
					Return(nil, errUnexpected).
					Once()
			},
			wantError: errUnexpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)

			got, err := a.SearchTracks(t.Context(), tt.artist, tt.title)

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

func TestYandexAdapter_searchAlbums(t *testing.T) {
	errUnexpected := errors.New("unexpected error")

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
					On("searchAlbums", "sample artist \u2013 sample album").
					Return([]searchAlbum{
						{
							ID:    123,
							Title: "First Album",
							Artists: []artist{
								{Name: "First Artist"},
							},
						},
						{
							ID:    456,
							Title: "Second Album",
							Artists: []artist{
								{Name: "Second Artist"},
							},
						},
					}, nil).
					Once()
			},
			want: []release.SearchAlbum{
				{
					ID:       "123",
					Title:    "First Album",
					Artist:   "First Artist",
					URL:      "https://music.yandex.com/album/123",
					Provider: release.Yandex,
				},
				{
					ID:       "456",
					Title:    "Second Album",
					Artist:   "Second Artist",
					URL:      "https://music.yandex.com/album/456",
					Provider: release.Yandex,
				},
			},
		},
		{
			name:   "no candidates",
			artist: "sample artist",
			title:  "sample album",
			mockClient: func(m *clientMock) {
				m.
					On("searchAlbums", "sample artist \u2013 sample album").
					Return([]searchAlbum{}, nil).
					Once()
			},
			want: []release.SearchAlbum{},
		},
		{
			name:   "not found returns no candidates",
			artist: "sample artist",
			title:  "sample album",
			mockClient: func(m *clientMock) {
				m.
					On("searchAlbums", "sample artist \u2013 sample album").
					Return(nil, errNotFound).
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
					On("searchAlbums", "sample artist \u2013 sample album").
					Return(nil, errUnexpected).
					Once()
			},
			wantError: errUnexpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientMock{}
			tt.mockClient(cm)

			a := NewAdapter(cm)

			got, err := a.SearchAlbums(t.Context(), tt.artist, tt.title)

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

func (m *clientMock) searchTracks(_ context.Context, query string) ([]searchTrack, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]searchTrack), args.Error(1)
}

func (m *clientMock) fetchAlbum(_ context.Context, id string) (album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return album{}, args.Error(1)
	}
	return args.Get(0).(album), args.Error(1)
}

func (m *clientMock) searchAlbums(_ context.Context, query string) ([]searchAlbum, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]searchAlbum), args.Error(1)
}
