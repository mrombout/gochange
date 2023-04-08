package changelog

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	// arrange
	release020 := Release{
		Name: "0.2.0",
		Date: "2018-08-14",
		Added: []Entry{
			{
				Description: "Some stuff.",
			},
		},
		PreviousRelease: &Release{
			Name: "NONE",
		},
	}
	release100 := Release{
		Name: "1.0.0",
		Date: "2018-12-28",
		Added: []Entry{
			{
				Description: "Some stuff.",
			},
		},
		PreviousRelease: &release020,
	}

	currentChangelog := Changelog{
		URL:         "http://github.com/mrombout/gochange/",
		Description: "Lorum ipsum dolor sit amet consectatur.",
		Unreleased: Release{
			Name: "Unreleased",
			Added: []Entry{
				{
					Description: "Some more stuff.",
				},
				{
					Description: "Even more stuff.",
				},
			},
			Removed: []Entry{
				{
					Description: "Easter egg.",
				},
				{
					Description: "Bitcoin Miner",
				},
			},
			Changed: []Entry{
				{
					Description: "Some things.",
				},
			},
		},
		LatestRelease: release100,
		Releases: []Release{
			release100,
			release020,
		},
	}
	expectedOutput, err := ioutil.ReadFile("testdata/rendered.md")
	if err != nil {
		t.Error(err)
	}
	actualOutput := strings.Builder{}

	// act
	Render(currentChangelog, &actualOutput)

	// assert
	actualOutputStr := actualOutput.String()
	actualExpectedOutputStr := string(expectedOutput)
	if actualOutputStr != actualExpectedOutputStr {
		t.Errorf("actual output does not match expected output")
	}
}
