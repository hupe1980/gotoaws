package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/spf13/cobra"
)

type sshOptions struct {
	target   string
	port     string
	user     string
	identity string
	fwd      string
}

func newSSHCmd() *cobra.Command {
	opts := &sshOptions{}
	cmd := &cobra.Command{
		Use:           "ssh [command]",
		Short:         "SSH over Session Manager",
		Example:       "gotoaws ssh -t myserver -i key.pem",
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

			if err := session.RunSSH(&internal.RunSSHInput{
				User:                opts.user,
				InstanceID:          instanceID,
				Identity:            opts.identity,
				LocalPortForwarding: opts.fwd,
				Command:             strings.Join(args, " "),
			}); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.target, "target", "t", "", "name|ID|IP|DNS of the instance (optional)")
	cmd.Flags().StringVarP(&opts.fwd, "lforward", "L", "", "local port forwarding (optional)")
	cmd.Flags().StringVarP(&opts.port, "port", "p", "22", "SSH port to us (optional)")
	cmd.Flags().StringVarP(&opts.user, "user", "l", "ec2-user", "SSH user to us (optional)")
	cmd.Flags().StringVarP(&opts.identity, "identity", "i", "", "file from which the identity (private key) for public key authentication is read (required)")

	if err := cmd.MarkFlagRequired("identity"); err != nil {
		panic(err)
	}

	return cmd
}
