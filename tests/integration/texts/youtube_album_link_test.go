package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_YoutubeAlbumLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         utils.FixturesMap
	}{
		{
			name:         "when youtube album regular link given and album found",
			input:        "https://www.youtube.com/playlist?list=PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Spotify",
							Data: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
						},
						{
							Text: "Yandex",
							Data: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ya",
						},
					},
				},
			},
			fixturesMap: utils.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read(
						"youtube/get_album_radiohead_amnesiac.json",
					),
				},
			},
		},
		{
			name:         "when youtube album short link given and album found",
			input:        "https://youtu.be/playlist?list=PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Spotify",
							Data: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
						},
						{
							Text: "Yandex",
							Data: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ya",
						},
					},
				},
			},
			fixturesMap: utils.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read(
						"youtube/get_album_radiohead_amnesiac.json",
					),
				},
			},
		},
		{
			name:                "when youtube album regular link given and album found",
			input:               "https://www.youtube.com/watch?v=notFound",
			expectedText:        "No supported link found",
			expectedReplyMarkup: nil,
			fixturesMap: utils.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"notFound": fixture.Read("youtube/not_found.json"),
				},
			},
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
