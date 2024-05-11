package telegram

import (
	"regexp"
	"strings"

	"github.com/tucnak/telebot"
)

type TextHandler struct {
	Re          *regexp.Regexp
	HandlerFunc TextHandlerFunc
}

type CallbackHandler struct {
	Route       string
	HandlerFunc CallbackHandlerFunc
}

type TextHandlerFunc func(inMsg *telebot.Message)

type CallbackHandlerFunc func(callback *Callback)

type Router struct {
	TextHandlers            []*TextHandler
	TextNotFoundHandler     TextHandlerFunc
	CallbackHandlers        []*CallbackHandler
	CallbackHandlerNotFound CallbackHandlerFunc
}

func (r *Router) RouteText(inMsg *telebot.Message) {
	for _, h := range r.TextHandlers {
		if h.Re.MatchString(inMsg.Text) {
			h.HandlerFunc(inMsg)
			return
		}
	}
	if r.TextNotFoundHandler != nil {
		r.TextNotFoundHandler(inMsg)
	}
}

func (r *Router) RouteCallback(callback *telebot.Callback) {
	splittedData := strings.Split(callback.Data, CallbackRouteSeparator)
	cb := Callback{
		Sender: callback.Sender,
		Data: &CallbackData{
			Route:  splittedData[0],
			Params: splittedData[1:],
		},
	}

	for _, h := range r.CallbackHandlers {
		if h.Route == cb.Data.Route {
			h.HandlerFunc(&cb)
			return
		}
	}
	if r.CallbackHandlerNotFound != nil {
		r.CallbackHandlerNotFound(&cb)
	}
}
