package streamnx

import (
	"github.com/GeorgeGorbanev/streamnx/internal/spotify"
	"github.com/GeorgeGorbanev/streamnx/internal/translator"
)

type Credentials struct {
	GoogleTranslatorAPIKeyJSON string
	GoogleTranslatorProjectID  string
	YoutubeAPIKey              string
	SpotifyClientID            string
	SpotifyClientSecret        string
}

func (c Credentials) google() *translator.GoogleCredentials {
	return &translator.GoogleCredentials{
		APIKeyJSON: c.GoogleTranslatorAPIKeyJSON,
		ProjectID:  c.GoogleTranslatorProjectID,
	}
}

func (c Credentials) spotify() *spotify.Credentials {
	return &spotify.Credentials{
		ClientID:     c.SpotifyClientID,
		ClientSecret: c.SpotifyClientSecret,
	}
}
