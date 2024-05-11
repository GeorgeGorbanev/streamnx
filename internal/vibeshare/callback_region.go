package vibeshare

import (
	"fmt"
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/yandex"
)

type regionParams struct {
	EntityID string
	Region   *yandex.Region
}

func (rp *regionParams) marshal() []string {
	return []string{
		rp.EntityID,
		rp.Region.DomainZone,
	}
}

func (rp *regionParams) unmarshal(s []string) error {
	if len(s) != 2 {
		return fmt.Errorf("invalid spicfy locale params: %s", s)
	}

	locale := yandex.FindRegionByDomainZone(s[1])
	if !locale.IsValid() {
		return fmt.Errorf("invalid locale: %s", s[1])
	}

	rp.EntityID = s[0]
	rp.Region = locale

	return nil
}

func (vs *Vibeshare) trackRegion(callback *telegram.Callback) {
	params := regionParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	track, err := vs.musicRegistry.Adapter(music.Yandex).GetTrack(params.EntityID)
	if err != nil {
		slog.Error("failed to search track", slog.Any("error", err))
		return
	}
	if track == nil {
		slog.Error("track not found")
		return
	}

	link := params.Region.LocalizeLink(track.URL)
	response := &telegram.Message{To: callback.Sender, Text: link}
	vs.send(response)
}

func (vs *Vibeshare) albumRegion(callback *telegram.Callback) {
	params := regionParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	album, err := vs.musicRegistry.Adapter(music.Yandex).GetAlbum(params.EntityID)
	if err != nil {
		slog.Error("failed to search album", slog.Any("error", err))
		return
	}
	if album == nil {
		slog.Error("album not found")
		return
	}

	link := params.Region.LocalizeLink(album.URL)
	response := &telegram.Message{To: callback.Sender, Text: link}
	vs.send(response)
}
