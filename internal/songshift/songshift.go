package songshift

import (
	"errors"
	"log"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
	"github.com/tucnak/telebot"
)

type Songshift struct {
	spotifyClient  *spotify.Client
	telegramSender telegram.Sender
	ymusicClient   *ymusic.Client
}

type Input struct {
	SpotifyClient  *spotify.Client
	TelegramSender telegram.Sender
	YmusicClient   *ymusic.Client
}

func NewSongshift(input *Input) *Songshift {
	return &Songshift{
		spotifyClient:  input.SpotifyClient,
		telegramSender: input.TelegramSender,
		ymusicClient:   input.YmusicClient,
	}
}

func (s *Songshift) HandleText(inMsg *telebot.Message) {
	trackID := spotify.DetectTrackID(inMsg.Text)
	if trackID == "" {
		outMsg, err := s.respond(inMsg, "no track link found")
		if err != nil {
			log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
			return
		}
		log.Printf("Sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
		return
	}

	track, err := s.spotifyClient.GetTrack(trackID)
	if err != nil {
		log.Printf("Error fetching track: %s", err)
		if errors.Is(err, spotify.TrackNotFoundError) {
			outMsg, err := s.respond(inMsg, "track not found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
		}
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

func (s *Songshift) respond(inMsg *telebot.Message, text string) (*telebot.Message, error) {
	return s.telegramSender.Send(
		&telegram.Message{
			To:   inMsg.Sender,
			Text: text,
		},
	)
}
