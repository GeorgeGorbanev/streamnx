package vibeshare

import (
	"fmt"
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/tucnak/telebot"
)

func (vs *Vibeshare) CallbackHandler(cb *telebot.Callback) {
	slog.Info("handling callback",
		slog.String("from", cb.Sender.Username),
		slog.String("data", cb.Data))
	vs.telegramRouter.RouteCallback(cb)
}

func (vs *Vibeshare) convertTrack(callback *telegram.Callback) {
	params := convertParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	sourceTrack, err := vs.musicRegistry.Adapter(params.Source).GetTrack(params.ID)
	if err != nil {
		slog.Error("failed to search track", slog.Any("error", err))
		return
	}
	if sourceTrack == nil {
		slog.Error("source track not found")
		return
	}

	track, err := vs.musicRegistry.Adapter(params.Target).SearchTrack(sourceTrack.Artist, sourceTrack.Title)
	if err != nil {
		slog.Error("failed to search track", slog.Any("error", err))
		return
	}
	if track == nil {
		text := fmt.Sprintf("Track not found in %s", params.Target.Name)
		vs.respond(&telegram.Message{To: callback.Sender, Text: text})
		return
	}

	vs.respond(&telegram.Message{To: callback.Sender, Text: track.URL})
}

func (vs *Vibeshare) convertAlbum(callback *telegram.Callback) {
	params := convertParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	sourceAlbum, err := vs.musicRegistry.Adapter(params.Source).GetAlbum(params.ID)
	if err != nil {
		slog.Error("failed to search album", slog.Any("error", err))
		return
	}
	if sourceAlbum == nil {
		slog.Error("source album not found")
		return
	}

	album, err := vs.musicRegistry.Adapter(params.Target).SearchAlbum(sourceAlbum.Artist, sourceAlbum.Title)
	if err != nil {
		slog.Error("failed to search album", slog.Any("error", err))
		return
	}
	if album == nil {
		text := fmt.Sprintf("Album not found in %s", params.Target.Name)
		vs.respond(&telegram.Message{To: callback.Sender, Text: text})
		return
	}

	vs.respond(&telegram.Message{To: callback.Sender, Text: album.URL})
}

func (vs *Vibeshare) callbackNotFoundHandler(callback *telegram.Callback) {
	slog.Warn("callback not found", slog.String("data", callback.Data.Route))
}
