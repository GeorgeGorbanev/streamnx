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

func TestText_YandexAlbumLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         utils.FixturesMap
	}{
		{
			name:         "when yandex music album link given and album found",
			input:        "https://music.yandex.com/album/3389008",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Spotify",
							Data: "convert_album/yandex/3389008/spotify",
						},
						{
							Text: "Youtube",
							Data: "convert_album/yandex/3389008/youtube",
						},
					},
				},
			},
			fixturesMap: utils.FixturesMap{
				YandexAlbums: map[string][]byte{
					"3389008": fixture.Read("yandex/get_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when yandex music album link given and yandex music album not found",
			input:        "https://music.yandex.com/album/0",
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
