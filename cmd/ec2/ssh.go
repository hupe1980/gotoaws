package ec2

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/ec2"
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
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := internal.NewConfigFromFlags()
			if err != nil {
				return err
			}

			inst, err := findInstance(cfg, opts.target)
			if err != nil {
				return err
			}

			docName := "AWS-StartSSHSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters:   map[string][]string{"portNumber": {opts.port}},
				Target:       &inst.ID,
			}
			session, err := ec2.NewSession(cfg, input)
			if err != nil {
				return err
			}
			defer session.Close()

			if err := session.RunSSH(&ec2.RunSSHInput{
				User:                opts.user,
				InstanceID:          inst.ID,
				Identity:            opts.identity,
				LocalPortForwarding: opts.fwd,
				Command:             strings.Join(args, " "),
			}); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.target, "target", "t", "", "name|ID|IP|DNS of the instance")
	cmd.Flags().StringVarP(&opts.fwd, "lforward", "L", "", "local port forwarding")
	cmd.Flags().StringVarP(&opts.port, "port", "p", "22", "SSH port to us")
	cmd.Flags().StringVarP(&opts.user, "user", "l", "ec2-user", "SSH user to us")
	cmd.Flags().StringVarP(&opts.identity, "identity", "i", "", "file from which the identity (private key) for public key authentication is read (required)")

	if err := cmd.MarkFlagRequired("identity"); err != nil {
		panic(err)
	}

	return cmd
}
