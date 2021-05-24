package cmd

import (
	"fmt"

	"github.com/hupe1980/awsconnect/internal"
	"github.com/spf13/cobra"
)

type runOptions struct {
	target string
	cmd    string
}

func newRunCmd() *cobra.Command {
	opts := &runOptions{}
	cmd := &cobra.Command{
		Use:           "run",
		Short:         "Run commands",
		Example:       "awsconnect ec2 run -t myserver -c 'cat /etc/passwd'",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := newConfig(cmd)
			if err != nil {
				return err
			}
			instanceID, err := findInstance(cfg, opts.target)
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

	cmd.Flags().StringVarP(&opts.target, "target", "t", "", "name|ID|IP|DNS of the instance (optional)")
	cmd.Flags().StringVarP(&opts.cmd, "cmd", "c", "", "command to exceute (required)")
	if err := cmd.MarkFlagRequired("cmd"); err != nil {
		panic(err)
	}

	return cmd
}
