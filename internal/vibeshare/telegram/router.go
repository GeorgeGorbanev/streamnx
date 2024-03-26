package telegram

import (
	"regexp"
	"strings"

	"github.com/tucnak/telebot"
)

type Router struct {
	textHandlers            []*textHandler
	textNotFoundHandler     TextHandlerFunc
	callbackHandlers        []*callbackHandler
	callbackHandlerNotFound CallbackHandlerFunc
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

func (r *Router) HandleCallback(route string, hf CallbackHandlerFunc) {
	r.callbackHandlers = append(r.callbackHandlers, &callbackHandler{
		route:       route,
		handlerFunc: hf,
	})
}

func (r *Router) HandleTextNotFound(hf TextHandlerFunc) {
	r.textNotFoundHandler = hf
}

func (r *Router) HandleCallbackNotFound(hf CallbackHandlerFunc) {
	r.callbackHandlerNotFound = hf
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
	splittedData := strings.Split(callback.Data, CallbackRouteSeparator)
	cb := Callback{
		Sender: callback.Sender,
		Data: &CallbackData{
			Route:  splittedData[0],
			Params: splittedData[1:],
		},
	}

	for _, h := range r.callbackHandlers {
		if h.route == cb.Data.Route {
			h.handlerFunc(&cb)
			return
		}
	}
	if r.callbackHandlerNotFound != nil {
		r.callbackHandlerNotFound(&cb)
	}
}
