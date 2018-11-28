package main

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

var appFs = afero.NewOsFs()

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
