package cmd

import (
	"github.com/spf13/cobra"
)

func newECSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ecs",
		Short:        "Connect to ecs",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newExecCmd(),
	)

	return cmd
}
