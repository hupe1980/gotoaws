package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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

type session struct {
	id         *string
	streamURL  *string
	tokenValue *string
	client     *ssm.Client
	input      *ssm.StartSessionInput
	profile    string
	plugin     string
	region     string
	timeout    time.Duration
}

type ECSSession interface {
	Close() error
	RunPlugin() error
}

type EC2Session interface {
	Close() error
	RunPlugin() error
	RunSSH(input *RunSSHInput) error
	RunSCP(input *RunSCPInput) error
}

func NewECSSession(cfg *Config, input *ecs.ExecuteCommandInput) (ECSSession, error) {
	client := ecs.NewFromConfig(cfg.awsCfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)

	defer cancel()

	output, err := client.ExecuteCommand(ctx, input)
	if err != nil {
		return nil, err
	}

	return &session{
		id:         output.Session.SessionId,
		streamURL:  output.Session.StreamUrl,
		tokenValue: output.Session.TokenValue,
		client:     ssm.NewFromConfig(cfg.awsCfg),
		input: &ssm.StartSessionInput{
			Target: aws.String(fmt.Sprintf("ecs:%s_%s_%s", *input.Cluster, *input.Task, *input.Container)),
		},
		profile: cfg.Profile,
		plugin:  cfg.plugin,
		region:  cfg.awsCfg.Region,
		timeout: cfg.timeout,
	}, nil
}

func NewEC2Session(cfg *Config, input *ssm.StartSessionInput) (EC2Session, error) {
	client := ssm.NewFromConfig(cfg.awsCfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)

	defer cancel()

	output, err := client.StartSession(ctx, input)
	if err != nil {
		return nil, err
	}

	return &session{
		id:         output.SessionId,
		streamURL:  output.StreamUrl,
		tokenValue: output.TokenValue,
		client:     client,
		input:      input,
		profile:    cfg.Profile,
		plugin:     cfg.plugin,
		region:     cfg.Region,
		timeout:    cfg.timeout,
	}, nil
}

func (sess *session) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), sess.timeout)
	defer cancel()

	_, err := sess.client.TerminateSession(ctx, &ssm.TerminateSessionInput{SessionId: sess.id})
	if err != nil {
		return err
	}

	return nil
}

func (sess *session) RunPlugin() error {
	sessJSON, err := json.Marshal(map[string]*string{
		"SessionId":  sess.id,
		"StreamUrl":  sess.streamURL,
		"TokenValue": sess.tokenValue,
	})
	if err != nil {
		return err
	}

	inputJSON, err := json.Marshal(sess.input)
	if err != nil {
		return err
	}

	return runSubprocess(sess.plugin, string(sessJSON), sess.region, "StartSession", sess.profile, string(inputJSON))
}

func (sess *session) RunSSH(input *RunSSHInput) error {
	pc, err := sess.proxyCommand()
	if err != nil {
		return err
	}

	args := []string{"-o", pc}

	for _, sep := range strings.Split(sshArgs(input), " ") {
		if sep != "" {
			args = append(args, sep)
		}
	}

	if err := runSubprocess("ssh", args...); err != nil {
		return err
	}

	return nil
}

func (sess *session) RunSCP(input *RunSCPInput) error {
	pc, err := sess.proxyCommand()
	if err != nil {
		return err
	}

	args := []string{"-o", pc}

	for _, sep := range strings.Split(scpArgs(input), " ") {
		if sep != "" {
			args = append(args, sep)
		}
	}

	fmt.Println(strings.Join(args, " "))

	if err := runSubprocess("scp", args...); err != nil {
		return err
	}

	return nil
}

func (sess *session) proxyCommand() (string, error) {
	sessJSON, err := json.Marshal(map[string]*string{
		"SessionId":  sess.id,
		"StreamUrl":  sess.streamURL,
		"TokenValue": sess.tokenValue,
	})
	if err != nil {
		return "", err
	}

	inputJSON, err := json.Marshal(sess.input)
	if err != nil {
		return "", err
	}

	pc := fmt.Sprintf("ProxyCommand=%s '%s' %s %s %s '%s'", sess.plugin, string(sessJSON), sess.region, "StartSession", sess.profile, string(inputJSON))

	return pc, nil
}

func runSubprocess(process string, args ...string) error {
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
