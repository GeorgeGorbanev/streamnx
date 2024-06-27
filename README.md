# Streaminx

Streaminx is a library that unifies interactions with various music streaming links into a single system.
With Streaminx, you can integrate platforms such as Apple Music, Spotify, YouTube and Yandex Music into your applications using a unified interface for searching and retrieving data about tracks and albums.

## Motivation

The main user of this library is the Telegram bot [Vibeshare](https://t.me/vibeshare_bot).
This bot helps users share song links from the streaming services they use, enabling others to listen to the same songs on different platforms.
Therefore, the primary use case for this library is *converting links* from one service to another.

## Providers Supported

The library supports the following music streaming services
- Apple Music
- Spotify
- Yandex Music
- YouTube (also YouTube Music)

## Installation

You can install the Streaminx library by running the following command in your terminal:

``` bash
go get github.com/GeorgeGorbanev/streaminx
```

To include lib in your project, import it in your Go source file:

``` golang
import "github.com/GeorgeGorbanev/streaminx"
```

## Configuration

To configure Streaminx, you need to set up the necessary credentials and API keys for the supported music streaming services.
Here are the steps to configure the library:

1) *Google Translator API.* Obtain the Google Translator API key and project ID from the [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2) *YouTube API*. Obtain the YouTube API key from the [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
3) *Spotify*. Register your application and obtain the Client ID with Client Secret on the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard).
4) *Init registry*. When you have all the necessary credentials, you can initialize the Streaminx *registry* with the following code.

``` golang
package main

import (
    "context"

    "github.com/GeorgeGorbanev/streaminx"
)

func main() {
    ctx := context.Background()
    registry, err := streaminx.NewRegistry(ctx, streaminx.Credentials{
        GoogleTranslatorAPIKeyJSON: "YOUR_GOOGLE_TRANSLATOR_API_KEY_JSON",
        GoogleTranslatorProjectID:  "YOUR_GOOGLE_TRANSLATOR_PROJECT_ID",
        YoutubeAPIKey:              "YOUR_YOUTUBE_API_KEY",
        SpotifyClientID:            "YOUR_SPOTIFY_CLIENT_ID",
        SpotifyClientSecret:        "YOUR_SPOTIFY_CLIENT_SECRET",
    })
    if err != nil {
        // Handle error
    }
    defer registry.Close()

    // Your code to use the registry
}
```
Replace the placeholders with your actual API keys and credentials.

## Usage

Here is an example of how to convert a link from Spotify to YouTube:

``` golang

func convert(link string) string {
    spotify := registry.Adapter(streaminx.Spotify)
    id, err := spotify.DetectTrackID(link)
    if err != nil {
        // Handle error
    }
    track, err := spotify.GetTrack(ctx, link)
    if err != nil {
        // Handle error
    }

    convertedTrack, err := registry.Adapter(streaminx.Youtube).SearchTrack(track.Artist, track.Name)
    if err != nil {
        // Handle error
    }

    return convertedTrack.URL
}

```

## API

#### Registry

The root of the library is the `Registry` struct.
The purpose of the `Registry` is to provide a unified interface for working with different streaming services.
Basically, `Registry` is a map, where the key is `Provider` and the value is the `Adapter` for this provider.
``` golang
apple := registry.Adapter(streaminx.Apple)
spotify := registry.Adapter(streaminx.Spotify)
yandex := registry.Adapter(streaminx.Yandex)
youtube := registry.Adapter(streaminx.Youtube)
```

#### Provider

Provider represents a music streaming service. All providers are accessible via the `Providers` enum:
``` golang
for _, provider := range streaminx.Providers {
    fmt.Println(provider.Name)
}
```

#### Adapter

Using *adapters* you can fetch and search for tracks and albums on the supported streaming services.
Each adapter implements the `Adapter` interface, which provides methods for working with tracks and albums.

``` golang
type Adapter interface {
	DetectTrackID(trackURL string) (string, error)
	GetTrack(id string) (*Track, error)
	SearchTrack(artistName, trackName string) (*Track, error)

	DetectAlbumID(albumURL string) (string, error)
	GetAlbum(id string) (*Album, error)
	SearchAlbum(artistName, albumName string) (*Album, error)
}
```

#### Track, Album

`Track` and `Album` are the main entities of the library. They represent a song and an album, respectively.
Each entity has a set of fields that contain information about the song or album, such as the name, artist, album, and URL.

``` golang
type Track struct {
    ID       string
    Title    string
    Artist   string
    URL      string
    Provider *Provider
}

type Album struct {
    ID       string
    Title    string
    Artist   string
    URL      string
    Provider *Provider
}
```

## Testing

For testing purposes, you can use the `RegistryOption`.

Implement `streaminx.Adapter` interface for provider you want to mock and pass it to `streaminx.NewRegistry` function:

``` golang
registry, err := streaminx.NewRegistry(
    ctx,
    streaminx.Credentials{},
    streaminx.WithProviderAdapter(streaminx.Apple, appleAdapterMock),
    streaminx.WithProviderAdapter(...)
}
```

Or if you want to go further and mock HTTP calls:

``` golang
registry, err := streaminx.NewRegistry(
    ctx,
    streaminx.Credentials{},
    streaminx.WithAppleAPIURL(appleAPIServerMock.URL),
    streaminx.WithAppleWebPlayerURL(appleWebPlayerServerMock.URL),
    streaminx.WithSpotifyAuthURL(spotifyAuthServerMock.URL),
    streaminx.WithSpotifyAPIURL(spotifyAPIServerMock.URL),
    streaminx.WithYandexAPIURL(yandexMockServer.URL),
    streaminx.WithYoutubeAPIURL(youtubeMockServer.URL),
}
```

To mock translator calls you also have option too:

``` golang
registry, err := streaminx.NewRegistry(
    ctx,
    streaminx.Credentials{},
    streaminx.WithGoogleTranslatorURL(translatorMockServer.URL),
}
```

## Why do we need translator?

Some streaming services translate or transliterate the names of artists.
So when we search for a track or album, we need to translate the artist's name to the language of the service.
For example Spotify doesn't allow non-latin characters in artist names. If we have a Yandex Music track by the artist "Дельфин" we need to make it "Dolphin" to find it on Spotify.


## Contribution and development

Contributions are welcome. It would be great if you could help us to add more providers or languages to the library.

To run the test and linter use the following commands:

```bash
make test
make lint
```

