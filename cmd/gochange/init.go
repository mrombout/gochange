package main

import (
	"os"

	"github.com/mrombout/gochange/changelog"
	"github.com/spf13/cobra"
)

// ForceInit indicates whether to overwrite the changelog if one already exists.
var ForceInit bool

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&ForceInit, "force", "f", false, "overwrite existing changelog if one already exists")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an empty changelog",
	Long:  "Generates an empty changelog template.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Initializing changelog...")

		file, err := os.Create("CHANGELOG.md")
		if err != nil {
			panic(err)
		}

		newChangelog := changelog.Changelog{
			Description: "Lorum ipsum dolor sit amet.",
			URL:         "http://github.com/",
			LatestRelease: changelog.Release{
				Name: "HEAD",
			},
		}
		changelog.Render(newChangelog, file)

		cmd.Println("Changelog has been initialized.")
	},
}
