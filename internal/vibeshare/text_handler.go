package vibeshare

import (
	"fmt"
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/templates"
	"github.com/tucnak/telebot"
)

func (vs *Vibeshare) TextHandler(inMsg *telebot.Message) {
	slog.Info("handling text message",
		slog.String("from", inMsg.Sender.Username),
		slog.String("text", inMsg.Text))
	vs.vibeshareRouter.RouteText(inMsg)
}

func (vs *Vibeshare) startCommand(inMsg *telebot.Message) {
	vs.send(&telegram.Message{To: inMsg.Sender, Text: templates.Start})
}

func (vs *Vibeshare) feedbackCommand(inMsg *telebot.Message) {
	text := fmt.Sprintf(templates.FeedbackCommand, vs.feedbackBotName)
	vs.send(&telegram.Message{To: inMsg.Sender, Text: text})
}

func (vs *Vibeshare) whatLinks(inMsg *telebot.Message) {
	vs.send(&telegram.Message{To: inMsg.Sender, Text: templates.WhatLinksResponse, ReplyMarkup: whatLinksMenu()})
}

func (vs *Vibeshare) skip(_ *telebot.Message) {
}

func (vs *Vibeshare) appleTrackLink(inMsg *telebot.Message) {
	vs.trackLink(music.Apple, inMsg)
}

func (vs *Vibeshare) appleAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(music.Apple, inMsg)
}

func (vs *Vibeshare) spotifyTrackLink(inMsg *telebot.Message) {
	vs.trackLink(music.Spotify, inMsg)
}

func (vs *Vibeshare) spotifyAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(music.Spotify, inMsg)
}

func (vs *Vibeshare) yandexTrackLink(inMsg *telebot.Message) {
	vs.trackLink(music.Yandex, inMsg)
}

func (vs *Vibeshare) yandexAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(music.Yandex, inMsg)
}

func (vs *Vibeshare) youtubeTrackLink(inMsg *telebot.Message) {
	vs.trackLink(music.Youtube, inMsg)
}

func (vs *Vibeshare) youtubeAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(music.Youtube, inMsg)
}

func (vs *Vibeshare) textNotFoundHandler(inMsg *telebot.Message) {
	vs.send(&telegram.Message{To: inMsg.Sender, Text: templates.NotFound, ReplyMarkup: notFoundMenu()})
}

func (vs *Vibeshare) trackLink(provider *music.Provider, inMsg *telebot.Message) {
	trackID := vs.musicRegistry.Adapter(provider).DetectTrackID(inMsg.Text)
	track, err := vs.musicRegistry.Adapter(provider).GetTrack(trackID)
	if err != nil {
		slog.Error("error fetching track", slog.Any("error", err))
		return
	}
	if track == nil {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "Link is invalid"})
		return
	}

	menu, err := convertTrackMenu(track)
	if err != nil {
		slog.Error("error creating menu", slog.Any("error", err))
		return
	}

	vs.send(&telegram.Message{To: inMsg.Sender, Text: "Select target link provider", ReplyMarkup: menu})
}

func (vs *Vibeshare) albumLink(provider *music.Provider, inMsg *telebot.Message) {
	albumID := vs.musicRegistry.Adapter(provider).DetectAlbumID(inMsg.Text)
	album, err := vs.musicRegistry.Adapter(provider).GetAlbum(albumID)
	if err != nil {
		slog.Error("error fetching album", slog.Any("error", err))
		return
	}
	if album == nil {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "Link is invalid"})
		return
	}

	menu, err := convertAlbumMenu(album)
	if err != nil {
		slog.Error("error creating menu", slog.Any("error", err))
		return
	}

	vs.send(&telegram.Message{To: inMsg.Sender, Text: "Select target link provider", ReplyMarkup: menu})
}
