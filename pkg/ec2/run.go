package ec2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/hupe1980/gotoaws/pkg/config"
)

type CommandRunner struct {
	client     *ssm.Client
	instanceID string
	output     *ssm.SendCommandOutput
	timeout    time.Duration
}

func NewCommandRunner(cfg *config.Config, inst *Instance, command []string) (*CommandRunner, error) {
	client := ssm.NewFromConfig(cfg.AWSConfig)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)

	defer cancel()

	// only linux support (window = "AWS-RunPowerShellScript")
	docName := "AWS-RunShellScript"

	if inst.Platform == "Windows" {
		docName = "AWS-RunPowerShellScript"
	}

	input := &ssm.SendCommandInput{
		DocumentName:   &docName,
		InstanceIds:    []string{inst.ID},
		TimeoutSeconds: int32(60), // 60 seconds
		CloudWatchOutputConfig: &types.CloudWatchOutputConfig{
			CloudWatchOutputEnabled: true,
		},
		Parameters: map[string][]string{
			"commands": {strings.Join(command, " ")},
		},
	}

	output, err := client.SendCommand(ctx, input)
	if err != nil {
		return nil, err
	}

	return &CommandRunner{
		client:     client,
		instanceID: inst.ID,
		output:     output,
		timeout:    cfg.Timeout,
	}, nil
}

func (cmd *CommandRunner) Result() (string, error) {
	input := &ssm.GetCommandInvocationInput{
		CommandId:  cmd.output.Command.CommandId,
		InstanceId: &cmd.instanceID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmd.timeout)
	defer cancel()

	for {
		time.Sleep(1 * time.Second)

		output, err := cmd.client.GetCommandInvocation(ctx, input)
		if err != nil {
			return "", err
		}

		switch output.Status {
		case types.CommandInvocationStatusPending, types.CommandInvocationStatusInProgress, types.CommandInvocationStatusDelayed:
		case types.CommandInvocationStatusSuccess:
			return *output.StandardOutputContent, nil
		case types.CommandInvocationStatusCancelled, types.CommandInvocationStatusCancelling, types.CommandInvocationStatusFailed, types.CommandInvocationStatusTimedOut:
			return "", fmt.Errorf(*output.StandardErrorContent)
		}
	}
}
