package telegram

import "github.com/tucnak/telebot"

type Sender interface {
	Send(msg *Message) (*telebot.Message, error)
}

type TelebotSender struct {
	bot *telebot.Bot
}

func NewTelebotSender(bot *telebot.Bot) *TelebotSender {
	return &TelebotSender{
		bot: bot,
	}
}

func (t *TelebotSender) Send(msg *Message) (*telebot.Message, error) {
	return t.bot.Send(msg.To, msg.Text)
}
