package vibeshare

import (
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
)

func (vs *Vibeshare) callbackNotFoundHandler(callback *telegram.Callback) {
	slog.Warn("callback not found", slog.String("data", callback.Data.Route))
}
