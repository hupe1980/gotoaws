package ssm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/internal"
)

type Session struct {
	ID         *string
	StreamURL  *string
	TokenValue *string
	SSMClient  *ssm.Client
	Input      *ssm.StartSessionInput
	Profile    string
	Plugin     string
	Region     string
	Timeout    time.Duration
}

func (sess *Session) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), sess.Timeout)
	defer cancel()

	_, err := sess.SSMClient.TerminateSession(ctx, &ssm.TerminateSessionInput{SessionId: sess.ID})
	if err != nil {
		return err
	}

	return nil
}

func (sess *Session) RunPlugin() error {
	sessJSON, err := json.Marshal(map[string]*string{
		"SessionId":  sess.ID,
		"StreamUrl":  sess.StreamURL,
		"TokenValue": sess.TokenValue,
	})
	if err != nil {
		return err
	}

	inputJSON, err := json.Marshal(sess.Input)
	if err != nil {
		return err
	}

	return internal.RunSubprocess(sess.Plugin, string(sessJSON), sess.Region, "StartSession", sess.Profile, string(inputJSON))
}

func (sess *Session) ProxyCommand() (string, error) {
	sessJSON, err := json.Marshal(map[string]*string{
		"SessionId":  sess.ID,
		"StreamUrl":  sess.StreamURL,
		"TokenValue": sess.TokenValue,
	})
	if err != nil {
		return "", err
	}

	inputJSON, err := json.Marshal(sess.Input)
	if err != nil {
		return "", err
	}

	pc := fmt.Sprintf("ProxyCommand=%s '%s' %s %s %s '%s'", sess.Plugin, string(sessJSON), sess.Region, "StartSession", sess.Profile, string(inputJSON))

	return pc, nil
}
