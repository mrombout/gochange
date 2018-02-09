package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mrombout/gochange/changelog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gochange [change]",
	Short: "Gochange helps updating changelogs.",
	Long:  `A tool that helps to create and update changelogs using a simple command-line interface.`,
	Args: func(cmd *cobra.Command, args []string) error {
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

		change := args[0]

		switch {
		case strings.HasPrefix(change, "Added"):
			currentChangelog.Unreleased.Added = append(currentChangelog.Unreleased.Added, changelog.Entry{
				Description: change,
			})
		case strings.HasPrefix(change, "Changed"):
			currentChangelog.Unreleased.Changed = append(currentChangelog.Unreleased.Changed, changelog.Entry{
				Description: change,
			})
		case strings.HasPrefix(change, "Deprecated"):
			currentChangelog.Unreleased.Deprecated = append(currentChangelog.Unreleased.Deprecated, changelog.Entry{
				Description: change,
			})
		case strings.HasPrefix(change, "Removed"):
			currentChangelog.Unreleased.Removed = append(currentChangelog.Unreleased.Removed, changelog.Entry{
				Description: change,
			})
		case strings.HasPrefix(change, "Fixed"):
			currentChangelog.Unreleased.Fixed = append(currentChangelog.Unreleased.Fixed, changelog.Entry{
				Description: change,
			})
		case strings.HasPrefix(change, "Security"):
			currentChangelog.Unreleased.Security = append(currentChangelog.Unreleased.Security, changelog.Entry{
				Description: change,
			})
		}

		file.Seek(0, 0)
		file.Truncate(0)
		changelog.Render(currentChangelog, file)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
