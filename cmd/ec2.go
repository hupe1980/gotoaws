package cmd

import (
	"github.com/spf13/cobra"
)

func newEC2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ec2",
		Short:        "Connect to ec2",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newRunCmd(),
		newFwdCmd(),
		newSCPCmd(),
		newSSHCmd(),
		newSessionCmd(),
	)

	return cmd
}
