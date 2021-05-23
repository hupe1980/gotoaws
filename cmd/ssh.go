package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/awsconnect/internal"
	"github.com/spf13/cobra"
)

type sshOptions struct {
	port     string
	user     string
	identity string
	cmd      string
}

func newSSHCmd() *cobra.Command {
	opts := &sshOptions{}
	cmd := &cobra.Command{
		Use:           "ssh [name|ID|IP|DNS| ]",
		Short:         "SSH over Session Manager",
		Example:       "awsconnect ssh myserver -i key.pem",
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
				User:       &opts.user,
				InstanceID: &instanceID,
				Identity:   &opts.identity,
				Command:    &opts.cmd,
			}); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "22", "SSH port to us (optional)")
	cmd.Flags().StringVarP(&opts.user, "user", "l", "ec2-user", "SSH user to us (optional)")
	cmd.Flags().StringVarP(&opts.cmd, "cmd", "c", "", "command to exceute (optional)")
	cmd.Flags().StringVarP(&opts.identity, "identity", "i", "file from which the identity (private key) for public key authentication is read", " (required)")
	if err := cmd.MarkFlagRequired("identity"); err != nil {
		panic(err)
	}

	return cmd
}
