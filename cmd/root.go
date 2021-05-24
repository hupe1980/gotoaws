package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hupe1980/awsconnect/internal"
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
	timeout time.Duration
}

func newRootCmd(version string) *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:     "awsconnect",
		Version: version,
		Short:   "awsconnect is an interactive CLI tool that you can use to connect to your AWS resources (EC2, ECS container)",
		Long: `awsconnect is an interactive CLI tool 
that you can use to connect to your AWS resources (EC2, ECS container) 
using the AWS Systems Manager Session Manager. 
It provides secure and auditable resource management without the need to open inbound ports, 
maintain bastion hosts, or manage SSH keys.`,
		SilenceErrors: true,
	}
	cmd.PersistentFlags().StringVar(&opts.profile, "profile", "default", "AWS profile (optional)")
	cmd.PersistentFlags().StringVar(&opts.region, "region", "", "AWS region (optional)")
	cmd.PersistentFlags().DurationVar(&opts.timeout, "timeout", time.Second*15, "timeout for network requests")
	cmd.AddCommand(
		newEC2Cmd(),
		newECSCmd(),
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

	timeout, err := cmd.Root().PersistentFlags().GetDuration("timeout")
	if err != nil {
		return nil, err
	}

	cfg, err := internal.NewConfig(profile, region, timeout)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Using profile %s (Region %s)\n", cfg.Profile, cfg.Region)

	return cfg, nil
}

func findInstance(cfg *internal.Config, identifier string) (string, error) {
	if identifier != "" {
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
