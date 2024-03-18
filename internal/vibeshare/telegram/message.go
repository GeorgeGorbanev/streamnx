package telegram

import "github.com/tucnak/telebot"

type Message struct {
	To   telebot.Recipient
	Text string
}
