package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Config struct {
	profile string
	plugin  string
	awsCfg  aws.Config
}

func NewConfig(profile string, region string) (*Config, error) {
	if profile == "" {
		profile = "default"
		if os.Getenv("AWS_PROFILE") != "" {
			profile = os.Getenv("AWS_PROFILE")
		}
	}

	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	if region != "" {
		awsCfg.Region = region
	}

	pluginPath, err := exec.LookPath("session-manager-plugin")
	if err != nil {
		return nil, fmt.Errorf(`session-manager-plugin not found

https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html
		`)
	}

	return &Config{
		profile: profile,
		plugin:  pluginPath,
		awsCfg:  awsCfg,
	}, nil
}

func LoadDefaultConfig(profile string) (aws.Config, error) {
	if profile == "" {
		profile = "default"
		if os.Getenv("AWS_PROFILE") != "" {
			profile = os.Getenv("AWS_PROFILE")
		}
	}

	return config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
}
