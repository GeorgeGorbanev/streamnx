package telegram

import (
	"encoding/json"
	"regexp"

	"github.com/tucnak/telebot"
)

type textHandler struct {
	re          *regexp.Regexp
	handlerFunc TextHandlerFunc
}

type callbackHandler struct {
	command     string
	handlerFunc CallbackHandlerFunc
}

type TextHandlerFunc func(inMsg *telebot.Message)

type CallbackHandlerFunc func(callback *Callback)

type callbackData struct {
	Command string          `json:"command"`
	Payload json.RawMessage `json:"payload"`
}

type Callback struct {
	Data *callbackData
}

type Message struct {
	To          telebot.Recipient
	Text        string
	ReplyMarkup *telebot.ReplyMarkup
}
