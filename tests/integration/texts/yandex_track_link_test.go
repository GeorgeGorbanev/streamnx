package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_YandexTrackLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         fixture.FixturesMap
	}{
		{
			name:         "when yandex music track link given and track found with .ru",
			input:        "prfx https://music.yandex.ru/album/3192570/track/354093?query=sample",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Apple",
							Data: "cnvtr/ya/354093/ap",
						},
						{
							Text: "Spotify",
							Data: "cnvtr/ya/354093/sf",
						},
						{
							Text: "Youtube",
							Data: "cnvtr/ya/354093/yt",
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				YandexTracks: map[string][]byte{
					"354093": fixture.Read("yandex/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when yandex music track link given and track found with .com",
			input:        "https://music.yandex.com/album/3192570/track/354093",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Apple",
							Data: "cnvtr/ya/354093/ap",
						},
						{
							Text: "Spotify",
							Data: "cnvtr/ya/354093/sf",
						},
						{
							Text: "Youtube",
							Data: "cnvtr/ya/354093/yt",
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				YandexTracks: map[string][]byte{
					"354093": fixture.Read("yandex/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when yandex music track link given and track not found",
			input:        "prfx https://music.yandex.ru/album/3192570/track/0",
			expectedText: "Link is invalid",
			fixturesMap:  fixture.FixturesMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixturesMap.Merge(&tt.fixturesMap)
			defer fixturesMap.Reset()
			defer senderMock.Reset()

			msg := &telebot.Message{
				Sender: user,
				Text:   tt.input,
			}

			vs.TextHandler(msg)

			require.NotNil(t, senderMock.Response)
			require.Equal(t, user, senderMock.Response.To)
			require.Equal(t, tt.expectedText, senderMock.Response.Text)
			require.Equal(t, tt.expectedReplyMarkup, senderMock.Response.ReplyMarkup)
		})
	}
}
