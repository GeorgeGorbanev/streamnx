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

func (vs *Vibeshare) TextHandler(inMsg *telebot.Message) {
	vs.telegramRouter.RouteText(inMsg)
}

func (vs *Vibeshare) CallbackHandler(callback *telebot.Callback) {
	vs.telegramRouter.RouteCallback(callback)
}

func (vs *Vibeshare) makeRouter() *telegram.Router {
	router := telegram.NewRouter()

	router.HandleText(spotify.TrackRe, vs.spotifyTrackLink)
	router.HandleText(spotify.AlbumRe, vs.spotifyAlbumLink)
	router.HandleText(yandex.TrackRe, vs.yandexTrackLink)
	router.HandleText(yandex.AlbumRe, vs.yandexAlbumLink)
	router.HandleTextNotFound(vs.textNotFoundHandler)

	router.HandleCallbackNotFound(vs.callbackNotFoundHandler)

	return router
}

func (vs *Vibeshare) send(response *telegram.Message) {
	_, err := vs.telegramSender.Send(response)
	if err != nil {
		slog.Error("failed to send message", slog.String("error", err.Error()))
		return
	}
	slog.Info("sent message", slog.String("to", response.To.Recipient()), slog.String("text", response.Text))
}

func (vs *Vibeshare) spotifyTrackLink(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertTrack(inMsg.Text, converter.Spotify, converter.Yandex)
	if err != nil {
		slog.Error("failed to convert track", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "failed to convert"})
		return
	}

	vs.send(&telegram.Message{To: inMsg.Sender, Text: link})
}

func (vs *Vibeshare) spotifyAlbumLink(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertAlbum(inMsg.Text, converter.Spotify, converter.Yandex)
	if err != nil {
		slog.Error("failed to convert album", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "failed to convert"})
		return
	}

	vs.send(&telegram.Message{To: inMsg.Sender, Text: link})
}

func (vs *Vibeshare) yandexTrackLink(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertTrack(inMsg.Text, converter.Yandex, converter.Spotify)
	if err != nil {
		slog.Error("failed to convert track", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "failed to convert"})
		return
	}

	vs.send(&telegram.Message{To: inMsg.Sender, Text: link})
}

func (vs *Vibeshare) yandexAlbumLink(inMsg *telebot.Message) {
	link, err := vs.converter.ConvertAlbum(inMsg.Text, converter.Yandex, converter.Spotify)
	if err != nil {
		slog.Error("failed to convert album", slog.String("error", err.Error()))
		return
	}
	if link == "" {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "failed to convert"})
		return
	}

	vs.send(&telegram.Message{To: inMsg.Sender, Text: link})
}

func (vs *Vibeshare) textNotFoundHandler(inMsg *telebot.Message) {
	vs.send(&telegram.Message{To: inMsg.Sender, Text: "no link found"})
}

func (vs *Vibeshare) callbackNotFoundHandler(callback *telegram.Callback) {
	slog.Warn("callback not found", slog.String("data", callback.Data.Command))
}
