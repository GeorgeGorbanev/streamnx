package telegram

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/tucnak/telebot"
)

const timeout = 10 * time.Second

type Bot struct {
	telebotBot *telebot.Bot
}

func NewBot(token string) (*Bot, error) {
	telebotBot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: timeout},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating telebot: %w", err)
	}
	return &Bot{
		telebotBot: telebotBot,
	}, nil
}

func (b *Bot) Sender() Sender {
	return NewTelebotSender(b.telebotBot)
}

func (b *Bot) HandleText(handler func(inMsg *telebot.Message)) {
	b.telebotBot.Handle(telebot.OnText, handler)
}

func (b *Bot) HandleCallback(handler func(callback *telebot.Callback)) {
	b.telebotBot.Handle(telebot.OnCallback, handler)
}

func (b *Bot) Name() string {
	return b.telebotBot.Me.Username
}

func (b *Bot) Run() {
	slog.Info("Bot started", slog.String("bot_name", b.Name()))
	b.telebotBot.Start()
}

func (b *Bot) Stop() {
	b.telebotBot.Stop()
}
