package streamnx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLink(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		want          *Link
		expectedError error
	}{
		{
			name: "Apple track",
			url:  "https://music.apple.com/us/album/song-name/1234567890?i=987654321",
			want: &Link{
				URL:        "https://music.apple.com/us/album/song-name/1234567890?i=987654321",
				Provider:   Apple,
				EntityID:   "us-987654321",
				EntityType: Track,
			},
		},
		{
			name: "Deezer cloak entity",
			url:  "https://link.deezer.com/s/30FbcgrctxIrNQImDnVEZ",
			want: &Link{
				URL:        "https://link.deezer.com/s/30FbcgrctxIrNQImDnVEZ",
				Provider:   Deezer,
				EntityID:   "30FbcgrctxIrNQImDnVEZ",
				EntityType: Cloak,
			},
		},
		{
			name: "Spotify album",
			url:  "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg",
			want: &Link{
				URL:        "https://open.spotify.com/album/7uv632EkfwYhXoqf8rhYrg",
				Provider:   Spotify,
				EntityID:   "7uv632EkfwYhXoqf8rhYrg",
				EntityType: Album,
			},
		},
		{
			name: "Yandex track",
			url:  "https://music.yandex.by/album/3192570/track/1197793",
			want: &Link{
				URL:        "https://music.yandex.by/album/3192570/track/1197793",
				Provider:   Yandex,
				EntityID:   "1197793",
				EntityType: Track,
			},
		},
		{
			name: "Youtube album",
			url:  "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
			want: &Link{
				URL:        "https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
				Provider:   Youtube,
				EntityID:   "PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj",
				EntityType: Album,
			},
		},
		{
			name: "Deezer track",
			url:  "https://www.deezer.com/track/123456789",
			want: &Link{
				URL:        "https://www.deezer.com/track/123456789",
				Provider:   Deezer,
				EntityID:   "123456789",
				EntityType: Track,
			},
		},
		{
			name: "Bandcamp track",
			url:  "https://artistname.bandcamp.com/track/song-name",
			want: &Link{
				URL:        "https://artistname.bandcamp.com/track/song-name",
				Provider:   Bandcamp,
				EntityID:   "artistname:song-name",
				EntityType: Track,
			},
		},
		{
			name: "Bandcamp album",
			url:  "https://artistname.bandcamp.com/album/album-name",
			want: &Link{
				URL:        "https://artistname.bandcamp.com/album/album-name",
				Provider:   Bandcamp,
				EntityID:   "artistname:album-name",
				EntityType: Album,
			},
		},
		{
			name:          "Unknown provider",
			url:           "https://example.com/track/123456789",
			expectedError: UnknownLinkError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLink(tt.url)
			if tt.expectedError != nil {
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
