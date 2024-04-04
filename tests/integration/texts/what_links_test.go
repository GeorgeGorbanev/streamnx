package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/templates"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_WhatLinks(t *testing.T) {
	user := &telebot.User{Username: "sample_username"}
	msg := &telebot.Message{Sender: user, Text: templates.WhatLinksButton}

	senderMock := utils.NewTelegramSenderMock()
	vs, err := vibeshare.NewVibeshare(&vibeshare.Input{}, vibeshare.WithVibeshareSender(senderMock))
	require.NoError(t, err)

	vs.TextHandler(msg)

	require.NotNil(t, senderMock.Response)
	require.Equal(t, user, senderMock.Response.To)
	require.Equal(t, templates.WhatLinksResponse, senderMock.Response.Text)
}
