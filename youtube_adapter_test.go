package streaminx

import (
	"context"
	"testing"
	"time"

	"github.com/GeorgeGorbanev/streaminx/internal/youtube"

	"github.com/stretchr/testify/require"
)

type youtubeClientMock struct {
	getVideo         map[string]*youtube.Video
	getPlaylist      map[string]*youtube.Playlist
	searchVideo      map[string]*youtube.Video
	searchPlaylist   map[string]*youtube.Playlist
	getPlaylistItems map[string][]youtube.Video
}

func (c *youtubeClientMock) GetVideo(_ context.Context, id string) (*youtube.Video, error) {
	return c.getVideo[id], nil
}

func (c *youtubeClientMock) SearchVideo(_ context.Context, query string) (*youtube.Video, error) {
	return c.searchVideo[query], nil
}

func (c *youtubeClientMock) GetPlaylist(_ context.Context, id string) (*youtube.Playlist, error) {
	return c.getPlaylist[id], nil
}

func (c *youtubeClientMock) SearchPlaylist(_ context.Context, query string) (*youtube.Playlist, error) {
	return c.searchPlaylist[query], nil
}

func (c *youtubeClientMock) GetPlaylistItems(_ context.Context, id string) ([]youtube.Video, error) {
	return c.getPlaylistItems[id], nil
}

func TestYoutubeAdapter_FetchTrack(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		youtubeClientMock youtubeClientMock
		expectedTrack     *Entity
	}{
		{
			name: "found ID",
			id:   "sampleID",
			youtubeClientMock: youtubeClientMock{
				getVideo: map[string]*youtube.Video{
					"sampleID": {
						ID:    "sampleID",
						Title: "sample artist – sample track",
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample track",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/watch?v=sampleID",
				Provider: Youtube,
				Type:     Track,
			},
		},
		{
			name: "found ID autogenerated track video",
			id:   "sampleID",
			youtubeClientMock: youtubeClientMock{
				getVideo: map[string]*youtube.Video{
					"sampleID": {
						ID:           "sampleID",
						Title:        "track name (remastered)",
						Description:  "bla bla bla. Auto-generated by YouTube",
						ChannelTitle: "sample artist - Topic",
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "track name",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/watch?v=sampleID",
				Provider: Youtube,
				Type:     Track,
			},
		},
		{
			name:              "not found ID",
			id:                "notFoundID",
			youtubeClientMock: youtubeClientMock{},
			expectedTrack:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYoutubeAdapter(&tt.youtubeClientMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.FetchTrack(ctx, tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYoutubeAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name              string
		artistName        string
		searchName        string
		youtubeClientMock youtubeClientMock
		expectedTrack     *Entity
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample track",
			youtubeClientMock: youtubeClientMock{
				searchVideo: map[string]*youtube.Video{
					"sample artist – sample track": {
						ID:    "sampleID",
						Title: "sample artist – sample track",
					},
				},
			},
			expectedTrack: &Entity{
				ID:       "sampleID",
				Title:    "sample track",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/watch?v=sampleID",
				Provider: Youtube,
				Type:     Track,
			},
		},
		{
			name:              "not found query",
			artistName:        "not found artist",
			searchName:        "not found name",
			youtubeClientMock: youtubeClientMock{},
			expectedTrack:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYoutubeAdapter(&tt.youtubeClientMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.SearchTrack(ctx, tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYoutubeAdapter_FetchAlbum(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		youtubeClientMock youtubeClientMock
		expectedAlbum     *Entity
	}{
		{
			name: "found ID",
			id:   "sampleID",
			youtubeClientMock: youtubeClientMock{
				getPlaylist: map[string]*youtube.Playlist{
					"sampleID": {
						ID:    "sampleID",
						Title: "sample artist – sample album",
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "sampleID",
				Title:    "sample album",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/playlist?list=sampleID",
				Provider: Youtube,
				Type:     Album,
			},
		},
		{
			name: "found ID autogenerated album playlist",
			id:   "sampleID",
			youtubeClientMock: youtubeClientMock{
				getPlaylist: map[string]*youtube.Playlist{
					"sampleID": {
						ID:    "sampleID",
						Title: "Album - sample album",
					},
				},
				getPlaylistItems: map[string][]youtube.Video{
					"sampleID": {
						{
							ID:           "sampleTrackID",
							Title:        "sample track",
							ChannelTitle: "sample artist - Topic",
							Description:  "bla bla bla Auto-generated by YouTube bla bla bla",
						},
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "sampleID",
				Title:    "sample album",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/playlist?list=sampleID",
				Provider: Youtube,
				Type:     Album,
			},
		},
		{
			name: "found ID autogenerated album playlist empty",
			id:   "sampleID",
			youtubeClientMock: youtubeClientMock{
				getPlaylist: map[string]*youtube.Playlist{
					"sampleID": {
						ID:    "sampleID",
						Title: "Album - sample album",
					},
				},
				getPlaylistItems: map[string][]youtube.Video{
					"sampleID": {},
				},
			},
			expectedAlbum: &Entity{
				ID:       "sampleID",
				Title:    "sample album",
				Artist:   "Album",
				URL:      "https://www.youtube.com/playlist?list=sampleID",
				Provider: Youtube,
				Type:     Album,
			},
		},
		{
			name: "found ID autogenerated album not autogenerated track",
			id:   "sampleID",
			youtubeClientMock: youtubeClientMock{
				getPlaylist: map[string]*youtube.Playlist{
					"sampleID": {
						ID:    "sampleID",
						Title: "Album - sample album",
					},
				},
				getPlaylistItems: map[string][]youtube.Video{
					"sampleID": {
						{
							ID:           "sampleTrackID",
							Title:        "sample track",
							ChannelTitle: "sample artist - Topic",
							Description:  "any not autogenerated description",
						},
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "sampleID",
				Title:    "sample album",
				Artist:   "Album",
				URL:      "https://www.youtube.com/playlist?list=sampleID",
				Provider: Youtube,
				Type:     Album,
			},
		},
		{
			name:              "not found ID",
			id:                "notFoundID",
			youtubeClientMock: youtubeClientMock{},
			expectedAlbum:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYoutubeAdapter(&tt.youtubeClientMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.FetchAlbum(ctx, tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedAlbum, result)
		})
	}
}

func TestYoutubeAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name              string
		artistName        string
		searchName        string
		youtubeClientMock youtubeClientMock
		expectedAlbum     *Entity
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample album",
			youtubeClientMock: youtubeClientMock{
				searchPlaylist: map[string]*youtube.Playlist{
					"sample artist – sample album": {
						ID:    "sampleID",
						Title: "sample artist – sample album",
					},
				},
			},
			expectedAlbum: &Entity{
				ID:       "sampleID",
				Title:    "sample album",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/playlist?list=sampleID",
				Provider: Youtube,
				Type:     Album,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			expectedAlbum: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYoutubeAdapter(&tt.youtubeClientMock)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := a.SearchAlbum(ctx, tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedAlbum, result)
		})
	}
}

func TestYoutubeAdapter_cleanAndSplitTitle(t *testing.T) {
	tests := []struct {
		title          string
		expectedArtist string
		expectedEntity string
	}{
		{
			title:          "rick astley - never gonna give you up",
			expectedArtist: "rick astley",
			expectedEntity: "never gonna give you up",
		},
		{
			title:          "radiohead amnesiac (2001)",
			expectedArtist: "radiohead",
			expectedEntity: "amnesiac",
		},
		{
			title:          "artist | title [official music video]",
			expectedArtist: "artist",
			expectedEntity: "title",
		},
		{
			title:          "band {live} - song",
			expectedArtist: "band",
			expectedEntity: "song",
		},
		{
			title:          "Michael Jackson - Billie Jean",
			expectedArtist: "Michael Jackson",
			expectedEntity: "Billie Jean",
		},
		{
			title:          "queen – bohemian rhapsody (official video)",
			expectedArtist: "queen",
			expectedEntity: "bohemian rhapsody",
		},
		{
			title:          "adele | someone like you",
			expectedArtist: "adele",
			expectedEntity: "someone like you",
		},
		{
			title:          "the beatles - hey jude [HQ]",
			expectedArtist: "the beatles",
			expectedEntity: "hey jude",
		},
		{
			title:          "coldplay – yellow (official video)",
			expectedArtist: "coldplay",
			expectedEntity: "yellow",
		},
	}
	for _, test := range tests {
		adapter := &YoutubeAdapter{}
		artist, entity := adapter.cleanAndSplitTitle(test.title)
		require.Equal(t, test.expectedArtist, artist)
		require.Equal(t, test.expectedEntity, entity)
	}
}
