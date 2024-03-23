package converter

import (
	"fmt"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/translit"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/ymusic"
)

type Converter struct {
	spotifyClient *spotify.Client
	yandexClient  *ymusic.Client
}

type Input struct {
	SpotifyClient *spotify.Client
	YandexClient  *ymusic.Client
}

func NewConverter(input *Input) *Converter {
	return &Converter{
		spotifyClient: input.SpotifyClient,
		yandexClient:  input.YandexClient,
	}
}

func (c *Converter) SpotifyTrackToYandex(link string) (string, error) {
	trackID := spotify.DetectTrackID(link)
	spotifyTrack, err := c.spotifyClient.GetTrack(trackID)
	if err != nil {
		return "", fmt.Errorf("error fetching track: %w", err)
	}
	if spotifyTrack == nil {
		return "", nil
	}

	yandexTrack, err := c.yandexTrackSearch(spotifyTrack)
	if err != nil {
		return "", fmt.Errorf("error searching yandex track: %w", err)
	}
	if yandexTrack == nil {
		return "", nil
	}

	return yandexTrack.URL(), nil
}

func (c *Converter) SpotifyAlbumToYandex(link string) (string, error) {
	albumID := spotify.DetectAlbumID(link)
	spotifyAlbum, err := c.spotifyClient.GetAlbum(albumID)
	if err != nil {
		return "", fmt.Errorf("error fetching album: %w", err)
	}
	if spotifyAlbum == nil {
		return "", nil
	}

	yandexAlbum, err := c.yandexClient.SearchAlbum(spotifyAlbum.Artists[0].Name, spotifyAlbum.Name)
	if err != nil {
		return "", fmt.Errorf("error searching album: %w", err)
	}

	if yandexAlbum == nil {
		return "", nil
	}

	return yandexAlbum.URL(), nil
}

func (c *Converter) YandexTrackToSpotify(link string) (string, error) {
	trackID := ymusic.DetectTrackID(link)
	yandexTrack, err := c.yandexClient.GetTrack(trackID)
	if err != nil {
		return "", fmt.Errorf("error fetching track: %w", err)
	}
	if yandexTrack == nil {
		return "", nil
	}

	spotifyTrack, err := c.spotifyClient.SearchTrack(yandexTrack.Artists[0].Name, yandexTrack.Title)
	if err != nil {
		return "", fmt.Errorf("error searching track: %w", err)
	}
	if spotifyTrack == nil {
		return "", nil
	}

	return spotifyTrack.URL(), nil
}

func (c *Converter) YandexAlbumToSpotify(link string) (string, error) {
	albumID := ymusic.DetectAlbumID(link)
	yandexAlbum, err := c.yandexClient.GetAlbum(albumID)
	if err != nil {
		return "", fmt.Errorf("error fetching album: %w", err)
	}
	if yandexAlbum == nil {
		return "", nil
	}

	spotifyAlbum, err := c.spotifyClient.SearchAlbum(yandexAlbum.Artists[0].Name, yandexAlbum.Title)
	if err != nil {
		return "", fmt.Errorf("error searching album: %w", err)
	}
	if spotifyAlbum == nil {
		return "", nil
	}

	return spotifyAlbum.URL(), nil
}

func (с *Converter) yandexTrackSearch(spotifyTrack *spotify.Track) (*ymusic.Track, error) {
	artistName := strings.ToLower(spotifyTrack.Artists[0].Name)
	trackName := strings.ToLower(spotifyTrack.Name)

	yandexTrack, err := с.yandexClient.SearchTrack(artistName, trackName)
	if err != nil {
		return nil, fmt.Errorf("error searching track: %w", err)
	}
	if yandexTrack != nil {
		foundLowcasedArtist := strings.ToLower(yandexTrack.Artists[0].Name)
		if artistName == foundLowcasedArtist {
			return yandexTrack, nil
		}

		translitedArtist := translit.CyrillicToLatin(foundLowcasedArtist)
		if artistName == translitedArtist {
			return yandexTrack, nil
		}
		return nil, nil
	}

	if spotifyTrack.NameContainsRussianLetters() {
		translitedArtist := translit.LatinToCyrillic(artistName)
		yandexTrack, err = с.yandexClient.SearchTrack(translitedArtist, trackName)
		if err != nil {
			return nil, fmt.Errorf("error searching yandex track: %w", err)
		}
	}

	return yandexTrack, nil
}
