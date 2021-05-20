package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

var (
	remotePortNumber string
	localPortNumber  string

	fwdCommand = &cobra.Command{
		Use:    "fwd [identifier]",
		Short:  "",
		Long:   "",
		PreRun: preRun,
		Run: func(cmd *cobra.Command, args []string) {
			docName := "AWS-StartPortForwardingSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters: map[string][]string{
					"portNumber":      {remotePortNumber},
					"localPortNumber": {localPortNumber},
				},
				Target: &instanceID,
			}
			session, err := internal.NewSession(cfg, input)
			if err != nil {
				panic(err)
			}
			defer session.Close()
			if err := session.RunPlugin(); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	fwdCommand.Flags().StringVarP(&remotePortNumber, "remote", "r", "", "Remote port to forward to (required)")
	fwdCommand.MarkFlagRequired("remote")
	fwdCommand.Flags().StringVarP(&localPortNumber, "local", "l", "", "Local port to use (required)")
	fwdCommand.MarkFlagRequired("local")
	rootCmd.AddCommand(fwdCommand)
}
