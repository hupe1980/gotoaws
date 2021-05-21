package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/ec2connect/internal"
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
		Use:           "ssh [identifier]",
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

			docName := "AWS-StartSSHSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters:   map[string][]string{"portNumber": {opts.port}},
				Target:       &instanceID,
			}
			session, err := internal.NewSession(cfg, input)
			if err != nil {
				panic(err)
			}
			defer session.Close()

			pc, err := session.ProxyCommand()
			if err != nil {
				panic(err)
			}
			sshArgs := []string{"-o", pc}
			for _, sep := range strings.Split(internal.SSHArgs(opts.user, instanceID, opts.identity, opts.cmd), " ") {
				if sep != "" {
					sshArgs = append(sshArgs, sep)
				}
			}
			if err := internal.RunSubprocess("ssh", sshArgs...); err != nil {
				panic(err)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.port, "port", "p", "22", "SSH port (optional)")
	cmd.Flags().StringVarP(&opts.user, "user", "l", "ec2-user", "SSH user to us (optional)")
	cmd.Flags().StringVarP(&opts.cmd, "cmd", "c", "", "(optional)")
	cmd.Flags().StringVarP(&opts.identity, "identity", "i", "", "(required)")
	cmd.MarkFlagRequired("identity")

	return cmd
}
