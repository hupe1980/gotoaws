package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/awsconnect/internal"
	"github.com/spf13/cobra"
)

func newSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "session [name|ID|IP|DNS|_]",
		Short:         "Start a session",
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

			input := &ssm.StartSessionInput{Target: &instanceID}
			session, err := internal.NewSession(cfg, input)
			if err != nil {
				panic(err)
			}
			defer session.Close()

			if err := session.RunPlugin(); err != nil {
				panic(err)
			}
			return nil
		},
	}

	return cmd
}
