package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/templates"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackSpotifyToYandex(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		translatorMock   utils.TranslatorMock
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when spotify track link given and yandex track found",
			input: "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.yandex.com/album/35627/track/354093",
				},
				{
					To:   user,
					Text: templates.SpecifyRegion,
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: [][]telebot.InlineButton{
							{
								{
									Text: "🇧🇾 Belarus",
									Data: "regtr/354093/by",
								},
							},
							{
								{
									Text: "🇰🇿 Kazakhstan",
									Data: "regtr/354093/kz",
								},
							},
							{
								{
									Text: "🇷🇺 Russia",
									Data: "regtr/354093/ru",
								},
							},
							{
								{
									Text: "🇺🇿 Uzbekistan",
									Data: "regtr/354093/uz",
								},
							},
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"massive attack – angel": fixture.Read("yandex/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:  "when spotify track link given and yandex track not found",
			input: "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Track not found in Yandex",
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				YandexSearchTracks: map[string][]byte{},
			},
		},
		{
			name:             "when yandex track not found",
			input:            "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ya",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks:      map[string][]byte{},
				YandexSearchTracks: map[string][]byte{},
			},
		},
		{
			name:  "when spotify track link given, track found and yandex track found, but artist name not match",
			input: "cnvtr/sf/7DSAEUvxU8FajXtRloy8M0/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Track not found in Yandex",
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7DSAEUvxU8FajXtRloy8M0": fixture.Read("spotify/get_track_miley_cyrus_flowers.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"miley cyrus – flowers": fixture.Read("yandex/search_track_miley_cyrus_flowers.json"),
				},
			},
		},
		{
			name:  "when spotify track link given, yandex track found and artist name not match, but match in translit",
			input: "cnvtr/sf/3NHSz1GyC5IeK1soZSjIIX/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.yandex.com/album/81431/track/732401",
				},
				{
					To:   user,
					Text: templates.SpecifyRegion,
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: [][]telebot.InlineButton{
							{
								{
									Text: "🇧🇾 Belarus",
									Data: "regtr/732401/by",
								},
							},
							{
								{
									Text: "🇰🇿 Kazakhstan",
									Data: "regtr/732401/kz",
								},
							},
							{
								{
									Text: "🇷🇺 Russia",
									Data: "regtr/732401/ru",
								},
							},
							{
								{
									Text: "🇺🇿 Uzbekistan",
									Data: "regtr/732401/uz",
								},
							},
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"3NHSz1GyC5IeK1soZSjIIX": fixture.Read("spotify/get_track_zemfira_iskala.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"zemfira – искала": fixture.Read("yandex/search_track_zemfira_iskala.json"),
				},
			},
		},
		{
			name:  "when spotify track link given, track found, yandex track not found, but found in translit",
			input: "cnvtr/sf/2sP5VgY8PWb6c9DhgZEpSv/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.yandex.com/album/4058886/track/33223088",
				},
				{
					To:   user,
					Text: templates.SpecifyRegion,
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: [][]telebot.InlineButton{
							{
								{
									Text: "🇧🇾 Belarus",
									Data: "regtr/33223088/by",
								},
							},
							{
								{
									Text: "🇰🇿 Kazakhstan",
									Data: "regtr/33223088/kz",
								},
							},
							{
								{
									Text: "🇷🇺 Russia",
									Data: "regtr/33223088/ru",
								},
							},
							{
								{
									Text: "🇺🇿 Uzbekistan",
									Data: "regtr/33223088/uz",
								},
							},
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"2sP5VgY8PWb6c9DhgZEpSv": fixture.Read("spotify/get_track_nadezhda_kadysheva_shiroka_reka.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"надежда кадышева – широка река": fixture.Read("yandex/search_track_nadezhda_kadysheva_shiroka_reka.json"),
				},
			},
		},
		{
			name:  "when spotify track link given, track found and artist name is translated",
			input: "cnvtr/sf/7m4RjyO6bsM4HlHCWZDUcM/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.yandex.com/album/3879329/track/31920275",
				},
				{
					To:   user,
					Text: templates.SpecifyRegion,
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: [][]telebot.InlineButton{
							{
								{
									Text: "🇧🇾 Belarus",
									Data: "regtr/31920275/by",
								},
							},
							{
								{
									Text: "🇰🇿 Kazakhstan",
									Data: "regtr/31920275/kz",
								},
							},
							{
								{
									Text: "🇷🇺 Russia",
									Data: "regtr/31920275/ru",
								},
							},
							{
								{
									Text: "🇺🇿 Uzbekistan",
									Data: "regtr/31920275/uz",
								},
							},
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7m4RjyO6bsM4HlHCWZDUcM": fixture.Read("spotify/get_track_dolphin_sumerki.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"dolphin – сумерки": fixture.Read("yandex/search_track_dolphin_sumerki.json"),
				},
			},
			translatorMock: utils.TranslatorMock{
				EnToRu: map[string]string{
					"dolphin": "дельфин",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixturesMap.Merge(&tt.fixturesMap)
			translatorMock.EnToRu = tt.translatorMock.EnToRu
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
