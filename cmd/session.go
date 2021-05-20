package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

var (
	sessionCommand = &cobra.Command{
		Use:    "session [identifier]",
		Short:  "",
		Long:   "",
		PreRun: preRun,
		Run: func(cmd *cobra.Command, args []string) {
			input := &ssm.StartSessionInput{Target: &instanceID}
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
	rootCmd.AddCommand(sessionCommand)
}
