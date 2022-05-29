package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/spf13/cobra"
)

type scpOptions struct {
	target    string
	port      string
	user      string
	identity  string
	receiving bool
}

func newSCPCmd() *cobra.Command {
	opts := &scpOptions{}
	cmd := &cobra.Command{
		Use:           "scp [source(s)] [target]",
		Short:         "SCP over Session Manager",
		Example:       "gotoaws ec2 scp file.txt /opt/ -t myserver -i key.pem",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := newConfig(cmd)
			if err != nil {
				return err
			}

			instanceID, err := findInstance(cfg, opts.target)
			if err != nil {
				return err
			}

			docName := "AWS-StartSSHSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters:   map[string][]string{"portNumber": {opts.port}},
				Target:       &instanceID,
			}
			session, err := internal.NewEC2Session(cfg, input)
			if err != nil {
				return err
			}
			defer session.Close()

			pos := len(args) - 1

			mode := internal.SCPModeSending
			if opts.receiving {
				mode = internal.SCPModeReceiving
			}

			if err := session.RunSCP(&internal.RunSCPInput{
				User:       opts.user,
				InstanceID: instanceID,
				Identity:   opts.identity,
				Sources:    args[:pos],
				Target:     args[pos],
				Mode:       mode,
			}); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.receiving, "recv", "R", false, "receive files from target (optional)")
	cmd.Flags().StringVarP(&opts.target, "target", "t", "", "name|ID|IP|DNS of the instance (optional)")
	cmd.Flags().StringVarP(&opts.port, "port", "p", "22", "SSH port to us (optional)")
	cmd.Flags().StringVarP(&opts.user, "user", "l", "ec2-user", "SCP user to us (optional)")
	cmd.Flags().StringVarP(&opts.identity, "identity", "i", "", "file from which the identity (private key) for public key authentication is read (required)")

	if err := cmd.MarkFlagRequired("identity"); err != nil {
		panic(err)
	}

	return cmd
}
