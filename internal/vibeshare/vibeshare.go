package vibeshare

import (
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/converter"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"

	"github.com/tucnak/telebot"
)

type Vibeshare struct {
	converter      *converter.Converter
	telegramRouter *telegram.Router
	telegramSender telegram.Sender
}

type Input struct {
	Converter      *converter.Converter
	SpotifyClient  *spotify.HTTPClient
	TelegramSender telegram.Sender
	YandexClient   *yandex.HTTPClient
}

func NewVibeshare(input *Input) *Vibeshare {
	vs := Vibeshare{
		converter:      input.Converter,
		telegramSender: input.TelegramSender,
	}
	vs.telegramRouter = vs.makeRouter()
	return &vs
}

func (vs *Vibeshare) HandleText(inMsg *telebot.Message) {
	vs.telegramRouter.RouteMessage(inMsg)
}

func (vs *Vibeshare) makeRouter() *telegram.Router {
	router := telegram.NewRouter()

	router.Register(spotify.TrackRe, vs.spotifyTrackHandler)
	router.Register(spotify.AlbumRe, vs.spotifyAlbumHandler)
	router.Register(yandex.TrackRe, vs.yandexTrackHandler)
	router.Register(yandex.AlbumRe, vs.yandexAlbumHandler)

	router.RegisterNotFound(vs.notFoundHandler)

	return router
}

func (vs *Vibeshare) respond(inMsg *telebot.Message, text string) {
	response := &telegram.Message{
		To:   inMsg.Sender,
		Text: text,
	}

	_, err := vs.telegramSender.Send(response)
	if err != nil {
		slog.Error("failed to send message", slog.String("error", err.Error()))
		return
	}
	slog.Info("sent message", slog.Int("to", inMsg.Sender.ID), slog.String("text", text))
}

func (vs *Vibeshare) spotifyTrackHandler(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertTrack(inMsg.Text, converter.Spotify, converter.Yandex)
	if err != nil {
		slog.Error("failed to convert track", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.respond(inMsg, "failed to convert")
		return
	}

	vs.respond(inMsg, link)
}

func (vs *Vibeshare) spotifyAlbumHandler(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertAlbum(inMsg.Text, converter.Spotify, converter.Yandex)
	if err != nil {
		slog.Error("failed to convert album", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.respond(inMsg, "failed to convert")
		return
	}

	vs.respond(inMsg, link)
}

func (vs *Vibeshare) yandexTrackHandler(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertTrack(inMsg.Text, converter.Yandex, converter.Spotify)
	if err != nil {
		slog.Error("failed to convert track", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.respond(inMsg, "failed to convert")
		return
	}

	vs.respond(inMsg, link)
}

func (vs *Vibeshare) yandexAlbumHandler(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertAlbum(inMsg.Text, converter.Yandex, converter.Spotify)
	if err != nil {
		slog.Error("failed to convert album", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.respond(inMsg, "failed to convert")
		return
	}

	vs.respond(inMsg, link)
}

func (vs *Vibeshare) notFoundHandler(inMsg *telebot.Message) {
	vs.respond(inMsg, "no link found")
}
