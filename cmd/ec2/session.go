package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/ec2"
	"github.com/spf13/cobra"
)

type sessionOptions struct {
	target string
}

func newSessionCmd() *cobra.Command {
	opts := &sessionOptions{}
	cmd := &cobra.Command{
		Use:           "session",
		Short:         "Start a session",
		Example:       "gotoaws ec2 session -t myserver",
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

			input := &ssm.StartSessionInput{Target: &inst.ID}
			session, err := ec2.NewSession(cfg, input)
			if err != nil {
				return err
			}
			defer session.Close()

			if err := session.RunPlugin(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.target, "target", "t", "", "name|ID|IP|DNS of the instance")

	return cmd
}
