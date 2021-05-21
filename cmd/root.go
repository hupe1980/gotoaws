package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hupe1980/ec2connect/internal"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func Execute(version string) {
	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type rootOptions struct {
	profile string
	region  string
}

func newRootCmd(version string) *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:     "ec2connect",
		Version: version,
		Short:   "ec2connect is an interactive CLI tool that you can use to connect to your EC2 instances",
		Long: `ec2connect is an interactive CLI tool 
that you can use to connect to your EC2 instances 
using the AWS Systems Manager Session Manager.`,
		SilenceErrors: true,
	}
	cmd.PersistentFlags().StringVarP(&opts.profile, "profile", "", "default", "AWS profile (optional)")
	cmd.PersistentFlags().StringVarP(&opts.region, "region", "", "", "AWS region (optional)")
	cmd.AddCommand(
		newRunCmd(),
		newFwdCmd(),
		newSCPCmd(),
		newSSHCmd(),
		newSessionCmd(),
		newCompletionCmd(),
	)
	return cmd
}

func newConfig(cmd *cobra.Command) (*internal.Config, error) {
	profile, err := cmd.Root().PersistentFlags().GetString("profile")
	if err != nil {
		return nil, err
	}

	region, err := cmd.Root().PersistentFlags().GetString("region")
	if err != nil {
		return nil, err
	}

	return internal.NewConfig(profile, region)
}

func findInstance(cfg *internal.Config, args []string) (string, error) {
	if len(args) > 0 {
		identifier := args[0]
		instances, err := internal.FindInstanceByIdentifier(cfg, identifier)
		if err != nil {
			panic(err)
		}
		if len(instances) > 1 {
			return chooseInstance(instances)
		}
		return instances[0].ID, nil
	}
	instances, err := internal.FindPossibleInstances(cfg)
	if err != nil {
		return "", err
	}
	return chooseInstance(instances)
}

func chooseInstance(instances []internal.Instance) (string, error) {
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf(`%s {{ .Name | cyan | bold }} ({{ .ID }})`, promptui.IconSelect),
		Inactive: `   {{ .Name | cyan }} ({{ .ID }})`,
		Selected: fmt.Sprintf(`%s {{ "Instance" | bold }}: {{ .Name | cyan }} ({{ .ID }})`, promptui.IconGood),
	}

	searcher := func(input string, index int) bool {
		instance := instances[index]
		name := strings.Replace(strings.ToLower(instance.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose an instance",
		Items:     instances,
		Templates: templates,
		Size:      15,
		Searcher:  searcher,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return instances[i].ID, nil
}
