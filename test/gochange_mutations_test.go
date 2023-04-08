package test

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/golden"
)

func TestGoChangeMutate(t *testing.T) {
	goldenFile, err := filepath.Abs("testdata/mutated_changelog.md")
	requireNoError(t, err)

	givenCoverageInstrumentedBinaryInPATH(t)
	runInTempDir(t, func(t *testing.T, dir string) {
		whenCalledWithArgs(t, "Added some new stuff")
		thenChangelogIsEqualTo(t, goldenFile)
	})
}

func thenChangelogIsEqualTo(t *testing.T, filename string) {
	content, err := os.ReadFile("CHANGELOG.md")
	requireNoError(t, err)

	golden.AssertBytes(t, content, filename)
}
