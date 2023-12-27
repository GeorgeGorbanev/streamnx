package songshift

import (
	"errors"
	"fmt"
	"log"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	"github.com/tucnak/telebot"
)

type Songshift struct {
	spotifyClient  *spotify.Client
	telegramSender telegram.Sender
}

func NewSongshift(spotifyClient *spotify.Client, telegramSender telegram.Sender) *Songshift {
	return &Songshift{
		spotifyClient:  spotifyClient,
		telegramSender: telegramSender,
	}
}

func (s *Songshift) HandleText(inMsg *telebot.Message) {
	trackID := spotify.DetectTrackID(inMsg.Text)
	if trackID == "" {
		outMsg, err := s.respond(inMsg, "no track link found")
		if err != nil {
			log.Printf("Error sending message to %s: %s", inMsg.Sender.Username, err)
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
				log.Printf("Error sending message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("Sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
		}
		return
	}

	outMsg, err := s.respond(inMsg, fmt.Sprintf(`Track: "%s – %s"`, track.ArtistsString(), track.Name))
	if err != nil {
		log.Printf("Error sending message to %s: %s", inMsg.Sender.Username, err)
		return
	}
	log.Printf("Sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
}

func (s *Songshift) respond(inMsg *telebot.Message, text string) (*telebot.Message, error) {
	return s.telegramSender.Send(
		&telegram.Message{
			To:   inMsg.Sender,
			Text: text,
		},
	)
}
