package music

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/translit"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
)

type YandexAdapter struct {
	client yandex.Client
}

func newYandexAdapter(client yandex.Client) *YandexAdapter {
	return &YandexAdapter{
		client: client,
	}
}

func (a *YandexAdapter) DetectTrackID(trackURL string) string {
	match := yandex.TrackRe.FindStringSubmatch(trackURL)
	if match == nil || len(match) < 3 {
		return ""
	}
	return match[2]
}

func (a *YandexAdapter) DetectAlbumID(albumURL string) string {
	match := yandex.AlbumRe.FindStringSubmatch(albumURL)
	if match == nil || len(match) < 3 {
		return ""
	}
	return match[2]
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

func (a *YandexAdapter) SearchTrack(artistName, trackName string) (*Track, error) {
	artist := strings.ToLower(artistName)
	track := strings.ToLower(trackName)

	yandexTrack, err := a.client.SearchTrack(artist, track)
	if err != nil {
		return nil, fmt.Errorf("error searching track: %w", err)
	}
	if yandexTrack != nil {
		foundLowcasedArtist := strings.ToLower(yandexTrack.Artists[0].Name)
		if artist == foundLowcasedArtist {
			return a.adaptTrack(yandexTrack), nil
		}

		translitedArtist := translit.CyrillicToLatin(artist)
		if artist == translitedArtist {
			return a.adaptTrack(yandexTrack), nil
		}
		return nil, nil
	}

	if translit.Translitable(artist) || translit.Translitable(track) {
		translitedArtist := translit.LatinToCyrillic(artist)
		yandexTrack, err = a.client.SearchTrack(translitedArtist, track)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex track: %w", err)
		}
	}
	if yandexTrack != nil {
		return a.adaptTrack(yandexTrack), nil
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
	yandexAlbum, err := a.client.SearchAlbum(artistName, albumName)
	if err != nil {
		return nil, fmt.Errorf("failed to search album on yandex music: %w", err)
	}
	if yandexAlbum == nil {
		return nil, nil
	}

	return a.adaptAlbum(yandexAlbum), nil
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
