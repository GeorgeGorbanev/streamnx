package songshift

import (
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
	log.Printf("Received message from %s: %s", inMsg.Sender.Username, inMsg.Text)

	token, err := s.spotifyClient.FetchToken()
	if err != nil {
		log.Printf("Error fetching token: %s", err)
		return
	}
	response := fmt.Sprintf("Received message: %s (token %s)", inMsg.Text, token.AccessToken)
	outMsg, err := s.telegramSender.Send(
		&telegram.Message{
			To:   inMsg.Sender,
			Text: response,
		},
	)
	if err != nil {
		log.Printf("Error sending message to %s: %s", inMsg.Sender.Username, err)
		return
	}
	log.Printf("Sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
}
