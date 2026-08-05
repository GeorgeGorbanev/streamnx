# AGENTS.md

This file applies to the entire `streamnx` repository. It is the working guide
for agents changing this module. Keep it aligned with the code when architecture,
public contracts, build commands, or provider conventions change.

## Project in one paragraph

Streamnx is a Go library that presents several music services as one catalog API.
It parses provider release links, fetches track and album metadata, searches for
track and album candidates, and resolves supported short/cloaked links. It does
not play media, manage users or playlists, recommend music, or decide which
cross-provider search result is the best match. The module path is
`github.com/GeorgeGorbanev/streamnx/v2`; the current Go version is declared in
`go.mod` and mirrored in CI.

The supported providers are:

| Provider | Code | Integration shape | Catalog release ID |
| --- | --- | --- | --- |
| Apple Music | `ap` | catalog API using a token scraped from the web player | composite `storefront-id` |
| Bandcamp | `bc` | release-page HTML/LD+JSON plus the public search endpoint | composite `artist-slug:release-slug` |
| Deezer | `dz` | public API plus Deezer cloak redirects | numeric string |
| Spotify | `sf` | Web API with client-credentials OAuth | opaque Spotify ID |
| SoundCloud | `sc` | hydrated web pages plus API calls using a scraped client ID | composite `user-slug:release-slug` |
| Yandex Music | `ya` | unofficial public API | albums use a numeric string; tracks use composite `album-id:track-id` |
| YouTube | `yt` | Data API; videos are tracks and playlists are albums | video or playlist ID |
| YouTube Music | `ym` | anonymous YouTube Music web-client Innertube API | video ID for tracks; composite `b:<browse-id>` or `p:<playlist-id>` for albums |

## Repository map

- `catalog.go` is the public orchestration layer. It owns the private adapter
  interface, explicit provider registration, input validation, dispatch, and
  public empty-result guarantees.
- `catalog_options.go` exposes provider-specific functional options. This is the
  composition root: public credentials and HTTP overrides are converted into
  internal client options, then a client and adapter are registered.
- `release.go` re-exports the canonical release models, provider/type constants,
  and shared errors from `internal/release` by aliases. Do not create parallel
  public model definitions.
- `internal/release/release.go` is the canonical domain model and provider list.
- `internal/release/compositekey` validates and serializes provider IDs that have
  multiple parts. Use it instead of ad hoc `strings.Split`/`Join` logic.
- `internal/release/duration` contains shared duration conversion helpers.
- `internal/release/extid` validates and normalizes ISRC recording identifiers
  and UPC-A/EAN-13 album identifiers. Keep identifier-specific functions and
  variables prefixed with `isrc` or `upc` within this shared package.
- `internal/providers/<provider>/` contains one isolated integration:
  - `client.go`: HTTP, authentication, provider query syntax, status handling,
    and decoding into provider-native response types.
  - `adapter.go`: implementation of the catalog adapter contract and mapping
    from provider-native values to `internal/release` values.
  - `response.go`: private structs mirroring external JSON/HTML data.
  - `link.go`: provider URL recognition and canonical URL construction.
  - `key.go`: optional composite-ID schema and helpers.
  - `client_test.go` and `adapter_test.go`: transport/parser tests and mapping
    tests respectively.
- `catalog_test.go` tests catalog-level validation, registration, dispatch, and
  result contracts with a mocked adapter.
- `tests/` contains public-API integration tests. They run complete catalog ->
  adapter -> client flows against local fixture servers.
- `tests/fixtures/server.go` is the shared embedded-fixture HTTP server;
  `tests/fixtures/responses/` stores captured provider responses.
- `README.md` is user-facing API documentation and part of the public contract.
- `Makefile`, `.golangci.yml`, and `.github/workflows/ci.yml` define the accepted
  validation path.

## Architecture and dependency direction

The intended call path is:

```text
consumer
  -> streamnx.Catalog
     -> private catalog adapter interface
        -> provider Adapter
           -> provider-private adapterClient interface
              -> provider Client
                 -> external HTTP service
```

The return path reverses the transformation:

```text
HTTP response -> response.go types -> Adapter normalization
              -> internal/release model -> public type alias
```

Keep these boundaries intact:

1. The root package knows the unified domain and the provider constructors, but
   not provider response schemas or endpoint behavior.
2. A provider client knows transport and provider-native data, but does not
   construct public `streamnx` values.
3. A provider adapter owns semantic mapping: IDs, artists, cover selection,
   dates, durations, track lists, not-found translation, and unsupported
   operations.
4. Provider packages do not import each other.
5. Shared code belongs under `internal/release` only when at least two providers
   genuinely share the same rule. Do not erase provider differences merely to
   deduplicate a few lines.
6. Interfaces live at the consuming boundary. `adapter` is declared by the
   catalog; `adapterClient` is declared by each adapter. Keep interfaces small
   and mockable rather than exporting internal abstractions.

## Product and domain boundaries

Streamnx is deliberately release-centric. Keep these decisions in consumers,
not in this library, unless the public design is explicitly changed:

- fuzzy matching and normalization across providers;
- candidate scoring, ranking, deduplication, or “best match” selection;
- product-specific fallbacks or confidence thresholds;
- playback, playlists, recommendations, and user-account data;
- a global result limit layered on top of provider ordering.

Search candidates stay in the order returned by the provider. Preserve that
order. Do not silently sort, trim, or filter legitimate candidates in an
adapter. A provider may use a provider-side limit; that is transport behavior,
not a cross-provider ranking policy.

## Public API contracts

Treat the following behavior as compatibility-sensitive.

### Catalog construction and registration

- `NewCatalog` starts empty. Providers are opt-in through `WithApple`,
  `WithBandcamp`, `WithDeezer`, `WithSpotify`, `WithSoundcloud`, `WithYandex`,
  `WithYoutube`, and `WithYoutubeMusic`.
- A catalog recognizes and dispatches only providers registered in that catalog.
- A nil `CatalogOption`, nil provider option, or nil public HTTP client override
  is ignored. Do not turn those cases into panics.
- Duplicate provider registration returns an error wrapping
  `ErrDuplicateProvider`.
- Unknown providers and nil adapters are rejected with an error wrapping
  `ErrInvalidProvider`.
- `Providers()` returns a clone. Never expose the mutable canonical provider
  slice.

### Links and IDs

- `Catalog.ParseLink` returns `ErrUnknownLink` when none of the registered
  adapters recognizes the input.
- Link parsing intentionally accepts provider URLs embedded in surrounding text
  where the provider regex permits it. Preserve tested query-string, locale,
  and alternate-link variants.
- `Link.URL` retains the original input; `ReleaseID` is the provider-normalized
  ID used by fetch/uncloak calls.
- Release types are track, album, and cloak. Only Deezer currently supports
  uncloaking; other adapters return `ErrUnsupportedOperation`.
- `FetchTrack`, `FetchAlbum`, and `Uncloak` reject empty or whitespace-only IDs
  with `ErrInvalidID` before calling an adapter. `FetchTracksByISRC` accepts a
  canonical 12-character ISRC or its hyphenated presentation form, normalizes
  it before dispatch, and returns `ErrInvalidID` for malformed values.
  `FetchAlbumsByUPC` accepts a check-digit-valid 12-digit UPC-A or 13-digit
  EAN-13. A 13-digit value with a leading zero is normalized to its equivalent
  12-digit UPC-A before dispatch.
- Composite IDs are an external contract. Changing their delimiter, part order,
  validation, or normalization is a breaking change. Route all creation and
  parsing through the package's `keyScheme`.

### ISRC fetch

- `FetchTracksByISRC` returns full `Track` objects because an ISRC can map to
  multiple provider catalog entries. Successful calls always return a non-nil
  slice, including zero results.
- Every adapter implements `FetchTracksByISRC`. Apple, Deezer, and Spotify
  support it; other adapters return `ErrUnsupportedOperation` before any
  provider request, following the same contract as unsupported `Uncloak` calls.
- ISRC identifies recordings, not albums. Do not add it to album models or
  album search; UPC/EAN is the corresponding product-level identifier.

### UPC fetch

- `FetchAlbumsByUPC` returns full `Album` objects because a UPC can map to
  multiple provider catalog entries. Successful calls always return a non-nil
  slice, including zero results.
- Every adapter implements `FetchAlbumsByUPC`. Apple, Deezer, and Spotify
  support it; other adapters return `ErrUnsupportedOperation` before any
  provider request.
- Spotify UPC search returns simplified album objects, so the client hydrates
  every result through the full album endpoint before returning to the adapter.
- Provider models may expose the same GTIN as 13-digit EAN with a leading zero
  or 12-digit UPC-A. Preserve provider metadata in model fields; catalog input
  normalization handles equivalent query representations.

### Search

- Search is valid when at least one of `Artist` or `Title` is non-empty after
  trimming for validation. Both empty yields `ErrInvalidSearchQuery`.
- Pass original artist/title query strings to the adapter; catalog validation
  does not rewrite them.
- Successful public searches always return a non-nil slice, including zero
  results. `Catalog` enforces this with `forceNonNil`.
- Adapter/client errors return a nil slice unless a provider explicitly maps a
  provider-level “not found” search response to a successful empty result.
- Track and album searches return `SearchTrack` and `SearchAlbum`, not fetched
  `Track` and `Album` objects.

### Release models

- Fetch methods and `ParseLink` return values, not pointers. Return the zero
  value alongside an error.
- `Track`, `Album`, `SearchTrack`, and `SearchAlbum` share the common fields
  `ID`, `Title`, `Artist`, `URL`, `AlternativeURL`, `CoverURL`, `Provider`,
  `Creator`, and `Description`.
- `Track` additionally owns `ISRC`, `AlbumID`, `AlbumTitle`, `Duration`, and
  `ReleaseDate`; `Album` owns `UPC`, `Label`, `ReleaseDate`, and `TrackIDs`;
  `SearchTrack` owns `ISRC`, `AlbumID`, and `AlbumTitle`; `SearchAlbum` owns
  `UPC`.
- `Duration` is integral seconds. Use shared conversion helpers and preserve
  their rounding semantics instead of truncating milliseconds or ISO-8601
  fractional seconds.
- `ReleaseDate` preserves partial precision. Unknown year/month/day components
  remain zero; do not invent dates. Spotify explicitly supports year-only and
  year-month values.
- Optional metadata is provider-dependent. Empty strings, zero dates/durations,
  and empty track lists are valid when the source omits data.
- `Creator` and `Description` are raw provider-side metadata, not normalized
  artist fields. Keep provider meaning intact.
- `TrackIDs` contains provider IDs in provider order. Skip only unusable empty
  IDs where the adapter already follows that policy.

### Errors

- Stable cross-provider sentinels are `ErrNotFound`, `ErrScrapingBlocked`,
  `ErrUnsupportedOperation`, `ErrInvalidProvider`, `ErrDuplicateProvider`,
  `ErrUnknownLink`, `ErrInvalidID`, and `ErrInvalidSearchQuery`.
- Clients may use private sentinels such as `errNotFound` and
  `errUnauthorized`. Adapters translate provider not-found conditions to
  `release.ErrNotFound`.
- Preserve error identity with `%w` and test stable classification using
  `errors.Is`/`require.ErrorIs`. Add context at layer boundaries, for example
  `failed to get track from spotify: %w`.
- Do not expose provider-private sentinels through the public API.
- Malformed success payloads and unexpected statuses are operational errors,
  not automatically `ErrNotFound`.

## Provider package conventions

### `client.go`

- Keep endpoint defaults and client functional options close to `Client`.
- Accept `context.Context` for every operation that can perform I/O and construct
  requests with `http.NewRequestWithContext`.
- Build query strings with `url.Values`; do not manually concatenate unescaped
  user input.
- Use the injected `*http.Client`. This supports tests, proxies, custom
  transports, and consumer timeouts. Do not call package-level `http.Get`.
- Close every non-nil response body. If a response is abandoned before a retry,
  close it before sending the next request.
- Check provider-specific status semantics explicitly. Some providers report
  missing releases in a successful JSON response rather than as HTTP 404.
- Decode only the response shape needed by the integration. Keep provider DTOs
  private and in `response.go` unless a small operation-local struct is clearer.
- Return provider-native structs to the adapter. Do not leak HTTP response bodies
  or root public models upward.
- Wrap errors with useful operation context and retain causes.

### Authentication and concurrency

Spotify tokens, the Apple web-player token, and the SoundCloud client ID are
shared mutable caches. Their locking and retry behavior is intentional:

- use `sync.RWMutex` for read-mostly credential state;
- double-check cached state after acquiring the write lock so concurrent callers
  do not all refresh it;
- on a 401, compare the rejected credential with the currently cached one before
  refreshing—another goroutine may already have replaced it;
- retry an unauthorized request at most once;
- propagate the caller's context into refresh requests;
- never hold or expose credentials in package globals, logs, fixtures, or error
  strings.

Changes in these areas require the provider's concurrency tests and a targeted
race run. Do not simplify the mutex/refetch sequence without proving equivalent
behavior under concurrent initial requests, expiry, and simultaneous 401s.

### `adapter.go`

- `Adapter` contains the smallest client interface required for the provider.
  Constructor injection enables mapping tests without HTTP.
- `ParseLink` tries supported release shapes and returns `("", "", false)` for
  non-matches. Parsing failures are not public errors.
- Fetch methods translate private not-found errors, then wrap all other client
  failures with provider and operation context.
- Construct domain structs in one visible literal. Populate every reliable field
  the provider exposes, but tolerate documented optional data.
- Put provider-specific fallback logic in small helpers: first artist, label,
  largest image, artwork template replacement, creator fallback, date parsing,
  and similar rules.
- Guard slices and nested optional data before indexing. A malformed upstream
  response must return a useful error or valid zero field, never panic.
- Preserve candidate and album-track order.
- Unsupported `Uncloak` methods return zero type, empty ID, and
  `release.ErrUnsupportedOperation`.

### `link.go` and `key.go`

- Compile stable regexes once as package variables.
- Keep recognition and canonical URL generation together.
- Test empty strings, wrong domains, missing IDs, query strings, locale paths,
  alternate official URL forms, and embedded URLs relevant to the provider.
- A composite key schema defines part names, regexes, delimiter, and order.
  Validate every part on both construction and parse.
- Use canonical provider URLs when an API response lacks a trustworthy public
  URL; otherwise preserve the authoritative URL returned by the provider.

### `response.go`

- Mirror external names with JSON tags while keeping Go names idiomatic.
- Avoid embedding unified-domain policy in DTO methods.
- Keep fields only when they are consumed or needed to model a tested response.
- Add fixture coverage when a schema change is based on an observed payload.

## Provider-specific behavior worth preserving

- Apple IDs include storefront because catalog calls are storefront-scoped.
  Artwork URL templates are expanded to a large JPEG. Token discovery and 401
  refresh scrape the public web player and are concurrency-safe.
- Bandcamp and SoundCloud IDs depend on URL slugs, not numeric values hidden in
  scraped payloads. Their link/key round trips are therefore central behavior.
- Deezer may represent not-found as a successful response with an error payload.
  It is also the sole cloak implementation: follow the redirect without losing
  `Location`, extract `dest`, then accept only track or album destinations.
- Spotify uses the first artist for the unified singular `Artist` field, selects
  the largest cover by pixel area, preserves partial release dates, and refreshes
  OAuth tokens safely under concurrency.
- SoundCloud prefers publisher artist metadata over uploader name, upgrades
  known artwork-size suffixes to `original`, prefers full duration, and refreshes
  a scraped API client ID after one unauthorized response.
- Yandex track IDs include the album ID as `album-id:track-id`; fetching selects
  that exact album association and its cover. Missing artists are tolerated;
  missing required album data is an adaptation error. Search-level not-found
  maps to an empty candidate slice.
- YouTube videos/playlists do not reliably provide normalized artist/album
  fields, so channel ownership belongs in `Creator`. Search hydrates each search
  result through a detail request; preserve result order and full fixture
  coverage for those follow-up calls.
- YouTube Music uses the anonymous `WEB_REMIX` Innertube client: `player` plus
  `next` for exact track and album metadata, filtered `search` for normalized
  candidates, and `browse` for album metadata and tracks. Album IDs preserve
  whether a link used a browse or playlist ID; available counterpart links are
  exposed through `AlternativeURL`. It owns `music.youtube.com` links; ordinary
  YouTube must not recognize that subdomain.

## Functional options and public configuration

Follow the existing two-level option pattern:

```go
type ExampleOption func(*exampleOptions)

type exampleOptions struct {
    client []example.ClientOption
}

func WithExample(opts ...ExampleOption) CatalogOption {
    return func(catalog *Catalog) error {
        cfg := exampleOptions{}
        for _, opt := range opts {
            if opt == nil {
                continue
            }
            opt(&cfg)
        }
        return catalog.register(Example, example.NewAdapter(example.NewClient(cfg.client...)))
    }
}
```

Public options should expose consumer needs—credentials, endpoint overrides, and
HTTP-client injection—without exporting provider internals. Ignore nil public
HTTP-client overrides so defaults remain usable. Keep real service URLs as
internal defaults; endpoint override options exist primarily for tests and
controlled deployments.

## Go style

- Follow the standard library first. The only regular third-party test
  dependency is `testify`.
- Formatting is stricter than plain `gofmt`: CI checks `gofumpt` and `gci`.
  Import groups are standard library, third-party, then the
  `github.com/GeorgeGorbanev` prefix.
- Use short, consistent receivers (`c` for client, `a` for adapter). Avoid a
  receiver name that conflicts with the type's role.
- Keep exported API surface small. Provider DTOs, client operations, sentinels,
  parsers, and helpers remain unexported unless consumers truly need them.
- Prefer early returns and explicit status/error branches. The enabled linters
  enforce receiver naming, error naming, indentation/early-return conventions,
  unused code, vet/static analysis, and preallocation opportunities.
- Preallocate result slices when size or capacity is known.
- Use zero values to represent absent optional metadata. Avoid pointer-helper
  layers solely to distinguish absent scalar values unless the public model is
  intentionally redesigned.
- Avoid hidden network calls, background goroutines, global mutable state,
  panics, and logging from library code.
- Comments should explain provider quirks, external constraints, or non-obvious
  invariants. Do not narrate straightforward code.
- Do not edit `.env`, read secrets from it for tests, or commit credentials.
  Tests must be hermetic.

## Testing strategy

Every behavioral change should be tested at the narrowest useful layer and, when
public behavior changes, at the public integration layer too.

### Catalog tests

Use `catalog_test.go` for provider registration, validation order, adapter
dispatch, error identity, registered-provider parsing, and public non-nil empty
search results. The local `adapterMock` uses `testify/mock`; assert expectations.

### Adapter tests

Adapter tests are in the provider package, not an external `_test` package, so
they can exercise private DTOs and sentinels. Mock `adapterClient` and verify:

- every recognized and rejected link form;
- exact provider-native method arguments;
- complete domain-model mapping;
- not-found translation and contextual wrapping;
- optional/malformed fields and fallback behavior;
- empty candidate behavior and order;
- `Uncloak` behavior.

Prefer table-driven tests when cases share setup. On errors, assert both the zero
or nil result contract and error classification/text as appropriate. Always call
`AssertExpectations` for mocks.

### Client tests

Use `httptest.Server` or a controlled transport. Assert method, path, query,
headers, body, redirect behavior, status mapping, response decoding, refresh,
and retry count. Use `t.Context()` rather than `context.Background()` in tests.
For cached credentials, retain tests for initial fetch, expiry/401 refresh,
network failure, concurrent callers, double-checking under the write lock, and a
single retry ceiling.

### Public fixture tests

Tests under `tests/` import `github.com/GeorgeGorbanev/streamnx/v2` as a consumer
would. They should use only public constructors/options and local embedded
fixtures. The shared fixture server matches routes by method, path, and expected
query values and can run extra request assertions.

- Never call live provider services in the test suite.
- Name fixtures `<provider>_<operation>_<case>_<status>.<ext>` consistently.
- Prefer realistic captured payloads, scrubbed of secrets and personal data.
- Add or update routes when an operation performs multiple requests, such as
  token acquisition or YouTube candidate hydration.
- Assert complete public structs when practical; this catches accidental field
  loss and normalization drift.
- Include not-found and unsupported-operation coverage for public behavior.

## Commands

Run commands from the repository root (`streamnx/`).

```bash
# Root catalog package only
go test .

# One provider while iterating
go test ./internal/providers/spotify

# One public fixture test
go test ./tests -run TestSpotifyCatalogFetchTrack

# Canonical full suite used by CI
make test

# Canonical lint/format check used by CI
make lint
```

CI currently installs `golangci-lint` v2.7.0 and runs `make lint` followed by
`make test`. The linter runs with module downloads in readonly mode and includes
`gci` and `gofumpt` format checks. Do not assume plain `gofmt` is the complete
formatting gate. The local linter binary must itself be built with a Go version
at least as new as the version targeted by `go.mod`.

For changes to credential caching or request concurrency, also run a focused
race check before the full suite:

```bash
go test -race ./internal/providers/spotify
go test -race ./internal/providers/apple
go test -race ./internal/providers/soundcloud
```

Use `go mod tidy` only when imports or dependencies actually change, and inspect
both `go.mod` and `go.sum` afterward. Do not introduce a dependency for behavior
that is small and clear with the standard library.

## Change workflows

### Fixing an existing provider

1. Reproduce the issue at the client or adapter boundary.
2. Add the smallest realistic client payload/fixture that demonstrates it.
3. Fix transport/schema handling in `client.go`/`response.go`, or unified mapping
   in `adapter.go`; do not blur the layers.
4. Verify link/key changes independently because IDs are public compatibility.
5. Run the provider's client and adapter tests, its public fixture tests, then
   `make test` and `make lint`.

### Adding a provider

1. Confirm it fits the release-link/track/album scope.
2. Add its provider code to `internal/release/release.go` and the canonical
   provider list; keep codes short, stable, and unique.
3. Create `internal/providers/<provider>` with client, adapter, response, link,
   optional key files, and focused tests following neighboring packages.
4. Implement the complete private catalog adapter interface, returning
   `ErrUnsupportedOperation` where necessary.
5. Add public constant aliases and shared errors only through root `release.go`.
6. Add public credentials/options and explicit registration in
   `catalog_options.go`, including endpoint and HTTP-client overrides needed for
   hermetic tests.
7. Extend catalog tests for registration, nil options, parsing, duplicates, and
   dispatch.
8. Add captured response fixtures and public integration tests for track/album
   fetch, not-found, both searches, and cloak/unsupported behavior.
9. Update the README provider list, configuration example, access/credential
   notes, and any affected API reference.
10. Run all validation commands. A provider is not complete with adapter-only or
    live-network-only tests.

### Changing the release model or public API

1. Treat root aliases, constructors, option signatures, provider codes, ID
   formats, sentinels, and documented result semantics as public API.
2. Change the canonical model under `internal/release`, then audit every adapter
   and every exact struct assertion.
3. Add catalog contract tests for validation or result-shape changes.
4. Update README examples and reference text in the same change.
5. Avoid provider-specific fields in the common model unless their meaning is
   coherent for multiple providers and useful to consumers.

### Changing link parsing or composite IDs

1. Preserve all currently tested official variants.
2. Add positive and negative table cases before changing regexes.
3. Verify parse -> ID -> fetch-argument behavior and canonical URL generation.
4. For composite IDs, test invalid part count, invalid values, unknown/missing
   parts, copy semantics, and keys from the wrong schema.
5. Call out compatibility impact if an already emitted ID could change.

## Final review checklist

Before handing off a change:

- inspect `git status` and `git diff`; preserve unrelated user changes;
- confirm code was changed in the correct layer;
- confirm provider order and provider-returned candidate/track order are intact;
- confirm contexts reach every HTTP request and all response bodies are closed;
- confirm errors wrap causes and public sentinel identity is preserved;
- confirm empty searches are non-nil at the public catalog boundary;
- confirm optional upstream data cannot cause an index panic;
- confirm auth/client-ID refresh remains concurrency-safe and retries once;
- confirm tests are local, deterministic, and credential-free;
- confirm README is updated for public behavior;
- run targeted tests, `make test`, and `make lint`;
- review formatting and import grouping, not just compilation.

If code and this guide disagree, first determine whether the code expresses a
provider-specific exception or the guide is stale. Preserve intentional,
well-tested exceptions and update this file when the stable project convention
has changed.
