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
			wantID:   "p:OLAK5uy_sample-123",
			wantOK:   true,
		},
		{
			name:     "browse album",
			input:    "https://music.youtube.com/browse/MPREb_sample-123?si=sample",
			wantType: release.TypeAlbum,
			wantID:   "b:MPREb_sample-123",
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
	tests := []struct {
		name           string
		id             string
		mockClient     func(m *clientMock)
		wantTrack      release.Track
		wantErr        error
		wantErrMessage string
	}{
		{
			name: "track metadata",
			id:   "track-id",
			mockClient: func(m *clientMock) {
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
						{
							Text: "Artist",
							NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
								BrowseID: "artist-id",
								Context: browseEndpointContext{Music: browseEndpointMusic{
									PageType: "MUSIC_PAGE_TYPE_ARTIST",
								}},
							}},
						},
						{Text: " • "},
						{
							Text: "Album title",
							NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
								BrowseID: "MPREalbum",
								Context: browseEndpointContext{Music: browseEndpointMusic{
									PageType: "MUSIC_PAGE_TYPE_ALBUM",
								}},
							}},
						},
						{Text: " • "},
						{Text: "1995"},
					}},
				}, nil).Once()
			},
			wantTrack: release.Track{
				ID:          "track-id",
				Title:       "Track title",
				Artist:      "Artist",
				AlbumID:     "b:MPREalbum",
				AlbumTitle:  "Album title",
				URL:         "https://music.youtube.com/watch?v=track-id",
				CoverURL:    "large",
				Duration:    259,
				ReleaseDate: release.Date{Year: 1995},
				Provider:    release.YoutubeMusic,
				Creator:     "Artist - Topic",
				Description: "Description",
			},
		},
		{
			name: "watch next fallbacks",
			id:   "track-id",
			mockClient: func(m *clientMock) {
				found := playerResponse{VideoDetails: videoDetails{VideoID: "track-id"}}
				found.Microformat.Renderer.VideoDetails.DurationSeconds = "259"
				m.On("fetchTrack", "track-id").Return(found, nil).Once()
				m.On("fetchWatchNext", "track-id").Return(playlistPanelVideoRenderer{
					VideoID:   "track-id",
					Title:     text{Runs: []run{{Text: "Track title"}}},
					Thumbnail: thumbnailList{Thumbnails: []thumbnail{{URL: "cover", Width: 544, Height: 544}}},
					LongBylineText: text{Runs: []run{
						{
							Text: "Artist",
							NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
								BrowseID: "artist-id",
								Context: browseEndpointContext{Music: browseEndpointMusic{
									PageType: "MUSIC_PAGE_TYPE_ARTIST",
								}},
							}},
						},
						{Text: " • "},
						{
							Text: "Album title",
							NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
								BrowseID: "MPREalbum",
								Context: browseEndpointContext{Music: browseEndpointMusic{
									PageType: "MUSIC_PAGE_TYPE_ALBUM",
								}},
							}},
						},
					}},
				}, nil).Once()
			},
			wantTrack: release.Track{
				ID:         "track-id",
				Title:      "Track title",
				Artist:     "Artist",
				AlbumID:    "b:MPREalbum",
				AlbumTitle: "Album title",
				URL:        "https://music.youtube.com/watch?v=track-id",
				CoverURL:   "cover",
				Duration:   259,
				Provider:   release.YoutubeMusic,
			},
		},
		{
			name: "not found",
			id:   "track-id",
			mockClient: func(m *clientMock) {
				m.On("fetchTrack", "track-id").Return(playerResponse{}, errNotFound).Once()
			},
			wantErr: release.ErrNotFound,
		},
		{
			name: "login required",
			id:   "track-id",
			mockClient: func(m *clientMock) {
				m.On("fetchTrack", "track-id").Return(playerResponse{}, errLoginRequired).Once()
			},
			wantErr: release.ErrScrapingBlocked,
		},
		{
			name: "track client error",
			id:   "track-id",
			mockClient: func(m *clientMock) {
				m.On("fetchTrack", "track-id").Return(playerResponse{}, errors.New("boom")).Once()
			},
			wantErrMessage: "failed to get track from youtube music: boom",
		},
		{
			name: "watch next client error",
			id:   "track-id",
			mockClient: func(m *clientMock) {
				m.On("fetchTrack", "track-id").Return(playerResponse{
					VideoDetails: videoDetails{VideoID: "track-id"},
				}, nil).Once()
				m.On("fetchWatchNext", "track-id").Return(
					playlistPanelVideoRenderer{}, errors.New("boom"),
				).Once()
			},
			wantErrMessage: "failed to get track metadata from youtube music: boom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &clientMock{}
			tt.mockClient(m)

			got, err := NewAdapter(m).FetchTrack(t.Context(), tt.id)

			if tt.wantErr != nil || tt.wantErrMessage != "" {
				require.Zero(t, got)
				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
				} else {
					require.EqualError(t, err, tt.wantErrMessage)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantTrack, got)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestAdapterFetchAlbum(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockClient func(m *clientMock)
		wantAlbum  release.Album
	}{
		{
			name: "album metadata",
			id:   "b:MPREalbum",
			mockClient: func(m *clientMock) {
				header := responsiveHeader{
					Title:            text{Runs: []run{{Text: "Album title"}}},
					StraplineTextOne: text{Runs: []run{{Text: "Artist"}}},
					Subtitle:         text{Runs: []run{{Text: "Album"}, {Text: " • "}, {Text: "1995"}}},
					Buttons: []responsiveHeaderButton{{
						PlayButton: playButtonRenderer{PlayNavigationEndpoint: navigationEndpoint{
							WatchPlaylistEndpoint: watchPlaylistEndpoint{PlaylistID: "OLAKalbum"},
						}},
					}},
				}
				header.Thumbnail.Renderer.Thumbnail.Thumbnails = []thumbnail{{URL: "cover", Width: 544, Height: 544}}
				header.Description.Shelf.Description.Runs = []run{{Text: "First "}, {Text: "second"}}
				m.On("fetchAlbum", "MPREalbum").Return(header, []responsiveListItem{
					{PlaylistItemData: playlistItemData{VideoID: "track-1"}},
					{
						FlexColumns: []flexColumn{
							{
								Renderer: flexColumnRenderer{
									Text: text{Runs: []run{
										{
											NavigationEndpoint: navigationEndpoint{
												WatchEndpoint: watchEndpoint{VideoID: "track-2"},
											},
										},
									}},
								},
							},
						},
					},
					{},
				}, nil).Once()
			},
			wantAlbum: release.Album{
				ID:             "b:MPREalbum",
				Title:          "Album title",
				Artist:         "Artist",
				URL:            "https://music.youtube.com/browse/MPREalbum",
				AlternativeURL: "https://music.youtube.com/playlist?list=OLAKalbum",
				CoverURL:       "cover",
				ReleaseDate:    release.Date{Year: 1995},
				Provider:       release.YoutubeMusic,
				Description:    "First second",
				TrackIDs:       []string{"track-1", "track-2"},
			},
		},
		{
			name: "playlist album metadata from tracks",
			id:   "p:OLAKalbum",
			mockClient: func(m *clientMock) {
				m.On("fetchAlbum", "OLAKalbum").Return(responsiveHeader{}, []responsiveListItem{{
					FlexColumns: []flexColumn{
						{Renderer: flexColumnRenderer{Text: text{Runs: []run{{Text: "Track title"}}}}},
						{Renderer: flexColumnRenderer{Text: text{Runs: []run{{
							Text: "Artist",
							NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
								Context: browseEndpointContext{Music: browseEndpointMusic{
									PageType: "MUSIC_PAGE_TYPE_ARTIST",
								}},
							}},
						}}}}},
						{Renderer: flexColumnRenderer{Text: text{Runs: []run{{
							Text: "Album title",
							NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
								BrowseID: "MPREalbum",
								Context: browseEndpointContext{Music: browseEndpointMusic{
									PageType: "MUSIC_PAGE_TYPE_ALBUM",
								}},
							}},
						}}}}},
					},
					Thumbnail: itemThumbnail{Renderer: thumbnailRenderer{
						Thumbnail: thumbnailList{Thumbnails: []thumbnail{{URL: "cover", Width: 120, Height: 120}}},
					}},
					PlaylistItemData: playlistItemData{VideoID: "track-1"},
				}}, nil).Once()
			},
			wantAlbum: release.Album{
				ID:             "p:OLAKalbum",
				Title:          "Album title",
				Artist:         "Artist",
				URL:            "https://music.youtube.com/playlist?list=OLAKalbum",
				AlternativeURL: "https://music.youtube.com/browse/MPREalbum",
				CoverURL:       "cover",
				Provider:       release.YoutubeMusic,
				TrackIDs:       []string{"track-1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &clientMock{}
			tt.mockClient(m)

			got, err := NewAdapter(m).FetchAlbum(t.Context(), tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.wantAlbum, got)
			m.AssertExpectations(t)
		})
	}
}

func TestAdapterSearchTracks(t *testing.T) {
	tests := []struct {
		name       string
		artist     string
		title      string
		mockClient func(m *clientMock)
		wantTracks []release.SearchTrack
	}{
		{
			name:   "track candidates",
			artist: "Artist",
			title:  "Track",
			mockClient: func(m *clientMock) {
				m.On("searchTracks", "Artist – Track").Return([]responsiveListItem{
					{
						FlexColumns: []flexColumn{
							{
								Renderer: flexColumnRenderer{Text: text{Runs: []run{{Text: "Track"}}}},
							},
							{
								Renderer: flexColumnRenderer{Text: text{Runs: []run{
									{
										Text: "Artist",
										NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
											BrowseID: "Artist-id",
											Context: browseEndpointContext{Music: browseEndpointMusic{
												PageType: "MUSIC_PAGE_TYPE_ARTIST",
											}},
										}},
									},
									{Text: " • "},
									{
										Text: "Album",
										NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
											BrowseID: "MPREalbum",
											Context: browseEndpointContext{Music: browseEndpointMusic{
												PageType: "MUSIC_PAGE_TYPE_ALBUM",
											}},
										}},
									},
								}}},
							},
						},
						Thumbnail: itemThumbnail{Renderer: thumbnailRenderer{
							Thumbnail: thumbnailList{Thumbnails: []thumbnail{{URL: "cover", Width: 120, Height: 120}}},
						}},
						PlaylistItemData: playlistItemData{VideoID: "video-1"},
					},
					{},
				}, nil).Once()
			},
			wantTracks: []release.SearchTrack{
				{
					ID:         "video-1",
					Title:      "Track",
					Artist:     "Artist",
					AlbumID:    "b:MPREalbum",
					AlbumTitle: "Album",
					URL:        "https://music.youtube.com/watch?v=video-1",
					CoverURL:   "cover",
					Provider:   release.YoutubeMusic,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &clientMock{}
			tt.mockClient(m)

			got, err := NewAdapter(m).SearchTracks(t.Context(), tt.artist, tt.title)

			require.NoError(t, err)
			require.Equal(t, tt.wantTracks, got)
			m.AssertExpectations(t)
		})
	}
}

func TestAdapterSearchAlbums(t *testing.T) {
	tests := []struct {
		name       string
		artist     string
		title      string
		mockClient func(m *clientMock)
		wantAlbums []release.SearchAlbum
	}{
		{
			name:   "album candidates",
			artist: "Artist",
			title:  "Album",
			mockClient: func(m *clientMock) {
				m.On("searchAlbums", "Artist – Album").Return([]responsiveListItem{
					{
						FlexColumns: []flexColumn{
							{
								Renderer: flexColumnRenderer{Text: text{Runs: []run{{Text: "Album"}}}},
							},
							{
								Renderer: flexColumnRenderer{Text: text{Runs: []run{
									{
										Text: "Artist",
										NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
											BrowseID: "Artist-id",
											Context: browseEndpointContext{Music: browseEndpointMusic{
												PageType: "MUSIC_PAGE_TYPE_ARTIST",
											}},
										}},
									},
									{Text: " • "},
									{
										NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{
											Context: browseEndpointContext{Music: browseEndpointMusic{
												PageType: "MUSIC_PAGE_TYPE_ALBUM",
											}},
										}},
									},
								}}},
							},
						},
						Thumbnail: itemThumbnail{Renderer: thumbnailRenderer{
							Thumbnail: thumbnailList{Thumbnails: []thumbnail{{URL: "cover", Width: 120, Height: 120}}},
						}},
						NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{BrowseID: "MPREalbum"}},
						Overlay: itemOverlay{Renderer: itemThumbnailOverlayRenderer{
							Content: itemOverlayContent{PlayButton: playButtonRenderer{
								PlayNavigationEndpoint: navigationEndpoint{WatchPlaylistEndpoint: watchPlaylistEndpoint{
									PlaylistID: "OLAKalbum",
								}},
							}},
						}},
					},
					{
						FlexColumns: []flexColumn{{
							Renderer: flexColumnRenderer{Text: text{Runs: []run{{Text: "Album without playlist"}}}},
						}},
						NavigationEndpoint: navigationEndpoint{BrowseEndpoint: browseEndpoint{BrowseID: "MPREfallback"}},
					},
					{},
				}, nil).Once()
			},
			wantAlbums: []release.SearchAlbum{
				{
					ID:             "b:MPREalbum",
					Title:          "Album",
					Artist:         "Artist",
					URL:            "https://music.youtube.com/browse/MPREalbum",
					AlternativeURL: "https://music.youtube.com/playlist?list=OLAKalbum",
					CoverURL:       "cover",
					Provider:       release.YoutubeMusic,
				},
				{
					ID:       "b:MPREfallback",
					Title:    "Album without playlist",
					URL:      "https://music.youtube.com/browse/MPREfallback",
					Provider: release.YoutubeMusic,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &clientMock{}
			tt.mockClient(m)

			got, err := NewAdapter(m).SearchAlbums(t.Context(), tt.artist, tt.title)

			require.NoError(t, err)
			require.Equal(t, tt.wantAlbums, got)
			m.AssertExpectations(t)
		})
	}
}

func TestAdapterUncloak(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType release.Type
		wantID   string
		wantErr  error
	}{
		{
			name:    "unsupported operation",
			input:   "sample",
			wantErr: release.ErrUnsupportedOperation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotID, err := NewAdapter(&clientMock{}).Uncloak(t.Context(), tt.input)

			require.Equal(t, tt.wantType, gotType)
			require.Equal(t, tt.wantID, gotID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAdapterFetchTracksByISRC(t *testing.T) {
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
