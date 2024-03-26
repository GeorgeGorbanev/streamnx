package vibeshare

import (
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"

	"github.com/tucnak/telebot"
)

func convertTrackMenu(track *music.Track) (*telebot.ReplyMarkup, error) {
	buttonsParams := make([]convertParams, 0, len(music.Providers)-1)
	for _, provider := range music.Providers {
		if provider != track.Provider {
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
			Text: string(buttonParams.Target),
			Data: cbData.Marshal(),
		})
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			buttons,
		},
	}, nil
}

func convertAlbumMenu(album *music.Album) (*telebot.ReplyMarkup, error) {
	buttonsParams := make([]convertParams, 0, len(music.Providers)-1)
	for _, provider := range music.Providers {
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
			Text: string(buttonParams.Target),
			Data: cbData.Marshal(),
		})
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			buttons,
		},
	}, nil
}
