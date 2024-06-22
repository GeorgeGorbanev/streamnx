package vibeshare

import (
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/streaminx"
	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"

	"github.com/tucnak/telebot"
)

func (vs *Vibeshare) appleTrackLink(inMsg *telebot.Message) {
	vs.trackLink(streaminx.Apple, inMsg)
}

func (vs *Vibeshare) appleAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(streaminx.Apple, inMsg)
}

func (vs *Vibeshare) spotifyTrackLink(inMsg *telebot.Message) {
	vs.trackLink(streaminx.Spotify, inMsg)
}

func (vs *Vibeshare) spotifyAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(streaminx.Spotify, inMsg)
}

func (vs *Vibeshare) yandexTrackLink(inMsg *telebot.Message) {
	vs.trackLink(streaminx.Yandex, inMsg)
}

func (vs *Vibeshare) yandexAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(streaminx.Yandex, inMsg)
}

func (vs *Vibeshare) youtubeTrackLink(inMsg *telebot.Message) {
	vs.trackLink(streaminx.Youtube, inMsg)
}

func (vs *Vibeshare) youtubeAlbumLink(inMsg *telebot.Message) {
	vs.albumLink(streaminx.Youtube, inMsg)
}

func (vs *Vibeshare) trackLink(provider *streaminx.Provider, inMsg *telebot.Message) {
	trackID, err := vs.streaminxRegistry.Adapter(provider).DetectTrackID(inMsg.Text)
	if err != nil {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "Link is invalid"})
		return
	}

	track, err := vs.streaminxRegistry.Adapter(provider).GetTrack(trackID)
	if err != nil {
		slog.Error("error fetching track", slog.Any("error", err))
		return
	}
	if track == nil {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "Link is invalid"})
		return
	}

	vs.send(&telegram.Message{
		To:          inMsg.Sender,
		Text:        "Select target link provider",
		ReplyMarkup: convertTrackMenu(track),
	})
}

func (vs *Vibeshare) albumLink(provider *streaminx.Provider, inMsg *telebot.Message) {
	albumID, err := vs.streaminxRegistry.Adapter(provider).DetectAlbumID(inMsg.Text)
	if err != nil {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "Link is invalid"})
		return
	}
	album, err := vs.streaminxRegistry.Adapter(provider).GetAlbum(albumID)
	if err != nil {
		slog.Error("error fetching album", slog.Any("error", err))
		return
	}
	if album == nil {
		vs.send(&telegram.Message{To: inMsg.Sender, Text: "Link is invalid"})
		return
	}

	vs.send(&telegram.Message{
		To:          inMsg.Sender,
		Text:        "Select target link provider",
		ReplyMarkup: convertAlbumMenu(album),
	})
}
