package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_YandexTrackLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         utils.FixturesMap
	}{
		{
			name:         "when yandex music track link given and track found with .ru",
			input:        "prfx https://music.yandex.ru/album/3192570/track/354093?query=sample",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
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
			fixturesMap: utils.FixturesMap{
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
			fixturesMap: utils.FixturesMap{
				YandexTracks: map[string][]byte{
					"354093": fixture.Read("yandex/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when yandex music track link given and track not found",
			input:        "prfx https://music.yandex.ru/album/3192570/track/0",
			expectedText: "Link is invalid",
			fixturesMap:  utils.FixturesMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			senderMock := utils.NewTelegramSenderMock()

			yandexMockServer := utils.NewYandexAPIServerMock(t, tt.fixturesMap)
			defer yandexMockServer.Close()

			yandexClient := yandex.NewHTTPClient(yandex.WithAPIURL(yandexMockServer.URL))

			vs := vibeshare.NewVibeshare(&vibeshare.Input{
				MusicRegistry: music.NewRegistry(&music.RegistryInput{
					YandexClient: yandexClient,
				}),
				TelegramSender: senderMock,
			})

			user := &telebot.User{
				Username: "sample_username",
			}
			msg := &telebot.Message{
				Sender: user,
				Text:   tt.input,
			}

			vs.TextHandler(msg)

			require.Equal(t, user, senderMock.Response.To)
			require.Equal(t, tt.expectedText, senderMock.Response.Text)
			require.Equal(t, tt.expectedReplyMarkup, senderMock.Response.ReplyMarkup)
		})
	}
}
