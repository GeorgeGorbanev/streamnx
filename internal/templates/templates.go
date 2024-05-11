package templates

import (
	"fmt"
	"regexp"
)

const (
	ExampleTrack      = "Convert example track https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	FeedbackCommand   = `Tap <a href="https://t.me/%s">here</a> to open feedback dialogue`
	FeedbackStart     = "Here you can send feedback to the author of this bot. Just type your message and it will be delivered to the author directly. 🙏"
	FeedbackThanks    = "Your feedback will be delivered to the human. Thank you! 🙏"
	FeedbackReceived  = "From @%s (%d): %s"
	NotFound          = `No supported link found`
	Skip              = "Skip"
	SkipDemonstration = "Skip demonstration"
	SpecifyRegion     = "This link has <code>.com</code> domain. You can specify country of recipient to be sure link accessible in recipients`s region"
	Start             = `🎶 Welcome to Vibeshare! 🎶

I'm here to help you share your favorite music tracks and albums with friends, no matter which streaming platforms they use. With just one click, you can provide them with direct access to the music you want to share, without the need for searching.

Here's how to get started:

<b>1) Share a Track or Album Link:</b> <i>Simply send me a link to a track or album from Spotify, Apple Music, Yandex Music, or YouTube.</i>

<b>2) Choose Your Platform:</b> <i>I'll provide you with options to convert the link to other platforms.</i>

<b>3) Share the Music:</b> <i>Once you've got the link for your desired platform, share it with your friends and enjoy the music together!</i>

Let's spread the vibes! ☮️`
	WhatLinksButton   = "What links do you support?"
	WhatLinksResponse = `I support the following music platforms links:
<code>
🍏 Apple Music
			- https://music.apple.com/*/album/track/*,
			- https://music.apple.com/*/album/track/*?i=*,
			- https://music.apple.com/*/song/angel/*
			- https://music.apple.com/*/album/*	
🎵 Spotify
			- https://open.spotify.com/track/*	
			- https://open.spotify.com/album/*	
🔊 Yandex Music
			- https://music.yandex.(com|ru)/album/*/track/*
			- https://music.yandex.(com|ru)/album/*	
📺 YouTube
			- https://www.youtube.com/watch?v=*
			- https://www.youtu.be/*"
			- https://www.youtube.com/playlist?list=*
			- https://www.youtu.be/playlist?list=*
</code>
Send me any of them and I'll convert it to the other platforms for you`
)

var (
	WhatLinksButtonRe   = regexp.MustCompile(strictEqual(WhatLinksButton))
	SkipRe              = regexp.MustCompile(strictEqual(Skip))
	SkipDemonstrationRe = regexp.MustCompile(strictEqual(SkipDemonstration))
)

func strictEqual(text string) string {
	escaped := regexp.QuoteMeta(text)
	return fmt.Sprintf("^%s$", escaped)
}
