package vibeshare

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/templates"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/youtube"
)

const (
	convertTrackCallbackRoute = "cnvtr"
	convertAlbumCallbackRoute = "cnval"
)

type convertParams struct {
	ID     string
	Source *music.Provider
	Target *music.Provider
}

var (
	startCommand    = regexp.MustCompile("/start")
	feedbackCommand = regexp.MustCompile("/feedback")
)

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

func (vs *Vibeshare) send(response *telegram.Message) {
	_, err := vs.vibeshareSender.Send(response)
	if err != nil {
		slog.Error("failed to send message",
			slog.Any("error", err))
		return
	}
	slog.Info("sent message",
		slog.String("to", response.To.Recipient()),
		slog.String("text", response.Text))
}

func (p *convertParams) marshal() []string {
	return []string{
		p.Source.Code,
		p.ID,
		p.Target.Code,
	}
}

func (p *convertParams) unmarshal(s []string) error {
	if len(s) != 3 {
		return fmt.Errorf("invalid convert params: %s", s)
	}
	sourceProvider := music.FindProviderByCode(s[0])
	if sourceProvider == nil {
		return fmt.Errorf("invalid source provider: %s", s[0])
	}
	targetProvider := music.FindProviderByCode(s[2])
	if targetProvider == nil {
		return fmt.Errorf("invalid target provider: %s", s[2])
	}

	p.Source = sourceProvider
	p.ID = s[1]
	p.Target = targetProvider

	return nil
}
