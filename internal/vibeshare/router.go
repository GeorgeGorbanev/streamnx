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
		TextHandlers: []*telegram.TextHandler{
			{Re: startCommand, HandlerFunc: vs.startCommand},
			{Re: feedbackCommand, HandlerFunc: vs.feedbackCommand},

			{Re: apple.TrackRe, HandlerFunc: vs.appleTrackLink},
			{Re: apple.AlbumRe, HandlerFunc: vs.appleAlbumLink},

			{Re: spotify.TrackRe, HandlerFunc: vs.spotifyTrackLink},
			{Re: spotify.AlbumRe, HandlerFunc: vs.spotifyAlbumLink},

			{Re: yandex.TrackRe, HandlerFunc: vs.yandexTrackLink},
			{Re: yandex.AlbumRe, HandlerFunc: vs.yandexAlbumLink},

			{Re: youtube.VideoRe, HandlerFunc: vs.youtubeTrackLink},
			{Re: youtube.PlaylistRe, HandlerFunc: vs.youtubeAlbumLink},

			{Re: templates.WhatLinksButtonRe, HandlerFunc: vs.whatLinks},
			{Re: templates.SkipRe, HandlerFunc: vs.skip},
			{Re: templates.SkipDemonstrationRe, HandlerFunc: vs.skip},
		},
		CallbackHandlers: []*telegram.CallbackHandler{
			{Route: convertTrackCallbackRoute, HandlerFunc: vs.convertTrack},
			{Route: convertAlbumCallbackRoute, HandlerFunc: vs.convertAlbum},
			{Route: trackRegionCallbackRoute, HandlerFunc: vs.trackRegion},
			{Route: albumRegionCallbackRoute, HandlerFunc: vs.albumRegion},
		},
		TextNotFoundHandler:     vs.textNotFoundHandler,
		CallbackHandlerNotFound: vs.callbackNotFoundHandler,
	}
}

func (vs *Vibeshare) setupFeedbackRouter() {
	vs.feedbackRouter = &telegram.Router{
		TextHandlers: []*telegram.TextHandler{
			{Re: startCommand, HandlerFunc: vs.feedbackStart},
		},
		TextNotFoundHandler: vs.feedback,
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
