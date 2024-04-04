# Vibeshare Bot

Vibeshare is a Telegram bot that helps you share music across different platforms effortlessly. Whether it's a track or an album, Vibeshare provides direct links to the same content on Spotify, Apple Music, Yandex Music, and YouTube, making music sharing seamless and enjoyable.

Check out the bot: [@vibeshare_bot](https://t.me/vibeshare_bot)

## Features

- **Cross-platform sharing:** Convert music links from one service to another with ease.
- **Support for major music platforms:** Spotify, Apple Music, Yandex Music, and YouTube.
- **User-friendly:** Simple and intuitive interface. All you need to do is send a link.

## How It Works

1. **Send a link:** Share a link to a track or album from any of the supported platforms.
2. **Choose your platform:** Bot will provide options to convert the link to other platforms.
3. **Share the music:** Get the link for your desired platform and share it with your friends

## Technologies

Bot is built using Go and leverages the APIs of various music platforms to fetch and convert music links.
- **No tracking:** Bot does not store any data in databases, ensuring there's no unwanted tracking of user data.
- **Easy hosting:** The absence of a database requirement makes it straightforward to host and maintain the bot.

## Running the Bot

To run Vibeshare, you will need to set up the following environment variables:

- `TELEGRAM_TOKEN`: Main Telegram bot token.
- `SPOTIFY_CLIENT_ID`: Your Spotify API client ID.
- `SPOTIFY_CLIENT_SECRET`: Your Spotify API client secret.
- `YOUTUBE_API_KEY`: Your YouTube API key.
- `FEEDBACK_TOKEN`: Telegram bot token for feedback messages.
- `FEEDBACK_RECEIVER_ID`: Telegram user ID to receive feedback messages.

Once you have these variables set up, you can run the bot using make:

```bash
make run
```

Also you can run tests using:

```bash
make test
```
Test coverage close to 100% to ensure the reliability and stability of the bot. 

## Future Enhancements

- **More platforms:** Support for additional music streaming services.
- **Playlist sharing:** Ability to share entire playlists across platforms.
- **Personalization:** Custom settings for preferred music platforms.