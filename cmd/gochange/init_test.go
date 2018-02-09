package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_WhenNoChangelog_CreatesChangelog(t *testing.T) {
	// arrange
	buf := bytes.Buffer{}
	initCmd.SetOutput(&buf)

	// act
	initCmd.Run(initCmd, []string{})

	// assert
	output := buf.String()

	if !strings.Contains(output, "Initializing changelog...") {
		t.Errorf("output '%s' did not container '%s'", output, "Initializing changelog...")
	}
	if !strings.Contains(output, "Changelog has been initialized.") {

	}
}
