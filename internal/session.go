package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type session struct {
	id         *string
	streamUrl  *string
	tokenValue *string
	client     *ssm.Client
	input      *ssm.StartSessionInput
	profile    string
	plugin     string
	region     string
}

type ECSSession interface {
	Close() error
	RunPlugin() error
}

type EC2Session interface {
	Close() error
	RunPlugin() error
	ProxyCommand() (string, error)
}

func NewECSSession(cfg *Config, input *ecs.ExecuteCommandInput) (ECSSession, error) {
	client := ecs.NewFromConfig(cfg.awsCfg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
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
	}, nil
}

func NewEC2Session(cfg *Config, input *ssm.StartSessionInput) (EC2Session, error) {
	client := ssm.NewFromConfig(cfg.awsCfg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
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
	}, nil
}

func (sess *session) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := sess.client.TerminateSession(ctx, &ssm.TerminateSessionInput{SessionId: sess.id})
	if err != nil {
		return err
	}
	return nil
}

func (sess *session) RunPlugin() error {
	outputJson, err := json.Marshal(map[string]*string{
		"SessionId":  sess.id,
		"StreamUrl":  sess.streamUrl,
		"TokenValue": sess.tokenValue,
	})
	inputJson, err := json.Marshal(sess.input)
	if err != nil {
		return err
	}
	return RunSubprocess(sess.plugin, string(outputJson), sess.region, "StartSession", sess.profile, string(inputJson))
}

func (sess *session) ProxyCommand() (string, error) {
	outputJson, err := json.Marshal(map[string]*string{
		"SessionId":  sess.id,
		"StreamUrl":  sess.streamUrl,
		"TokenValue": sess.tokenValue,
	})
	if err != nil {
		return "", err
	}
	inputJson, err := json.Marshal(sess.input)
	if err != nil {
		return "", err
	}
	pc := fmt.Sprintf("ProxyCommand=%s '%s' %s %s %s '%s'", sess.plugin, string(outputJson), sess.region, "StartSession", sess.profile, string(inputJson))
	return pc, nil
}
