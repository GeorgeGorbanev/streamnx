package vibeshare

import (
	"fmt"
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/streaminx"
	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/templates"
)

type convertParams struct {
	ID     string
	Source *streaminx.Provider
	Target *streaminx.Provider
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
	sourceProvider := streaminx.FindProviderByCode(s[0])
	if sourceProvider == nil {
		return fmt.Errorf("invalid source provider: %s", s[0])
	}
	targetProvider := streaminx.FindProviderByCode(s[2])
	if targetProvider == nil {
		return fmt.Errorf("invalid target provider: %s", s[2])
	}

	p.Source = sourceProvider
	p.ID = s[1]
	p.Target = targetProvider

	return nil
}

func (vs *Vibeshare) convertTrack(callback *telegram.Callback) {
	params := convertParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	sourceTrack, err := vs.streaminxRegistry.Adapter(params.Source).GetTrack(params.ID)
	if err != nil {
		slog.Error("failed to search track", slog.Any("error", err))
		return
	}
	if sourceTrack == nil {
		slog.Error("source track not found")
		return
	}

	track, err := vs.streaminxRegistry.Adapter(params.Target).SearchTrack(sourceTrack.Artist, sourceTrack.Title)
	if err != nil {
		slog.Error("failed to search track", slog.Any("error", err))
		return
	}

	response := []*telegram.Message{}

	if track == nil {
		text := fmt.Sprintf("Track not found in %s", params.Target.Name)
		response = append(response, &telegram.Message{To: callback.Sender, Text: text})
	} else {
		response = append(response, &telegram.Message{To: callback.Sender, Text: track.URL})

		if params.Target == streaminx.Yandex {
			response = append(response, &telegram.Message{
				To:          callback.Sender,
				Text:        templates.SpecifyRegion,
				ReplyMarkup: trackRegionMenu(track),
			})
		}
	}

	vs.send(response...)
}

func (vs *Vibeshare) convertAlbum(callback *telegram.Callback) {
	params := convertParams{}
	if err := params.unmarshal(callback.Data.Params); err != nil {
		slog.Error("failed to unmarshal params", slog.Any("error", err))
		return
	}

	sourceAlbum, err := vs.streaminxRegistry.Adapter(params.Source).GetAlbum(params.ID)
	if err != nil {
		slog.Error("failed to search album", slog.Any("error", err))
		return
	}
	if sourceAlbum == nil {
		slog.Error("source album not found")
		return
	}

	album, err := vs.streaminxRegistry.Adapter(params.Target).SearchAlbum(sourceAlbum.Artist, sourceAlbum.Title)
	if err != nil {
		slog.Error("failed to search album", slog.Any("error", err))
		return
	}

	response := []*telegram.Message{}

	if album == nil {
		text := fmt.Sprintf("Album not found in %s", params.Target.Name)
		response = append(response, &telegram.Message{To: callback.Sender, Text: text})
	} else {
		response = append(response, &telegram.Message{To: callback.Sender, Text: album.URL})

		if params.Target == streaminx.Yandex {
			response = append(response, &telegram.Message{
				To:          callback.Sender,
				Text:        templates.SpecifyRegion,
				ReplyMarkup: albumRegionMenu(album),
			})
		}
	}

	vs.send(response...)
}
