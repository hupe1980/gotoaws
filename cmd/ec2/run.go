package ec2

import (
	"errors"
	"fmt"
	"os"

	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/ec2"
	"github.com/spf13/cobra"
)

type runOptions struct {
	target string
}

func newRunCmd() *cobra.Command {
	opts := &runOptions{}
	cmd := &cobra.Command{
		Use:   "run [flags] -- COMMAND [args...]",
		Short: "Run commands",
		Example: `gotoaws ec2 run -- date
gotoaws ec2 run -t myserver -- date`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := internal.NewConfigFromFlags()
			if err != nil {
				return err
			}

			inst, err := findInstance(cfg, opts.target)
			if err != nil {
				return err
			}

			command := []string{}
			if i := cmd.ArgsLenAtDash(); i != -1 {
				command = args[i:]
			}

			if len(command) == 0 {
				return errors.New("command is missing")
			}

			runner, err := ec2.NewCommandRunner(cfg, inst, command)
			if err != nil {
				return err
			}

			res, err := runner.Result()
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, res)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.target, "target", "t", "", "name|ID|IP|DNS of the instance")

	return cmd
}
