package utils

import (
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/tucnak/telebot"
)

type TelegramSenderMock struct {
	Response *telegram.Message
}

func NewTelegramSenderMock() *TelegramSenderMock {
	return &TelegramSenderMock{}
}

func (t *TelegramSenderMock) Send(msg *telegram.Message) (*telebot.Message, error) {
	t.Response = msg
	return &telebot.Message{}, nil
}
