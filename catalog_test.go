package streamnx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCatalog_ParseLink(t *testing.T) {
	catalog := mustCatalog(t,
		WithApple(),
		WithBandcamp(),
		WithDeezer(),
		WithSpotify(SpotifyCredentials{}),
		WithSoundcloud(),
		WithYandex(),
		WithYoutube(YoutubeCredentials{}),
	)

	tests := []struct {
		name          string
		url           string
		want          Link
		expectedError error
	}{
		{
			name: "Apple track",
			url:  "https://music.apple.com/us/album/song-name/1234567890?i=987654321",
			want: Link{
				URL:         "https://music.apple.com/us/album/song-name/1234567890?i=987654321",
				Provider:    Apple,
				ReleaseID:   "us-987654321",
				ReleaseType: ReleaseTypeTrack,
			},
		},
		{
			name: "Deezer cloak entity",
			url:  "https://link.deezer.com/s/30FbcgrctxIrNQImDnVEZ",
			want: Link{
				URL:         "https://link.deezer.com/s/30FbcgrctxIrNQImDnVEZ",
				Provider:    Deezer,
				ReleaseID:   "30FbcgrctxIrNQImDnVEZ",
				ReleaseType: ReleaseTypeCloak,
			},
		},
		{
			name: "Spotify album",
			url:  "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg",
			want: Link{
				URL:         "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg",
				Provider:    Spotify,
				ReleaseID:   "7uv632EkfwYhXoqf8rhYrg",
				ReleaseType: ReleaseTypeAlbum,
			},
		},
		{
			name: "Soundcloud track",
			url:  "https://soundcloud.com/forss/flickermood",
			want: Link{
				URL:         "https://soundcloud.com/forss/flickermood",
				Provider:    Soundcloud,
				ReleaseID:   "forss:flickermood",
				ReleaseType: ReleaseTypeTrack,
			},
		},
		{
			name: "Soundcloud album",
			url:  "https://soundcloud.com/forss/sets/soulhack",
			want: Link{
				URL:         "https://soundcloud.com/forss/sets/soulhack",
				Provider:    Soundcloud,
				ReleaseID:   "forss:soulhack",
				ReleaseType: ReleaseTypeAlbum,
			},
		},
		{
			name: "Yandex track",
			url:  "https://music.yandex.by/album/3192570/track/1197793",
			want: Link{
				URL:         "https://music.yandex.by/album/3192570/track/1197793",
				Provider:    Yandex,
				ReleaseID:   "1197793",
				ReleaseType: ReleaseTypeTrack,
			},
		},
		{
			name: "Youtube album",
			url:  "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			want: Link{
				URL:         "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
				Provider:    Youtube,
				ReleaseID:   "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
				ReleaseType: ReleaseTypeAlbum,
			},
		},
		{
			name: "Deezer track",
			url:  "https://www.deezer.com/track/123456789",
			want: Link{
				URL:         "https://www.deezer.com/track/123456789",
				Provider:    Deezer,
				ReleaseID:   "123456789",
				ReleaseType: ReleaseTypeTrack,
			},
		},
		{
			name: "Bandcamp track",
			url:  "https://artistname.bandcamp.com/track/song-name",
			want: Link{
				URL:         "https://artistname.bandcamp.com/track/song-name",
				Provider:    Bandcamp,
				ReleaseID:   "artistname:song-name",
				ReleaseType: ReleaseTypeTrack,
			},
		},
		{
			name: "Bandcamp album",
			url:  "https://artistname.bandcamp.com/album/album-name",
			want: Link{
				URL:         "https://artistname.bandcamp.com/album/album-name",
				Provider:    Bandcamp,
				ReleaseID:   "artistname:album-name",
				ReleaseType: ReleaseTypeAlbum,
			},
		},
		{
			name:          "Unknown provider",
			url:           "https://example.com/track/123456789",
			expectedError: ErrUnknownLink,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := catalog.ParseLink(tt.url)
			if tt.expectedError != nil {
				require.Zero(t, got)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewCatalogRegistersExplicitProviders(t *testing.T) {
	am := &adapterMock{}

	catalog, err := NewCatalog(withAdapter(Apple, am))

	require.NoError(t, err)
	require.Same(t, am, catalog.adapters[Apple])
	require.Nil(t, catalog.adapters[Spotify])
}

func TestNewCatalogIgnoresNilOptions(t *testing.T) {
	catalog, err := NewCatalog(
		nil,
		WithApple(nil, WithAppleHTTPClient(nil)),
		WithBandcamp(nil, WithBandcampHTTPClient(nil)),
		WithDeezer(nil, WithDeezerHTTPClient(nil)),
		WithSpotify(SpotifyCredentials{}, nil, WithSpotifyHTTPClient(nil)),
		WithSoundcloud(nil, WithSoundcloudHTTPClient(nil)),
		WithYandex(nil, WithYandexHTTPClient(nil)),
		WithYoutube(YoutubeCredentials{}, nil, WithYoutubeHTTPClient(nil)),
	)

	require.NoError(t, err)
	require.Len(t, catalog.adapters, len(Providers()))
}

func TestProvidersReturnsCopy(t *testing.T) {
	providers := Providers()
	require.NotEmpty(t, providers)

	providers[0] = "modified"

	require.NotEqual(t, ReleaseProvider("modified"), Providers()[0])
}

func TestNewCatalogRegistersMultipleExplicitProviders(t *testing.T) {
	appleAdapter := &adapterMock{}
	spotifyAdapter := &adapterMock{}

	catalog, err := NewCatalog(
		withAdapter(Apple, appleAdapter),
		withAdapter(Spotify, spotifyAdapter),
	)

	require.NoError(t, err)
	require.Same(t, appleAdapter, catalog.adapters[Apple])
	require.Same(t, spotifyAdapter, catalog.adapters[Spotify])
}

func TestNewCatalogRejectsDuplicateProvider(t *testing.T) {
	_, err := NewCatalog(
		withAdapter(Apple, &adapterMock{}),
		withAdapter(Apple, &adapterMock{}),
	)

	require.ErrorIs(t, err, ErrDuplicateProvider)
}

func TestCatalog_FetchTrack(t *testing.T) {
	sampleProvider := Apple
	am := &adapterMock{}
	am.
		On("FetchTrack", "1").
		Return(Track{
			ID:       "1",
			Title:    "Sample Track",
			Artist:   "Sample Artist",
			URL:      "https://music.example/track/1",
			Provider: sampleProvider,
		}, nil).
		Once()

	catalog := mustCatalog(t, withAdapter(sampleProvider, am))

	result, err := catalog.FetchTrack(t.Context(), sampleProvider, "1")

	require.NoError(t, err)
	require.Equal(t, Track{
		ID:       "1",
		Title:    "Sample Track",
		Artist:   "Sample Artist",
		URL:      "https://music.example/track/1",
		Provider: sampleProvider,
	}, result)
	am.AssertExpectations(t)
}

func TestCatalog_FetchAlbum(t *testing.T) {
	sampleProvider := Apple
	am := &adapterMock{}
	am.
		On("FetchAlbum", "1").
		Return(Album{
			ID:       "1",
			Title:    "Sample Album",
			Artist:   "Sample Artist",
			URL:      "https://music.example/album/1",
			Provider: sampleProvider,
		}, nil).
		Once()

	catalog := mustCatalog(t, withAdapter(sampleProvider, am))

	result, err := catalog.FetchAlbum(t.Context(), sampleProvider, "1")

	require.NoError(t, err)
	require.Equal(t, Album{
		ID:       "1",
		Title:    "Sample Album",
		Artist:   "Sample Artist",
		URL:      "https://music.example/album/1",
		Provider: sampleProvider,
	}, result)
	am.AssertExpectations(t)
}

func TestCatalog_FetchTrackRejectsUnregisteredProvider(t *testing.T) {
	catalog := mustCatalog(t, withAdapter(Apple, &adapterMock{}))

	result, err := catalog.FetchTrack(t.Context(), Spotify, "1")

	require.Zero(t, result)
	require.ErrorIs(t, err, ErrInvalidProvider)
}

func TestCatalogRejectsEmptyIDs(t *testing.T) {
	catalog := mustCatalog(t, withAdapter(Apple, &adapterMock{}))

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "fetch track",
			call: func() error {
				_, err := catalog.FetchTrack(t.Context(), Apple, " \t")
				return err
			},
		},
		{
			name: "fetch album",
			call: func() error {
				_, err := catalog.FetchAlbum(t.Context(), Apple, "")
				return err
			},
		},
		{
			name: "uncloak",
			call: func() error {
				_, _, err := catalog.Uncloak(t.Context(), Apple, "\n")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, tt.call(), ErrInvalidID)
		})
	}
}

func TestCatalog_SearchTracks(t *testing.T) {
	sampleProvider := Apple
	am := &adapterMock{}
	am.
		On("SearchTracks", "artist", "name").
		Return([]SearchTrack{
			{
				ID:       "1",
				Title:    "First Track",
				Artist:   "First Artist",
				URL:      "https://music.example/track/1",
				Provider: sampleProvider,
			},
			{
				ID:       "2",
				Title:    "Second Track",
				Artist:   "Second Artist",
				URL:      "https://music.example/track/2",
				Provider: sampleProvider,
			},
		}, nil).
		Once()

	catalog := mustCatalog(t, withAdapter(sampleProvider, am))

	result, err := catalog.SearchTracks(t.Context(), sampleProvider, SearchQuery{
		Artist: "artist",
		Title:  "name",
	})

	require.NoError(t, err)
	require.Equal(t, []SearchTrack{
		{
			ID:       "1",
			Title:    "First Track",
			Artist:   "First Artist",
			URL:      "https://music.example/track/1",
			Provider: sampleProvider,
		},
		{
			ID:       "2",
			Title:    "Second Track",
			Artist:   "Second Artist",
			URL:      "https://music.example/track/2",
			Provider: sampleProvider,
		},
	}, result)
	am.AssertExpectations(t)
}

func TestCatalog_SearchAlbums(t *testing.T) {
	sampleProvider := Apple
	am := &adapterMock{}
	am.
		On("SearchAlbums", "artist", "name").
		Return([]SearchAlbum{
			{
				ID:       "1",
				Title:    "First Album",
				Artist:   "First Artist",
				URL:      "https://music.example/album/1",
				Provider: sampleProvider,
			},
			{
				ID:       "2",
				Title:    "Second Album",
				Artist:   "Second Artist",
				URL:      "https://music.example/album/2",
				Provider: sampleProvider,
			},
		}, nil).
		Once()

	catalog := mustCatalog(t, withAdapter(sampleProvider, am))

	result, err := catalog.SearchAlbums(t.Context(), sampleProvider, SearchQuery{
		Artist: "artist",
		Title:  "name",
	})

	require.NoError(t, err)
	require.Equal(t, []SearchAlbum{
		{
			ID:       "1",
			Title:    "First Album",
			Artist:   "First Artist",
			URL:      "https://music.example/album/1",
			Provider: sampleProvider,
		},
		{
			ID:       "2",
			Title:    "Second Album",
			Artist:   "Second Artist",
			URL:      "https://music.example/album/2",
			Provider: sampleProvider,
		},
	}, result)
	am.AssertExpectations(t)
}

func TestCatalogSearchReturnsNonNilEmptyResults(t *testing.T) {
	adapter := &adapterMock{}
	adapter.On("SearchTracks", "artist", "track").Return(nil, nil).Once()
	adapter.On("SearchAlbums", "artist", "album").Return(nil, nil).Once()

	catalog := mustCatalog(t, withAdapter(Apple, adapter))

	tracks, err := catalog.SearchTracks(t.Context(), Apple, SearchQuery{Artist: "artist", Title: "track"})
	require.NoError(t, err)
	require.NotNil(t, tracks)
	require.Empty(t, tracks)

	albums, err := catalog.SearchAlbums(t.Context(), Apple, SearchQuery{Artist: "artist", Title: "album"})
	require.NoError(t, err)
	require.NotNil(t, albums)
	require.Empty(t, albums)
	adapter.AssertExpectations(t)
}

func TestCatalogSearchAllowsQueryWithOneField(t *testing.T) {
	adapter := &adapterMock{}
	adapter.On("SearchTracks", "", "title").Return([]SearchTrack{}, nil).Once()
	adapter.On("SearchAlbums", "artist", "").Return([]SearchAlbum{}, nil).Once()

	catalog := mustCatalog(t, withAdapter(Apple, adapter))

	_, err := catalog.SearchTracks(t.Context(), Apple, SearchQuery{Title: "title"})
	require.NoError(t, err)

	_, err = catalog.SearchAlbums(t.Context(), Apple, SearchQuery{Artist: "artist"})
	require.NoError(t, err)
	adapter.AssertExpectations(t)
}

func TestCatalogRejectsEmptySearchQuery(t *testing.T) {
	catalog := mustCatalog(t, withAdapter(Apple, &adapterMock{}))

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "track",
			call: func() error {
				_, err := catalog.SearchTracks(t.Context(), Apple, SearchQuery{Artist: " ", Title: "\t"})
				return err
			},
		},
		{
			name: "album",
			call: func() error {
				_, err := catalog.SearchAlbums(t.Context(), Apple, SearchQuery{})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, tt.call(), ErrInvalidSearchQuery)
		})
	}
}

func TestCatalog_Uncloak(t *testing.T) {
	sampleProvider := Deezer
	am := &adapterMock{}
	am.
		On("Uncloak", "cloak-id").
		Return(ReleaseTypeTrack, "track-id", nil).
		Once()

	catalog := mustCatalog(t, withAdapter(sampleProvider, am))

	releaseType, id, err := catalog.Uncloak(t.Context(), sampleProvider, "cloak-id")

	require.NoError(t, err)
	require.Equal(t, ReleaseTypeTrack, releaseType)
	require.Equal(t, "track-id", id)
	am.AssertExpectations(t)
}

func TestCatalog_ParseLinkUsesRegisteredProvidersOnly(t *testing.T) {
	appleAdapter := &adapterMock{}
	appleAdapter.
		On("ParseLink", "https://music.example/track/1").
		Return(ReleaseTypeTrack, "1", true).
		Once()

	catalog := mustCatalog(t, withAdapter(Apple, appleAdapter))

	link, err := catalog.ParseLink("https://music.example/track/1")

	require.NoError(t, err)
	require.Equal(t, Link{
		URL:         "https://music.example/track/1",
		Provider:    Apple,
		ReleaseID:   "1",
		ReleaseType: ReleaseTypeTrack,
	}, link)
	appleAdapter.AssertExpectations(t)
}

func TestCatalog_ParseLinkReturnsUnknownForUnregisteredProvider(t *testing.T) {
	appleAdapter := &adapterMock{}
	appleAdapter.
		On("ParseLink", "https://open.spotify.com/track/1").
		Return(ReleaseType(""), "", false).
		Once()

	catalog := mustCatalog(t, withAdapter(Apple, appleAdapter))

	link, err := catalog.ParseLink("https://open.spotify.com/track/1")

	require.Zero(t, link)
	require.ErrorIs(t, err, ErrUnknownLink)
	appleAdapter.AssertExpectations(t)
}

func mustCatalog(t *testing.T, opts ...CatalogOption) *Catalog {
	t.Helper()

	catalog, err := NewCatalog(opts...)
	require.NoError(t, err)
	return catalog
}

func withAdapter(provider ReleaseProvider, adapter adapter) CatalogOption {
	return func(r *Catalog) error {
		return r.register(provider, adapter)
	}
}

type adapterMock struct {
	mock.Mock
}

func (m *adapterMock) FetchTrack(_ context.Context, id string) (Track, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return Track{}, args.Error(1)
	}
	return args.Get(0).(Track), args.Error(1)
}

func (m *adapterMock) FetchAlbum(_ context.Context, id string) (Album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return Album{}, args.Error(1)
	}
	return args.Get(0).(Album), args.Error(1)
}

func (m *adapterMock) SearchTracks(_ context.Context, artist, title string) ([]SearchTrack, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SearchTrack), args.Error(1)
}

func (m *adapterMock) SearchAlbums(_ context.Context, artist, title string) ([]SearchAlbum, error) {
	args := m.Called(artist, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SearchAlbum), args.Error(1)
}

func (m *adapterMock) ParseLink(rawURL string) (ReleaseType, string, bool) {
	args := m.Called(rawURL)
	return args.Get(0).(ReleaseType), args.String(1), args.Bool(2)
}

func (m *adapterMock) Uncloak(_ context.Context, cloakID string) (ReleaseType, string, error) {
	args := m.Called(cloakID)
	return args.Get(0).(ReleaseType), args.String(1), args.Error(2)
}
