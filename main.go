package main

import (
	"github.com/hupe1980/ec2connect/cmd"
)

const (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
