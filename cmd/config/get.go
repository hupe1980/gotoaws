package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getOptions struct {
	key string
}

func newGetCmd() *cobra.Command {
	opts := &getOptions{}
	cmd := &cobra.Command{
		Use:           "get",
		Short:         "Print a config value",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := viper.AllSettings()

			value, ok := cfg[opts.key]
			if !ok {
				return fmt.Errorf("unknown config key: %s", opts.key)
			}

			fmt.Fprintln(os.Stdout, value)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.key, "key", "k", "", "key of the config value (required)")

	if err := cmd.MarkFlagRequired("key"); err != nil {
		panic(err)
	}

	return cmd
}
