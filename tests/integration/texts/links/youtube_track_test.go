package links

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_YoutubeTrackLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         fixture.FixturesMap
	}{
		{
			name:         "when youtube track regular link given and track found",
			input:        "https://www.youtube.com/watch?v=hbe3CQamF8k",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Apple",
							Data: "cnvtr/yt/hbe3CQamF8k/ap",
						},
						{
							Text: "Spotify",
							Data: "cnvtr/yt/hbe3CQamF8k/sf",
						},
						{
							Text: "Yandex",
							Data: "cnvtr/yt/hbe3CQamF8k/ya",
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when youtube track short link given and track found",
			input:        "https://www.youtu.be/hbe3CQamF8k",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Apple",
							Data: "cnvtr/yt/hbe3CQamF8k/ap",
						},
						{
							Text: "Spotify",
							Data: "cnvtr/yt/hbe3CQamF8k/sf",
						},
						{
							Text: "Yandex",
							Data: "cnvtr/yt/hbe3CQamF8k/ya",
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:                "when youtube track regular link given and track found",
			input:               "https://www.youtube.com/watch?v=notFound",
			expectedText:        "No supported link found",
			expectedReplyMarkup: nil,
			fixturesMap: fixture.FixturesMap{
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
