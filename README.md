# Streamnx

Streamnx — a unified music streaming catalog library.

Streamnx provides a common interface for music streaming services.
Use a single API to parse release links, search tracks and albums, 
and retrieve release metadata across multiple streaming providers.

The library focuses on release-level data and cross-service links. 
It does not handle playback, playlists, recommendations, or user account data.

Typical use cases include:
* link sharing apps
* playlist migration tools
* release matching across streaming services
* music catalog enrichment pipelines

- [Streaming providers list](#streaming-providers-list)
- [Installation](#installation)
- [Streaming API access](#streaming-api-access)
- [Credentials](#credentials)
- [Configuration](#configuration)
- [Usage](#usage)
- [API reference](#api-reference)
    - [Catalog](#catalog)
    - [Track, album, and search API](#track-album-and-search-api)
    - [Provider](#provider)
    - [Link](#link)
- [Testing](#testing)
- [Contribution and development](#contribution-and-development)

## Streaming providers list

The library supports the following music streaming providers:

- Apple Music
- Bandcamp
- Deezer
- Soundcloud
- Spotify
- Yandex Music
- YouTube (compatible with YouTube Music)

## Installation

You can install the Streamnx library by running the following command in your terminal:

``` bash
go get github.com/GeorgeGorbanev/streamnx/v2
```

To include lib in your project, import it in your Go source file:

``` golang
import "github.com/GeorgeGorbanev/streamnx/v2"
```

## Streaming API access

Music streaming services differ widely in how openly they expose release data.
Some provide a free public API without registration, some require application
registration, some charge for access, some provide access only by individual
agreement, and some provide no public API at all.

Because of that, using every provider may require credentials, paid requests, or
a proxy for HTML scraping or unofficial APIs. Streamnx is oriented toward
amateur and enthusiast projects: for free or low-cost APIs it asks for
credentials, while expensive or unavailable APIs are handled with web scraping
where possible. For scraping integrations, you may still need a proxy and
anti-bot handling depending on your environment.

Provider availability, from easiest to most restricted:

1. Deezer: official API, no registration, free.
2. Spotify and YouTube: official APIs, application registration required, free
   within quota limits.
3. Apple Music: official API with developer registration and a paid Apple
   Developer Program membership; Streamnx currently uses the free web-player
   token from the public web player, with its limits.
4. Yandex Music: unofficial API only, no registration.
5. Soundcloud and Bandcamp: official API access is closed or granted by
   individual agreement, so Streamnx uses HTML scraping; these integrations may
   require proxies or anti-bot handling.

(!) Contributions are welcome. If you have access to closed APIs, or can add
alternative integrations for streaming services with free official APIs, that is
useful for the project.

## Credentials

You can get the following credentials for free:

1. YouTube API key from the [Google Cloud Console](https://console.cloud.google.com/apis/credentials).
2. Spotify Client ID and Client Secret from the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard).

If you are not going to work with these streaming providers, you can skip
registration. All streaming providers are optional.

## Configuration

To use integrations, initialize `Catalog` with only the providers your
application needs:

``` golang
catalog, err := streamnx.NewCatalog(
    streamnx.WithApple(),
    streamnx.WithBandcamp(),
    streamnx.WithDeezer(),
    streamnx.WithSpotify(streamnx.SpotifyCredentials{
        ClientID:     "[your spotify client id]",
        ClientSecret: "[your spotify client secret]",
    }),
    streamnx.WithSoundcloud(),
    streamnx.WithYandex(),
    streamnx.WithYoutube(streamnx.YoutubeCredentials{
        APIKey: "[your youtube api key]",
    }),
)
if err != nil {
    // Handle error
}
```

## Usage

Here is an example of how to get Spotify track candidates for an Apple Music link.
The application decides how to rank the candidates and what fallback to show when
there are no confident matches.

``` golang
func convertAppleTrackToSpotify(ctx context.Context, catalog *streamnx.Catalog, link string) (streamnx.SearchTrack, bool, error) {
    parsedLink, err := catalog.ParseLink(link)
    if err != nil {
        return streamnx.SearchTrack{}, false, err
    }

    if parsedLink.ReleaseType != streamnx.ReleaseTypeTrack {
        return streamnx.SearchTrack{}, false, fmt.Errorf("expected track link, got %s", parsedLink.ReleaseType)
    }

    track, err := catalog.FetchTrack(ctx, parsedLink.Provider, parsedLink.ReleaseID)
    if err != nil {
        return streamnx.SearchTrack{}, false, err
    }

    candidates, err := catalog.SearchTracks(ctx, streamnx.Spotify, streamnx.SearchQuery{
        Artist: track.Artist,
        Title:  track.Title,
    })
    if err != nil {
        return streamnx.SearchTrack{}, false, err
    }

    if len(candidates) == 0 {
        return streamnx.SearchTrack{}, false, nil
    }

    return candidates[0], true, nil
}
```

## API reference

#### Catalog

`Catalog` is the main entry point of the library.

The purpose of the `Catalog` is to provide a unified interface for working with streaming services by HTTP API.

#### Track, album, and search API

The fetch API uses separate models for tracks and albums. Search uses separate
candidate models, so applications can keep fetched source releases distinct from
destination search results:

``` golang
track, err := catalog.FetchTrack(ctx, provider, trackID)

trackCandidates, err := catalog.SearchTracks(ctx, provider, streamnx.SearchQuery{
    Artist: "Artist Name",
    Title:  "Track Title",
})

album, err := catalog.FetchAlbum(ctx, provider, albumID)

albumCandidates, err := catalog.SearchAlbums(ctx, provider, streamnx.SearchQuery{
    Artist: "Artist Name",
    Title:  "Album Title",
})
```

`FetchTrack` returns `Track`, `FetchAlbum` returns `Album`, `SearchTracks`
returns `[]SearchTrack`, and `SearchAlbums` returns `[]SearchAlbum`.

All four models expose `ID`, `Title`, `Artist`, `URL`, `CoverURL`, `Provider`,
`Creator`, and `Description`. The release type is implied by the concrete Go
type. The models also expose fields specific to their role:

- `Track` adds `AlbumID`, `AlbumTitle`, `Duration`, and `ReleaseDate`.
- `Album` adds `Label`, `ReleaseDate`, and `TrackIDs`.
- `SearchTrack` adds `AlbumID` and `AlbumTitle`.
- `SearchAlbum` has only the common fields.

`Duration` is expressed in seconds. `ReleaseDate` contains separate `Year`,
`Month`, and `Day` components; unavailable components remain zero. `TrackIDs`
contains provider-side track IDs when the provider exposes an album track list.
These metadata fields are optional and provider-dependent.

`Creator` contains raw provider-side creator, owner, uploader, channel, or label
values when a provider exposes them. `Description` contains raw provider
descriptions. These fields are optional and provider-dependent: they can be
empty, and they are not normalized artist metadata. Consumers that need
conversion-specific matching should interpret these raw fields in application
code.

Search methods return candidates in provider-defined order. If a search
completes but finds no results, it returns a non-nil empty slice and a nil
error. `SearchQuery` requires at least one non-empty field: `Artist` or
`Title`. They do not choose the best conversion target. Matching, ranking,
fuzzy search decisions, and user-facing fallback behavior belong in consumers
of the library. Streamnx does not apply a cross-provider client-side result
limit; consumers can trim candidate lists after search if their product flow
needs it.

`FetchTrack`, `FetchAlbum`, and `Uncloak` require a non-empty ID. They return
`ErrInvalidID` for an empty or whitespace-only ID. A fetch returns
`ErrNotFound` when the provider reports that no release exists for the supplied
ID. Searches use `ErrInvalidSearchQuery` only when both query fields are empty.

#### Provider

`ReleaseProvider` represents a music streaming service code.

``` golang
for _, provider := range streamnx.Providers() {
    fmt.Println(provider)
}
```

`Providers` returns a copy of the supported-provider list. The caller may
modify the returned slice without affecting Streamnx.

Known providers are exported as constants:

``` golang
provider := streamnx.Apple
// => "ap"
```

#### Link

`Link` struct represents a parsed link to a track or album on a streaming service. 

Useful to extract the ID and provider from the link. 

``` golang
type Link struct {
    URL         string
    Provider    ReleaseProvider
    ReleaseID   string
    ReleaseType ReleaseType
}
```

It is returned by the catalog-scoped `ParseLink` method. Parsing only recognizes
providers that were registered in that catalog:

``` golang          
link, err := catalog.ParseLink("https://music.apple.com/us/album/song-name/1234?i=4567")
// => Link{
//      URL: "https://music.apple.com/us/album/song-name/1234?i=4567", 
//      Provider: streamnx.Apple,
//      ReleaseID: "us-4567",
//      ReleaseType: streamnx.ReleaseTypeTrack,
//  }, nil
```

`ReleaseType` has three values: `ReleaseTypeTrack`, `ReleaseTypeAlbum`, and
`ReleaseTypeCloak`. `ParseLink` returns `Link` by value.

Use `FetchTrack` for parsed track links and `FetchAlbum` for parsed album links.
Use `Uncloak` for parsed cloak links; it resolves the cloak and returns the
resolved `ReleaseType` plus release ID:

``` golang
link, err := catalog.ParseLink(rawURL)
if err != nil {
    // Handle error
}

switch link.ReleaseType {
case streamnx.ReleaseTypeTrack:
    track, err := catalog.FetchTrack(ctx, link.Provider, link.ReleaseID)
case streamnx.ReleaseTypeAlbum:
    album, err := catalog.FetchAlbum(ctx, link.Provider, link.ReleaseID)
case streamnx.ReleaseTypeCloak:
    releaseType, releaseID, err := catalog.Uncloak(ctx, link.Provider, link.ReleaseID)
}
```

## Testing

For testing purposes, use `CatalogOption` values and HTTP-client
overrides to point provider clients at mock servers:

``` golang
catalog, err := streamnx.NewCatalog(
    streamnx.WithApple(
        streamnx.WithAppleAPIURL(appleAPIServerMock.URL),
        streamnx.WithAppleWebPlayerURL(appleWebPlayerServerMock.URL),
        streamnx.WithAppleHTTPClient(appleHTTPClient),
    ),
    streamnx.WithSoundcloud(
        streamnx.WithSoundcloudAPIURL(soundcloudAPIServerMock.URL),
        streamnx.WithSoundcloudWebURL(soundcloudWebServerMock.URL),
    ),
    streamnx.WithSpotify(
        streamnx.SpotifyCredentials{},
        streamnx.WithSpotifyAuthURL(spotifyAuthServerMock.URL),
        streamnx.WithSpotifyAPIURL(spotifyAPIServerMock.URL),
    ),
    streamnx.WithYandex(streamnx.WithYandexAPIURL(yandexMockServer.URL)),
    streamnx.WithYoutube(
        streamnx.YoutubeCredentials{},
        streamnx.WithYoutubeAPIURL(youtubeMockServer.URL),
        streamnx.WithYoutubeHTTPClient(youtubeHTTPClient),
    ),
)
```

## Contribution and development

Contributions are welcome, especially fixes for provider integrations, docs,
tests, and new release-data providers that fit the library's focus on links plus
track and album metadata.

To run the tests and linter, use the following commands:

```bash
make test
make lint
```
