package vibeshare

import (
	"log/slog"
	"regexp"

	"github.com/GeorgeGorbanev/vibeshare/internal/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/templates"
	"github.com/GeorgeGorbanev/vibeshare/internal/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/youtube"

	"github.com/tucnak/telebot"
)

const (
	convertTrackCallbackRoute = "cnvtr"
	convertAlbumCallbackRoute = "cnval"
	trackRegionCallbackRoute  = "regtr"
	albumRegionCallbackRoute  = "regal"
)

var (
	startCommand    = regexp.MustCompile("/start")
	feedbackCommand = regexp.MustCompile("/feedback")
)

func (vs *Vibeshare) TextHandler(inMsg *telebot.Message) {
	slog.Info("handling text message",
		slog.String("from", inMsg.Sender.Username),
		slog.String("text", inMsg.Text))
	vs.vibeshareRouter.RouteText(inMsg)
}

func (vs *Vibeshare) CallbackHandler(cb *telebot.Callback) {
	slog.Info("handling callback",
		slog.String("from", cb.Sender.Username),
		slog.String("data", cb.Data))
	vs.vibeshareRouter.RouteCallback(cb)
}

func (vs *Vibeshare) FeedbackTextHandler(inMsg *telebot.Message) {
	slog.Info("handling feedback message",
		slog.String("from", inMsg.Sender.Username),
		slog.String("text", inMsg.Text))
	vs.feedbackRouter.RouteText(inMsg)
}

func (vs *Vibeshare) setupVibeshareRouter() {
	vs.vibeshareRouter = &telegram.Router{
		TextRoutes: []*telegram.TextRoute{
			{Pattern: startCommand, Handler: vs.startCommand},
			{Pattern: feedbackCommand, Handler: vs.feedbackCommand},

			{Pattern: apple.TrackRe, Handler: vs.appleTrackLink},
			{Pattern: apple.AlbumRe, Handler: vs.appleAlbumLink},

			{Pattern: spotify.TrackRe, Handler: vs.spotifyTrackLink},
			{Pattern: spotify.AlbumRe, Handler: vs.spotifyAlbumLink},

			{Pattern: yandex.TrackRe, Handler: vs.yandexTrackLink},
			{Pattern: yandex.AlbumRe, Handler: vs.yandexAlbumLink},

			{Pattern: youtube.VideoRe, Handler: vs.youtubeTrackLink},
			{Pattern: youtube.PlaylistRe, Handler: vs.youtubeAlbumLink},

			{Pattern: templates.WhatLinksButtonRe, Handler: vs.whatLinks},
			{Pattern: templates.SkipRe, Handler: vs.skip},
			{Pattern: templates.SkipDemonstrationRe, Handler: vs.skip},
		},
		CallbackRoutes: []*telegram.CallbackRoute{
			{Address: convertTrackCallbackRoute, Handler: vs.convertTrack},
			{Address: convertAlbumCallbackRoute, Handler: vs.convertAlbum},
			{Address: trackRegionCallbackRoute, Handler: vs.trackRegion},
			{Address: albumRegionCallbackRoute, Handler: vs.albumRegion},
		},
		TextNotFound:     vs.textNotFoundHandler,
		CallbackNotFound: vs.callbackNotFoundHandler,
	}
}

func (vs *Vibeshare) setupFeedbackRouter() {
	vs.feedbackRouter = &telegram.Router{
		TextRoutes: []*telegram.TextRoute{
			{Pattern: startCommand, Handler: vs.feedbackStart},
		},
		TextNotFound: vs.feedback,
	}
}

func (vs *Vibeshare) send(messages ...*telegram.Message) {
	for _, message := range messages {
		_, err := vs.vibeshareSender.Send(message)
		if err != nil {
			slog.Error("failed to send message",
				slog.Any("error", err))
			return
		}
		slog.Info("sent message",
			slog.String("to", message.To.Recipient()),
			slog.String("text", message.Text))
	}
}
