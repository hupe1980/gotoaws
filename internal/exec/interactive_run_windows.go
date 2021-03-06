//go:build windows
// +build windows

package exec

import (
	"os"
	"os/signal"
)

// InteractiveRun runs the input command that starts a child process.
func (c *Cmd) InteractiveRun(name string, args ...string) error {
	sig := make(chan os.Signal, 1)

	// See https://golang.org/pkg/os/signal/#hdr-Windows
	signal.Notify(sig, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	cmd := c.command(name, args, Stdout(os.Stdout), Stdin(os.Stdin), Stderr(os.Stderr))

	return cmd.Run()
}
