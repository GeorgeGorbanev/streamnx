package vibeshare

import (
	"fmt"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
)

type Input struct {
	MusicRegistry      *music.Registry
	VibeshareBotToken  string
	FeedbackBotToken   string
	FeedbackReceiverID int
}

type Vibeshare struct {
	musicRegistry *music.Registry

	vibeshareBot    *telegram.Bot
	vibeshareSender telegram.Sender
	vibeshareRouter *telegram.Router

	feedbackBot        *telegram.Bot
	feedbackSender     telegram.Sender
	feedbackRouter     *telegram.Router
	feedbackBotName    string
	feedbackReceiverID int
}

func NewVibeshare(input *Input, opts ...Option) (*Vibeshare, error) {
	vs := &Vibeshare{}
	vs.musicRegistry = input.MusicRegistry
	vs.feedbackReceiverID = input.FeedbackReceiverID

	if err := vs.setupVibeshareBot(input); err != nil {
		return vs, err
	}
	if err := vs.setupFeedbackBot(input); err != nil {
		return vs, err
	}

	for _, opt := range opts {
		opt(vs)
	}

	vs.setupVibeshareRouter()
	vs.setupFeedbackRouter()

	return vs, nil
}

func (vs *Vibeshare) Run() {
	if vs.feedbackBot != nil {
		go vs.feedbackBot.Run()
	}
	if vs.vibeshareBot != nil {
		vs.vibeshareBot.Run()
	}
}

func (vs *Vibeshare) Stop() {
	if vs.feedbackBot != nil {
		vs.feedbackBot.Stop()
	}
	if vs.vibeshareBot != nil {
		vs.vibeshareBot.Stop()
	}
	if vs.musicRegistry != nil {
		vs.musicRegistry.Close()
	}
}

func (vs *Vibeshare) setupVibeshareBot(input *Input) error {
	if input.VibeshareBotToken == "" {
		return nil
	}
	bot, err := telegram.NewBot(input.VibeshareBotToken)
	if err != nil {
		return fmt.Errorf("failed to create vibeshare bot: %w", err)
	}
	bot.HandleText(vs.TextHandler)
	bot.HandleCallback(vs.CallbackHandler)

	vs.vibeshareSender = bot.Sender()
	vs.vibeshareBot = bot
	return nil
}

func (vs *Vibeshare) setupFeedbackBot(input *Input) error {
	if input.FeedbackBotToken == "" {
		return nil
	}
	bot, err := telegram.NewBot(input.FeedbackBotToken)
	if err != nil {
		return fmt.Errorf("failed to create feedback bot: %w", err)
	}
	bot.HandleText(vs.FeedbackTextHandler)

	vs.feedbackSender = bot.Sender()
	vs.feedbackBotName = bot.Name()
	vs.feedbackBot = bot
	return nil
}
