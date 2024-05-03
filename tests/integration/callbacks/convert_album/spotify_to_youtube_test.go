package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumSpotifyToYoutube(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when spotify album link given and youtube album found",
			input: "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/yt",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://www.youtube.com/playlist?list=PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C",
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyAlbums: map[string][]byte{
					"1HrMmB5useeZ0F5lHrMvl0": fixture.Read("spotify/get_album_radiohead_amnesiac.json"),
				},
				YoutubeSearchAlbums: map[string][]byte{
					"Radiohead – Amnesiac": fixture.Read("youtube/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:  "when spotify album link given and youtube album not found",
			input: "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/yt",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Album not found in Youtube",
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyAlbums: map[string][]byte{
					"1HrMmB5useeZ0F5lHrMvl0": fixture.Read("spotify/get_album_radiohead_amnesiac.json"),
				},
				YoutubeSearchAlbums: map[string][]byte{},
			},
		},
		{
			name:             "when spotify album not found",
			input:            "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/yt",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				SpotifyAlbums:       map[string][]byte{},
				YoutubeSearchAlbums: map[string][]byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixturesMap.Merge(&tt.fixturesMap)
			defer fixturesMap.Reset()
			defer senderMock.Reset()

			callback := telebot.Callback{
				Sender: user,
				Data:   tt.input,
			}

			vs.CallbackHandler(&callback)

			require.Equal(t, tt.expectedMessages, senderMock.AllSent)
		})
	}
}
