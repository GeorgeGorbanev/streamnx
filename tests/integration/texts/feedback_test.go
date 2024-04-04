package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_Feedback(t *testing.T) {
	feedbackReceiverID := 42
	feedbackAuthor := telebot.User{
		Username: "feedback_author",
	}

	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
	}{
		{
			name:  "when /start command is sent",
			input: "/start",
			expectedMessages: []*telegram.Message{
				{
					To: &feedbackAuthor,
					Text: "Here you can send feedback to the author of this bot. " +
						"Just type your message and it will be delivered to the author directly. 🙏",
				},
			},
		},
		{
			name:  "when any other text is sent",
			input: "some feedback",
			expectedMessages: []*telegram.Message{
				{
					To: &telebot.User{
						ID: feedbackReceiverID,
					},
					Text: "From @feedback_author (0): some feedback",
				},
				{
					To:   &feedbackAuthor,
					Text: "Your feedback will be delivered to the author. Thank you! 🙏",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			senderMock := utils.NewTelegramSenderMock()
			vs, err := vibeshare.NewVibeshare(&vibeshare.Input{
				FeedbackReceiverID: feedbackReceiverID,
			}, vibeshare.WithFeedbackSender(senderMock))
			require.NoError(t, err)

			vs.FeedbackTextHandler(&telebot.Message{
				Text:   tt.input,
				Sender: &feedbackAuthor,
			})

			require.Equal(t, tt.expectedMessages, senderMock.AllSent)
		})
	}
}
