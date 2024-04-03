package links

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_AppleAlbumLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         fixture.FixturesMap
	}{
		{
			name:         "when apple album link given and album found",
			input:        "https://music.apple.com/us/album/amnesiac/1097864180",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Spotify",
							Data: "cnval/ap/1097864180/sf",
						},
						{
							Text: "Yandex",
							Data: "cnval/ap/1097864180/ya",
						},
						{
							Text: "Youtube",
							Data: "cnval/ap/1097864180/yt",
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				AppleAlbums: map[string][]byte{
					"1097864180": fixture.Read("apple/get_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when apple album link given and spotify album not found",
			input:        "https://music.apple.com/us/album/amnesiac/0",
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
