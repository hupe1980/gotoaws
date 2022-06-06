package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hupe1980/gotoaws/cmd/config"
	"github.com/hupe1980/gotoaws/cmd/ec2"
	"github.com/hupe1980/gotoaws/cmd/ecs"
	"github.com/hupe1980/gotoaws/cmd/eks"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute(version string) {
	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			internal.PrintError(ae.ErrorMessage())
		} else {
			internal.PrintError(err)
		}

		os.Exit(1)
	}
}

func newRootCmd(version string) *cobra.Command {
	var cfgFile string

	cobra.OnInitialize(initConfig(cfgFile))

	cmd := &cobra.Command{
		Use:     "gotoaws",
		Version: version,
		Short:   "Connect to your EC2 instance or ECS container without the need to open inbound ports, maintain bastion hosts, or manage SSH keys",
		Long: `gotoaws is an interactive CLI tool that you can use to connect to your AWS resources 
(EC2, ECS container) using the AWS Systems Manager Session Manager. 
It provides secure and auditable resource management without the need to open inbound 
ports, maintain bastion hosts, or manage SSH keys.`,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().String("profile", "", "AWS profile")
	cmd.PersistentFlags().String("region", "", "AWS region")
	cmd.PersistentFlags().Duration("timeout", time.Second*15, "timeout for network requests")
	cmd.PersistentFlags().Bool("silent", false, "run gotoaws without printing logs")
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default \"$HOME/.config/configstore/gotoaws.json\")")

	cmd.AddCommand(
		ec2.NewEC2Cmd(),
		ecs.NewECSCmd(),
		eks.NewEKSCmd(),
		config.NewConfigCmd(),
		newCompletionCmd(),
	)

	err := viper.BindPFlag("profile", cmd.PersistentFlags().Lookup("profile"))
	cobra.CheckErr(err)

	err = viper.BindPFlag("region", cmd.PersistentFlags().Lookup("region"))
	cobra.CheckErr(err)

	err = viper.BindPFlag("timeout", cmd.PersistentFlags().Lookup("timeout"))
	cobra.CheckErr(err)

	err = viper.BindPFlag("silent", cmd.PersistentFlags().Lookup("silent"))
	cobra.CheckErr(err)

	return cmd
}

func initConfig(cfgFile string) func() {
	return func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			viper.AddConfigPath(filepath.Join(home, ".config", "configstore"))
			viper.SetConfigType("json")
			viper.SetConfigName("gotoaws")
		}

		if err := viper.ReadInConfig(); err == nil {
			internal.PrintInfof("Using config file: %s", viper.ConfigFileUsed())
		}
	}
}
