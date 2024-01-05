package songshift

import (
	"fmt"
	"log"
	"strings"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"

	"github.com/essentialkaos/translit/v2"
	"github.com/tucnak/telebot"
)

type Songshift struct {
	spotifyClient  *spotify.Client
	telegramRouter *telegram.Router
	telegramSender telegram.Sender
	ymusicClient   *ymusic.Client
}

type Input struct {
	SpotifyClient  *spotify.Client
	TelegramSender telegram.Sender
	YmusicClient   *ymusic.Client
}

func NewSongshift(input *Input) *Songshift {
	s := Songshift{
		spotifyClient:  input.SpotifyClient,
		telegramSender: input.TelegramSender,
		ymusicClient:   input.YmusicClient,
	}

	s.telegramRouter = s.makeRouter()

	return &s
}

func (s *Songshift) HandleText(inMsg *telebot.Message) {
	s.telegramRouter.RouteMessage(inMsg)
}

func (s *Songshift) makeRouter() *telegram.Router {
	router := telegram.NewRouter()

	router.Register(spotify.OpenTrackRe, s.spotifyTrack())
	router.Register(ymusic.TrackURLRegExp, s.yMusicTrack())
	router.RegisterNotFound(s.notFound())

	return router
}

func (s *Songshift) respond(inMsg *telebot.Message, text string) (*telebot.Message, error) {
	return s.telegramSender.Send(
		&telegram.Message{
			To:   inMsg.Sender,
			Text: text,
		},
	)
}

func (s *Songshift) yMusicSearch(spotifyTrack *spotify.Track) (*ymusic.Track, error) {
	artistName := strings.ToLower(spotifyTrack.Artists[0].Name)
	trackName := strings.ToLower(spotifyTrack.Name)

	yMusicTrack, err := s.ymusicClient.SearchTrack(artistName, trackName)
	if err != nil {
		return nil, fmt.Errorf("failed to find ymusic track: %w", err)
	}
	if yMusicTrack != nil {
		foundLowcased := strings.ToLower(yMusicTrack.Artists[0].Name)
		if artistName == foundLowcased {
			return yMusicTrack, nil
		}

		translited := translit.ICAO(foundLowcased)
		if artistName == translited {
			return yMusicTrack, nil
		}
	}

	return nil, nil
}

func (s *Songshift) spotifyTrack() telegram.HandlerFunc {
	return func(inMsg *telebot.Message) {
		trackID := spotify.DetectTrackID(inMsg.Text)
		spotifyTrack, err := s.spotifyClient.GetTrack(trackID)
		if err != nil {
			log.Printf("error fetching track: %s", err)
			return
		}
		if spotifyTrack == nil {
			outMsg, err := s.respond(inMsg, "track not found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		yMusicTrack, err := s.yMusicSearch(spotifyTrack)
		if err != nil {
			log.Printf("failed to search ymusic: %s", err)
			return
		}
		if yMusicTrack == nil {
			outMsg, err := s.respond(inMsg, "no ym track found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		outMsg, err := s.respond(inMsg, yMusicTrack.URL())
		if err != nil {
			log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
			return
		}
		log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
	}
}

func (s *Songshift) yMusicTrack() telegram.HandlerFunc {
	return func(inMsg *telebot.Message) {
		trackID := ymusic.ParseTrackID(inMsg.Text)
		yMusicTrack, err := s.ymusicClient.GetTrack(trackID)
		if err != nil {
			log.Printf("error fetching track: %s", err)
			return
		}
		if yMusicTrack == nil {
			outMsg, err := s.respond(inMsg, "track not found in yandex music")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		spotifyTrack, err := s.spotifyClient.SearchTrack(yMusicTrack.Artists[0].Name, yMusicTrack.Title)
		if err != nil {
			log.Printf("failed to search spotify: %s", err)
			return
		}
		if spotifyTrack == nil {
			outMsg, err := s.respond(inMsg, "no track found in spotify")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		outMsg, err := s.respond(inMsg, spotifyTrack.URL())
		if err != nil {
			log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
			return
		}
		log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
	}
}

func (s *Songshift) notFound() telegram.HandlerFunc {
	return func(inMsg *telebot.Message) {
		outMsg, err := s.respond(inMsg, "no track link found")
		if err != nil {
			log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
			return
		}
		log.Printf("Sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
	}
}
