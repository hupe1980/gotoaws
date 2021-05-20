package internal

import (
	"context"
	"os"

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

	return &Config{
		profile: profile,
		plugin:  "session-manager-plugin",
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
