package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_SpotifyAlbumLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         utils.FixturesMap
	}{
		{
			name:         "when spotify album link given and album found",
			input:        "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Yandex",
							Data: "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/ya",
						},
						{
							Text: "Youtube",
							Data: "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/yt",
						},
					},
				},
			},
			fixturesMap: utils.FixturesMap{
				SpotifyAlbums: map[string][]byte{
					"1HrMmB5useeZ0F5lHrMvl0": fixture.Read("spotify/get_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when spotify album link given and spotify album not found",
			input:        "https://open.spotify.com/album/0",
			expectedText: "Link is invalid",
			fixturesMap:  utils.FixturesMap{},
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
