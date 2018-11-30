# go-appcast

[![Travis (.org)](https://img.shields.io/travis/victorpopkov/go-appcast.svg)](https://travis-ci.org/victorpopkov/go-appcast)
[![Codecov](https://img.shields.io/codecov/c/github/victorpopkov/go-appcast.svg)](https://codecov.io/gh/victorpopkov/go-appcast)
[![Go Report Card](https://goreportcard.com/badge/github.com/victorpopkov/go-appcast)](https://goreportcard.com/report/github.com/victorpopkov/go-appcast)
[![GoDoc](https://godoc.org/github.com/victorpopkov/go-appcast?status.svg)](https://godoc.org/github.com/victorpopkov/go-appcast)

An extendable library which provides functionality for working with different
appcasts. It can work in both ways: retrieve data from an already existing
appcast (unmarshal) or generate an appcast from the provided data (marshal).

- [What "appcast" means?](#what-appcast-means)
- [What this library does?](#what-this-library-does)
- [Providers](#providers)
  - [Sparkle RSS Feed](#sparkle-rss-feed)
  - [SourceForge RSS Feed](#sourceforge-rss-feed)
  - [GitHub Atom Feed](#github-atom-feed)
- [Sources](#sources)
- [Outputs](#outputs)

## What "appcast" means?

The word "appcast" is usually referred to a remote web page providing
information about software updates. This kind of pages is usually created for
different software update frameworks like [Sparkle][] or generated by different
services that distribute applications ([GitHub][], [SourceForge][] and etc.).

There are plenty of different methods available, but originally "appcasting" was
the practice of using an [RSS enclosure][] to distribute updates and release
notes.

## What this library does?

As you can find plenty of different ways how vendors distribute their software
updates it becomes pretty tedious to extract useful information from their
appcasts. Especially, considering the idea that sometimes you would like to
transpile one appcast type into another or simply make changes into an already
existing one.

Here comes this library handy. It provides the core functionality for working
with the supported providers in a reliable and consistent way. In addition, it
was designed to be extendable which enables you adding more features or extend
the supported providers, sources and even outputs to match your needs.

### Features

- [x] Designed to be extendable
- [x] Detect release stability from the semantic version
- [x] Different outputs to save to
- [x] Different sources to load from
- [x] Filter releases by stability, title, media type or download URL
- [x] Guess the supported provider
- [x] Sort releases by version
- [ ] Transpilation from one provider into another

## Providers

Out of the box, 3 providers are supported:

- [GitHub Atom Feed](#github-atom-feed)
- [SourceForge RSS Feed](#sourceforge-rss-feed)
- [Sparkle RSS Feed](#sparkle-rss-feed)

Each provider can be used separately by explicitly importing only those packages
you are going to use. This is useful when you don't need any extra stuff in your
project and you know which appcast provider you are dealing with.

In other cases importing a single `github.com/victorpopkov/go-appcast` is the
best option. This will give you all the necessary functionality to work with the
supported providers as it will automatically detect which is used and then call
the appropriate methods.

### GitHub Atom Feed

Each project that uses [GitHub][] releases to distribute applications has its
own Atom Feed available that can be considered as an appcast.

Example URL: <https://github.com/atom/atom/releases.atom>

### SourceForge RSS Feed

Each project hosted on [SourceForge][] has its own releases RSS feed available
that can be considered to be an appcast.

Example URL: <https://sourceforge.net/projects/filezilla/rss>

### Sparkle RSS Feed

Appcasts, created by the [Sparkle Framework][]. Originally, Sparkle was created
to distribute software updates for macOS applications. However, for Windows,
there is a [WinSparkle](https://winsparkle.org/) framework which uses the same
[RSS enclosure][] technique to distribute updates and release notes.

A good example of a Sparkle appcast is how [Adium](https://adium.im/)
distributes software updates: <https://www.adium.im/sparkle/appcast-release.xml>.
You can find a [GoDoc][] examples here:

- [`import "github.com/victorpopkov/go-appcast"`](https://godoc.org/github.com/victorpopkov/go-appcast)
- [`import "github.com/victorpopkov/go-appcast/provider/sparkle"`](https://godoc.org/github.com/victorpopkov/go-appcast/provider/sparkle)

## Sources

Out of the box, 2 sources are supported:

- [`source.Local`](#sourcelocal) (load from the local file)
- [`source.Remote`](#sourceremote) (load from the remote location)

This means that you have by default 2 options from where an appcast can be
loaded. You can just choose the appropriate one from the `source` package or
create your own.

You can find a [GoDoc][] standalone package example here:
[`import "github.com/victorpopkov/go-appcast/source"`](https://godoc.org/github.com/victorpopkov/go-appcast/source).
In addition, by digging into the ["Providers"](#providers) you can see how to
use them alongside with an appcast in their examples:

- [`import "github.com/victorpopkov/go-appcast/provider/github"`](https://godoc.org/github.com/victorpopkov/go-appcast/provider/github)
- [`import "github.com/victorpopkov/go-appcast/provider/sourceforge"`](https://godoc.org/github.com/victorpopkov/go-appcast/provider/sourceforge)
- [`import "github.com/victorpopkov/go-appcast/provider/sparkle"`](https://godoc.org/github.com/victorpopkov/go-appcast/provider/sparkle)

### `source.Remote`

This was designed to retrieve an appcast data from the remote location by URL.
It should cover most use cases when the appcast is available remotely.

For convenience purposes an `Appcast.LoadFromRemoteSource` can be used when
using the default `appcast` package. It sets the `Appcast` to use the
`source.Remote`, loads the source and unmarshals it.

### `source.Local`

This was designed to retrieve an appcast data from the local file by path.

For convenience purposes an `Appcast.LoadFromLocalSource` can be used when using
the default `appcast` package. It sets the `Appcast` to use the `source.Local`,
loads the source and unmarshals it.

## Outputs

Out of the box, only a single `output.Local` is available to save an appcast to
the local file.

Just like the ["Sources"](#sources), outputs are designed in the same way
meaning that you can extend the list of the supported outputs by creating your
own.

You can find a [GoDoc][] standalone package example here:
[`import "github.com/victorpopkov/go-appcast/output"`](https://godoc.org/github.com/victorpopkov/go-appcast/output).
In addition, by digging into the ["Providers"](#providers) you can see how to
use them alongside with an appcast in their "Marshal" examples:

- [`import "github.com/victorpopkov/go-appcast/provider/github"`](https://godoc.org/github.com/victorpopkov/go-appcast/provider/github)
- [`import "github.com/victorpopkov/go-appcast/provider/sourceforge"`](https://godoc.org/github.com/victorpopkov/go-appcast/provider/sourceforge)
- [`import "github.com/victorpopkov/go-appcast/provider/sparkle"`](https://godoc.org/github.com/victorpopkov/go-appcast/provider/sparkle)

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).

[github]: https://github.com/
[godoc]: https://godoc.org/
[rss enclosure]: https://en.wikipedia.org/wiki/RSS_enclosure
[sourceforge]: https://sourceforge.net/
[sparkle framework]: https://sparkle-project.org/
[sparkle]: https://sparkle-project.org/
