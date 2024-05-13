package links

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_AppleTrackLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         fixture.FixturesMap
	}{
		{
			name:         "when apple track link given and track found",
			input:        "https://music.apple.com/us/album/angel/724466069?i=724466660",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Spotify",
							Data: "cnvtr/ap/us-724466660/sf",
						},
						{
							Text: "Yandex",
							Data: "cnvtr/ap/us-724466660/ya",
						},
						{
							Text: "Youtube",
							Data: "cnvtr/ap/us-724466660/yt",
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				AppleTracks: map[string][]byte{
					"us-724466660": fixture.Read("apple/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when apple track link given and track not found",
			input:        "https://music.apple.com/us/album/angel/724466069?i=724466660123123123123",
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
