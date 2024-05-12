package vibeshare

import (
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
)

func (vs *Vibeshare) callbackNotFoundHandler(callback *telegram.Callback) {
	slog.Warn("callback not found", slog.String("data", callback.Data.Route))
	vs.send(&telegram.Message{
		To:   callback.Sender,
		Text: "This button not available anymore. Try to request a new button 🙏",
	})
}
