package vibeshare

import (
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/streaminx"
	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
)

type regionParams struct {
	EntityID string
	Region   string
}

func (rp *regionParams) marshal() []string {
	return []string{
		rp.EntityID,
		rp.Region,
	}
}

func (rp *regionParams) unmarshal(s []string) error {
	if len(s) != 2 {
		return fmt.Errorf("invalid spicfy locale params: %s", s)
	}

	if !slices.Contains(streaminx.Yandex.Regions, s[1]) {
		return fmt.Errorf("invalid locale: %s", s[1])
	}

	rp.EntityID = s[0]
	rp.Region = s[1]

	return nil
}

func (vs *Vibeshare) trackRegion(callback *telegram.Callback) {
	params := regionParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	track, err := vs.streaminxRegistry.Adapter(streaminx.Yandex).GetTrack(params.EntityID)
	if err != nil {
		slog.Error("failed to search track", slog.Any("error", err))
		return
	}
	if track == nil {
		slog.Error("track not found")
		return
	}

	link := changeDomain(track.URL, params.Region)
	response := &telegram.Message{To: callback.Sender, Text: link}
	vs.send(response)
}

func (vs *Vibeshare) albumRegion(callback *telegram.Callback) {
	params := regionParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	album, err := vs.streaminxRegistry.Adapter(streaminx.Yandex).GetAlbum(params.EntityID)
	if err != nil {
		slog.Error("failed to search album", slog.Any("error", err))
		return
	}
	if album == nil {
		slog.Error("album not found")
		return
	}

	link := changeDomain(album.URL, params.Region)
	response := &telegram.Message{To: callback.Sender, Text: link}
	vs.send(response)
}

func changeDomain(link, domain string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}

	hostParts := strings.Split(u.Host, ".")
	hostParts[len(hostParts)-1] = domain
	u.Host = strings.Join(hostParts, ".")

	return u.String()
}
