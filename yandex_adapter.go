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

func (a *YandexAdapter) SearchTrack(ctx context.Context, artist, title string) (*Entity, error) {
	foundTrack, err := a.findTrack(ctx, artist, title)
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

func (a *YandexAdapter) SearchAlbum(ctx context.Context, artist, title string) (*Entity, error) {
	foundAlbum, err := a.findAlbum(ctx, artist, title)
	if err != nil {
		if errors.Is(err, yandex.NotFoundError) {
			return nil, EntityNotFoundError
		}
		return nil, err
	}

	return a.adaptAlbum(foundAlbum), nil
}

func (a *YandexAdapter) findTrack(ctx context.Context, artist, title string) (*yandex.Track, error) {
	track, err := a.searchTrackRequest(ctx, artist, title)
	if err != nil && !errors.Is(err, yandex.NotFoundError) {
		return nil, fmt.Errorf("error searching track: %w", err)
	}
	if track != nil {
		artistMatch, err := a.artistMatch(ctx, track.Artists[0].Name, artist)
		if err != nil {
			return nil, fmt.Errorf("failed to check artist match: %w", err)
		}
		if artistMatch {
			return track, nil
		}
	}

	if translator.HasCyrillic(title) {
		translited := translator.TranslitLatToCyr(artist)
		track, err = a.searchTrackRequest(ctx, translited, title)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex track: %w", err)
		}
		if track != nil {
			return track, nil
		}
	}

	return nil, yandex.NotFoundError
}

func (a *YandexAdapter) findAlbum(ctx context.Context, artist, title string) (*yandex.Album, error) {
	album, err := a.searchAlbumRequest(ctx, artist, title)
	if err != nil && !errors.Is(err, yandex.NotFoundError) {
		return nil, fmt.Errorf("error searching album: %w", err)
	}
	if album != nil {
		artistMatch, err := a.artistMatch(ctx, album.Artists[0].Name, artist)
		if err != nil {
			return nil, fmt.Errorf("failed to check artist match: %w", err)
		}
		if artistMatch {
			return album, nil
		}
	}

	if translator.HasCyrillic(title) {
		translited := translator.TranslitLatToCyr(artist)
		album, err = a.searchAlbumRequest(ctx, translited, title)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex album: %w", err)
		}
		if album != nil {
			return album, nil
		}
	}

	return nil, yandex.NotFoundError
}

func (a *YandexAdapter) searchTrackRequest(ctx context.Context, artist, title string) (*yandex.Track, error) {
	query := a.prepareQuery(artist, title)
	return a.client.SearchTrack(ctx, query)
}

func (a *YandexAdapter) searchAlbumRequest(ctx context.Context, artist, title string) (*yandex.Album, error) {
	query := a.prepareQuery(artist, title)
	return a.client.SearchAlbum(ctx, query)
}

func (a *YandexAdapter) prepareQuery(artist, title string) string {
	query := entityFullTitle(artist, title)
	return strings.ToLower(query)
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

func (a *YandexAdapter) artistMatch(ctx context.Context, found, query string) (bool, error) {
	lcFound := strings.ToLower(found)
	lcQuery := strings.ToLower(query)
	if lcQuery == lcFound {
		return true, nil
	}

	translitedFoundArtist := translator.TranslitCyrToLat(lcFound)
	if lcQuery == translitedFoundArtist {
		return true, nil
	}

	if translator.HasCyrillic(found) {
		translatedArtist, err := a.translator.TranslateEnToRu(ctx, lcQuery)
		if err != nil {
			return false, fmt.Errorf("failed to translate artist name: %w", err)
		}
		if translatedArtist == lcFound {
			return true, nil
		}
	}

	return false, nil
}
