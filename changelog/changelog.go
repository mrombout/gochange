package changelog

// Changelog represents a projects changelog.
type Changelog struct {
	URL         string
	Description string
	Unreleased  Release
	Releases    []Release

	LatestRelease Release
}

// Release represents a single release of a project.
type Release struct {
	Name string
	Date string

	PreviousRelease *Release

	Added      []Entry
	Changed    []Entry
	Deprecated []Entry
	Removed    []Entry
	Fixed      []Entry
	Security   []Entry
}

// Entry represents a single entry for a projects release. An entry must be part
// of one of the sections "Added", "Changed", "Deprecated", "Removed", "Fixed"
// or "Security".
type Entry struct {
	Description string
}

func newChangelog() Changelog {
	return Changelog{
		URL: "http://github.com/",
		LatestRelease: Release{
			Name: "HEAD",
		},
	}
}
