package youtubemusic

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

func TestAdapterParseLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType release.Type
		wantID   string
		wantOK   bool
	}{
		{
			name:     "track",
			input:    "listen here: https://music.youtube.com/watch?si=sample&v=5PgdZDXg0z0",
			wantType: release.TypeTrack,
			wantID:   "5PgdZDXg0z0",
			wantOK:   true,
		},
		{
			name:     "track without scheme",
			input:    "music.youtube.com/watch?v=5PgdZDXg0z0",
			wantType: release.TypeTrack,
			wantID:   "5PgdZDXg0z0",
			wantOK:   true,
		},
		{
			name:     "playlist album",
			input:    "https://music.youtube.com/playlist?list=OLAK5uy_sample-123",
			wantType: release.TypeAlbum,
			wantID:   "OLAK5uy_sample-123",
			wantOK:   true,
		},
		{
			name:     "browse album",
			input:    "https://music.youtube.com/browse/MPREb_sample-123?si=sample",
			wantType: release.TypeAlbum,
			wantID:   "MPREb_sample-123",
			wantOK:   true,
		},
		{
			name:   "regular youtube belongs to youtube provider",
			input:  "https://www.youtube.com/watch?v=5PgdZDXg0z0",
			wantOK: false,
		},
		{
			name:   "missing id",
			input:  "https://music.youtube.com/watch?v=",
			wantOK: false,
		},
		{
			name:   "empty",
			input:  "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter(&clientMock{})

			gotType, gotID, gotOK := a.ParseLink(tt.input)

			require.Equal(t, tt.wantType, gotType)
			require.Equal(t, tt.wantID, gotID)
			require.Equal(t, tt.wantOK, gotOK)
		})
	}
}

func TestAdapterFetchTrack(t *testing.T) {
	m := &clientMock{}
	found := playerResponse{
		VideoDetails: videoDetails{
			VideoID:       "track-id",
			Title:         "Track title",
			Author:        "Artist",
			LengthSeconds: "259",
			Thumbnail: thumbnailList{Thumbnails: []thumbnail{
				{URL: "small", Width: 60, Height: 60},
				{URL: "large", Width: 544, Height: 544},
			}},
		},
	}
	found.Microformat.Renderer.URLCanonical = "https://music.youtube.com/watch?v=track-id"
	found.Microformat.Renderer.Description = "Description"
	found.Microformat.Renderer.PageOwnerDetails.Name = "Artist - Topic"
	m.On("fetchTrack", "track-id").Return(found, nil).Once()
	m.On("fetchWatchNext", "track-id").Return(playlistPanelVideoRenderer{
		VideoID: "track-id",
		LongBylineText: text{Runs: []run{
			metadataRun("Artist", "artist-id", "MUSIC_PAGE_TYPE_ARTIST"),
			{Text: " • "},
			metadataRun("Album title", "MPREalbum", "MUSIC_PAGE_TYPE_ALBUM"),
			{Text: " • "},
			{Text: "1995"},
		}},
	}, nil).Once()

	got, err := NewAdapter(m).FetchTrack(t.Context(), "track-id")

	require.NoError(t, err)
	require.Equal(t, release.Track{
		ID:          "track-id",
		Title:       "Track title",
		Artist:      "Artist",
		AlbumID:     "MPREalbum",
		AlbumTitle:  "Album title",
		URL:         "https://music.youtube.com/watch?v=track-id",
		CoverURL:    "large",
		Duration:    259,
		ReleaseDate: release.Date{Year: 1995},
		Provider:    release.YoutubeMusic,
		Creator:     "Artist - Topic",
		Description: "Description",
	}, got)
	m.AssertExpectations(t)
}

func TestAdapterFetchTrackUsesWatchNextFallbacks(t *testing.T) {
	m := &clientMock{}
	found := playerResponse{VideoDetails: videoDetails{VideoID: "track-id"}}
	found.Microformat.Renderer.VideoDetails.DurationSeconds = "259"
	m.On("fetchTrack", "track-id").Return(found, nil).Once()
	m.On("fetchWatchNext", "track-id").Return(playlistPanelVideoRenderer{
		VideoID:   "track-id",
		Title:     text{Runs: []run{{Text: "Track title"}}},
		Thumbnail: thumbnailList{Thumbnails: []thumbnail{{URL: "cover", Width: 544, Height: 544}}},
		LongBylineText: text{Runs: []run{
			metadataRun("Artist", "artist-id", "MUSIC_PAGE_TYPE_ARTIST"),
			{Text: " • "},
			metadataRun("Album title", "MPREalbum", "MUSIC_PAGE_TYPE_ALBUM"),
		}},
	}, nil).Once()

	got, err := NewAdapter(m).FetchTrack(t.Context(), "track-id")

	require.NoError(t, err)
	require.Equal(t, release.Track{
		ID:         "track-id",
		Title:      "Track title",
		Artist:     "Artist",
		AlbumID:    "MPREalbum",
		AlbumTitle: "Album title",
		URL:        "https://music.youtube.com/watch?v=track-id",
		CoverURL:   "cover",
		Duration:   259,
		Provider:   release.YoutubeMusic,
	}, got)
	m.AssertExpectations(t)
}

func TestAdapterFetchTrackErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   error
		want    error
		message string
	}{
		{name: "not found", input: errNotFound, want: release.ErrNotFound},
		{name: "login required", input: errLoginRequired, want: release.ErrScrapingBlocked},
		{name: "client error", input: errors.New("boom"), message: "failed to get track from youtube music: boom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &clientMock{}
			m.On("fetchTrack", "track-id").Return(playerResponse{}, tt.input).Once()

			got, err := NewAdapter(m).FetchTrack(t.Context(), "track-id")

			require.Zero(t, got)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.EqualError(t, err, tt.message)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestAdapterFetchTrackWatchNextError(t *testing.T) {
	m := &clientMock{}
	m.On("fetchTrack", "track-id").Return(playerResponse{
		VideoDetails: videoDetails{VideoID: "track-id"},
	}, nil).Once()
	m.On("fetchWatchNext", "track-id").Return(playlistPanelVideoRenderer{}, errors.New("boom")).Once()

	got, err := NewAdapter(m).FetchTrack(t.Context(), "track-id")

	require.Zero(t, got)
	require.EqualError(t, err, "failed to get track metadata from youtube music: boom")
	m.AssertExpectations(t)
}

func TestAdapterFetchAlbum(t *testing.T) {
	m := &clientMock{}
	header := responsiveHeader{
		Title:            text{Runs: []run{{Text: "Album title"}}},
		StraplineTextOne: text{Runs: []run{{Text: "Artist"}}},
		Subtitle:         text{Runs: []run{{Text: "Album"}, {Text: " • "}, {Text: "1995"}}},
	}
	header.Thumbnail.Renderer.Thumbnail.Thumbnails = []thumbnail{{URL: "cover", Width: 544, Height: 544}}
	header.Description.Shelf.Description.Runs = []run{{Text: "First "}, {Text: "second"}}
	m.On("fetchAlbum", "MPREalbum").Return(header, []responsiveListItem{
		itemWithVideoID("track-1"),
		itemWithWatchEndpoint("track-2"),
		{},
	}, nil).Once()

	got, err := NewAdapter(m).FetchAlbum(t.Context(), "MPREalbum")

	require.NoError(t, err)
	require.Equal(t, release.Album{
		ID:          "MPREalbum",
		Title:       "Album title",
		Artist:      "Artist",
		URL:         "https://music.youtube.com/browse/MPREalbum",
		CoverURL:    "cover",
		ReleaseDate: release.Date{Year: 1995},
		Provider:    release.YoutubeMusic,
		Description: "First second",
		TrackIDs:    []string{"track-1", "track-2"},
	}, got)
	m.AssertExpectations(t)
}

func TestAdapterSearchTracks(t *testing.T) {
	m := &clientMock{}
	m.On("searchTracks", "Artist – Track").Return([]responsiveListItem{
		searchItem("video-1", "Track", "Artist", "MPREalbum", "Album", "cover"),
		{},
	}, nil).Once()

	got, err := NewAdapter(m).SearchTracks(t.Context(), "Artist", "Track")

	require.NoError(t, err)
	require.Equal(t, []release.SearchTrack{{
		ID:         "video-1",
		Title:      "Track",
		Artist:     "Artist",
		AlbumID:    "MPREalbum",
		AlbumTitle: "Album",
		URL:        "https://music.youtube.com/watch?v=video-1",
		CoverURL:   "cover",
		Provider:   release.YoutubeMusic,
	}}, got)
	m.AssertExpectations(t)
}

func TestAdapterSearchAlbums(t *testing.T) {
	m := &clientMock{}
	item := searchItem("", "Album", "Artist", "", "", "cover")
	item.NavigationEndpoint.BrowseEndpoint.BrowseID = "MPREalbum"
	m.On("searchAlbums", "Artist – Album").Return([]responsiveListItem{item, {}}, nil).Once()

	got, err := NewAdapter(m).SearchAlbums(t.Context(), "Artist", "Album")

	require.NoError(t, err)
	require.Equal(t, []release.SearchAlbum{{
		ID:       "MPREalbum",
		Title:    "Album",
		Artist:   "Artist",
		URL:      "https://music.youtube.com/browse/MPREalbum",
		CoverURL: "cover",
		Provider: release.YoutubeMusic,
	}}, got)
	m.AssertExpectations(t)
}

func TestAdapterUncloak(t *testing.T) {
	gotType, gotID, err := NewAdapter(&clientMock{}).Uncloak(t.Context(), "sample")

	require.Empty(t, gotType)
	require.Empty(t, gotID)
	require.ErrorIs(t, err, release.ErrUnsupportedOperation)
}

func searchItem(videoID, title, artist, albumID, albumTitle, cover string) responsiveListItem {
	artistRun := metadataRun(artist, artist+"-id", "MUSIC_PAGE_TYPE_ARTIST")
	albumRun := metadataRun(albumTitle, albumID, "MUSIC_PAGE_TYPE_ALBUM")
	item := responsiveListItem{
		FlexColumns: []flexColumn{{}, {}},
	}
	item.Thumbnail.Renderer.Thumbnail.Thumbnails = []thumbnail{{URL: cover, Width: 120, Height: 120}}
	item.PlaylistItemData.VideoID = videoID
	item.FlexColumns[0].Renderer.Text.Runs = []run{{Text: title}}
	item.FlexColumns[1].Renderer.Text.Runs = []run{artistRun, {Text: " • "}, albumRun}
	return item
}

func metadataRun(textValue, browseID, pageType string) run {
	value := run{
		Text: textValue,
		NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
			BrowseID: browseID,
		}},
	}
	value.NavigationEndpoint.BrowseEndpoint.Context.Music.PageType = pageType
	return value
}

func itemWithVideoID(id string) responsiveListItem {
	item := responsiveListItem{}
	item.PlaylistItemData.VideoID = id
	return item
}

func itemWithWatchEndpoint(id string) responsiveListItem {
	item := responsiveListItem{FlexColumns: []flexColumn{{}}}
	item.FlexColumns[0].Renderer.Text.Runs = []run{{
		NavigationEndpoint: navigationEndpoint{WatchEndpoint: watchEndpoint{VideoID: id}},
	}}
	return item
}

type clientMock struct {
	mock.Mock
}

func (m *clientMock) fetchTrack(_ context.Context, id string) (playerResponse, error) {
	args := m.Called(id)
	return args.Get(0).(playerResponse), args.Error(1)
}

func (m *clientMock) fetchWatchNext(_ context.Context, id string) (playlistPanelVideoRenderer, error) {
	args := m.Called(id)
	return args.Get(0).(playlistPanelVideoRenderer), args.Error(1)
}

func (m *clientMock) fetchAlbum(_ context.Context, id string) (responsiveHeader, []responsiveListItem, error) {
	args := m.Called(id)
	return args.Get(0).(responsiveHeader), args.Get(1).([]responsiveListItem), args.Error(2)
}

func (m *clientMock) searchTracks(_ context.Context, query string) ([]responsiveListItem, error) {
	args := m.Called(query)
	return args.Get(0).([]responsiveListItem), args.Error(1)
}

func (m *clientMock) searchAlbums(_ context.Context, query string) ([]responsiveListItem, error) {
	args := m.Called(query)
	return args.Get(0).([]responsiveListItem), args.Error(1)
}
