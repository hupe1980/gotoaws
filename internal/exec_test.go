package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSSHArgs(t *testing.T) {
	expected := "-i key.pem ec2-user@i-123456789"
	actual := SSHArgs("ec2-user", "i-123456789", "key.pem", "")
	require.Equal(t, expected, actual)
}

func TestSSHArgsWithCmd(t *testing.T) {
	expected := "-i key.pem ec2-user@i-123456789 uname -a"
	actual := SSHArgs("ec2-user", "i-123456789", "key.pem", "uname -a")
	require.Equal(t, expected, actual)
}

func TestSCPArgs(t *testing.T) {
	expected := "-i key.pem ./test.txt ec2-user@i-123456789:/opt/"
	actual := SCPArgs("ec2-user", "i-123456789", "key.pem", "./test.txt", "/opt/")
	require.Equal(t, expected, actual)
}
