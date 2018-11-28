package cmd

import (
	"os/exec"
	"testing"
)

func TestInit(t *testing.T) {
	// Arrange
	assertGoChangeInstalled(t)
	cmd := exec.Command("gochange", "init")

	// Act
	output, err := cmd.CombinedOutput()

	// Assert
	if err != nil {
		t.Fatalf("expected result to be nil, but got %t", err)
	}
	if string(output) != "" {
		t.Fatalf("wron!")
	}
	// TODO: Assert that a CHANGELOG.md has been created
	// TODO: Assert that created CHANGELOG.md complies to oracle
	// TODO: Assert output is correct according to oracle
}

func TestAddedEntry(t *testing.T) {
	// Arrange
	assertGoChangeInstalled(t)
	cmd := exec.Command("gochange", "Added a bug.")
	// TODO: Make sure that CHANGELOG.md is in a known state

	// Act
	output, err := cmd.CombinedOutput()

	// Assert
	if err != nil {
		t.Fatalf("expected result to be nil, but got %t", err)
	}
	if string(output) != "" {
		t.Fatalf("wron!")
	}
	// TODO: Assert that new Added entry has been added
}

func TestRelease(t *testing.T) {
	// Arrange
	assertGoChangeInstalled(t)
	cmd := exec.Command("gochange", "release", "v1.0.0")
	// TODO: Make sure that CHANGELOG.md is in a known state and contains all sections

	// Act
	output, err := cmd.CombinedOutput()

	// Assert
	if err != nil {
		t.Fatalf("expected result to be nil, but got %t", err)
	}
	if string(output) != "" {
		t.Fatalf("wron!")
	}
	// TODO: Assert that all unreleased changed have been moved to the correct section
}

func assertGoChangeInstalled(t *testing.T) {
	_, err := exec.LookPath("gochange")
	if err != nil {
		t.Fatal("Application 'gochange' can not be found in $PATH")
	}
}
