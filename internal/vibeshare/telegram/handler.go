package telegram

import (
	"regexp"

	"github.com/tucnak/telebot"
)

type textHandler struct {
	re          *regexp.Regexp
	handlerFunc TextHandlerFunc
}

type callbackHandler struct {
	route       string
	handlerFunc CallbackHandlerFunc
}

type TextHandlerFunc func(inMsg *telebot.Message)

type CallbackHandlerFunc func(callback *Callback)
