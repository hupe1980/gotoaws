package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

type scpOptions struct {
	port     string
	user     string
	identity string
	source   string
	target   string
}

func newSCPCmd() *cobra.Command {
	opts := &scpOptions{}
	cmd := &cobra.Command{
		Use:           "scp [name|ID|IP|DNS|_]",
		Short:         "Tunnel scp",
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
				return err
			}
			defer session.Close()

			pc, err := session.ProxyCommand()
			if err != nil {
				return err
			}
			scpArgs := []string{"-o", pc}
			for _, sep := range strings.Split(internal.SCPArgs(opts.user, instanceID, opts.identity, opts.source, opts.target), " ") {
				if sep != "" {
					scpArgs = append(scpArgs, sep)
				}
			}
			if err := internal.RunSubprocess("scp", scpArgs...); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "22", "SSH port to us (optional)")
	cmd.Flags().StringVarP(&opts.user, "user", "l", "ec2-user", "SCP user to us (optional)")
	cmd.Flags().StringVarP(&opts.source, "source", "s", "", "Source in the local host (required)")
	cmd.MarkFlagRequired("source")
	cmd.Flags().StringVarP(&opts.target, "target", "t", "", "Target in the remote host (required)")
	cmd.MarkFlagRequired("target")
	cmd.Flags().StringVarP(&opts.identity, "identity", "i", "", "(required)")
	cmd.MarkFlagRequired("identity")

	return cmd
}
