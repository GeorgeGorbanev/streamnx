package music

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/stretchr/testify/require"
)

func TestYandexAdapter_DetectTrackID(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		wantID string
	}{
		{
			name:   "Valid Track URL",
			url:    "https://music.yandex.ru/album/3192570/track/1197793",
			wantID: "1197793",
		},
		{
			name:   "Invalid URL - Missing track ID",
			url:    "https://music.yandex.ru/album/3192570/track/",
			wantID: "",
		},
		{
			name:   "Invalid URL - Non-numeric track ID",
			url:    "https://music.yandex.ru/album/3192570/track/abc",
			wantID: "",
		},
		{
			name:   "Invalid URL - Incorrect format",
			url:    "https://example.com/album/3192570/track/1197793",
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newYandexAdapter(nil)
			result := adapter.DetectTrackID(tt.url)
			require.Equal(t, tt.wantID, result)
		})
	}
}

func TestYandexAdapter_DetectAlbumID(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		wantID string
	}{
		{
			name:   "Valid album URL",
			url:    "https://music.yandex.ru/album/1197793",
			wantID: "1197793",
		},
		{
			name:   "Invalid URL - Missing album ID",
			url:    "https://music.yandex.ru/album/",
			wantID: "",
		},
		{
			name:   "Invalid URL - Non-numeric album ID",
			url:    "https://music.yandex.ru/album/letters",
			wantID: "",
		},
		{
			name:   "Invalid URL - Incorrect host",
			url:    "https://example.com/album/3192570",
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newYandexAdapter(nil)
			result := adapter.DetectAlbumID(tt.url)
			require.Equal(t, tt.wantID, result)
		})
	}
}

func TestYandexAdapter_GetTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedTrack *Track
	}{
		{
			name: "found ID",
			id:   "42",
			expectedTrack: &Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
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
			a := &YandexAdapter{
				client: &yandexClientMock{},
			}

			result, err := a.GetTrack(tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYandexAdapter_SearchTrack(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		expectedTrack *Track
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			expectedTrack: &Track{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/41/track/42",
				Provider: Yandex,
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
			a := &YandexAdapter{
				client: &yandexClientMock{},
			}

			result, err := a.SearchTrack(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYandexAdapter_GetAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedTrack *Album
	}{
		{
			name: "found id",
			id:   "42",
			expectedTrack: &Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
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
			a := &YandexAdapter{
				client: &yandexClientMock{},
			}

			result, err := a.GetAlbum(tt.id)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

func TestYandexAdapter_SearchAlbum(t *testing.T) {
	tests := []struct {
		name          string
		artistName    string
		searchName    string
		expectedTrack *Album
	}{
		{
			name:       "found query",
			artistName: "sample artist",
			searchName: "sample name",
			expectedTrack: &Album{
				ID:       "42",
				Title:    "sample name",
				Artist:   "sample artist",
				URL:      "https://music.yandex.com/album/42",
				Provider: Yandex,
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
			a := &YandexAdapter{
				client: &yandexClientMock{},
			}

			result, err := a.SearchAlbum(tt.artistName, tt.searchName)

			require.NoError(t, err)
			require.Equal(t, tt.expectedTrack, result)
		})
	}
}

type yandexClientMock struct{}

func (c *yandexClientMock) GetTrack(id string) (*yandex.Track, error) {
	if id != "42" {
		return nil, nil
	}
	return &yandex.Track{
		ID:    42,
		Title: "sample name",
		Artists: []yandex.Artist{
			{Name: "sample artist"},
		},
		Albums: []yandex.Album{
			{ID: 41},
		},
	}, nil
}

func (c *yandexClientMock) SearchTrack(artistName, trackName string) (*yandex.Track, error) {
	if artistName != "sample artist" || trackName != "sample name" {
		return nil, nil
	}

	return &yandex.Track{
		ID:    42,
		Title: "sample name",
		Artists: []yandex.Artist{
			{Name: "sample artist"},
		},
		Albums: []yandex.Album{
			{ID: 41},
		},
	}, nil
}

func (c *yandexClientMock) GetAlbum(id string) (*yandex.Album, error) {
	if id != "42" {
		return nil, nil
	}
	return &yandex.Album{
		ID:    42,
		Title: "sample name",
		Artists: []yandex.Artist{
			{Name: "sample artist"},
		},
	}, nil
}

func (c *yandexClientMock) SearchAlbum(artistName, albumName string) (*yandex.Album, error) {
	if artistName != "sample artist" || albumName != "sample name" {
		return nil, nil
	}
	return &yandex.Album{
		ID:    42,
		Title: "sample name",
		Artists: []yandex.Artist{
			{Name: "sample artist"},
		},
	}, nil
}
