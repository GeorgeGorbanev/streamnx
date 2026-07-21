package streamnx

import (
	"slices"

	"github.com/GeorgeGorbanev/streamnx/v2/internal/release"
)

type (
	Track           = release.Track
	Album           = release.Album
	SearchTrack     = release.SearchTrack
	SearchAlbum     = release.SearchAlbum
	ReleaseDate     = release.Date
	ReleaseType     = release.Type
	ReleaseProvider = release.Provider
)

const (
	ReleaseTypeTrack = release.TypeTrack
	ReleaseTypeAlbum = release.TypeAlbum
	ReleaseTypeCloak = release.TypeCloak

	Apple        = release.Apple
	Bandcamp     = release.Bandcamp
	Deezer       = release.Deezer
	Spotify      = release.Spotify
	Soundcloud   = release.Soundcloud
	Yandex       = release.Yandex
	Youtube      = release.Youtube
	YoutubeMusic = release.YoutubeMusic
)

var (
	ErrNotFound             = release.ErrNotFound
	ErrScrapingBlocked      = release.ErrScrapingBlocked
	ErrUnsupportedOperation = release.ErrUnsupportedOperation
)

func Providers() []ReleaseProvider {
	return slices.Clone(release.Providers)
}
