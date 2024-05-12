package music

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/youtube"
	"github.com/stretchr/testify/require"
)

type youtubeClientMock struct{}

func (c *youtubeClientMock) GetVideo(id string) (*youtube.Video, error) {
	if id != "sampleID" {
		return nil, nil
	}
	return &youtube.Video{
		ID:    "sampleID",
		Title: "sample artist – sample track",
	}, nil
}

func (c *youtubeClientMock) SearchVideo(query string) (*youtube.Video, error) {
	if query != "sample artist – sample track" {
		return nil, nil
	}

	return &youtube.Video{
		ID:    "sampleID",
		Title: "sample artist – sample track",
	}, nil
}

func (c *youtubeClientMock) GetPlaylist(id string) (*youtube.Playlist, error) {
	if id != "sampleID" {
		return nil, nil
	}
	return &youtube.Playlist{
		ID:    "sampleID",
		Title: "sample artist – sample album",
	}, nil
}

func (c *youtubeClientMock) SearchPlaylist(query string) (*youtube.Playlist, error) {
	if query != "sample artist – sample album" {
		return nil, nil
	}
	return &youtube.Playlist{
		ID:    "sampleID",
		Title: "sample artist – sample album",
	}, nil
}

func TestYoutubeAdapter_DetectTrackID(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError error
	}{
		{
			name:     "Short URL",
			input:    "https://youtu.be/dQw4w9WgXcQ",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "Long URL",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:     "URL with extra parameters",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=youtu.be",
			expected: "dQw4w9WgXcQ",
		},
		{
			name:          "Invalid URL",
			input:         "https://www.youtube.com/watch?v=",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Non-YouTube URL",
			input:         "https://www.example.com/watch?v=dQw4w9WgXcQ",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Empty string",
			input:         "",
			expected:      "",
			expectedError: IDNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newYoutubeAdapter(nil)
			result, err := adapter.DetectTrackID(tt.input)
			require.Equal(t, tt.expected, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestYoutubeAdapter_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError error
	}{
		{
			name:     "Standard URL",
			input:    "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expected: "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
		},
		{
			name:     "Shortened URL",
			input:    "https://youtu.be/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expected: "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
		},
		{
			name:     "URL with extra parameters",
			input:    "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj&feature=share",
			expected: "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
		},
		{
			name:          "Invalid URL",
			input:         "https://www.example.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			expected:      "",
			expectedError: IDNotFoundError,
		},
		{
			name:          "Empty string",
			input:         "",
			expected:      "",
			expectedError: IDNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newYoutubeAdapter(nil)
			result, err := adapter.DetectAlbumID(tt.input)
			require.Equal(t, tt.expected, result)

			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestYoutubeAdapter_GetTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedTrack *Track
	}{
		{
			name: "found ID",
			id:   "sampleID",
			expectedTrack: &Track{
				ID:       "sampleID",
				Title:    "sample track",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/watch?v=sampleID",
				Provider: Youtube,
			},
		},
		{
			name:          "not found ID",
			id:            "notFoundID",
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYoutubeAdapter(&youtubeClientMock{})

			result, err := a.GetTrack(tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYoutubeAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		expectedTrack *Track
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample track",
			expectedTrack: &Track{
				ID:       "sampleID",
				Title:    "sample track",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/watch?v=sampleID",
				Provider: Youtube,
			},
		},
		{
			name:          "not found query",
			artistName:    "not found artist",
			searchName:    "not found name",
			expectedTrack: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYoutubeAdapter(&youtubeClientMock{})

			result, err := a.SearchTrack(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYoutubeAdapter_GetAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedAlbum *Album
	}{
		{
			name: "found ID",
			id:   "sampleID",
			expectedAlbum: &Album{
				ID:       "sampleID",
				Title:    "sample album",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/playlist?list=sampleID",
				Provider: Youtube,
			},
		},
		{
			name:          "not found ID",
			id:            "notFoundID",
			expectedAlbum: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newYoutubeAdapter(&youtubeClientMock{})

			result, err := a.GetAlbum(tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedAlbum, result)
		})
	}
}

func TestYoutubeAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		expectedAlbum *Album
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample album",
			expectedAlbum: &Album{
				ID:       "sampleID",
				Title:    "sample album",
				Artist:   "sample artist",
				URL:      "https://www.youtube.com/playlist?list=sampleID",
				Provider: Youtube,
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
			a := newYoutubeAdapter(&youtubeClientMock{})

			result, err := a.SearchAlbum(tt.artistName, tt.searchName)

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
