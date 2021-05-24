package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSSHArgs(t *testing.T) {
	expected := "-i key.pem ec2-user@i-123456789"
	actual := sshArgs("ec2-user", "i-123456789", "key.pem", "", "")
	require.Equal(t, expected, actual)
}

func TestSSHArgsWithCmd(t *testing.T) {
	expected := "-i key.pem ec2-user@i-123456789 uname -a"
	actual := sshArgs("ec2-user", "i-123456789", "key.pem", "", "uname -a")
	require.Equal(t, expected, actual)
}

func TestSSHArgsWithFwd(t *testing.T) {
	expected := "-L 80:intra.example.com:80 -i key.pem ec2-user@i-123456789"
	actual := sshArgs("ec2-user", "i-123456789", "key.pem", "80:intra.example.com:80", "")
	require.Equal(t, expected, actual)
}

func TestSCPArgs(t *testing.T) {
	expected := "-i key.pem ./test.txt ec2-user@i-123456789:/opt/"
	actual := scpArgs("ec2-user", "i-123456789", "key.pem", SCPModeSending, []string{"./test.txt"}, "/opt/")
	require.Equal(t, expected, actual)
}
