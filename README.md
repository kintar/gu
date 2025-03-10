# go-utils

This module provides common functionality for applications written in Go.

## MakeVersion
`cmd/makeversion` contains a command which can be run with a `//go:generate` directive in order to read version
information from a git repository's tags and create a `version.txt` file for embedding into compiled binaries.

## datautil
Package `datautil` currently contains a single function, `CreateStructByUserInput`, which makes filling configuration
structures from user input much simpler. Full documentation is in the godocs for the function.