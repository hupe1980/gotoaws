package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

type fwdOptions struct {
	remotePortNumber string
	localPortNumber  string
}

func newFwdCmd() *cobra.Command {
	opts := &fwdOptions{}
	cmd := &cobra.Command{
		Use:           "fwd [identifier]",
		Short:         "",
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

			docName := "AWS-StartPortForwardingSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters: map[string][]string{
					"portNumber":      {opts.remotePortNumber},
					"localPortNumber": {opts.localPortNumber},
				},
				Target: &instanceID,
			}
			session, err := internal.NewSession(cfg, input)
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

	cmd.Flags().StringVarP(&opts.remotePortNumber, "remote", "r", "", "Remote port to forward to (required)")
	cmd.MarkFlagRequired("remote")
	cmd.Flags().StringVarP(&opts.localPortNumber, "local", "l", "", "Local port to use (required)")
	cmd.MarkFlagRequired("local")

	return cmd
}
