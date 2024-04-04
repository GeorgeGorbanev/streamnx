package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_Skip(t *testing.T) {
	user := &telebot.User{Username: "sample_username"}
	msg := &telebot.Message{Sender: user, Text: `Skip`}

	senderMock := utils.NewTelegramSenderMock()
	vs, err := vibeshare.NewVibeshare(&vibeshare.Input{}, vibeshare.WithVibeshareSender(senderMock))
	require.NoError(t, err)

	vs.TextHandler(msg)

	require.Nil(t, senderMock.Response)
}

func TestText_SkipDemonstration(t *testing.T) {
	user := &telebot.User{Username: "sample_username"}
	msg := &telebot.Message{Sender: user, Text: `Skip demonstration`}

	senderMock := utils.NewTelegramSenderMock()
	vs, err := vibeshare.NewVibeshare(&vibeshare.Input{}, vibeshare.WithVibeshareSender(senderMock))
	require.NoError(t, err)

	vs.TextHandler(msg)

	require.Nil(t, senderMock.Response)
}
