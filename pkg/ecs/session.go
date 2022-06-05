package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	aws_ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/hupe1980/gotoaws/pkg/ssm"
)

type Session interface {
	Close() error
	RunPlugin() error
}

func NewSession(cfg *config.Config, input *aws_ecs.ExecuteCommandInput) (Session, error) {
	client := aws_ecs.NewFromConfig(cfg.AWSConfig)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)

	defer cancel()

	output, err := client.ExecuteCommand(ctx, input)
	if err != nil {
		return nil, err
	}

	return &ssm.Session{
		ID:         output.Session.SessionId,
		StreamURL:  output.Session.StreamUrl,
		TokenValue: output.Session.TokenValue,
		SSMClient:  aws_ssm.NewFromConfig(cfg.AWSConfig),
		Input: &aws_ssm.StartSessionInput{
			Target: aws.String(fmt.Sprintf("ecs:%s_%s_%s", *input.Cluster, *input.Task, *input.Container)),
		},
		Profile: cfg.Profile,
		Plugin:  cfg.Plugin,
		Region:  cfg.AWSConfig.Region,
		Timeout: cfg.Timeout,
	}, nil
}
