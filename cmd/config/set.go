package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type setOptions struct {
	values map[string]string
}

func newSetCmd() *cobra.Command {
	opts := &setOptions{}
	cmd := &cobra.Command{
		Use:           "set",
		Short:         "Create a new config value",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(opts.values)
			cfg := viper.AllSettings()

			for k, v := range opts.values {
				cfg[k] = v
			}

			return writeConfig(cfg)
		},
	}

	cmd.Flags().StringToStringVarP(&opts.values, "key", "k", nil, "new config value (required)")

	if err := cmd.MarkFlagRequired("key"); err != nil {
		panic(err)
	}

	return cmd
}
