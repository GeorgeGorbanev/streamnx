package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/templates"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumYoutubeToYandex(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when youtube album link given and yandex album found",
			input: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.yandex.com/album/3389008",
				},
				{
					To:   user,
					Text: templates.SpecifyRegion,
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: [][]telebot.InlineButton{
							{
								{
									Text: "🇧🇾 Belarus",
									Data: "regal/3389008/by",
								},
							},
							{
								{
									Text: "🇰🇿 Kazakhstan",
									Data: "regal/3389008/kz",
								},
							},
							{
								{
									Text: "🇷🇺 Russia",
									Data: "regal/3389008/ru",
								},
							},
							{
								{
									Text: "🇺🇿 Uzbekistan",
									Data: "regal/3389008/uz",
								},
							},
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				YandexSearchAlbums: map[string][]byte{
					"radiohead – amnesiac": fixture.Read("yandex/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:  "when youtube album link given and yandex album not found",
			input: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Album not found in Yandex",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				YandexSearchAlbums: map[string][]byte{},
			},
		},
		{
			name:             "when youtube album not found",
			input:            "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ya",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums:      map[string][]byte{},
				YandexSearchAlbums: map[string][]byte{},
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
