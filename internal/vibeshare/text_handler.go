package vibeshare

import (
	"fmt"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/templates"

	"github.com/tucnak/telebot"
)

func (vs *Vibeshare) startCommand(inMsg *telebot.Message) {
	vs.send(&telegram.Message{To: inMsg.Sender, Text: templates.Start})
}

func (vs *Vibeshare) feedbackCommand(inMsg *telebot.Message) {
	text := fmt.Sprintf(templates.FeedbackCommand, vs.feedbackBotName)
	vs.send(&telegram.Message{To: inMsg.Sender, Text: text})
}

func (vs *Vibeshare) whatLinks(inMsg *telebot.Message) {
	vs.send(&telegram.Message{To: inMsg.Sender, Text: templates.WhatLinksResponse, ReplyMarkup: whatLinksMenu()})
}

func (vs *Vibeshare) skip(_ *telebot.Message) {
}

func (vs *Vibeshare) textNotFoundHandler(inMsg *telebot.Message) {
	vs.send(&telegram.Message{To: inMsg.Sender, Text: templates.NotFound, ReplyMarkup: notFoundMenu()})
}
