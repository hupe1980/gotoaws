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

type RunSSHInput struct {
	User       *string
	InstanceID *string
	Identity   *string
	Command    *string
}

type RunSCPInput struct {
	User       *string
	InstanceID *string
	Identity   *string
	Source     *string
	Target     *string
}

type session struct {
	id         *string
	streamUrl  *string
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
		streamUrl:  output.Session.StreamUrl,
		tokenValue: output.Session.TokenValue,
		client:     ssm.NewFromConfig(cfg.awsCfg),
		input: &ssm.StartSessionInput{
			Target: aws.String(fmt.Sprintf("ecs:%s_%s_%s", *input.Cluster, *input.Task, *input.Container)),
		},
		profile: cfg.profile,
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
		streamUrl:  output.StreamUrl,
		tokenValue: output.TokenValue,
		client:     client,
		input:      input,
		profile:    cfg.profile,
		plugin:     cfg.plugin,
		region:     cfg.awsCfg.Region,
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
		"StreamUrl":  sess.streamUrl,
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
	for _, sep := range strings.Split(sshArgs(*input.User, *input.InstanceID, *input.Identity, *input.Command), " ") {
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
	for _, sep := range strings.Split(scpArgs(*input.User, *input.InstanceID, *input.Identity, *input.Source, *input.Target), " ") {
		if sep != "" {
			args = append(args, sep)
		}
	}
	if err := runSubprocess("scp", args...); err != nil {
		return err
	}
	return nil
}

func (sess *session) proxyCommand() (string, error) {
	sessJSON, err := json.Marshal(map[string]*string{
		"SessionId":  sess.id,
		"StreamUrl":  sess.streamUrl,
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

func sshArgs(user string, instanceID string, identity string, cmd string) string {
	ssh := fmt.Sprintf("-i %s %s@%s", identity, user, instanceID)
	if cmd == "" {
		return ssh
	}
	return fmt.Sprintf("%s %s", ssh, cmd)
}

func scpArgs(user string, instanceID string, identity string, source string, target string) string {
	return fmt.Sprintf("-i %s %s %s@%s:%s", identity, source, user, instanceID, target)
}
