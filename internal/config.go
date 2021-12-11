package internal

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
	Account string
	Profile string
	Region  string
	plugin  string
	timeout time.Duration
	awsCfg  aws.Config
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
		Account: *output.Account,
		Profile: profile,
		Region:  awsCfg.Region,
		plugin:  pluginPath,
		timeout: timeout,
		awsCfg:  awsCfg,
	}, nil
}
