package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "exec",
		Short:         "Exec into container",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := newConfig(cmd)
			if err != nil {
				return err
			}
			fmt.Println(cfg)

			return nil
		},
	}

	return cmd
}
