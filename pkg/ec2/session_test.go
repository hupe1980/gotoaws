package ec2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHArgs(t *testing.T) {
	expected := "-i key.pem ec2-user@i-123456789"
	actual := sshArgs(&RunSSHInput{
		User:       "ec2-user",
		Identity:   "key.pem",
		InstanceID: "i-123456789",
	})
	assert.Equal(t, expected, actual)
}

func TestSSHArgsWithCmd(t *testing.T) {
	expected := "-i key.pem ec2-user@i-123456789 uname -a"
	actual := sshArgs(&RunSSHInput{
		User:       "ec2-user",
		Identity:   "key.pem",
		InstanceID: "i-123456789",
		Command:    "uname -a",
	})
	assert.Equal(t, expected, actual)
}

func TestSSHArgsWithFwd(t *testing.T) {
	expected := "-L 80:intra.example.com:80 -i key.pem ec2-user@i-123456789"
	actual := sshArgs(&RunSSHInput{
		User:                "ec2-user",
		Identity:            "key.pem",
		InstanceID:          "i-123456789",
		LocalPortForwarding: "80:intra.example.com:80",
	})
	assert.Equal(t, expected, actual)
}

func TestSCPArgs(t *testing.T) {
	expected := "-i key.pem ./test.txt ec2-user@i-123456789:/opt/"
	actual := scpArgs(&RunSCPInput{
		User:       "ec2-user",
		Identity:   "key.pem",
		InstanceID: "i-123456789",
		Mode:       SCPModeSending,
		Sources:    []string{"./test.txt"},
		Target:     "/opt/",
	})
	assert.Equal(t, expected, actual)
}
