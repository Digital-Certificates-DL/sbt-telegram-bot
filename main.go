package main

import (
	"os"

	"gitlab.com/tokend/course-certificates/sbt-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
