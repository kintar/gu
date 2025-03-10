//go:build generate

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	c := exec.Command("git", "describe", "--always", "--tags", "--match", "'v*'")
	l := log.New(os.Stderr, "", 0)

	if output, err := c.Output(); err != nil {
		l.Printf("failed to execute git command: %s\n", err.Error())
		os.Exit(1)
	} else {
		version := string(output)
		if len(version) == 0 {
			l.Printf("failed to identify a version: git returned empty data\n")
			os.Exit(1)
		}

		version = strings.ReplaceAll(version, "\r", "")
		version = strings.ReplaceAll(version, "\n", "")

		if version[0] != 'v' {
			version = fmt.Sprintf("%s-dev", version)
		}

		if f, err := os.Create("version.txt"); err != nil {
			l.Printf("failed to create version.txt: %s", err.Error())
			os.Exit(2)
		} else {
			_, err = f.WriteString(version)
			if err != nil {
				l.Printf("failed to write to version.txt: %s", err.Error())
				os.Exit(3)
			}
			fmt.Println(version)
		}
	}

}
