package appcast

import "regexp"

// Provider holds different supported providers.
type Provider int

const (
	// Unknown represents an unknown appcast provider.
	Unknown Provider = iota

	// Sparkle represents an RSS feed that is generated by the Sparkle
	// framework.
	Sparkle

	// SourceForgeRSSFeed represents an RSS feed of the releases generated by
	// SourceForge.
	SourceForgeRSSFeed

	// GitHubAtomFeed represents an Atom feed of the releases generated by
	// GitHub.
	GitHubAtomFeed
)

var providerNames = [...]string{
	"-",
	"Sparkle RSS Feed",
	"SourceForge RSS Feed",
	"GitHub Atom Feed",
}

// GuessProviderByContent attempts to guess the supported provider from the
// passed content. By default returns Provider.Unknown.
func GuessProviderByContent(content []byte) Provider {
	regexSparkleRSSFeed := regexp.MustCompile(`(?s)(<rss.*xmlns:sparkle)|(?s)(<rss.*<enclosure)`)
	regexSourceForgeRSSFeed := regexp.MustCompile(`(?s)(<rss.*xmlns:sf)|(?s)(<channel.*xmlns:sf)`)
	regexGitHubAtomFeed := regexp.MustCompile(`(?s)<feed.*<id>tag:github.com`)

	if regexSparkleRSSFeed.Match(content) {
		return Sparkle
	}

	if regexSourceForgeRSSFeed.Match(content) {
		return SourceForgeRSSFeed
	}

	if regexGitHubAtomFeed.Match(content) {
		return GitHubAtomFeed
	}

	return Unknown
}

// GuessProviderByContentString attempts to guess the supported provider from
// the passed content string. By default returns Provider.Unknown.
func GuessProviderByContentString(content string) Provider {
	return GuessProviderByContent([]byte(content))
}

// GuessProviderByUrl attempts to guess the supported provider from the passed
// URL. Only appcasts that are web-service specific can be guessed. By default
// returns Provider.Unknown.
func GuessProviderByUrl(url string) Provider {
	regexSourceForgeRSSFeed := regexp.MustCompile(`.*sourceforge.net/projects/.*/rss`)
	regexGitHubAtomFeed := regexp.MustCompile(`.*github\.com/(?P<user>.*?)/(?P<repo>.*?)/releases\.atom`)

	if regexSourceForgeRSSFeed.MatchString(url) {
		return SourceForgeRSSFeed
	}

	if regexGitHubAtomFeed.MatchString(url) {
		return GitHubAtomFeed
	}

	return Unknown
}

// String returns the string representation of the Provider.
func (p Provider) String() string {
	return providerNames[p]
}
