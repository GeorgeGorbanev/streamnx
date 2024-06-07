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
		name                     string
		inMsg                    *telebot.Message
		handleNotFound           bool
		expectStartHandlerCalled bool
		expectNotFoundCalled     bool
	}{
		{
			name: "when msg matched",
			inMsg: &telebot.Message{
				Text: "/start",
			},
			handleNotFound:           false,
			expectStartHandlerCalled: true,
			expectNotFoundCalled:     false,
		},
		{
			name: "when msg not matched",
			inMsg: &telebot.Message{
				Text: "not matched",
			},
			handleNotFound:           false,
			expectStartHandlerCalled: false,
			expectNotFoundCalled:     false,
		},
		{
			name: "when msg not matched and not found given",
			inMsg: &telebot.Message{
				Text: "not matched",
			},
			handleNotFound:           true,
			expectStartHandlerCalled: false,
			expectNotFoundCalled:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startCalled := false
			notFoundCalled := false

			router := Router{
				TextRoutes: []*TextRoute{
					{
						Pattern: startRe,
						Handler: func(inMsg *telebot.Message) {
							startCalled = true
						},
					},
				},
			}

			if tt.handleNotFound {
				router.TextNotFound = func(inMsg *telebot.Message) {
					notFoundCalled = true
				}
			}

			router.RouteText(tt.inMsg)
			require.Equal(t, tt.expectStartHandlerCalled, startCalled)
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
			name:                 "when route matched",
			callbackData:         "sample_callback/any/route/params",
			handleNotFound:       false,
			expectHandlerCalled:  true,
			expectNotFoundCalled: false,
		},
		{
			name:                 "when msg not matched",
			callbackData:         "not_found_callback/any/route/params",
			handleNotFound:       false,
			expectHandlerCalled:  false,
			expectNotFoundCalled: false,
		},
		{
			name:                 "when msg not matched and not found given",
			callbackData:         "not_found_callback/any/route/params",
			handleNotFound:       true,
			expectHandlerCalled:  false,
			expectNotFoundCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampleCallbackCalled := false
			notFoundCalled := false

			router := Router{
				CallbackRoutes: []*CallbackRoute{
					{
						Address: "sample_callback",
						Handler: func(callback *Callback) {
							sampleCallbackCalled = true
						},
					},
				},
			}

			if tt.handleNotFound {
				router.CallbackNotFound = func(callback *Callback) {
					notFoundCalled = true
				}
			}

			router.RouteCallback(&telebot.Callback{
				Data: tt.callbackData,
			})

			require.Equal(t, tt.expectHandlerCalled, sampleCallbackCalled)
			require.Equal(t, tt.expectNotFoundCalled, notFoundCalled)
		})
	}
}
