package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

var (
	scpPort     string
	scpUser     string
	scpIdentity string
	scpSource   string
	scpTarget   string

	scpCommand = &cobra.Command{
		Use:    "scp [identifier]",
		Short:  "",
		Long:   "",
		PreRun: preRun,
		Run: func(cmd *cobra.Command, args []string) {
			docName := "AWS-StartSSHSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters:   map[string][]string{"portNumber": {scpPort}},
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
			scpArgs := []string{"-o", pc}
			for _, sep := range strings.Split(internal.SCPCommand(scpUser, instanceID, scpIdentity, scpSource, scpTarget), " ") {
				if sep != "" {
					scpArgs = append(scpArgs, sep)
				}
			}
			if err := internal.RunSubprocess("scp", scpArgs...); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	scpCommand.Flags().StringVarP(&scpPort, "port", "p", "22", "(optional)")
	scpCommand.Flags().StringVarP(&scpUser, "user", "l", "ec2-user", "SCP user to us (optional)")
	scpCommand.Flags().StringVarP(&scpSource, "source", "s", "", "(required)")
	scpCommand.MarkFlagRequired("source")
	scpCommand.Flags().StringVarP(&scpTarget, "target", "t", "", "(required)")
	scpCommand.MarkFlagRequired("target")
	scpCommand.Flags().StringVarP(&scpIdentity, "identity", "i", "", "(required)")
	scpCommand.MarkFlagRequired("identity")
	rootCmd.AddCommand(scpCommand)
}
