package main

import (
	"bufio"
	"errors"
	"os"
	"time"

	"github.com/mrombout/gochange/changelog"
	"github.com/spf13/cobra"
)

// Force indicates whether to overwrite an existing release if one already exists.
var Force bool

// Merge indiacates whether to merge with an existing release if one already exists.
var Merge bool

func init() {
	rootCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().BoolVarP(&Force, "force", "f", false, "overwrite an existing release if it already exists")
	releaseCmd.Flags().BoolVarP(&Merge, "merge", "m", false, "merge with existing release if it already exists")
}

var releaseCmd = &cobra.Command{
	Use:   "release <version>",
	Short: "Move all unreleased changes to a release",
	Long:  "Moves all unreleased changes to a release.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one argument")
		}

		// parse changelog
		// move unreleased to release
		// write changelog

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.OpenFile("CHANGELOG.md", os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)

		tokens, err := changelog.Lex(scanner)
		if err != nil {
			panic(err)
		}

		currentChangelog, err := changelog.Parse(tokens)
		if err != nil {
			panic(err)
		}

		newRelease := currentChangelog.Unreleased
		newRelease.Name = args[0]
		newRelease.Date = time.Now().Format("2006-01-02")

		if len(currentChangelog.Releases) > 0 {
			newRelease.PreviousRelease = &currentChangelog.Releases[0]
		}

		currentChangelog.Releases = append([]changelog.Release{newRelease}, currentChangelog.Releases...)

		currentChangelog.Unreleased = changelog.Release{}

		file.Seek(0, 0)
		file.Truncate(0)
		changelog.Render(currentChangelog, file)
	},
}
