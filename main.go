package main

import (
	"github.com/hupe1980/gotoaws/cmd"
)

// nolint: gochecknoglobals
var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
