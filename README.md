# go-appcast

[![Report Card](https://goreportcard.com/badge/github.com/victorpopkov/go-appcast)](https://goreportcard.com/badge/github.com/victorpopkov/go-appcast)
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
information about software updates. This kind of pages are usually generated by
different frameworks ([Sparkle](https://sparkle-project.org/)) or services that
distribute applications ([SourceForge](https://sourceforge.net/), [GitHub](https://github.com/)
and etc.). There are plenty of different methods available, but originally
"appcasting" was the practice of using an [RSS enclosure](https://en.wikipedia.org/wiki/RSS_enclosure)
to distribute updates and release notes.

## Why is this library needed?

Since today you can find plenty of different ways how vendors distribute their
software updates, this library attempts to provide a universal way for supported
providers of analyzing and retrieving useful information from appcasts.

## Supported providers

At the moment, only 3 providers are supported:

- [Sparkle RSS Feed](https://sparkle-project.org/)
- [SourceForge RSS Feed](https://sourceforge.net/)
- [GitHub Atom Feed](https://github.com/)

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).
