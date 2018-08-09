# go-appcast

[![Build Status](https://travis-ci.org/victorpopkov/go-appcast.svg?branch=master)](https://travis-ci.org/victorpopkov/go-appcast)
[![Coverage Status](https://coveralls.io/repos/github/victorpopkov/go-appcast/badge.svg?branch=master)](https://coveralls.io/github/victorpopkov/go-appcast?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/victorpopkov/go-appcast)](https://goreportcard.com/report/github.com/victorpopkov/go-appcast)
[![GoDoc](https://godoc.org/github.com/victorpopkov/go-appcast?status.svg)](https://godoc.org/github.com/victorpopkov/go-appcast)

This library provides functionality for working with appcasts. It retrieves
versions alongside with the download URLs if available and other useful
information.

- [What the heck is appcast?](#what-the-heck-is-appcast)
- [Why is this library needed and what it does?](#why-is-this-library-needed-and-what-it-does)
  - [Features](#features)
- [Supported providers](#supported-providers)
  - [Sparkle RSS Feed](#sparkle-rss-feed)
  - [SourceForge RSS Feed](#sourceforge-rss-feed)
  - [GitHub Atom Feed](#github-atom-feed)
- [Sources](#sources)
  - [RemoteSource](#remotesource)
  - [LocalSource](#localsource)

## What the heck is appcast?

The word "appcast" is usually referred to a remote web page providing
information about software updates. This kind of pages is usually created for
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

### Features

- [x] Load an appcast from different sources:
  - [x] remote source (from the remote location)
  - [x] local source (from the local file)
- [x] Detect which one of the supported providers is used
- [x] Unmarshal (extract) releases
- [x] Sort releases by version
- [x] Try to guess the release stability
- [x] Filter releases:
  - [x] by title, media type or download URL (using RegExp)
  - [x] by stability

## Supported providers

At the moment, only 3 providers are supported:

- [Sparkle RSS Feed](#sparkle-rss-feed)
- [SourceForge RSS Feed](#sourceforge-rss-feed)
- [GitHub Atom Feed](#github-atom-feed)

### Sparkle RSS Feed

Appcasts, created by the [Sparkle Framework](https://sparkle-project.org/).
Originally, Sparkle was created to distribute software updates for macOS
applications. However, for Windows, there is a [WinSparkle](https://winsparkle.org/)
framework which uses the same [RSS enclosure](https://en.wikipedia.org/wiki/RSS_enclosure)
technique to distribute updates and release notes.

Example URL: [https://www.adium.im/sparkle/appcast-release.xml](https://www.adium.im/sparkle/appcast-release.xml)

#### "Sparkle RSS Feed" example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	a := appcast.New()
	a.LoadFromRemoteSource("https://www.adium.im/sparkle/appcast-release.xml")
	a.SortReleasesByVersions(appcast.DESC)

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())
	fmt.Printf("Releases: %d total\n\n", len(a.Releases()))

	release := a.Releases()[0]
	fmt.Printf("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", release.Version())
	fmt.Printf("%12s %s\n", "Build:", release.Build())
	fmt.Printf("%12s %s\n", "Title:", release.Title())
	fmt.Printf("%12s %v\n", "Downloads:", release.Downloads())
	fmt.Printf("%12s %v\n", "Published:", release.PublishedDateTime())
	fmt.Printf("%12s %v\n", "Pre-release:", release.IsPreRelease())

	// Output:
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
	//
	// First release details:
	//
	//     Version: 1.5.10.4
	//       Build: 1.5.10.4
	//       Title: Adium 1.5.10.4
	//   Downloads: [{https://adiumx.cachefly.net/Adium_1.5.10.4.dmg application/octet-stream 21140435}]
	//   Published: 2017-05-14 12:04:01 +0000 UTC
	// Pre-release: false
}
```

### SourceForge RSS Feed

Each project hosted on [SourceForge](https://sourceforge.net/) has its own
releases RSS feed available that can be considered as an appcast.

Example URL: [https://sourceforge.net/projects/filezilla/rss](https://sourceforge.net/projects/filezilla/rss)

#### "SourceForge RSS Feed" example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	a := appcast.New()
	a.LoadFromRemoteSource("https://sourceforge.net/projects/filezilla/rss")

	// apply some filters
	a.FilterReleasesByMediaType("application/x-bzip2")
	a.FilterReleasesByTitle("FileZilla_Client_Unstable", true)
	a.FilterReleasesByURL("macosx")
	defer a.ResetFilters() // reset

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())
	fmt.Printf("Releases: %d total\n\n", len(a.Releases()))

	release := a.Releases()[0]
	fmt.Printf("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", release.Version())
	fmt.Printf("%12s %s\n", "Title:", release.Title())
	fmt.Printf("%12s %v\n", "Downloads:", release.Downloads())
	fmt.Printf("%12s %v\n", "Published:", release.PublishedDateTime())
	fmt.Printf("%12s %v\n", "Pre-release:", release.IsPreRelease())

	// Output:
	// Checksum: 69886b91a041ce9d742218a77317cd99f87a14199c3f8ba094042dd9d430f7fd
	// Provider: SourceForge RSS Feed
	// Releases: 5 total
	//
	// First release details:
	//
	//     Version: 3.25.2
	//       Title: /FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2
	//   Downloads: [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8453714}]
	//   Published: 2017-04-30 12:07:25 +0000 UTC
	// Pre-release: false
}
```

### GitHub Atom Feed

Each project that uses [GitHub](https://github.com/) releases to distribute
applications has its own Atom Feed available that can be considered as an
appcast.

Example URL: [https://github.com/atom/atom/releases.atom](https://github.com/atom/atom/releases.atom)

#### "GitHub Atom Feed" example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	a := appcast.New()
	a.LoadFromRemoteSource("https://github.com/atom/atom/releases.atom")

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())
	fmt.Printf("Releases: %d total\n\n", len(a.Releases()))

	release := a.Releases()[0]
	fmt.Printf("First release details:\n\n")
	fmt.Printf("%12s %s\n", "Version:", release.Version())
	fmt.Printf("%12s %s\n", "Title:", release.Title())
	fmt.Printf("%12s %v\n", "Downloads:", release.Downloads())
	fmt.Printf("%12s %v\n", "Published:", release.PublishedDateTime())
	fmt.Printf("%12s %v\n", "Pre-release:", release.IsPreRelease())

	// Output:
	// Checksum: 03b6d9b8199ea377036caafa5358512295afa3c740edf9031dc6739b89e3ba05
	// Provider: GitHub Atom Feed
	// Releases: 10 total
	//
	// First release details:
	//
	//     Version: 1.28.0-beta3
	//       Title: 1.28.0-beta3
	//   Downloads: []
	//   Published: 2018-06-06 17:09:54 +0000 UTC
	// Pre-release: true
}
```

## Sources

Out of the box, two main sources where an appcast can be retrieved are
supported:

- [RemoteSource](#remotesource) (get from the remote location)
- [LocalSource](#localsource) (get from the local file)

> Manually setting the source gives you more control over when the load and/or
> unmarshaling should happen. This also allows modifying the source content if
> needed.

### `RemoteSource`

`RemoteSource` is designed to retrieve an appcast data from the remote location
by URL. It should cover most use cases when the appcast is available remotely.

> For convenience purposes an `Appcast.LoadFromRemoteSource` can be used. It
> sets the `Appcast` to use the `RemoteSource`, calls the appropriate source
> load method and unmarshals releases.

#### `RemoteSource` example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	src, _ := appcast.NewRemoteSource("https://www.adium.im/sparkle/appcast-release.xml")

	a := appcast.New(src)
	a.LoadSource()
	a.UnmarshalReleases()

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())
	fmt.Printf("Releases: %d total\n", len(a.Releases()))

	// Output:
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
}
```

### `LocalSource`

`LocalSource` is designed to retrieve an appcast data from the local file by
path.

> For convenience purposes an `Appcast.LoadFromLocalSource` can be used. It sets
> the `Appcast` to use the `LocalSource`, calls the appropriate source load
> method and unmarshals releases.

#### `LocalSource` example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	src := appcast.NewLocalSource("/path/to/file.xml")

	a := appcast.New(src)
	a.LoadSource()
	a.UnmarshalReleases()

	fmt.Println("Checksum:", a.Source().Checksum())
	fmt.Println("Provider:", a.Source().Provider())
	fmt.Printf("Releases: %d total\n", len(a.Releases()))

	// Output:
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Releases: 5 total
}
```

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).
