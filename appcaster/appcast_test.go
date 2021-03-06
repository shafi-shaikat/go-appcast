package appcaster

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/release"
)

type TestSource struct {
	Source
}

type TestSourceError struct {
	Source
}

func (s *TestSource) Load() error {
	return nil
}

func (s *TestSourceError) Load() error {
	return fmt.Errorf("error")
}

// newTestAppcast creates a new Appcast instance for testing purposes and
// returns its pointer. By default the content is []byte("content"). However,
// own content can be provided as an argument.
func newTestAppcast(content ...interface{}) *Appcast {
	var resultContent []byte

	d := new(release.Download)
	d.SetUrl("https://example.com/app_2.0.0-beta.dmg")
	d.SetFiletype("application/octet-stream")
	d.SetLength(100000)

	r := new(release.Release)
	r.SetBuild("200")
	r.SetTitle("Release 2.0.0-beta")
	r.SetDescription("Release 2.0.0-beta Description")
	r.SetReleaseNotesLink("https://example.com/changelogs/2.0.0-beta.html")
	r.SetMinimumSystemVersion("10.10")
	r.SetDownloads([]release.Download{*d})
	r.SetIsPreRelease(true)

	// r1
	d1 := d
	r1 := *r
	r1.SetDownloads([]release.Download{*d1})

	err := r1.SetVersionString("2.0.0-beta")
	if err != nil {
		panic(err)
	}

	t, _ := time.Parse(time.RFC1123Z, "Fri, 13 May 2016 12:00:00 +0200")
	r1.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	// r2
	d2 := d
	d2.SetUrl("https://example.com/app_1.1.0.dmg")

	r2 := *r
	r2.SetBuild("110")
	r2.SetTitle("Release 1.1.0")
	r2.SetDescription("Release 1.1.0 Description")
	r2.SetReleaseNotesLink("https://example.com/changelogs/1.1.0.html")
	r2.SetMinimumSystemVersion("10.9")
	r2.SetDownloads([]release.Download{*d2})
	r2.SetIsPreRelease(false)

	err = r2.SetVersionString("1.1.0")
	if err != nil {
		panic(err)
	}

	t, _ = time.Parse(time.RFC1123Z, "Thu, 12 May 2016 12:00:00 +0200")
	r2.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	// r3
	d3 := d
	d3.SetUrl("https://example.com/app_1.0.1.dmg")

	r3 := *r
	r3.SetBuild("101")
	r3.SetTitle("Release 1.0.1")
	r3.SetDescription("Release 1.0.1 Description")
	r3.SetReleaseNotesLink("https://example.com/changelogs/1.0.1.html")
	r3.SetMinimumSystemVersion("10.9")
	r3.SetDownloads([]release.Download{*d3})
	r3.SetIsPreRelease(false)

	err = r3.SetVersionString("1.0.1")
	if err != nil {
		panic(err)
	}

	t, _ = time.Parse(time.RFC1123Z, "Wed, 11 May 2016 12:00:00 +0200")
	r3.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	// r4
	d4 := d
	d4.SetUrl("https://example.com/app_1.0.0.dmg")

	r4 := *r
	r4.SetBuild("100")
	r4.SetTitle("Release 1.0.0")
	r4.SetDescription("Release 1.0.0 Description")
	r4.SetReleaseNotesLink("https://example.com/changelogs/1.0.0.html")
	r4.SetMinimumSystemVersion("10.9")
	r4.SetDownloads([]release.Download{*d3})
	r4.SetIsPreRelease(false)

	err = r4.SetVersionString("1.0.0")
	if err != nil {
		panic(err)
	}

	t, _ = time.Parse(time.RFC1123Z, "Tue, 10 May 2016 12:00:00 +0200")
	r4.SetPublishedDateTime(release.NewPublishedDateTime(&t))

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("content")
	}

	return &Appcast{
		source: &Source{
			content:  resultContent,
			provider: Provider(0),
		},
		output: &Output{
			content: resultContent,
			checksum: &Checksum{
				algorithm: SHA256,
				source:    resultContent,
				result:    []byte("test"),
			},
			provider: Provider(0),
		},
		releases: release.NewReleases([]release.Releaser{&r1, &r2, &r3, &r4}),
	}
}

func TestNew(t *testing.T) {
	// test (without source)
	a := New()
	assert.IsType(t, Appcast{}, *a)
	assert.Nil(t, a.source)

	// test (with source)
	a = New(
		&Source{
			content:  []byte("content"),
			provider: Provider(0),
		},
	)

	assert.IsType(t, Appcast{}, *a)
	assert.NotNil(t, a.source)
}

func TestExtractSemanticVersions(t *testing.T) {
	testCases := map[string][]string{
		// single
		"Version 1":           nil,
		"Version 1.0":         nil,
		"Version 1.0.2":       {"1.0.2"},
		"Version 1.0.2-alpha": {"1.0.2-alpha"},
		"Version 1.0.2-beta":  {"1.0.2-beta"},
		"Version 1.0.2-dev":   {"1.0.2-dev"},
		"Version 1.0.2-rc1":   {"1.0.2-rc1"},

		// multiples
		"First is v1.0.1, second is v1.0.2, third is v1.0.3": {"1.0.1", "1.0.2", "1.0.3"},
	}

	// test
	for data, versions := range testCases {
		actual, err := ExtractSemanticVersions(data)
		if versions == nil {
			assert.Error(t, err)
			assert.EqualError(t, err, "no semantic versions found")
		} else {
			assert.Nil(t, err)
			assert.Equal(t, versions, actual)
		}
	}
}

func TestAppcast_GenerateSourceChecksum(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.Nil(t, a.source.Checksum())

	// test
	result := a.GenerateSourceChecksum(MD5)
	assert.Equal(t, result.String(), a.source.Checksum().String())
	assert.Equal(t, "9a0364b9e99bb480dd25e1f0284c8555", result.String())
	assert.Equal(t, MD5, a.source.Checksum().Algorithm())
}

func TestAppcast_LoadSource(t *testing.T) {
	// preparations
	a := newTestAppcast()

	// test (successful)
	a.SetSource(&TestSource{})
	err := a.LoadSource()
	assert.Nil(t, err)

	// test (error)
	a.SetSource(&TestSourceError{})
	err = a.LoadSource()
	assert.Error(t, err)
}

func TestAppcast_GuessSourceProvider(t *testing.T) {
	a := newTestAppcast()
	assert.Panics(t, func() {
		a.GuessSourceProvider()
	})
}

func TestAppcast_Unmarshal(t *testing.T) {
	a := newTestAppcast()
	assert.Panics(t, func() {
		a.Unmarshal()
	})
}

func TestAppcast_Uncomment(t *testing.T) {
	a := newTestAppcast()
	assert.Panics(t, func() {
		_ = a.Uncomment()
	})
}

func TestAppcast_Source(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.source, a.Source())
}

func TestAppcast_SetSource(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.NotNil(t, a.source)

	// test
	a.SetSource(nil)
	assert.Nil(t, a.source)
}

func TestAppcast_Output(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.output, a.Output())
}

func TestAppcast_SetOutput(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.NotNil(t, a.output)

	// test
	a.SetOutput(nil)
	assert.Nil(t, a.output)
}

func TestAppcast_Releases(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.releases, a.Releases())
}

func TestAppcast_SetReleases(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.NotNil(t, a.releases)

	// test
	a.SetReleases(nil)
	assert.Nil(t, a.releases)
}

func TestAppcast_FirstRelease(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.releases.First(), a.FirstRelease())
}
