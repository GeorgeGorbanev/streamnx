package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_NotFound(t *testing.T) {
	sampleSender := &telebot.User{Username: "sample_username"}
	senderMock := utils.NewTelegramSenderMock()

	vs, err := vibeshare.NewVibeshare(&vibeshare.Input{}, vibeshare.WithVibeshareSender(senderMock))
	require.NoError(t, err)

	msg := telebot.Message{Sender: sampleSender, Text: "just message with no link"}
	vs.TextHandler(&msg)

	require.Equal(t, "No supported link found", senderMock.Response.Text)
}
