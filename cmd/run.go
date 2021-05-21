package cmd

import (
	"fmt"

	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

type runOptions struct {
	cmd string
}

func newRunCmd() *cobra.Command {
	opts := &runOptions{}
	cmd := &cobra.Command{
		Use:           "run [name|ID|IP|DNS|_]",
		Short:         "Run commands",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := newConfig(cmd)
			if err != nil {
				return err
			}
			instanceID, err := findInstance(cfg, args)
			if err != nil {
				return err
			}
			command, err := internal.NewLinuxCommand(cfg, instanceID, opts.cmd)
			if err != nil {
				return err
			}
			res, err := command.Result()
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.cmd, "cmd", "c", "", "command to exceute (required)")
	cmd.MarkFlagRequired("cmd")

	return cmd
}
