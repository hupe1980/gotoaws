package config

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Config struct {
	// The Amazon Web Services account ID number of the account that owns or contains the calling entity
	Account string

	// The SharedConfigProfile that is used
	Profile string

	// The region to send requests to.
	Region  string
	Timeout time.Duration

	// A Config provides service configuration for aws service clients
	AWSConfig aws.Config

	// Path to the session-manager-plugin
	Plugin string
}

func NewConfig(profile string, region string, timeout time.Duration) (*Config, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
		config.WithAssumeRoleCredentialOptions(func(aro *stscreds.AssumeRoleOptions) {
			aro.TokenProvider = stscreds.StdinTokenProvider
		}),
	)
	if err != nil {
		return nil, err
	}

	pluginPath, err := exec.LookPath("session-manager-plugin")
	if err != nil {
		return nil, fmt.Errorf(`session-manager-plugin not found

https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html
		`)
	}

	client := sts.NewFromConfig(awsCfg)

	output, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return &Config{
		Account:   *output.Account,
		Profile:   profile,
		Region:    awsCfg.Region,
		Plugin:    pluginPath,
		Timeout:   timeout,
		AWSConfig: awsCfg,
	}, nil
}
