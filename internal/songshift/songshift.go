package songshift

import (
	"fmt"
	"log"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	"github.com/tucnak/telebot"
)

type Songshift struct {
	telegramSender telegram.Sender
}

func NewSongshift(telegramSender telegram.Sender) *Songshift {
	return &Songshift{
		telegramSender: telegramSender,
	}
}

func (s *Songshift) HandleText(inMsg *telebot.Message) {
	log.Printf("Received message from %s: %s", inMsg.Sender.Username, inMsg.Text)

	response := fmt.Sprintf("Received message: %s", inMsg.Text)
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
