package vibeshare

import (
	"testing"

	telegram_utils "github.com/GeorgeGorbanev/vibeshare/tests/utils"
	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestVibeshare_HandleText(t *testing.T) {
	sampleSender := &telebot.User{Username: "sample_username"}
	senderMock := telegram_utils.NewTelegramSenderMock()
	vs := NewVibeshare(&Input{
		TelegramSender: senderMock,
	})

	vs.TextHandler(&telebot.Message{
		Sender: sampleSender,
		Text:   "just message with no link",
	})

	require.Equal(t, "No supported link found", senderMock.Response.Text)
}
