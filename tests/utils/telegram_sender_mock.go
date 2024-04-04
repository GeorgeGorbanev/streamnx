package utils

import (
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"

	"github.com/tucnak/telebot"
)

type TelegramSenderMock struct {
	Response *telegram.Message
	AllSent  []*telegram.Message
}

func NewTelegramSenderMock() *TelegramSenderMock {
	return &TelegramSenderMock{}
}

func (t *TelegramSenderMock) Send(msg *telegram.Message) (*telebot.Message, error) {
	t.Response = msg
	t.AllSent = append(t.AllSent, msg)
	return nil, nil
}

func (t *TelegramSenderMock) Reset() {
	t.Response = nil
}
