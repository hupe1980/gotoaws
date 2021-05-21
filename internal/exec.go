package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func RunSubprocess(process string, args ...string) error {
	call := exec.Command(process, args...)
	call.Stderr = os.Stderr
	call.Stdout = os.Stdout
	call.Stdin = os.Stdin

	if err := call.Run(); err != nil {
		return err
	}
	return nil
}

func SSHArgs(user string, instanceID string, identity string, cmd string) string {
	ssh := fmt.Sprintf("-i %s %s@%s", identity, user, instanceID)
	if cmd == "" {
		return ssh
	}
	return fmt.Sprintf("%s %s", ssh, cmd)
}

func SCPArgs(user string, instanceID string, identity string, source string, target string) string {
	return fmt.Sprintf("-i %s %s %s@%s:%s", identity, source, user, instanceID, target)
}
