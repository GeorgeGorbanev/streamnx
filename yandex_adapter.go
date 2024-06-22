package streaminx

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/GeorgeGorbanev/streaminx/internal/translator"
	"github.com/GeorgeGorbanev/streaminx/internal/yandex"
)

type YandexAdapter struct {
	client     yandex.Client
	translator translator.Translator
}

func newYandexAdapter(c yandex.Client, t translator.Translator) *YandexAdapter {
	return &YandexAdapter{
		client:     c,
		translator: t,
	}
}

func (a *YandexAdapter) DetectTrackID(trackURL string) (string, error) {
	match := yandex.TrackRe.FindStringSubmatch(trackURL)
	if match == nil || len(match) < 3 {
		return "", IDNotFoundError
	}
	return match[2], nil
}

func (a *YandexAdapter) DetectAlbumID(albumURL string) (string, error) {
	match := yandex.AlbumRe.FindStringSubmatch(albumURL)
	if match == nil || len(match) < 3 {
		return "", IDNotFoundError
	}
	return match[2], nil
}

func (a *YandexAdapter) GetTrack(id string) (*Track, error) {
	yandexTrack, err := a.client.GetTrack(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get track from yandex music: %w", err)
	}
	if yandexTrack == nil {
		return nil, nil
	}

	return a.adaptTrack(yandexTrack), nil
}

func (a *YandexAdapter) SearchTrack(artist, track string) (*Track, error) {
	lowcasedArtist := strings.ToLower(artist)
	lowcasedTrack := strings.ToLower(track)

	foundTrack, err := a.findTrack(lowcasedArtist, lowcasedTrack)
	if err != nil {
		return nil, err
	}
	if foundTrack != nil {
		return a.adaptTrack(foundTrack), nil
	}

	return nil, nil
}

func (a *YandexAdapter) GetAlbum(id string) (*Album, error) {
	yandexAlbum, err := a.client.GetAlbum(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get album from yandex music: %w", err)
	}
	if yandexAlbum == nil {
		return nil, nil
	}

	return a.adaptAlbum(yandexAlbum), nil
}

func (a *YandexAdapter) SearchAlbum(artistName, albumName string) (*Album, error) {
	lowcasedArtist := strings.ToLower(artistName)
	lowcasedAlbum := strings.ToLower(albumName)

	foundAlbum, err := a.findAlbum(lowcasedArtist, lowcasedAlbum)
	if err != nil {
		return nil, err
	}
	if foundAlbum != nil {
		return a.adaptAlbum(foundAlbum), nil
	}

	return nil, nil
}

func (a *YandexAdapter) findTrack(artist, track string) (*yandex.Track, error) {
	foundTrack, err := a.client.SearchTrack(artist, track)
	if err != nil {
		return nil, fmt.Errorf("error searching track: %w", err)
	}
	if foundTrack != nil {
		artistMatch, err := a.artistMatch(foundTrack.Artists[0].Name, artist)
		if err != nil {
			return nil, fmt.Errorf("failed to check artist match: %w", err)
		}
		if artistMatch {
			return foundTrack, nil
		}
	}

	if translator.HasCyrillic(track) {
		translited := translator.TranslitLatToCyr(artist)
		foundTranslitedTrack, err := a.client.SearchTrack(translited, track)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex track: %w", err)
		}
		if foundTranslitedTrack != nil {
			return foundTranslitedTrack, nil
		}
	}

	return nil, nil
}

func (a *YandexAdapter) findAlbum(artist, album string) (*yandex.Album, error) {
	foundAlbum, err := a.client.SearchAlbum(artist, album)
	if err != nil {
		return nil, fmt.Errorf("error searching album: %w", err)
	}
	if foundAlbum != nil {
		artistMatch, err := a.artistMatch(foundAlbum.Artists[0].Name, artist)
		if err != nil {
			return nil, fmt.Errorf("failed to check artist match: %w", err)
		}
		if artistMatch {
			return foundAlbum, nil
		}
	}

	if translator.HasCyrillic(album) {
		translited := translator.TranslitLatToCyr(artist)
		foundTranslitedAlbum, err := a.client.SearchAlbum(translited, album)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex album: %w", err)
		}
		if foundTranslitedAlbum != nil {
			return foundTranslitedAlbum, nil
		}
	}

	return nil, nil
}

func (a *YandexAdapter) adaptTrack(yandexTrack *yandex.Track) *Track {
	return &Track{
		ID:       yandexTrack.IDString(),
		Title:    yandexTrack.Title,
		Artist:   yandexTrack.Artists[0].Name,
		URL:      yandexTrack.URL(),
		Provider: Yandex,
	}
}

func (a *YandexAdapter) adaptAlbum(yandexAlbum *yandex.Album) *Album {
	return &Album{
		ID:       strconv.Itoa(yandexAlbum.ID),
		Title:    yandexAlbum.Title,
		Artist:   yandexAlbum.Artists[0].Name,
		URL:      yandexAlbum.URL(),
		Provider: Yandex,
	}
}

func (a *YandexAdapter) artistMatch(foundArtist, artistName string) (bool, error) {
	lowcasedFoundArtist := strings.ToLower(foundArtist)
	if artistName == lowcasedFoundArtist {
		return true, nil
	}

	translitedFoundArtist := translator.TranslitCyrToLat(lowcasedFoundArtist)
	if artistName == translitedFoundArtist {
		return true, nil
	}

	if translator.HasCyrillic(foundArtist) {
		translatedArtist, err := a.translator.TranslateEnToRu(context.Background(), artistName)
		if err != nil {
			return false, fmt.Errorf("failed to translate artist name: %w", err)
		}
		if translatedArtist == lowcasedFoundArtist {
			return true, nil
		}
	}

	return false, nil
}
