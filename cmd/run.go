package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hupe1980/ec2connect/internal"
	"github.com/spf13/cobra"
)

type runOptions struct {
	cmd string
}

func newRunCmd() *cobra.Command {
	opts := &runOptions{}
	cmd := &cobra.Command{
		Use:           "run [identifier]",
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
			command, err := internal.NewLinuxCommand(cfg, instanceID, opts.cmd)
			if err != nil {
				return err
			}
			res, err := command.Result()
			if err != nil {
				return err
			}
			fmt.Println(color.GreenString(res))
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.cmd, "cmd", "c", "", "Command to exceute (required)")
	cmd.MarkFlagRequired("cmd")

	return cmd
}

// var (
// 	linuxCmd string

// 	cmdCommand = &cobra.Command{
// 		Use:    "cmd [identifier]",
// 		Short:  "",
// 		Long:   "",
// 		PreRun: preRun,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			command, err := internal.NewLinuxCommand(cfg, instanceID, linuxCmd)
// 			if err != nil {
// 				panic(err)
// 			}
// 			res, err := command.Result()
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Println(color.GreenString(res))
// 		},
// 	}
// )

// func init() {
// 	cmdCommand.Flags().StringVarP(&linuxCmd, "cmd", "c", "", "Command to exceute (required)")
// 	cmdCommand.MarkFlagRequired("cmd")
// 	rootCmd.AddCommand(cmdCommand)
// }
