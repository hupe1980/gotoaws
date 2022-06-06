package config

import (
	"bytes"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "config",
		Short:        "Manage your local gotoaws CLI config file",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newGetCmd(),
		newSetCmd(),
		newUnsetCmd(),
	)

	return cmd
}

func writeConfig(cfg map[string]interface{}) error {
	b, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	if err = viper.ReadConfig(bytes.NewReader(b)); err != nil {
		return err
	}

	if err := viper.SafeWriteConfig(); err != nil {
		return viper.WriteConfig()
	}

	return nil
}
