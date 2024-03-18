package telegram

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

var startRe = regexp.MustCompile(`^/start$`)

func TestRouter(t *testing.T) {
	tests := []struct {
		name                 string
		inMsg                *telebot.Message
		handleNotFound       bool
		expectHandlerCalled  bool
		expectNotFoundCalled bool
	}{
		{
			name: "when msg matched",
			inMsg: &telebot.Message{
				Text: "/start",
			},
			handleNotFound:       false,
			expectHandlerCalled:  true,
			expectNotFoundCalled: false,
		},
		{
			name: "when msg not matched",
			inMsg: &telebot.Message{
				Text: "not matched",
			},
			handleNotFound:       false,
			expectHandlerCalled:  false,
			expectNotFoundCalled: false,
		},
		{
			name: "when msg not matched and not found given",
			inMsg: &telebot.Message{
				Text: "not matched",
			},
			handleNotFound:       true,
			expectHandlerCalled:  false,
			expectNotFoundCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()

			handlerCalled := false
			router.Register(startRe, func(inMsg *telebot.Message) {
				handlerCalled = true
			})

			notFoundCalled := false
			if tt.handleNotFound {
				router.RegisterNotFound(func(inMsg *telebot.Message) {
					notFoundCalled = true
				})
			}

			router.RouteMessage(tt.inMsg)
			require.Equal(t, tt.expectHandlerCalled, handlerCalled)
			require.Equal(t, tt.expectNotFoundCalled, notFoundCalled)
		})
	}
}
