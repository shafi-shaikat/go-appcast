# go-appcast

[![Build Status](https://travis-ci.org/victorpopkov/go-appcast.svg?branch=master)](https://travis-ci.org/victorpopkov/go-appcast)
[![Coverage Status](https://coveralls.io/repos/github/victorpopkov/go-appcast/badge.svg?branch=master)](https://coveralls.io/github/victorpopkov/go-appcast?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/victorpopkov/go-appcast)](https://goreportcard.com/report/github.com/victorpopkov/go-appcast)
[![GoDoc](https://godoc.org/github.com/victorpopkov/go-appcast?status.svg)](https://godoc.org/github.com/victorpopkov/go-appcast)

**NOTICE:** Currently in development.

This library provides functionality for working with appcasts. It retrieves
versions alongside with the download URLs if available and other useful
information.

- [What the heck is appcast?](#what-the-heck-is-appcast)
- [Why is this library needed?](#why-is-this-library-needed)
- [Supported providers](#supported-providers)

## What the heck is appcast?

The word "appcast" is usually referred to a remote web page providing
information about software updates. This kind of pages are usually created for
different software update frameworks like ([Sparkle](https://sparkle-project.org/))
or generated by different services that distribute applications
([SourceForge](https://sourceforge.net/), [GitHub](https://github.com/)
and etc.). There are plenty of different methods available, but originally
"appcasting" was the practice of using an [RSS enclosure](https://en.wikipedia.org/wiki/RSS_enclosure)
to distribute updates and release notes.

## Why is this library needed and what it does?

Since today you can find plenty of different ways how vendors distribute their
software updates, this library attempts to provide a universal way of analyzing
and retrieving the useful information from appcasts of the supported providers.

Basically, at the moment it knows how to:

- [x] Load the appcast from remote URL
- [ ] Load the appcast from local file _(not yet implemented)_
- [x] Detect which one of the supported providers is used
- [x] Extract and sort releases
- [ ] Try to guess which release is stable and which is not _(not yet implemented)_

## Supported providers

At the moment, only 3 providers are supported:

- [Sparkle RSS Feed](#sparkle-rss-feed)
- [SourceForge RSS Feed](#sourceforge-rss-feed) _(not yet implemented)_
- [GitHub Atom Feed](#github-atom-feed) _(not yet implemented)_

### Sparkle RSS Feed

Appcasts, created for the [Sparkle Framework](https://sparkle-project.org/).
Originally, Sparkle was created to distribute software updates for macOS
applications. However, for Windows, there is a [WinSparkle](https://winsparkle.org/)
framework which uses the same [RSS enclosure](https://en.wikipedia.org/wiki/RSS_enclosure)
technique to distribute updates and release notes.

Example URL: [https://www.adium.im/sparkle/appcast-release.xml](https://www.adium.im/sparkle/appcast-release.xml)

#### Example

```go
a := New()
a.LoadFromURL("https://www.adium.im/sparkle/appcast-release.xml")
a.GenerateChecksum(Sha256)
a.ExtractReleases()
a.SortReleasesByVersions(DESC)

fmt.Println("Checksum:", a.GetChecksum())
fmt.Println("Provider:", a.Provider)

for i, release := range a.Releases {
  fmt.Println(fmt.Sprintf("Release #%d:", i+1), release)
}

// Output:
// Checksum: bfee8d59301ff64b44d72572e973a8354d0397657552c91db61bedac353b04b8
// Provider: Sparkle RSS Feed
// Release #1: {1.5.10.4 1.5.10.4 Adium 1.5.10.4  [{https://adiumx.cachefly.net/Adium_1.5.10.4.dmg application/octet-stream 21140435}] 2017-05-14 05:04:01 -0700 -0700 false}
// Release #2: {1.5.10 1.5.10 Adium 1.5.10  [{https://adiumx.cachefly.net/Adium_1.5.10.dmg application/octet-stream 24595712}] 0001-01-01 00:00:00 +0000 UTC false}
// Release #3: {1.4.5 1.4.5 Adium 1.4.5  [{https://adiumx.cachefly.net/Adium_1.4.5.dmg application/octet-stream 23065688}] 0001-01-01 00:00:00 +0000 UTC false}
// Release #4: {1.3.10 1.3.10 Adium 1.3.10  [{https://adiumx.cachefly.net/Adium_1.3.10.dmg application/octet-stream 22369877}] 0001-01-01 00:00:00 +0000 UTC false}
// Release #5: {1.0.6 1.0.6 Adium 1.0.6  [{https://adiumx.cachefly.net/Adium_1.0.6.dmg application/octet-stream 13795246}] 0001-01-01 00:00:00 +0000 UTC false}
```

### SourceForge RSS Feed

_(not yet implemented)_

Each project hosted on [SourceForge](https://sourceforge.net/) has its own
releases RSS feed available that can be considered as an appcast.

Example URL: [https://sourceforge.net/projects/filezilla/rss](https://sourceforge.net/projects/filezilla/rss)

### GitHub Atom Feed

_(not yet implemented)_

Each project that uses [GitHub](https://github.com/) releases to distribute
applications has its own Atom Feed available that can be considered as an
appcast.

Example URL: [https://github.com/atom/atom/releases.atom](https://github.com/atom/atom/releases.atom)

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).
