package vibeshare

import "github.com/GeorgeGorbanev/vibeshare/internal/telegram"

type Option func(vs *Vibeshare)

func WithVibeshareSender(sender telegram.Sender) Option {
	return func(vs *Vibeshare) {
		vs.vibeshareSender = sender
	}
}

func WithFeedbackSender(sender telegram.Sender) Option {
	return func(vs *Vibeshare) {
		vs.feedbackSender = sender
	}
}

func WithFeedbackBotName(name string) Option {
	return func(vs *Vibeshare) {
		vs.feedbackBotName = name
	}
}
