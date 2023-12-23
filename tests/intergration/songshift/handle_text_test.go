package songshift

import (
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	"github.com/GeorgeGorbanev/songshift/tests/utils"
	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestSongshift_HandleText(t *testing.T) {
	tests := []struct {
		name             string
		inMsg            *telebot.Message
		expectedResponse *telegram.Message
	}{
		{
			name: "should send a response to the user",
			inMsg: &telebot.Message{
				Sender: &telebot.User{
					Username: "sample_username",
				},
				Text: "sample_text",
			},
			expectedResponse: &telegram.Message{
				To:   &telebot.User{Username: "sample_username"},
				Text: "Received message: sample_text",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			senderMock := utils.NewTelegramSenderMock()
			ss := songshift.NewSongshift(senderMock)
			ss.HandleText(tt.inMsg)

			require.Equal(t, senderMock.Response, tt.expectedResponse)
		})
	}
}
