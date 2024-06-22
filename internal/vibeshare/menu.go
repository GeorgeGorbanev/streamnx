package vibeshare

import (
	"github.com/GeorgeGorbanev/vibeshare/internal/streaminx"
	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/templates"
	"github.com/tucnak/telebot"
)

var (
	regionsLabels = map[string]string{
		"by": "🇧🇾 Belarus",
		"kz": "🇰🇿 Kazakhstan",
		"ru": "🇷🇺 Russia",
		"uz": "🇺🇿 Uzbekistan",
	}
)

func convertTrackMenu(track *streaminx.Track) *telebot.ReplyMarkup {
	buttonsParams := make([]convertParams, 0, len(streaminx.Providers)-1)
	for _, provider := range streaminx.Providers {
		if provider.Code != track.Provider.Code {
			buttonsParams = append(buttonsParams, convertParams{
				ID:     track.ID,
				Source: track.Provider,
				Target: provider,
			})
		}
	}

	buttons := make([]telebot.InlineButton, 0, len(buttonsParams))
	for _, buttonParams := range buttonsParams {
		cbData := telegram.CallbackData{
			Route:  convertTrackCallbackRoute,
			Params: buttonParams.marshal(),
		}

		buttons = append(buttons, telebot.InlineButton{
			Text: buttonParams.Target.Name,
			Data: cbData.Marshal(),
		})
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			buttons,
		},
	}
}

func convertAlbumMenu(album *streaminx.Album) *telebot.ReplyMarkup {
	buttonsParams := make([]convertParams, 0, len(streaminx.Providers)-1)
	for _, provider := range streaminx.Providers {
		if provider != album.Provider {
			buttonsParams = append(buttonsParams, convertParams{
				ID:     album.ID,
				Source: album.Provider,
				Target: provider,
			})
		}
	}

	buttons := make([]telebot.InlineButton, 0, len(buttonsParams))
	for _, buttonParams := range buttonsParams {
		cbData := telegram.CallbackData{
			Route:  convertAlbumCallbackRoute,
			Params: buttonParams.marshal(),
		}

		buttons = append(buttons, telebot.InlineButton{
			Text: buttonParams.Target.Name,
			Data: cbData.Marshal(),
		})
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			buttons,
		},
	}
}

func trackRegionMenu(track *streaminx.Track) *telebot.ReplyMarkup {
	buttonsParams := make([]regionParams, 0, len(streaminx.Yandex.Regions))
	for _, locale := range streaminx.Yandex.Regions {
		buttonsParams = append(buttonsParams, regionParams{
			EntityID: track.ID,
			Region:   locale,
		})
	}

	buttons := make([][]telebot.InlineButton, 0, len(buttonsParams))
	for _, buttonParams := range buttonsParams {
		cbData := telegram.CallbackData{
			Route:  trackRegionCallbackRoute,
			Params: buttonParams.marshal(),
		}

		buttons = append(buttons, []telebot.InlineButton{
			{
				Text: regionsLabels[buttonParams.Region],
				Data: cbData.Marshal(),
			},
		})
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: buttons,
	}
}

func albumRegionMenu(album *streaminx.Album) *telebot.ReplyMarkup {
	buttonsParams := make([]regionParams, 0, len(streaminx.Yandex.Regions))
	for _, r := range streaminx.Yandex.Regions {
		buttonsParams = append(buttonsParams, regionParams{
			EntityID: album.ID,
			Region:   r,
		})
	}

	buttons := make([][]telebot.InlineButton, 0, len(buttonsParams))
	for _, buttonParams := range buttonsParams {
		cbData := telegram.CallbackData{
			Route:  albumRegionCallbackRoute,
			Params: buttonParams.marshal(),
		}

		buttons = append(buttons, []telebot.InlineButton{
			{
				Text: regionsLabels[buttonParams.Region],
				Data: cbData.Marshal(),
			},
		})
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: buttons,
	}
}

func notFoundMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		OneTimeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{
					Text: templates.WhatLinksButton,
				},
				{
					Text: templates.Skip,
				},
			},
		},
	}
}

func whatLinksMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		OneTimeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{
					Text: templates.ExampleTrack,
				},
				{
					Text: templates.SkipDemonstration,
				},
			},
		},
	}
}
