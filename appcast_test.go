package appcast

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
	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/victorpopkov/go-appcast/release"
)

var testdataPath = "./testdata/"

// getWorkingDir returns a current working directory path. If it's not available
// prints an error to os.Stdout and exits with error status 1.
func getWorkingDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return pwd
}

// getTestdata returns a file content as a byte slice from the provided testdata
// paths. If the file is not found, prints an error to os.Stdout and exits with
// exit status 1.
func getTestdata(paths ...string) []byte {
	path := getTestdataPath(paths...)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
		os.Exit(1)
	}

	return content
}

// getTestdataPath returns a full path for the provided testdata paths.
func getTestdataPath(paths ...string) string {
	return filepath.Join(getWorkingDir(), testdataPath, filepath.Join(paths...))
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

// getLineFromString returns a specified line from the passed string content.
func getLineFromString(lineNum int, content string) (line string, err error) {
	return getLine(lineNum, []byte(content))
}

// newTestAppcast creates a new Appcast instance for testing purposes and
// returns its pointer. By default the content is []byte("test"). However, own
// content can be provided as an argument.
func newTestAppcast(content ...interface{}) *Appcast {
	var resultContent []byte

	if len(content) > 0 {
		resultContent = content[0].([]byte)
	} else {
		resultContent = []byte("test")
	}

	url := "https://example.com/appcast.xml"
	r, _ := NewRequest(url)

	s := &Appcast{
		source: &RemoteSource{
			Source: &Source{
				content:  resultContent,
				provider: Unknown,
			},
			request: r,
			url:     url,
		},
		output: &LocalOutput{
			Output: &Output{
				content: resultContent,
				checksum: &Checksum{
					algorithm: SHA256,
					source:    resultContent,
					result:    []byte("test"),
				},
				provider: Unknown,
			},
			filepath:    "/tmp/test.txt",
			permissions: 0777,
		},
	}

	return s
}

func TestNew(t *testing.T) {
	// test (without source)
	a := New()
	assert.IsType(t, Appcast{}, *a)
	assert.Nil(t, a.source)

	// test (with source)
	a = New(NewLocalSource(getTestdataPath("sparkle/default.xml")))
	assert.IsType(t, Appcast{}, *a)
	assert.NotNil(t, a.source)
}

func TestAppcast_LoadFromRemoteSource(t *testing.T) {
	// mock the request
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"https://example.com/appcast.xml",
		httpmock.NewBytesResponder(200, getTestdata("sparkle/default.xml")),
	)
	defer httpmock.DeactivateAndReset()

	// test (successful) [URL]
	a := New()
	p, err := a.LoadFromRemoteSource("https://example.com/appcast.xml")
	assert.Nil(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &SparkleRSSFeedAppcast{}, p)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())
	assert.IsType(t, &SparkleRSSFeedAppcast{}, a.source.Appcast())

	// test (successful) [Request]
	a = New()
	r, _ := NewRequest("https://example.com/appcast.xml")
	p, err = a.LoadFromRemoteSource(r)
	assert.Nil(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &SparkleRSSFeedAppcast{}, p)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())
	assert.IsType(t, &SparkleRSSFeedAppcast{}, a.source.Appcast())

	// test "Invalid URL" error
	a = New()
	url := "http://192.168.0.%31/"
	p, err = a.LoadFromRemoteSource(url)
	assert.Error(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.EqualError(t, err, fmt.Sprintf("parse %s: invalid URL escape \"%%31\"", url))
	assert.Nil(t, a.Source())

	// test "Invalid request" error
	a = New()
	p, err = a.LoadFromRemoteSource("invalid")
	assert.Error(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.EqualError(t, err, "Get invalid: no responder found")
	assert.Nil(t, a.Source())
}

func TestAppcast_LoadFromLocalSource(t *testing.T) {
	// test (successful)
	path := getTestdataPath("sparkle/default.xml")
	content := getTestdata("sparkle/default.xml")

	localSourceReadFile = func(filename string) ([]byte, error) {
		return content, nil
	}

	a := New()
	p, err := a.LoadFromLocalSource(path)
	assert.IsType(t, &Appcast{}, a)
	assert.IsType(t, &SparkleRSSFeedAppcast{}, p)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Source().Content())
	assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
	assert.NotNil(t, a.Source().Checksum())
	assert.IsType(t, &SparkleRSSFeedAppcast{}, a.source.Appcast())

	// test (error)
	localSourceReadFile = func(filename string) ([]byte, error) {
		return nil, fmt.Errorf("error")
	}

	a = New()
	p, err = a.LoadFromLocalSource(path)
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Nil(t, a.Source())

	// test (unmarshalling error)
	localSourceReadFile = func(filename string) ([]byte, error) {
		return []byte("invalid"), nil
	}

	a = New()
	p, err = a.LoadFromLocalSource(path)
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.Error(t, err)
	assert.EqualError(t, err, "releases can't be unmarshaled from the \"Unknown\" provider")

	localSourceReadFile = ioutil.ReadFile
}

func TestAppcast_GenerateSourceChecksum(t *testing.T) {
	// preparations
	a := newTestSparkleRSSFeedAppcast()
	assert.Nil(t, a.Source().Checksum())

	// test
	result := a.GenerateSourceChecksum(MD5)
	assert.Equal(t, result.String(), a.Source().Checksum().String())
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", result.String())
	assert.Equal(t, MD5, a.Source().Checksum().Algorithm())
}

func TestAppcast_LoadSource(t *testing.T) {
	// preparations
	a := New(NewLocalSource(getTestdataPath("sparkle/default.xml")))
	assert.Nil(t, a.source.Content())

	// test
	a.LoadSource()
	assert.NotNil(t, a.source.Content())
}

func TestAppcast_UnmarshalReleases_Unknown(t *testing.T) {
	// preparations
	a := newTestAppcast()

	// provider "Unknown"
	p, err := a.UnmarshalReleases()
	assert.Error(t, err)
	assert.IsType(t, &Appcast{}, a)
	assert.Nil(t, p)
	assert.EqualError(t, err, "releases can't be unmarshaled from the \"Unknown\" provider")
	assert.Nil(t, a.source.Appcast())
}

func TestAppcast_UnmarshalReleases_SparkleRSSFeed(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"sparkle/attributes_as_elements.xml": {
			"checksum": "d59d258ce0b06d4c6216f6589aefb36e2bd37fbd647f175741cc248021e0e8b4",
			"releases": 4,
		},
		"sparkle/default_asc.xml": {
			"checksum": "9f8d8eb4c8acfdd53e3084fe5f59aa679bf141afc0c3887141cd2bdfe1427b41",
			"releases": 4,
		},
		"sparkle/default.xml": {
			"checksum": "0cb017e2dfd65e07b54580ca8d4eedbfcf6cef5824bcd9539a64afb72fa9ce8c",
			"releases": 4,
		},
		"sparkle/incorrect_namespace.xml": {
			"checksum": "ff464014dc6a2f6868aca7c3b42521930f791de5fc993d1cc19d747598bcd760",
			"releases": 4,
		},
		"sparkle/invalid_pubdate.xml": {
			"checksum": "9a59f9d0ccd08b317cf784656f6a5bd0e5a1868103ec56d3364baec175dd0da1",
			"releases": 4,
		},
		// "sparkle/multiple_enclosure.xml": {
		// 	"checksum": "48fc8531b253c5d3ed83abfe040edeeafb327d103acbbacf12c2288769dc80b9",
		// 	"releases": 4,
		// },
		"sparkle/no_releases.xml": {
			"checksum": "befd99d96be280ca7226c58ef1400309905ad20d2723e69e829cf050e802afcf",
			"releases": 0,
		},
		"sparkle/only_version.xml": {
			"checksum": "ee5a775fec4d7b95843e284bff6f35f7df30d76af2d1d7c26fc02f735383ef7f",
			"releases": 4,
		},
		"sparkle/prerelease.xml": {
			"checksum": "8e44fccf005ad4720bcc75b9afffb035befade81bdf9f587984c26842dd7c759",
			"releases": 4,
		},
		"sparkle/single.xml": {
			"checksum": "c59ec641579c6bad98017db7e1076a2997cdef7fff315323dd7f0cabed638d50",
			"releases": 1,
		},
		"sparkle/without_namespaces.xml": {
			"checksum": "888494294fc74990e4354689a02e50ff425cfcbd498162fdffd5b3d1cd096fa1",
			"releases": 4,
		},
	}

	errorTestCases := map[string]string{
		"sparkle/invalid_version.xml": "Malformed version: invalid",
		"sparkle/with_comments.xml":   "version is required, but it's not specified in release #1",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// test (successful)
	for filename, data := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		assert.Nil(t, a.Source())
		assert.Len(t, a.releases, 0)

		// load from URL
		src, err := NewRemoteSource("https://example.com/appcast.xml")
		a.SetSource(src)
		a.Source().Load()
		assert.Nil(t, err)
		assert.Equal(t, SparkleRSSFeed, a.Source().Provider())
		assert.NotEmpty(t, a.Source().Content())
		assert.NotNil(t, a.Source().Checksum())
		assert.Equal(t, data["checksum"].(string), a.Source().Checksum().String())
		assert.Len(t, a.releases, 0)

		// releases
		p, err := a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.IsType(t, &SparkleRSSFeedAppcast{}, p)
		assert.Len(t, a.releases, data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", filename))
		assert.IsType(t, &SparkleRSSFeedAppcast{}, a.source.Appcast())
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")

		// test
		p, err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, p)
		assert.EqualError(t, err, errorMsg)
		assert.Nil(t, a.source.Appcast())
	}
}

func TestAppcast_UnmarshalReleases_SourceForgeRSSFeed(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"sourceforge/default.xml": {
			"checksum": "d4afcf95e193a46b7decca76786731c015ee0954b276e4c02a37fa2661a6a5d0",
			"releases": 4,
		},
		"sourceforge/empty.xml": {
			"checksum": "569cb5c8fa66b2bae66e7c0d45e6fbbeb06a5f965fc7e6884ff45aab4f17b407",
			"releases": 0,
		},
		"sourceforge/invalid_pubdate.xml": {
			"checksum": "160885aaaa2f694b5306e91ea20d08ef514f424e51704947c9f07fffec787cf6",
			"releases": 4,
		},
		"sourceforge/single.xml": {
			"checksum": "5384ed38515985f60f990c125f1cceed0261c2c5c2b85181ebd4214a7bc709de",
			"releases": 1,
		},
	}

	errorTestCases := map[string]string{
		"sourceforge/invalid_version.xml": "version is required, but it's not specified in release #2",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// test (successful)
	for filename, data := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		assert.Nil(t, a.Source())
		assert.Len(t, a.releases, 0)

		// load from URL
		src, err := NewRemoteSource("https://example.com/appcast.xml")
		a.SetSource(src)
		a.Source().Load()
		assert.Nil(t, err)
		assert.Equal(t, SourceForgeRSSFeed, a.Source().Provider())
		assert.NotEmpty(t, a.Source().Content())
		assert.NotNil(t, a.Source().Checksum())
		assert.Equal(t, data["checksum"].(string), a.Source().Checksum().String())
		assert.Len(t, a.releases, 0)

		// releases
		p, err := a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.IsType(t, &SourceForgeRSSFeedAppcast{}, p)
		assert.Len(t, a.releases, data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", filename))
		assert.IsType(t, &SourceForgeRSSFeedAppcast{}, a.source.Appcast())
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")

		// test
		p, err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, p)
		assert.EqualError(t, err, errorMsg)
		assert.Nil(t, a.source.Appcast())
	}
}

func TestAppcast_UnmarshalReleases_GitHubAtomFeed(t *testing.T) {
	testCases := map[string]map[string]interface{}{
		"github/default.xml": {
			"checksum": "c28ff87daf2c02471fd2c836b7ed3776d927a8febbb6b8961daf64ce332f6185",
			"releases": 4,
		},
		"github/invalid_pubdate.xml": {
			"checksum": "52f87bba760a4e5f8ee418cdbc3806853d79ad10d3f961e5c54d1f5abf09b24b",
			"releases": 4,
		},
	}

	errorTestCases := map[string]string{
		"github/invalid_version.xml": "Malformed version: invalid",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// test (successful)
	for filename, data := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		assert.Nil(t, a.Source())
		assert.Len(t, a.releases, 0)

		// load from URL
		src, err := NewRemoteSource("https://example.com/appcast.xml")
		a.SetSource(src)
		a.Source().Load()
		assert.Nil(t, err)
		assert.Equal(t, GitHubAtomFeed, a.Source().Provider())
		assert.NotEmpty(t, a.Source().Content())
		assert.NotNil(t, a.Source().Checksum())
		assert.Equal(t, data["checksum"].(string), a.Source().Checksum().String())
		assert.Len(t, a.releases, 0)

		// releases
		p, err := a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.IsType(t, &GitHubAtomFeedAppcast{}, p)
		assert.Len(t, a.releases, data["releases"].(int), fmt.Sprintf("%s: number of releases doesn't match", filename))
		assert.IsType(t, &GitHubAtomFeedAppcast{}, a.source.Appcast())
	}

	// test (error)
	for filename, errorMsg := range errorTestCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")

		// test
		p, err := a.UnmarshalReleases()
		assert.Error(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.Nil(t, p)
		assert.EqualError(t, err, errorMsg)
		assert.Nil(t, a.source.Appcast())
	}
}

func TestAppcast_Uncomment_Unknown(t *testing.T) {
	// preparations
	a := newTestAppcast()

	// test
	err := a.Uncomment()
	assert.EqualError(t, err, "uncommenting is not available for the \"Unknown\" provider")
	a.SetSource(nil)
	err = a.Uncomment()
	assert.EqualError(t, err, "no source")
}

func TestAppcast_Uncomment_SparkleRSSFeed(t *testing.T) {
	// preparations
	a := newTestAppcast(getTestdata("sparkle/with_comments.xml"))
	a.source.SetProvider(SparkleRSSFeed)

	regexCommentStart := regexp.MustCompile(`<!--([[:space:]]*)?<`)
	regexCommentEnd := regexp.MustCompile(`>([[:space:]]*)?-->`)

	// test
	err := a.Uncomment()
	assert.Nil(t, err)
	for _, commentLine := range []int{13, 20} {
		line, _ := getLine(commentLine, a.Source().Content())
		check := regexCommentStart.MatchString(line) && regexCommentEnd.MatchString(line)
		assert.False(t, check)
	}
}

func TestAppcast_Uncomment_SourceForgeRSSFeed(t *testing.T) {
	// preparations
	a := newTestAppcast(getTestdata("sourceforge/default.xml"))
	a.source.SetProvider(SourceForgeRSSFeed)

	// test
	err := a.Uncomment()
	assert.Error(t, err)
	assert.EqualError(t, err, "uncommenting is not available for the \"SourceForge RSS Feed\" provider")
}

func TestAppcast_Uncomment_GitHubAtomFeed(t *testing.T) {
	// preparations
	a := newTestAppcast(getTestdata("github/default.xml"))
	a.source.SetProvider(GitHubAtomFeed)

	// test
	err := a.Uncomment()
	assert.Error(t, err)
	assert.EqualError(t, err, "uncommenting is not available for the \"GitHub Atom Feed\" provider")
}

func TestAppcast_SortReleasesByVersions(t *testing.T) {
	testCases := []string{
		"sparkle/attributes_as_elements.xml",
		"sparkle/default_asc.xml",
		"sparkle/default.xml",
		"sparkle/incorrect_namespace.xml",
		// "sparkle/multiple_enclosure.xml",
		"sparkle/without_namespaces.xml",
	}

	// preparations for mocking the request
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, filename := range testCases {
		// mock the request
		httpmock.RegisterResponder(
			"GET",
			"https://example.com/appcast.xml",
			httpmock.NewBytesResponder(200, getTestdata(filename)),
		)

		// preparations
		a := New()
		a.LoadFromRemoteSource("https://example.com/appcast.xml")
		p, err := a.UnmarshalReleases()
		assert.Nil(t, err)
		assert.IsType(t, &Appcast{}, a)
		assert.IsType(t, &SparkleRSSFeedAppcast{}, p)

		// test (ASC)
		a.SortReleasesByVersions(ASC)
		assert.Equal(t, "1.0.0", a.releases[0].Version().String())

		// test (DESC)
		a.SortReleasesByVersions(DESC)
		assert.Equal(t, "2.0.0", a.releases[0].Version().String())
	}
}

func TestAppcast_Filters(t *testing.T) {
	// mock the request
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"https://example.com/appcast.xml",
		httpmock.NewBytesResponder(200, getTestdata("sparkle/prerelease.xml")),
	)
	defer httpmock.DeactivateAndReset()

	// preparations
	a := New()
	a.LoadFromRemoteSource("https://example.com/appcast.xml")
	a.UnmarshalReleases()

	// Appcast.FilterReleasesByTitle
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByTitle("Release 1.0")
	assert.Len(t, a.releases, 2)
	a.FilterReleasesByTitle("Release 1.0.0", true)
	assert.Len(t, a.releases, 1)
	assert.Equal(t, "Release 1.0.1", a.releases[0].Title())
	a.ResetFilters()

	// Appcast.FilterReleasesByMediaType
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByMediaType("application/octet-stream")
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByMediaType("test", true)
	assert.Len(t, a.releases, 4)
	a.ResetFilters()

	// Appcast.FilterReleasesByURL
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByURL(`app_1.*dmg$`)
	assert.Len(t, a.releases, 3)
	a.FilterReleasesByURL(`app_1.0.*dmg$`, true)
	assert.Len(t, a.releases, 1)
	a.ResetFilters()

	// Appcast.FilterReleasesByPrerelease
	assert.Len(t, a.releases, 4)
	a.FilterReleasesByPrerelease()
	assert.Len(t, a.releases, 1)
	a.ResetFilters()

	assert.Len(t, a.releases, 4)
	a.FilterReleasesByPrerelease(true)
	assert.Len(t, a.releases, 3)
	a.ResetFilters()
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
	assert.Nil(t, a.originalReleases)

	// test
	a.SetReleases([]release.Releaser{&release.Release{}})
	assert.Len(t, a.releases, 1)
}

func TestAppcast_FirstRelease(t *testing.T) {
	// preparations
	a := newTestSparkleRSSFeedAppcast(getTestdata("sparkle/default.xml"))
	a.UnmarshalReleases()

	// test
	assert.Equal(t, a.releases[0].Version().String(), a.FirstRelease().Version().String())
}

func TestAppcast_OriginalReleases(t *testing.T) {
	a := newTestAppcast()
	assert.Equal(t, a.originalReleases, a.OriginalReleases())
}

func TestAppcast_SetOriginalReleases(t *testing.T) {
	// preparations
	a := newTestAppcast()
	assert.Nil(t, a.originalReleases)

	// test
	a.SetOriginalReleases([]release.Releaser{&release.Release{}})
	assert.Len(t, a.originalReleases, 1)
}
