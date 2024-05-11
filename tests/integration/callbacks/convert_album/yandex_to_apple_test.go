package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumYandexToApple(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when yandex album link given and apple album found",
			input: "cnval/ya/3389008/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.apple.com/us/album/amnesiac/1097864180",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YandexAlbums: map[string][]byte{
					"3389008": fixture.Read("yandex/get_album_radiohead_amnesiac.json"),
				},
				AppleSearchAlbums: map[string][]byte{
					"Radiohead Amnesiac": fixture.Read("apple/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:  "when yandex album link given and apple album not found",
			input: "cnval/ya/3389008/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Album not found in Apple",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YandexAlbums: map[string][]byte{
					"3389008": fixture.Read("yandex/get_album_radiohead_amnesiac.json"),
				},
				AppleSearchAlbums: map[string][]byte{},
			},
		},
		{
			name:             "when yandex album not found",
			input:            "cnval/ya/3389008/ap",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				YandexAlbums:      map[string][]byte{},
				AppleSearchAlbums: map[string][]byte{},
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
