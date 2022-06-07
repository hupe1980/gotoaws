package ec2

import (
	"context"
	"fmt"
	"strings"

	aws_ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/internal/exec"
	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/hupe1980/gotoaws/pkg/ssm"
)

type SCPMode string

const (
	SCPModeSending   SCPMode = "sending"
	SCPModeReceiving SCPMode = "receiving"
)

type RunSSHInput struct {
	User                string
	InstanceID          string
	Identity            string
	LocalPortForwarding string
	Command             string
}

type RunSCPInput struct {
	User       string
	InstanceID string
	Identity   string
	Sources    []string
	Target     string
	Mode       SCPMode
}

type Session interface {
	Close() error
	RunPlugin() error
	RunSSH(input *RunSSHInput) error
	RunSCP(input *RunSCPInput) error
}

type session struct {
	ssmSession *ssm.Session
}

func NewSession(cfg *config.Config, input *aws_ssm.StartSessionInput) (Session, error) {
	ssmClient := aws_ssm.NewFromConfig(cfg.AWSConfig)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)

	defer cancel()

	output, err := ssmClient.StartSession(ctx, input)
	if err != nil {
		return nil, err
	}

	return &session{
		ssmSession: &ssm.Session{
			ID:         output.SessionId,
			StreamURL:  output.StreamUrl,
			TokenValue: output.TokenValue,
			SSMClient:  ssmClient,
			Input:      input,
			Profile:    cfg.Profile,
			Plugin:     cfg.Plugin,
			Region:     cfg.AWSConfig.Region,
			Timeout:    cfg.Timeout,
		},
	}, nil
}

func (sess *session) Close() error {
	return sess.ssmSession.Close()
}

func (sess *session) RunPlugin() error {
	return sess.ssmSession.RunPlugin()
}

func (sess *session) RunSSH(input *RunSSHInput) error {
	pc, err := sess.ssmSession.ProxyCommand()
	if err != nil {
		return err
	}

	args := []string{"-o", pc}

	for _, sep := range strings.Split(sshArgs(input), " ") {
		if sep != "" {
			args = append(args, sep)
		}
	}

	cmd := exec.NewCmd()

	return cmd.InteractiveRun("ssh", args...)
}

func (sess *session) RunSCP(input *RunSCPInput) error {
	pc, err := sess.ssmSession.ProxyCommand()
	if err != nil {
		return err
	}

	args := []string{"-o", pc}

	for _, sep := range strings.Split(scpArgs(input), " ") {
		if sep != "" {
			args = append(args, sep)
		}
	}

	cmd := exec.NewCmd()

	return cmd.InteractiveRun("scp", args...)
}

func sshArgs(input *RunSSHInput) string {
	ssh := fmt.Sprintf("-i %s %s@%s", input.Identity, input.User, input.InstanceID)
	if input.LocalPortForwarding != "" {
		ssh = fmt.Sprintf("-L %s %s", input.LocalPortForwarding, ssh)
	}

	if input.Command != "" {
		ssh = fmt.Sprintf("%s %s", ssh, input.Command)
	}

	return ssh
}

func scpArgs(input *RunSCPInput) string {
	if input.Mode == SCPModeSending {
		return fmt.Sprintf("-i %s %s %s@%s:%s", input.Identity, strings.Join(input.Sources, " "), input.User, input.InstanceID, input.Target)
	}

	s := input.Sources[0]

	if len(input.Sources) > 1 {
		s = fmt.Sprintf("{%s}", strings.Join(input.Sources, ","))
	}

	return fmt.Sprintf("-i %s %s@%s:%s %s", input.Identity, input.User, input.InstanceID, s, input.Target)
}
