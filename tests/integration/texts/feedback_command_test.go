package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_FeedbackCommand(t *testing.T) {
	user := &telebot.User{Username: "sample_username"}
	msg := &telebot.Message{Sender: user, Text: `/feedback`}

	senderMock := utils.NewTelegramSenderMock()
	vs, err := vibeshare.NewVibeshare(&vibeshare.Input{},
		vibeshare.WithVibeshareSender(senderMock),
		vibeshare.WithFeedbackBotName("feedback_bot"))
	require.NoError(t, err)

	vs.TextHandler(msg)

	require.NotNil(t, senderMock.Response)
	require.Equal(t, user, senderMock.Response.To)
	require.Equal(t, "You can send feedback to @feedback_bot", senderMock.Response.Text)
}
