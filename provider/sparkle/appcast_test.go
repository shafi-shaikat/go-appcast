package sparkle

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/victorpopkov/go-appcast/appcaster"
)

// workingDir returns a current working directory path. If it's not available
// prints an error to os.Stdout and exits with error status 1.
func workingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return pwd
}

// testdata returns a file content as a byte slice from the provided testdata
// paths. If the file is not found, prints an error to os.Stdout and exits with
// exit status 1.
func testdata(paths ...string) []byte {
	path := testdataPath(paths...)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		os.Exit(1)
	}

	return content
}

// testdataPath returns a full path for the provided testdata paths.
func testdataPath(paths ...string) string {
	return filepath.Join(workingDir(), "./testdata/", filepath.Join(paths...))
}

// ReadLine reads a provided line number from io.Reader and returns it alongside
// with an error.
func readLine(r io.Reader, lineNum int) (line string, err error) {
	var lastLine int

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), nil
		}
	}

	return "", fmt.Errorf("there is no line \"%d\" in specified io.Reader", lineNum)
}

// getLine returns a specified line from the passed content.
func getLine(lineNum int, content []byte) (line string, err error) {
	return readLine(bytes.NewReader(content), lineNum)
}

// newTestAppcast creates a new Appcast instance for testing purposes and
// returns its pointer. By default the source is LocalSource and points to the
// "Sparkle RSS Feed" default.xml testdata.
func newTestAppcast(paths ...string) *Appcast {
	var content []byte

	if len(paths) > 0 {
		content = testdata(paths...)
	} else {
		content = testdata("unmarshal", "default.xml")
	}

	s := new(appcaster.Source)
	s.SetContent(content)
	s.GenerateChecksum(appcaster.SHA256)
	s.SetProvider(appcaster.Provider(0))

	a := new(Appcast)
	a.SetSource(s)

	return a
}

func TestNew(t *testing.T) {
	// test (without source)
	a := New()
	assert.IsType(t, Appcast{}, *a)
	assert.Nil(t, a.Source())

	// test (with source)
	src := new(appcaster.Source)
	src.SetContent([]byte("content"))
	src.SetProvider(appcaster.Provider(0))

	a = New(src)
	assert.IsType(t, Appcast{}, *a)
	assert.NotNil(t, a.Source())
}

func TestAppcast_Unmarshal(t *testing.T) {
	type testCase struct {
		path     string
		appcast  appcaster.Appcaster
		releases map[string][]string
		errors   []string
	}

	testCases := []testCase{
		{
			path:    "attributes_as_elements.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
				"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			},
		},
		{
			path:    "default.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
				"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			},
		},
		{
			path:    "default_asc.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			},
		},
		{
			path:    "incorrect_namespace.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
				"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			},
		},
		{
			path:    "invalid_pubdate.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
				"1.1.0": {"", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			},
			errors: []string{
				"release #2 (parsing of the published datetime failed)",
			},
		},
		{
			path:    "no_releases.xml",
			appcast: &Appcast{},
		},
		{
			path: "invalid_tag.xml",
			errors: []string{
				"XML syntax error on line 14: element <enclosure> closed by </item>",
			},
		},
		{
			path:    "invalid_version.xml",
			appcast: &Appcast{},
			errors: []string{
				"release #2 (malformed version: invalid)",
			},
		},
		//{
		//	path:    "multiple_enclosure.xml",
		//	appcast: &Appcast{},
		//},
		{
			path:    "only_version.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "2.0.0", "https://example.com/app_2.0.0.dmg", "10.10"},
				"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "1.1.0", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "1.0.1", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "1.0.0", "https://example.com/app_1.0.0.dmg", "10.9"},
			},
		},
		{
			path:    "prerelease.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0-beta": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0_beta.dmg", "10.10"},
				"1.1.0":      {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1":      {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0":      {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			},
		},
		{
			path:    "single.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
			},
		},
		{
			path:    "without_namespaces.xml",
			appcast: &Appcast{},
			releases: map[string][]string{
				"2.0.0": {"Fri, 13 May 2016 12:00:00 +0200", "200", "https://example.com/app_2.0.0.dmg", "10.10"},
				"1.1.0": {"Thu, 12 May 2016 12:00:00 +0200", "110", "https://example.com/app_1.1.0.dmg", "10.9"},
				"1.0.1": {"Wed, 11 May 2016 12:00:00 +0200", "101", "https://example.com/app_1.0.1.dmg", "10.9"},
				"1.0.0": {"Tue, 10 May 2016 12:00:00 +0200", "100", "https://example.com/app_1.0.0.dmg", "10.9"},
			},
		},
		{
			path:    "with_comments.xml",
			appcast: &Appcast{},
			errors: []string{
				"release #1 (no version)",
				"release #2 (no version)",
			},
		},
	}

	// test
	for _, testCase := range testCases {
		// preparations
		a := newTestAppcast("unmarshal", testCase.path)

		// test
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, a.Source().Appcast())
		assert.Nil(t, a.channel)
		assert.Empty(t, a.Releases())

		appcast, errors := a.Unmarshal()

		if testCase.appcast != nil {
			assert.IsType(t, testCase.appcast, appcast, fmt.Sprintf("%s: appcast type mismatch", testCase.path))
			assert.IsType(t, testCase.appcast, a.Source().Appcast())
		} else {
			assert.Equal(t, testCase.appcast, appcast, fmt.Sprintf("%s: appcast type mismatch", testCase.path))
		}

		if len(testCase.errors) == 0 {
			// successful
			assert.Nil(t, errors, fmt.Sprintf("%s: errors not nil", testCase.path))

			assert.IsType(t, &Channel{}, a.channel)
			assert.Equal(t, "App", a.channel.Title)
			assert.Equal(t, "https://example.com/app/", a.channel.Link)
			assert.Equal(t, "App Description", a.channel.Description)
			assert.Equal(t, "en", a.channel.Language)

			releases := testCase.releases
			assert.Len(t, releases, a.Releases().Len())

			for _, r := range a.Releases().Filtered() {
				v := r.Version().String()
				assert.Equal(t, fmt.Sprintf("Release %s", v), r.Title())
				assert.Equal(t, fmt.Sprintf("Release %s Description", v), r.Description())
				assert.Equal(t, releases[v][0], r.PublishedDateTime().String())
				assert.Equal(t, releases[v][1], r.Build())
				assert.Equal(t, releases[v][3], r.MinimumSystemVersion())

				// downloads
				assert.Equal(t, releases[v][2], r.Downloads()[0].Url())
				assert.Equal(t, "application/octet-stream", r.Downloads()[0].Filetype())
				assert.Equal(t, 100000, r.Downloads()[0].Length())
			}
		} else {
			// error (unmarshalling failure)
			assert.Len(t, errors, len(testCase.errors), fmt.Sprintf("%s: errors length mismatch", testCase.path))

			for i, errorMsg := range testCase.errors {
				err := errors[i]
				assert.EqualError(t, err, errorMsg)
			}
		}
	}

	// test (error) [no source]
	a := new(Appcast)

	p, errors := a.Unmarshal()

	assert.Len(t, errors, 1)
	err := errors[0]

	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, p)
	assert.Nil(t, a.Source())
	assert.Nil(t, a.channel)
}

func TestAppcast_Uncomment(t *testing.T) {
	testCases := map[string][]int{
		"attributes_as_elements.xml": nil,
		"default_asc.xml":            nil,
		"default.xml":                nil,
		"incorrect_namespace.xml":    nil,
		"multiple_enclosure.xml":     nil,
		"single.xml":                 nil,
		"with_comments.xml":          {13, 20},
		"without_namespaces.xml":     nil,
	}

	regexCommentStart := regexp.MustCompile(`<!--([[:space:]]*)?<`)
	regexCommentEnd := regexp.MustCompile(`>([[:space:]]*)?-->`)

	// test (successful)
	for filename, commentLines := range testCases {
		// preparations
		a := newTestAppcast("unmarshal", filename)

		// before
		for _, commentLine := range commentLines {
			line, _ := getLine(commentLine, a.Source().Content())
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.True(t, check, fmt.Sprintf("\"%s\" doesn't have a commented out line", filename))
		}

		err := a.Uncomment()

		// after
		assert.Nil(t, err)

		for _, commentLine := range commentLines {
			line, _ := getLine(commentLine, a.Source().Content())
			check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
			assert.False(t, check, fmt.Sprintf("\"%s\" didn't uncomment a \"%d\" line", filename, commentLine))
		}
	}

	// test (error) [no source]
	a := new(Appcast)

	err := a.Uncomment()
	assert.Error(t, err)
	assert.EqualError(t, err, "no source")
	assert.Nil(t, a.Source())
	assert.Nil(t, a.channel)
}

func TestAppcast_Channel(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.channel, a.Channel())
}

func TestAppcast_SetChannel(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.Nil(t, a.channel)

	// test
	a.SetChannel(&Channel{})
	assert.NotNil(t, a.channel)
}
