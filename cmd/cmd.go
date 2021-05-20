package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

var (
	linuxCmd string

	cmdCommand = &cobra.Command{
		Use:    "cmd [identifier]",
		Short:  "",
		Long:   "",
		PreRun: preRun,
		Run: func(cmd *cobra.Command, args []string) {
			command, err := internal.NewLinuxCommand(cfg, instanceID, linuxCmd)
			if err != nil {
				panic(err)
			}
			if res, err := command.Result(); err != nil {
				panic(err)
			} else {
				fmt.Println(color.GreenString(res))
			}
		},
	}
)

func init() {
	cmdCommand.Flags().StringVarP(&linuxCmd, "cmd", "c", "", "Command to exceute (required)")
	cmdCommand.MarkFlagRequired("cmd")
	rootCmd.AddCommand(cmdCommand)
}
