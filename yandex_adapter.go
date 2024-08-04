package streamnx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/GeorgeGorbanev/streamnx/internal/translator"
	"github.com/GeorgeGorbanev/streamnx/internal/yandex"
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

func (a *YandexAdapter) FetchTrack(ctx context.Context, id string) (*Entity, error) {
	yandexTrack, err := a.client.FetchTrack(ctx, id)
	if err != nil {
		if errors.Is(err, yandex.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get track from yandex music: %w", err)
	}

	return a.adaptTrack(yandexTrack), nil
}

func (a *YandexAdapter) SearchTrack(ctx context.Context, artist, track string) (*Entity, error) {
	lowcasedArtist := strings.ToLower(artist)
	lowcasedTrack := strings.ToLower(track)

	foundTrack, err := a.findTrack(ctx, lowcasedArtist, lowcasedTrack)
	if err != nil {
		if errors.Is(err, yandex.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, err
	}

	return a.adaptTrack(foundTrack), nil
}

func (a *YandexAdapter) FetchAlbum(ctx context.Context, id string) (*Entity, error) {
	yandexAlbum, err := a.client.FetchAlbum(ctx, id)
	if err != nil {
		if errors.Is(err, yandex.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, fmt.Errorf("failed to get album from yandex music: %w", err)
	}

	return a.adaptAlbum(yandexAlbum), nil
}

func (a *YandexAdapter) SearchAlbum(ctx context.Context, artistName, albumName string) (*Entity, error) {
	lowcasedArtist := strings.ToLower(artistName)
	lowcasedAlbum := strings.ToLower(albumName)

	foundAlbum, err := a.findAlbum(ctx, lowcasedArtist, lowcasedAlbum)
	if err != nil {
		if errors.Is(err, yandex.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, err
	}

	return a.adaptAlbum(foundAlbum), nil
}

func (a *YandexAdapter) findTrack(ctx context.Context, artist, track string) (*yandex.Track, error) {
	foundTrack, err := a.client.SearchTrack(ctx, artist, track)
	if err != nil && !errors.Is(err, yandex.NotFoundError) {
		return nil, fmt.Errorf("error searching track: %w", err)
	}
	if foundTrack != nil {
		artistMatch, err := a.artistMatch(ctx, foundTrack.Artists[0].Name, artist)
		if err != nil {
			return nil, fmt.Errorf("failed to check artist match: %w", err)
		}
		if artistMatch {
			return foundTrack, nil
		}
	}

	if translator.HasCyrillic(track) {
		translited := translator.TranslitLatToCyr(artist)
		foundTranslitedTrack, err := a.client.SearchTrack(ctx, translited, track)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex track: %w", err)
		}
		if foundTranslitedTrack != nil {
			return foundTranslitedTrack, nil
		}
	}

	return nil, yandex.NotFoundError
}

func (a *YandexAdapter) findAlbum(ctx context.Context, artist, album string) (*yandex.Album, error) {
	foundAlbum, err := a.client.SearchAlbum(ctx, artist, album)
	if err != nil && !errors.Is(err, yandex.NotFoundError) {
		return nil, fmt.Errorf("error searching album: %w", err)
	}
	if foundAlbum != nil {
		artistMatch, err := a.artistMatch(ctx, foundAlbum.Artists[0].Name, artist)
		if err != nil {
			return nil, fmt.Errorf("failed to check artist match: %w", err)
		}
		if artistMatch {
			return foundAlbum, nil
		}
	}

	if translator.HasCyrillic(album) {
		translited := translator.TranslitLatToCyr(artist)
		foundTranslitedAlbum, err := a.client.SearchAlbum(ctx, translited, album)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex album: %w", err)
		}
		if foundTranslitedAlbum != nil {
			return foundTranslitedAlbum, nil
		}
	}

	return nil, yandex.NotFoundError
}

func (a *YandexAdapter) adaptTrack(yandexTrack *yandex.Track) *Entity {
	return &Entity{
		ID:       yandexTrack.IDString(),
		Title:    yandexTrack.Title,
		Artist:   yandexTrack.Artists[0].Name,
		URL:      yandexTrack.URL(),
		Provider: Yandex,
		Type:     Track,
	}
}

func (a *YandexAdapter) adaptAlbum(yandexAlbum *yandex.Album) *Entity {
	return &Entity{
		ID:       strconv.Itoa(yandexAlbum.ID),
		Title:    yandexAlbum.Title,
		Artist:   yandexAlbum.Artists[0].Name,
		URL:      yandexAlbum.URL(),
		Provider: Yandex,
		Type:     Album,
	}
}

func (a *YandexAdapter) artistMatch(ctx context.Context, foundArtist, artistName string) (bool, error) {
	lowcasedFoundArtist := strings.ToLower(foundArtist)
	if artistName == lowcasedFoundArtist {
		return true, nil
	}

	translitedFoundArtist := translator.TranslitCyrToLat(lowcasedFoundArtist)
	if artistName == translitedFoundArtist {
		return true, nil
	}

	if translator.HasCyrillic(foundArtist) {
		translatedArtist, err := a.translator.TranslateEnToRu(ctx, artistName)
		if err != nil {
			return false, fmt.Errorf("failed to translate artist name: %w", err)
		}
		if translatedArtist == lowcasedFoundArtist {
			return true, nil
		}
	}

	return false, nil
}
