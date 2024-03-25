package telegram

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestRouter_RouteText(t *testing.T) {
	var startRe = regexp.MustCompile(`^/start$`)

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
			router.HandleText(startRe, func(inMsg *telebot.Message) {
				handlerCalled = true
			})

			notFoundCalled := false
			if tt.handleNotFound {
				router.HandleTextNotFound(func(inMsg *telebot.Message) {
					notFoundCalled = true
				})
			}

			router.RouteText(tt.inMsg)
			require.Equal(t, tt.expectHandlerCalled, handlerCalled)
			require.Equal(t, tt.expectNotFoundCalled, notFoundCalled)
		})
	}
}

func TestRouter_RouteCallback(t *testing.T) {
	tests := []struct {
		name                 string
		callbackData         string
		handleNotFound       bool
		expectHandlerCalled  bool
		expectNotFoundCalled bool
	}{
		{
			name: "when command matched",
			callbackData: `{
				"command":"sample_callback",
				"payload":{
					"sample":"payload"
				}
			}`,
			handleNotFound:       false,
			expectHandlerCalled:  true,
			expectNotFoundCalled: false,
		},
		{
			name: "when msg not matched",
			callbackData: `{
				"command":"not_found_callback",
				"payload":{
					"sample":"payload"
				}
			}`,
			handleNotFound:       false,
			expectHandlerCalled:  false,
			expectNotFoundCalled: false,
		},
		{
			name: "when msg not matched and not found given",
			callbackData: `{
				"command":"not_found_callback",
				"payload":{
					"sample":"payload"
				}
			}`,
			handleNotFound:       true,
			expectHandlerCalled:  false,
			expectNotFoundCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()

			handlerCalled := false
			router.HandleCallback("sample_callback", func(callback *Callback) {
				handlerCalled = true
			})

			notFoundCalled := false
			if tt.handleNotFound {
				router.HandleCallbackNotFound(func(callback *Callback) {
					notFoundCalled = true
				})
			}

			router.RouteCallback(&telebot.Callback{
				Data: tt.callbackData,
			})
			require.Equal(t, tt.expectHandlerCalled, handlerCalled)
			require.Equal(t, tt.expectNotFoundCalled, notFoundCalled)
		})
	}
}
