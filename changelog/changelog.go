package changelog

import (
	"strings"
	"time"
)

const sectionAdded = "Added"
const sectionChanged = "Changed"
const sectionDeprecated = "Deprecated"
const sectionRemoved = "Removed"
const sectionFixed = "Fixed"
const sectionSecurity = "Security"

// Changelog represents a projects changelog.
type Changelog struct {
	URL         string
	Description string
	Unreleased  Release
	Releases    []Release

	LatestRelease Release
}

// AddEntry adds a new entry to the unreleased section of the changelog.
func (c *Changelog) AddEntry(entry string) {
	switch {
	case strings.HasPrefix(entry, sectionAdded):
		c.Unreleased.Added = append(c.Unreleased.Added, Entry{
			Description: entry,
		})
	case strings.HasPrefix(entry, sectionChanged):
		c.Unreleased.Changed = append(c.Unreleased.Changed, Entry{
			Description: entry,
		})
	case strings.HasPrefix(entry, sectionDeprecated):
		c.Unreleased.Deprecated = append(c.Unreleased.Deprecated, Entry{
			Description: entry,
		})
	case strings.HasPrefix(entry, sectionRemoved):
		c.Unreleased.Removed = append(c.Unreleased.Removed, Entry{
			Description: entry,
		})
	case strings.HasPrefix(entry, sectionFixed):
		c.Unreleased.Fixed = append(c.Unreleased.Fixed, Entry{
			Description: entry,
		})
	case strings.HasPrefix(entry, sectionSecurity):
		c.Unreleased.Security = append(c.Unreleased.Security, Entry{
			Description: entry,
		})
	}
}

// Release creates a new release and moves all unreleased entries to the new release.
func (c *Changelog) Release(name string) {
	newRelease := c.Unreleased
	newRelease.Name = name
	newRelease.Date = time.Now().Format("2006-01-02")

	if len(c.Releases) > 0 {
		newRelease.PreviousRelease = &c.Releases[0]
	}

	c.Releases = append([]Release{newRelease}, c.Releases...)

	c.Unreleased = Release{}
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
