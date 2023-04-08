package test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGoChangeInit(t *testing.T) {
	givenCoverageInstrumentedBinaryInPATH(t)
	runInTempDir(t, func(t *testing.T, dir string) {
		whenCalledWithArgs(t, "init")
		thenEmptyChangelogFileIsCreated(t)
	})
}

func givenCoverageInstrumentedBinaryInPATH(t *testing.T) {
	t.Helper()

	// TODO: Fail if not built?

	absBinaryPath, err := filepath.Abs("../gochange")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	absBinaryDir := filepath.Dir(absBinaryPath)
	os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), absBinaryDir))
}

func whenCalledWithArgs(t *testing.T, args ...string) []byte {
	t.Log()

	gochangeCmd := exec.Command("gochange", "init")
	gochangeCmd.Env = []string{
		"GOCOVERDIR=coverage_data",
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}
	output, err := gochangeCmd.Output()
	if err != nil {
		t.Log(string(output))
		t.Error(err)
		t.FailNow()
	}

	return output
}

func thenEmptyChangelogFileIsCreated(t *testing.T) {
	changelogFilepath, err := filepath.Abs("CHANGELOG.md")
	requireNoError(t, err)

	_, err = os.Stat(changelogFilepath)
	requireNoError(t, err)
}

func runInTempDir(t *testing.T, f func(t *testing.T, dir string)) {
	dir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.ReadFile(dir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	os.Chdir(dir)
	defer os.Chdir(originalDir)

	f(t, dir)
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func requireNoError(t *testing.T, err error) {
	if err == nil {
		return
	}

	t.Error(err)
	t.FailNow()
}
