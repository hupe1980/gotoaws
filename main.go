package main

import (
	"github.com/hupe1980/gotoaws/cmd"
)

// nolint: gochecknoglobals // ok
var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
