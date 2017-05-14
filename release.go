package appcast

import "github.com/hashicorp/go-version"

// A Release represents an application release.
type Release struct {
	// Version specifies the release version. It should follow the SemVer
	// specification.
	Version *version.Version

	// Build specifies the release build. This could have any value.
	Build string

	// DownloadUrls specify an array of download URLs for release.
	DownloadURLs []string

	// IsPrerelease specifies if the current release is not stable.
	//
	// By default, each release is considered to be stable, so the default value
	// is "false". If the release version, build or any other provider specific
	// value point that the release is not stable, the value should become "true".
	IsPrerelease bool
}

// NewRelease returns a new Release instance pointer. Requires both version and
// build strings. By default, Release.IsPrerelease is set to "false", so the
// release will be considered as stable.
func NewRelease(version string, build string) (*Release, error) {
	r := &Release{
		IsPrerelease: false,
	}

	// add version
	err := r.SetVersion(version)
	if err != nil {
		return nil, err
	}

	// add build, if its not empty
	if build != "" {
		r.Build = build
	}

	return r, nil
}

// SetVersion sets the Release.Version from the provided version value string.
// Returns an error, if the provided version string value doesn't follow SemVer
// specification.
func (r *Release) SetVersion(value string) error {
	v, err := version.NewVersion(value)
	if err != nil {
		return err
	}

	r.Version = v

	return nil
}
