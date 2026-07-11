package release

import (
	"errors"
	"time"
)

type Track struct {
	ID          string
	Title       string
	Artist      string
	AlbumID     string
	AlbumTitle  string
	URL         string
	CoverURL    string
	Duration    int
	ReleaseDate Date
	Provider    Provider
	Creator     string
	Description string
}

type Album struct {
	ID          string
	Title       string
	Artist      string
	Label       string
	URL         string
	CoverURL    string
	ReleaseDate Date
	Provider    Provider
	Creator     string
	Description string
	TrackIDs    []string
}

type SearchTrack struct {
	ID          string
	Title       string
	Artist      string
	AlbumID     string
	AlbumTitle  string
	URL         string
	CoverURL    string
	Provider    Provider
	Creator     string
	Description string
}

type SearchAlbum struct {
	ID          string
	Title       string
	Artist      string
	URL         string
	CoverURL    string
	Provider    Provider
	Creator     string
	Description string
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

type Type string

const (
	TypeTrack Type = "track"
	TypeAlbum Type = "album"
	TypeCloak Type = "cloak"
)

type Provider string

const (
	Apple      Provider = "ap"
	Bandcamp   Provider = "bc"
	Deezer     Provider = "dz"
	Spotify    Provider = "sf"
	Soundcloud Provider = "sc"
	Yandex     Provider = "ya"
	Youtube    Provider = "yt"
)

var Providers = []Provider{
	Apple,
	Bandcamp,
	Deezer,
	Spotify,
	Soundcloud,
	Yandex,
	Youtube,
}

var (
	ErrNotFound             = errors.New("release not found")
	ErrUnsupportedOperation = errors.New("unsupported operation")
)
