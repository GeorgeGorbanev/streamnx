package telegram

import (
	"strings"

	"github.com/tucnak/telebot"
)

const CallbackRouteSeparator = "/"

type Callback struct {
	Sender *telebot.User
	Data   *CallbackData
}

type CallbackData struct {
	Route  string
	Params []string
}

func (cd *CallbackData) Marshal() string {
	return strings.Join(
		append([]string{cd.Route}, cd.Params...),
		CallbackRouteSeparator,
	)
}

func (cd *CallbackData) Unmarshal(data string) {
	parts := strings.Split(data, CallbackRouteSeparator)
	cd.Route = parts[0]
	cd.Params = parts[1:]
}
