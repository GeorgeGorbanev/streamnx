package telegram

import (
	"regexp"
	"strings"

	"github.com/tucnak/telebot"
)

type TextRoute struct {
	Pattern *regexp.Regexp
	Handler TextHandler
}

type CallbackRoute struct {
	Address string
	Handler CallbackHandler
}

type TextHandler func(inMsg *telebot.Message)

type CallbackHandler func(callback *Callback)

type Router struct {
	TextRoutes   []*TextRoute
	TextNotFound TextHandler

	CallbackRoutes   []*CallbackRoute
	CallbackNotFound CallbackHandler
}

func (r *Router) RouteText(inMsg *telebot.Message) {
	for _, route := range r.TextRoutes {
		if route.Pattern.MatchString(inMsg.Text) {
			route.Handler(inMsg)
			return
		}
	}
	if r.TextNotFound != nil {
		r.TextNotFound(inMsg)
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

	for _, route := range r.CallbackRoutes {
		if route.Address == cb.Data.Route {
			route.Handler(&cb)
			return
		}
	}
	if r.CallbackNotFound != nil {
		r.CallbackNotFound(&cb)
	}
}
