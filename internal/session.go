package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Session struct {
	ID      string
	client  *ssm.Client
	output  *ssm.StartSessionOutput
	input   *ssm.StartSessionInput
	profile string
	plugin  string
	region  string
}

func NewSession(cfg *Config, input *ssm.StartSessionInput) (*Session, error) {
	client := ssm.NewFromConfig(cfg.awsCfg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	output, err := client.StartSession(ctx, input)
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:      *output.SessionId,
		client:  client,
		output:  output,
		input:   input,
		profile: cfg.profile,
		plugin:  cfg.plugin,
		region:  cfg.awsCfg.Region,
	}, nil
}

func (sess *Session) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := sess.client.TerminateSession(ctx, &ssm.TerminateSessionInput{SessionId: &sess.ID})
	if err != nil {
		return err
	}
	return nil
}

func (sess *Session) RunPlugin() error {
	outputJson, err := json.Marshal(sess.output)
	if err != nil {
		return err
	}
	inputJson, err := json.Marshal(sess.input)
	if err != nil {
		return err
	}
	return RunSubprocess(sess.plugin, string(outputJson), sess.region, "StartSession", sess.profile, string(inputJson))
}

func (sess *Session) ProxyCommand() (string, error) {
	outputJson, err := json.Marshal(sess.output)
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
