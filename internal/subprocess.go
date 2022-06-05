package internal

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func RunSubprocess(process string, args ...string) error {
	call := exec.Command(process, args...)
	call.Stderr = os.Stderr
	call.Stdout = os.Stdout
	call.Stdin = os.Stdin

	signal.Ignore(syscall.SIGINT)

	if err := call.Run(); err != nil {
		return err
	}

	return nil
}
