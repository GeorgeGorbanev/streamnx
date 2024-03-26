package vibeshare

import (
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
)

const (
	convertTrackCallbackRoute = "convert_track"
	convertAlbumCallbackRoute = "convert_album"
)

type Vibeshare struct {
	musicRegistry  *music.Registry
	telegramRouter *telegram.Router
	telegramSender telegram.Sender
}

type Input struct {
	MusicRegistry  *music.Registry
	TelegramSender telegram.Sender
}

func NewVibeshare(input *Input) *Vibeshare {
	vs := Vibeshare{
		musicRegistry:  input.MusicRegistry,
		telegramSender: input.TelegramSender,
	}
	vs.telegramRouter = vs.makeRouter()
	return &vs
}

func (vs *Vibeshare) makeRouter() *telegram.Router {
	router := telegram.NewRouter()

	router.HandleText(spotify.TrackRe, vs.spotifyTrackLink)
	router.HandleText(spotify.AlbumRe, vs.spotifyAlbumLink)
	router.HandleText(yandex.TrackRe, vs.yandexTrackLink)
	router.HandleText(yandex.AlbumRe, vs.yandexAlbumLink)
	router.HandleTextNotFound(vs.textNotFoundHandler)

	router.HandleCallback(convertTrackCallbackRoute, vs.convertTrack)
	router.HandleCallback(convertAlbumCallbackRoute, vs.convertAlbum)
	router.HandleCallbackNotFound(vs.callbackNotFoundHandler)

	return router
}

func (vs *Vibeshare) respond(response *telegram.Message) {
	_, err := vs.telegramSender.Send(response)
	if err != nil {
		slog.Error("failed to send message", slog.Any("error", err))
		return
	}
	slog.Info("sent message",
		slog.String("to", response.To.Recipient()),
		slog.String("text", response.Text))
}
