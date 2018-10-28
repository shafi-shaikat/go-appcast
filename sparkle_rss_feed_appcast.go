package appcast

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/victorpopkov/go-appcast/release"
)

// SparkleRSSFeedAppcaster is the interface that wraps the SparkleRSSFeedAppcast
// methods.
type SparkleRSSFeedAppcaster interface {
	Appcaster
	Channel() *SparkleRSSFeedAppcastChannel
	SetChannel(channel *SparkleRSSFeedAppcastChannel)
}

// SparkleRSSFeedAppcast represents the "Sparkle RSS Feed" appcast which is
// generated by Sparkle Framework.
type SparkleRSSFeedAppcast struct {
	Appcast
	channel *SparkleRSSFeedAppcastChannel
}

// SparkleRSSFeedAppcastChannel represents the "Sparkle RSS Feed" appcast
// channel data.
type SparkleRSSFeedAppcastChannel struct {
	Title       string
	Link        string
	Description string
	Language    string
}

// unmarshalSparkleRSSFeed represents an RSS itself for unmarshal purposes.
type unmarshalSparkleRSSFeed struct {
	Channel unmarshalSparkleRSSFeedChannel `xml:"channel"`
}

// unmarshalSparkleRSSFeedChannel represents the "Sparkle RSS Feed" channel for
// unmarshal purposes.
type unmarshalSparkleRSSFeedChannel struct {
	Title       string                        `xml:"title"`
	Link        string                        `xml:"link"`
	Description string                        `xml:"description"`
	Language    string                        `xml:"language"`
	Items       []unmarshalSparkleRSSFeedItem `xml:"item"`
}

// unmarshalSparkleRSSFeedItem represents an RSS item.
type unmarshalSparkleRSSFeedItem struct {
	Title                string                           `xml:"title"`
	Description          string                           `xml:"description"`
	PubDate              string                           `xml:"pubDate"`
	ReleaseNotesLink     string                           `xml:"releaseNotesLink"`
	MinimumSystemVersion string                           `xml:"minimumSystemVersion"`
	Enclosure            unmarshalSparkleRSSFeedEnclosure `xml:"enclosure"`
	Version              string                           `xml:"version"`
	ShortVersionString   string                           `xml:"shortVersionString"`
}

// unmarshalSparkleRSSFeedEnclosure represents the "Sparkle RSS Feed" item
// enclosure for unmarshal purposes.
type unmarshalSparkleRSSFeedEnclosure struct {
	DsaSignature       string `xml:"dsaSignature,attr"`
	Version            string `xml:"version,attr"`
	ShortVersionString string `xml:"shortVersionString,attr"`
	URL                string `xml:"url,attr"`
	Length             int    `xml:"length,attr"`
	Type               string `xml:"type,attr"`
}

// UnmarshalReleases unmarshals the Appcast.source.content into the
// Appcast.releases for the "Sparkle RSS Feed" provider.
//
// It returns both: the supported provider-specific appcast implementing the
// Appcaster interface and an error.
func (a *SparkleRSSFeedAppcast) UnmarshalReleases() (Appcaster, error) {
	var x unmarshalSparkleRSSFeed
	var version, build string

	xml.Unmarshal(a.source.Content(), &x)

	a.channel = &SparkleRSSFeedAppcastChannel{
		Title:       x.Channel.Title,
		Link:        x.Channel.Link,
		Description: x.Channel.Description,
		Language:    x.Channel.Language,
	}

	items := make([]release.Releaser, len(x.Channel.Items))
	for i, item := range x.Channel.Items {
		if item.Enclosure.ShortVersionString == "" && item.ShortVersionString != "" {
			version = item.ShortVersionString
		} else {
			version = item.Enclosure.ShortVersionString
		}

		if item.Enclosure.Version == "" && item.Version != "" {
			build = item.Version
		} else {
			build = item.Enclosure.Version
		}

		if version == "" && build == "" {
			return nil, fmt.Errorf("version is required, but it's not specified in release #%d", i+1)
		} else if version == "" && build != "" {
			version = build
		}

		// new release
		r, err := release.New(version, build)
		if err != nil {
			return nil, err
		}

		r.SetTitle(item.Title)
		r.SetDescription(item.Description)
		r.SetReleaseNotesLink(item.ReleaseNotesLink)
		r.SetMinimumSystemVersion(item.MinimumSystemVersion)

		// publishedDateTime
		p := release.NewPublishedDateTime()
		p.Parse(item.PubDate)
		r.SetPublishedDateTime(p)

		// prerelease
		if r.Version().Prerelease() != "" {
			r.SetIsPreRelease(true)
		}

		// downloads
		d := release.NewDownload(
			item.Enclosure.URL,
			item.Enclosure.Type,
			item.Enclosure.Length,
			item.Enclosure.DsaSignature,
		)

		r.AddDownload(*d)

		items[i] = r
	}

	a.releases = items

	return a, nil
}

// Uncomment uncomments XML tags in SparkleRSSFeedAppcast.source.content.
func (a *SparkleRSSFeedAppcast) Uncomment() error {
	if a.source == nil || len(a.source.Content()) == 0 {
		return fmt.Errorf("no source")
	}

	regex := regexp.MustCompile(`(<!--([[:space:]]*)?)|(([[:space:]]*)?-->)`)
	if regex.Match(a.source.Content()) {
		a.source.SetContent(regex.ReplaceAll(a.source.Content(), []byte("")))
	}

	return nil
}

// Channel is a SparkleRSSFeedAppcast.channel getter.
func (a *SparkleRSSFeedAppcast) Channel() *SparkleRSSFeedAppcastChannel {
	return a.channel
}

// SetChannel is a SparkleRSSFeedAppcast.channel setter.
func (a *SparkleRSSFeedAppcast) SetChannel(channel *SparkleRSSFeedAppcastChannel) {
	a.channel = channel
}
