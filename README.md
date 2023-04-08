# GoChange
Go tool that helps writing changelogs.

[![master](https://github.com/mrombout/gochange/actions/workflows/master.yml/badge.svg)](https://github.com/mrombout/gochange/actions/workflows/master.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrombout/gochange)](https://goreportcard.com/report/github.com/mrombout/gochange)
[![codecov](https://codecov.io/gh/mrombout/gochange/branch/master/graph/badge.svg)](https://codecov.io/gh/mrombout/gochange)

## Usage

To initialize an empty changelog for your project use the command described below.

    gochange init

To add a new entry to the unreleased section of the changelog use the commands described below. The section is determined by the first world in your sentence.

    gochange "Added links to navigational page."
    gochange "Changed way navigation is handled."
    gochange "Deprecated the use of the navigational page."
    gochange "Removed links to navigational page."
    gochange "Fixed link to navigation page."
    gochange "Security navigational page is not longer a threat."

To bump all changes in the unreleased section up to a specific version use the command described below.

    gochange release 0.1.0

