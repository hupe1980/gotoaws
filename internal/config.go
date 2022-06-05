package internal

import (
	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/spf13/viper"
)

func NewConfigFromFlags() (*config.Config, error) {
	profile := viper.GetString("profile")
	region := viper.GetString("region")
	timeout := viper.GetDuration("timeout")

	cfg, err := config.NewConfig(profile, region, timeout)
	if err != nil {
		return nil, err
	}

	PrintInfof("Account: %s (%s)", cfg.Account, cfg.Region)

	return cfg, nil
}
