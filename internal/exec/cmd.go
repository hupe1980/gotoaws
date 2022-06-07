package exec

import (
	"io"
	"os"
	"os/exec"
)

// CmdOption is a type alias to configure a command.
type CmdOption func(cmd *exec.Cmd)

// Stdin sets the internal *exec.Cmd's Stdin field.
func Stdin(r io.Reader) CmdOption {
	return func(c *exec.Cmd) {
		c.Stdin = r
	}
}

// Stdout sets the internal *exec.Cmd's Stdout field.
func Stdout(writer io.Writer) CmdOption {
	return func(c *exec.Cmd) {
		c.Stdout = writer
	}
}

// Stderr sets the internal *exec.Cmd's Stderr field.
func Stderr(writer io.Writer) CmdOption {
	return func(c *exec.Cmd) {
		c.Stderr = writer
	}
}

type cmdRunner interface {
	Run() error
}

// Cmd runs external commands, it wraps the exec.Command function from the stdlib so that
// running external commands can be unit tested.
type Cmd struct {
	command func(name string, args []string, opts ...CmdOption) cmdRunner
}

// NewCmd returns a Cmd that can run external commands.
func NewCmd() *Cmd {
	return &Cmd{
		command: func(name string, args []string, opts ...CmdOption) cmdRunner {
			cmd := exec.Command(name, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			for _, opt := range opts {
				opt(cmd)
			}
			return cmd
		},
	}
}

// Run starts the named command and waits until it finishes.
func (c *Cmd) Run(name string, args []string, opts ...CmdOption) error {
	cmd := c.command(name, args, opts...)
	return cmd.Run()
}
