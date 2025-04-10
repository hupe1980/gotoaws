package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type unsetOptions struct {
	key string
}

func newUnsetCmd() *cobra.Command {
	opts := &unsetOptions{}
	cmd := &cobra.Command{
		Use:           "unset",
		Short:         "Remove a config value",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := viper.AllSettings()

			delete(cfg, opts.key)

			return writeConfig(cfg)
		},
	}

	cmd.Flags().StringVarP(&opts.key, "key", "k", "", "key of the config value (required)")

	if err := cmd.MarkFlagRequired("key"); err != nil {
		panic(err)
	}

	return cmd
}
