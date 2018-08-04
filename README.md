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

### Features

- [x] Load an appcast from both remote URL or local file
- [x] Detect which one of the supported providers is used
- [x] Extract releases
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

Appcasts, created for the [Sparkle Framework](https://sparkle-project.org/).
Originally, Sparkle was created to distribute software updates for macOS
applications. However, for Windows, there is a [WinSparkle](https://winsparkle.org/)
framework which uses the same [RSS enclosure](https://en.wikipedia.org/wiki/RSS_enclosure)
technique to distribute updates and release notes.

Example URL: [https://www.adium.im/sparkle/appcast-release.xml](https://www.adium.im/sparkle/appcast-release.xml)

#### Example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	a := appcast.New()
	a.LoadFromURL("https://www.adium.im/sparkle/appcast-release.xml")
	a.GenerateChecksum(appcast.SHA256)
	a.ExtractReleases()
	a.SortReleasesByVersions(appcast.DESC)

	fmt.Println("Checksum:", a.GetChecksum())
	fmt.Println("Provider:", a.GetProvider())

	for i, release := range a.Releases {
		fmt.Println(fmt.Sprintf("Release #%d:", i+1), release)
	}

	// Output:
	// Checksum: 6ec7c5abcaa78457cc4bf3c2196584446cca1461c65505cbaf0382a2f62128db
	// Provider: Sparkle RSS Feed
	// Release #1: {1.5.10.4 1.5.10.4 Adium 1.5.10.4  [{https://adiumx.cachefly.net/Adium_1.5.10.4.dmg application/octet-stream 21140435}] 2017-05-14 12:04:01 +0000 UTC false}
	// Release #2: {1.5.10 1.5.10 Adium 1.5.10  [{https://adiumx.cachefly.net/Adium_1.5.10.dmg application/octet-stream 24595712}] 2014-05-19 21:25:14 +0000 UTC false}
	// Release #3: {1.4.5 1.4.5 Adium 1.4.5  [{https://adiumx.cachefly.net/Adium_1.4.5.dmg application/octet-stream 23065688}] 2012-03-20 20:30:00 +0000 UTC false}
	// Release #4: {1.3.10 1.3.10 Adium 1.3.10  [{https://adiumx.cachefly.net/Adium_1.3.10.dmg application/octet-stream 22369877}] 2010-01-12 23:30:00 +0000 UTC false}
	// Release #5: {1.0.6 1.0.6 Adium 1.0.6  [{https://adiumx.cachefly.net/Adium_1.0.6.dmg application/octet-stream 13795246}] 2007-08-13 22:12:45 +0000 UTC false}
}
```

### SourceForge RSS Feed

Each project hosted on [SourceForge](https://sourceforge.net/) has its own
releases RSS feed available that can be considered as an appcast.

Example URL: [https://sourceforge.net/projects/filezilla/rss](https://sourceforge.net/projects/filezilla/rss)

#### Example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	a := appcast.New()
	a.LoadFromURL("https://sourceforge.net/projects/filezilla/rss")
	a.GenerateChecksum(appcast.SHA256)
	a.ExtractReleases()

	// apply some filters
	a.FilterReleasesByMediaType("application/x-bzip2")
	a.FilterReleasesByTitle("FileZilla_Client_Unstable", true)
	a.FilterReleasesByURL("macosx")
	defer a.ResetFilters() // reset

	fmt.Println("Checksum:", a.GetChecksum())
	fmt.Println("Provider:", a.GetProvider())

	for i, release := range a.Releases {
		fmt.Println(fmt.Sprintf("Release #%d:", i+1), release)
	}

	// Output:
	// Checksum: 69886b91a041ce9d742218a77317cd99f87a14199c3f8ba094042dd9d430f7fd
	// Provider: SourceForge RSS Feed
	// Release #1: {3.25.2  /FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2 /FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.2/FileZilla_3.25.2_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8453714}] 2017-04-30 12:07:25 +0000 UTC false}
	// Release #2: {3.25.1  /FileZilla_Client/3.25.1/FileZilla_3.25.1_macosx-x86.app.tar.bz2 /FileZilla_Client/3.25.1/FileZilla_3.25.1_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.1/FileZilla_3.25.1_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8460741}] 2017-03-20 17:11:09 +0000 UTC false}
	// Release #3: {3.25.0  /FileZilla_Client/3.25.0/FileZilla_3.25.0_macosx-x86.app.tar.bz2 /FileZilla_Client/3.25.0/FileZilla_3.25.0_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.25.0/FileZilla_3.25.0_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8461936}] 2017-03-13 14:36:41 +0000 UTC false}
	// Release #4: {3.24.1  /FileZilla_Client/3.24.1/FileZilla_3.24.1_macosx-x86.app.tar.bz2 /FileZilla_Client/3.24.1/FileZilla_3.24.1_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.24.1/FileZilla_3.24.1_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8764178}] 2017-02-21 22:00:38 +0000 UTC false}
	// Release #5: {3.24.0  /FileZilla_Client/3.24.0/FileZilla_3.24.0_macosx-x86.app.tar.bz2 /FileZilla_Client/3.24.0/FileZilla_3.24.0_macosx-x86.app.tar.bz2 [{https://sourceforge.net/projects/filezilla/files/FileZilla_Client/3.24.0/FileZilla_3.24.0_macosx-x86.app.tar.bz2/download application/x-bzip2; charset=binary 8765941}] 2017-01-13 20:20:31 +0000 UTC false}
}
```

### GitHub Atom Feed

Each project that uses [GitHub](https://github.com/) releases to distribute
applications has its own Atom Feed available that can be considered as an
appcast.

Example URL: [https://github.com/atom/atom/releases.atom](https://github.com/atom/atom/releases.atom)

#### Example

```go
package main

import (
	"fmt"

	"github.com/victorpopkov/go-appcast"
)

func main() {
	a := appcast.New()
	a.LoadFromURL("https://github.com/atom/atom/releases.atom")
	a.GenerateChecksum(appcast.SHA256)
	a.ExtractReleases()

	fmt.Println("Checksum:", a.GetChecksum())
	fmt.Println("Provider:", a.GetProvider())

	for i, release := range a.Releases {
		fmt.Println(fmt.Sprintf("Release #%d:", i+1), release.Version, release.Title, release.PublishedDateTime, release.IsPrerelease)
	}

	fmt.Println("Release #1 description:", a.Releases[0].Description)

  // Output:
	// Checksum: 03b6d9b8199ea377036caafa5358512295afa3c740edf9031dc6739b89e3ba05
	// Provider: GitHub Atom Feed
	// Release #1: 1.28.0-beta3 1.28.0-beta3 2018-06-06 17:09:54 +0000 UTC true
	// Release #2: 1.28.0-beta2 1.28.0-beta2 2018-05-31 13:55:54 +0000 UTC true
	// Release #3: 1.27.2 1.27.2 2018-05-31 13:55:49 +0000 UTC false
	// Release #4: 1.28.0-beta1 1.28.0-beta1 2018-05-21 14:46:23 +0000 UTC true
	// Release #5: 1.27.1 1.27.1 2018-05-21 14:46:10 +0000 UTC false
	// Release #6: 1.28.0-beta0 1.28.0-beta0 2018-05-15 17:46:08 +0000 UTC true
	// Release #7: 1.27.0 1.27.0 2018-05-15 17:46:03 +0000 UTC false
	// Release #8: 1.27.0-beta1 1.27.0-beta1 2018-04-26 19:40:51 +0000 UTC true
	// Release #9: 1.26.1 1.26.1 2018-04-26 19:40:40 +0000 UTC false
	// Release #10: 1.26.0 1.26.0 2018-04-18 23:00:10 +0000 UTC false
}
```

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).
