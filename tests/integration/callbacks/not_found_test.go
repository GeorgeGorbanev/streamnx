package callbacks

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_NotFound(t *testing.T) {
	senderMock := utils.NewTelegramSenderMock()
	app, err := vibeshare.NewVibeshare(&vibeshare.Input{}, vibeshare.WithVibeshareSender(senderMock))
	require.NoError(t, err)

	user := telebot.User{Username: "sample_username"}
	app.CallbackHandler(&telebot.Callback{
		Sender: &user,
		Data:   "sample_not_found_route",
	})

	require.Equal(t, &telegram.Message{
		To:   &user,
		Text: "This button not available anymore. Try to request a new button 🙏",
	}, senderMock.Response)
}
