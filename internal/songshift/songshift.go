package songshift

import (
	"log"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
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

func (s *Songshift) spotifyTrack() telegram.HandlerFunc {
	return func(inMsg *telebot.Message) {
		trackID := spotify.DetectTrackID(inMsg.Text)
		track, err := s.spotifyClient.GetTrack(trackID)
		if err != nil {
			log.Printf("error fetching track: %s", err)
			return
		}
		if track == nil {
			outMsg, err := s.respond(inMsg, "track not found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		searchResponse, err := s.ymusicClient.Search(track.Title())
		if err != nil {
			log.Printf("failed to search ymusic: %s", err)
			return
		}

		if !searchResponse.Result.AnyTracksFound() {
			outMsg, err := s.respond(inMsg, "no ym track found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		trackURL := searchResponse.Result.Tracks.Results[0].URL()

		outMsg, err := s.respond(inMsg, trackURL)
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
		track, err := s.ymusicClient.GetTrack(trackID)
		if err != nil {
			log.Printf("error fetching track: %s", err)
			return
		}
		if track == nil {
			outMsg, err := s.respond(inMsg, "track not found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		outMsg, err := s.respond(inMsg, track.FullTitle())
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

func (s *Songshift) respond(inMsg *telebot.Message, text string) (*telebot.Message, error) {
	return s.telegramSender.Send(
		&telegram.Message{
			To:   inMsg.Sender,
			Text: text,
		},
	)
}
