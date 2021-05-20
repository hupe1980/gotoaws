package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

var (
	sshPort     string
	sshUser     string
	sshIdentity string
	sshCmd      string

	sshCommand = &cobra.Command{
		Use:    "ssh [identifier]",
		Short:  "",
		Long:   "",
		PreRun: preRun,
		Run: func(cmd *cobra.Command, args []string) {
			docName := "AWS-StartSSHSession"
			input := &ssm.StartSessionInput{
				DocumentName: &docName,
				Parameters:   map[string][]string{"portNumber": {sshPort}},
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
			for _, sep := range strings.Split(internal.SSHCommand(sshUser, instanceID, sshIdentity, sshCmd), " ") {
				if sep != "" {
					sshArgs = append(sshArgs, sep)
				}
			}
			fmt.Println(sshArgs)
			if err := internal.RunSubprocess("ssh", sshArgs...); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	sshCommand.Flags().StringVarP(&sshPort, "port", "p", "22", "SSH port (optional)")
	sshCommand.Flags().StringVarP(&sshUser, "user", "l", "ec2-user", "SSH user to us (optional)")
	sshCommand.Flags().StringVarP(&sshCmd, "cmd", "c", "", "(optional)")
	sshCommand.Flags().StringVarP(&sshIdentity, "identity", "i", "", "(required)")
	sshCommand.MarkFlagRequired("identity")
	rootCmd.AddCommand(sshCommand)
}
