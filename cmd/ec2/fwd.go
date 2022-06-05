package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/ec2"
	"github.com/spf13/cobra"
)

type fwdOptions struct {
	target           string
	remotePortNumber string
	remoteHost       string
	localPortNumber  string
}

func newFwdCmd() *cobra.Command {
	opts := &fwdOptions{}
	cmd := &cobra.Command{
		Use:   "fwd",
		Short: "Port forwarding",
		Example: `gotoaws fwd run -t myserver -l 8080 -r 8080
gotoaws fwd run -t myserver -l 5432 -r 5432 -H xxx.rds.amazonaws.com`,
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

			docName := "AWS-StartPortForwardingSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters: map[string][]string{
					"portNumber":      {opts.remotePortNumber},
					"localPortNumber": {opts.localPortNumber},
				},
				Target: &inst.ID,
			}

			if opts.remoteHost != "" {
				docName = "AWS-StartPortForwardingSessionToRemoteHost"
				input = &ssm.StartSessionInput{
					DocumentName: &docName,
					Parameters: map[string][]string{
						"portNumber":      {opts.remotePortNumber},
						"localPortNumber": {opts.localPortNumber},
						"host":            {opts.remoteHost},
					},
					Target: &inst.ID,
				}
			}

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
	cmd.Flags().StringVarP(&opts.remotePortNumber, "remote", "r", "", "remote port to forward to (required)")

	if err := cmd.MarkFlagRequired("remote"); err != nil {
		panic(err)
	}

	cmd.Flags().StringVarP(&opts.remoteHost, "host", "H", "", "remote host to forward to")

	cmd.Flags().StringVarP(&opts.localPortNumber, "local", "l", "", "local port to use (required)")

	if err := cmd.MarkFlagRequired("local"); err != nil {
		panic(err)
	}

	return cmd
}
