package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"
	"github.com/stretchr/testify/require"

	"github.com/tucnak/telebot"
)

func TestText_Start(t *testing.T) {

	user := &telebot.User{
		Username: "sample_username",
	}
	msg := &telebot.Message{
		Sender: user,
		Text:   `/start`,
	}

	senderMock := utils.NewTelegramSenderMock()
	vs := vibeshare.NewVibeshare(&vibeshare.Input{
		TelegramSender: senderMock,
	})

	vs.TextHandler(msg)

	require.NotNil(t, senderMock.Response)
	require.Equal(t, user, senderMock.Response.To)
	require.Equal(t,
		`🎶 Welcome to Vibeshare! 🎶

I'm here to help you share your favorite music tracks and albums with friends, no matter which streaming platforms they use. With just one click, you can provide them with direct access to the music you want to share, without the need for searching.

Here's how to get started:

<b>1) Share a Track or Album Link:</b> <i>Simply send me a link to a track or album from Spotify, Apple Music, Yandex Music, or YouTube.</i>

<b>2) Choose Your Platform:</b> <i>I'll provide you with options to convert the link to other platforms.</i>

<b>3) Share the Music:</b> <i>Once you've got the link for your desired platform, share it with your friends and enjoy the music together!</i>

Let's spread the vibes! ☮️`,
		senderMock.Response.Text)
}
