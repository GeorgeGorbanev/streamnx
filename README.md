# Streamnx

Streamnx is a library that unifies interactions with various music streaming links into a single system.
With Streamnx, you can integrate platforms such as Apple Music, Spotify, YouTube and Yandex Music into your applications using a unified interface for searching and retrieving data about tracks and albums.

- [Motivation](#motivation)
- [Providers Supported](#providers-supported)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API reference](#api-reference)
    - [Registry](#registry)
    - [Provider](#provider)
    - [EntityType](#entitytype)
    - [Entity](#entity)
    - [Link](#link)
- [Testing](#testing)
- [Why do we need translator?](#why-do-we-need-translator)
- [Contribution and development](#contribution-and-development)

## Motivation

The main user of this library is the Telegram bot [Vibeshare](https://t.me/vibeshare_bot).
This bot helps users to share song links from the streaming services they use, enabling others to listen to the same songs on different platforms.
Therefore, the primary use case for this library is *converting links* from one service to another.

## Supported services

The library supports the following music streaming services
- Apple Music
- Spotify
- Yandex Music
- YouTube (also YouTube Music)

## Installation

You can install the Streamnx library by running the following command in your terminal:

``` bash
go get github.com/GeorgeGorbanev/streamnx
```

To include lib in your project, import it in your Go source file:

``` golang
import "github.com/GeorgeGorbanev/streamnx"
```

## Configuration

To configure streamnx, you need to set up the necessary credentials and API keys for the supported music streaming services.
Here are the steps to configure the library:

1) *Google Translator API.* Obtain the Google Translator API key and project ID from the [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2) *YouTube API*. Obtain the YouTube API key from the [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
3) *Spotify*. Register your application and obtain the Client ID with Client Secret on the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard).
4) *Build registry*. When you have all the necessary credentials, you can initialize the streamnx *registry* with the following code.

``` golang
package main

import (
    "context"

    "github.com/GeorgeGorbanev/streamnx"
)

func main() {
    ctx := context.Background()
    registry, err := streamnx.NewRegistry(ctx, streamnx.Credentials{
        GoogleTranslatorAPIKeyJSON: "[your google translator api key json]",
        GoogleTranslatorProjectID:  "[your google translator project id]",
        YoutubeAPIKey:              "[your youtube api key]",
        SpotifyClientID:            "[your spotify client id]",
        SpotifyClientSecret:        "[your spotify client secret]",
    })
    if err != nil {
        // Handle error
    }
    defer registry.Close()

    // use the registry to fetch or search for tracks and albums
}
```
Replace the placeholders with your actual API keys and credentials.

## Usage

Here is an example of how to convert a link from Apple to Spotify:

``` golang

func appleTrackToSpotify(ctx context.Context, link string) (string, error) {
    parsedLink, err := streamnx.ParseLink(link)
    if err != nil {
        // Handle error
    }
        
    track, err := registry.Fetch(ctx, streamnx.Apple, streamnx.Track, parsedLink.ID)
    if err != nil {
        // Handle error
    }

    converted, err := registry.Search(ctx, streamnx.Spotify, streamnx.Track, track.Artist, track.Name)
    if err != nil {
        // Handle error
    }

    return converted.URL
}

```

## API reference

#### Registry

`Registry` struct is the main entry point of the library. 

The purpose of the `Registry` is to provide a unified interface for working with streaming services by HTTP API.

It implements two main methods:
``` golang
// Fetch(...) – allows to get entities by their ID
entity, err := registry.Fetch(ctx, provider, entityType, entityID)

// Search(...) – allows to search for entities by name and artist
entity, err := registry.Search(ctx, provider, entityType, entityArtist, entityTitle) 
```

This methods requires to specify the *provider*, the *entity type* and *identifiers* explained below. 

#### Provider

`Provider` represents a music streaming service, implemented as an enum. 

All providers are accessible via the `Providers` enum:
``` golang
for _, provider := range streamnx.Providers {
    fmt.Println(provider.Name())
}

```

It implements following methods:

``` golang
p := streamnx.Apple

p.Name()  
// => "Apple" 
// human readable name of the provider

code := p.Code()
// => "ap"
// short code of the provider, useful for runtime provider definition

regions := p.Regions()
// => []string{"us", "es", "fr", "ru", ... }
// optional region codes for the provider, used for region-specific requests and referenced in the URL 
   
trackID, err := p.DetectTrackID("https://music.apple.com/us/album/song-name/1234?i=4567")
// => "4567", nil
// extract track ID from the link

albumID, err := p.DetectAlbumID("https://music.apple.com/us/album/album-name/1234")
// => "1234", nil
// extract album ID from the link


```

When you need to define a provider in runtime, you can use the `FindProviderByCode(string)` method:

``` golang
provider := streamnx.FindProviderByCode("ap")
// => streamnx.Apple

```

#### EntityType

`EntityType` simple string enum that represents the type of entity you want to fetch or search for. 

For now, it has two values: `Track` and `Album`.

``` golang
streamnx.Track
// => "track"

streamnx.Album
// => "album"
```

#### Entity

`Entity` struct implements unified representation of tracks and albums. 

This struct is returned by the `Fetch` and `Search` methods of the `Registry`.

``` golang
type Entity struct {
	ID       string
	Title    string
	Artist   string
	URL      string
	Provider *Provider
	Type     EntityType
}
```

#### Link

`Link` struct represents a parsed link to a track or album on a streaming service. 

Useful to extract the ID and provider from the link. 

``` golang
type Link struct {
    URL      string
    Provider *Provider
    Type     EntityType
    ID       string
}
```

It is returned by the `ParseLink` method:

``` golang          
link, err := streamnx.ParseLink("https://music.apple.com/us/album/song-name/1234?i=4567")
// => Link{
//      URL: "https://music.apple.com/us/album/song-name/1234?i=4567", 
//      Provider: streamnx.Apple,
//      Type: streamnx.Track,
//      ID: "4567",
//  }, nil
```

## Testing

For testing purposes, you can use the `RegistryOption`.

Implement `streamnx.Adapter` interface for provider you want to mock and pass it to `streamnx.NewRegistry` function:

``` golang
registry, err := streamnx.NewRegistry(
    ctx,
    streamnx.Credentials{},
    streamnx.WithProviderAdapter(streamnx.Apple, appleAdapterMock),
    streamnx.WithProviderAdapter(...)
}
```

Or if you want to go further and mock HTTP calls:

``` golang
registry, err := streamnx.NewRegistry(
    ctx,
    streamnx.Credentials{},
    streamnx.WithAppleAPIURL(appleAPIServerMock.URL),
    streamnx.WithAppleWebPlayerURL(appleWebPlayerServerMock.URL),
    streamnx.WithSpotifyAuthURL(spotifyAuthServerMock.URL),
    streamnx.WithSpotifyAPIURL(spotifyAPIServerMock.URL),
    streamnx.WithYandexAPIURL(yandexMockServer.URL),
    streamnx.WithYoutubeAPIURL(youtubeMockServer.URL),
}
```

To mock translator calls you also have option too:

``` golang
registry, err := streamnx.NewRegistry(
    ctx,
    streamnx.Credentials{},
    streamnx.WithGoogleTranslatorURL(translatorMockServer.URL),
}
```

## Why do we need translator?

Some streaming services translate or transliterate the names of artists.
So when we search for a track or album, we need to translate the artist's name to the language of the service.
For example Spotify doesn't allow non-latin characters in artist names. If we have a Yandex Music track by the artist "Дельфин" we need to make it "Dolphin" to find it on Spotify.


## Contribution and development

Contributions are welcome. It would be great if you could help us to add more providers (e.g. Deezer) or entities (e.g. artist, playlist) to the library.

To run the test and linter use the following commands:

```bash
make test
make lint
```

