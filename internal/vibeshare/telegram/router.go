package telegram

import (
	"encoding/json"
	"regexp"

	"github.com/tucnak/telebot"
)

type Router struct {
	textHandlers            []*textHandler
	textNotFoundHandler     TextHandlerFunc
	callbackHandlers        []*callbackHandler
	callbackHandlerNotFound CallbackHandlerFunc
	callbackHandlerError    CallbackHandlerFunc
}

func NewRouter() *Router {
	return &Router{
		textHandlers:     make([]*textHandler, 0),
		callbackHandlers: make([]*callbackHandler, 0),
	}
}

func (r *Router) HandleText(re *regexp.Regexp, hf TextHandlerFunc) {
	r.textHandlers = append(r.textHandlers, &textHandler{
		re:          re,
		handlerFunc: hf,
	})
}

func (r *Router) HandleCallback(command string, hf CallbackHandlerFunc) {
	r.callbackHandlers = append(r.callbackHandlers, &callbackHandler{
		command:     command,
		handlerFunc: hf,
	})
}

func (r *Router) HandleTextNotFound(hf TextHandlerFunc) {
	r.textNotFoundHandler = hf
}

func (r *Router) HandleCallbackNotFound(hf CallbackHandlerFunc) {
	r.callbackHandlerNotFound = hf
}

func (r *Router) HandleCallbackError(hf CallbackHandlerFunc) {
	r.callbackHandlerError = hf
}

func (r *Router) RouteText(inMsg *telebot.Message) {
	for _, h := range r.textHandlers {
		if h.re.MatchString(inMsg.Text) {
			h.handlerFunc(inMsg)
			return
		}
	}
	if r.textNotFoundHandler != nil {
		r.textNotFoundHandler(inMsg)
	}
}

func (r *Router) RouteCallback(callback *telebot.Callback) {
	cb := Callback{Data: &callbackData{}}
	if err := json.Unmarshal([]byte(callback.Data), cb.Data); err != nil {
		if r.callbackHandlerError != nil {
			r.callbackHandlerError(&cb)
		}
		return
	}

	for _, h := range r.callbackHandlers {
		if h.command == cb.Data.Command {
			h.handlerFunc(&cb)
			return
		}
	}
	if r.callbackHandlerNotFound != nil {
		r.callbackHandlerNotFound(&cb)
	}
}
