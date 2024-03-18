package telegram

import (
	"regexp"

	"github.com/tucnak/telebot"
)

type Router struct {
	handlers        []*handler
	notFoundHandler HandlerFunc
}

type HandlerFunc func(*telebot.Message)

type handler struct {
	re          *regexp.Regexp
	handlerFunc HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		handlers: make([]*handler, 0),
	}
}

func (r *Router) RouteMessage(inMsg *telebot.Message) {
	for _, h := range r.handlers {
		if h.re.MatchString(inMsg.Text) {
			h.handlerFunc(inMsg)
			return
		}
	}
	if r.notFoundHandler != nil {
		r.notFoundHandler(inMsg)
	}
}

func (r *Router) Register(re *regexp.Regexp, hf HandlerFunc) {
	r.handlers = append(r.handlers, &handler{
		re:          re,
		handlerFunc: hf,
	})
}

func (r *Router) RegisterNotFound(hf HandlerFunc) {
	r.notFoundHandler = hf
}
