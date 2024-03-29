package music

import (
	"slices"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	Yandex  Provider = "yandex"
	Spotify Provider = "spotify"
	Youtube Provider = "youtube"
)

type Provider string

var Providers = []Provider{
	Spotify,
	Yandex,
	Youtube,
}

func IsValidProvider(p Provider) bool {
	return slices.Contains(Providers, p)
}

func (p Provider) Name() string {
	return cases.Title(language.English).String(string(p))
}
