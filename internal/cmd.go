package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type Command struct {
	client     *ssm.Client
	instanceID string
	output     *ssm.SendCommandOutput
	timeout    time.Duration
}

func NewLinuxCommand(cfg *Config, instanceID string, command string) (*Command, error) {
	client := ssm.NewFromConfig(cfg.awsCfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	// only linux support (window = "AWS-RunPowerShellScript")
	docName := "AWS-RunShellScript"

	input := &ssm.SendCommandInput{
		DocumentName:   &docName,
		InstanceIds:    []string{instanceID},
		TimeoutSeconds: int32(60), // 60 seconds
		CloudWatchOutputConfig: &types.CloudWatchOutputConfig{
			CloudWatchOutputEnabled: true,
		},
		Parameters: map[string][]string{
			"commands": {command},
		},
	}

	output, err := client.SendCommand(ctx, input)
	if err != nil {
		return nil, err
	}

	return &Command{
		client:     client,
		instanceID: instanceID,
		output:     output,
		timeout:    cfg.timeout,
	}, nil
}

func (cmd *Command) Result() (string, error) {
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
		default:
			return "", fmt.Errorf(*output.StandardErrorContent)
		}
	}
}
