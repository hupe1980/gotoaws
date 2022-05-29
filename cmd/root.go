package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hupe1980/gotoaws/internal"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func Execute(version string) {
	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, promptui.IconBad, err)
		os.Exit(1)
	}
}

type rootOptions struct {
	profile string
	region  string
	timeout time.Duration
	silent  bool
}

func newRootCmd(version string) *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:     "gotoaws",
		Version: version,
		Short:   "Connect to your EC2 instance or ECS container without the need to open inbound ports, maintain bastion hosts, or manage SSH keys",
		Long: `gotoaws is an interactive CLI tool 
that you can use to connect to your AWS resources (EC2, ECS container) 
using the AWS Systems Manager Session Manager. 
It provides secure and auditable resource management without the need to open inbound ports, 
maintain bastion hosts, or manage SSH keys.`,
		SilenceErrors: true,
	}
	cmd.PersistentFlags().StringVar(&opts.profile, "profile", "", "AWS profile (optional)")
	cmd.PersistentFlags().StringVar(&opts.region, "region", "", "AWS region (optional)")
	cmd.PersistentFlags().DurationVar(&opts.timeout, "timeout", time.Second*15, "timeout for network requests")
	cmd.PersistentFlags().BoolVar(&opts.silent, "silent", false, "run gotoaws without printing logs")
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

	silent, err := cmd.Root().PersistentFlags().GetBool("silent")
	if err != nil {
		return nil, err
	}

	if !silent {
		fmt.Fprintf(os.Stdout, "%s Account: %s (%s)\n", promptui.IconGood, cfg.Account, cfg.Region)
	}

	return cfg, nil
}

func findInstance(cfg *internal.Config, identifier string) (*internal.Instance, error) {
	finder := internal.NewInstanceFinder(cfg)
	if identifier != "" {
		instances, err := finder.FindByIdentifier(identifier)
		if err != nil {
			return nil, err
		}

		if len(instances) > 1 {
			return chooseInstance(instances)
		}

		return &instances[0], nil
	}

	instances, err := finder.Find()
	if err != nil {
		return nil, err
	}

	return chooseInstance(instances)
}

func chooseInstance(instances []internal.Instance) (*internal.Instance, error) {
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
		return nil, err
	}

	return &instances[i], nil
}

func findContainer(cfg *internal.Config, cluster string, task string, cname string) (string, string, error) {
	finder := internal.NewContainerFinder(cfg)
	if task != "" {
		containers, err := finder.FindByIdentifier(cluster, task, cname)
		if err != nil {
			return "", "", err
		}

		if len(containers) > 1 {
			return chooseContainer(containers)
		}

		return containers[0].Task, containers[0].Name, nil
	}

	containers, err := finder.Find(cluster)
	if err != nil {
		return "", "", err
	}

	return chooseContainer(containers)
}

func chooseContainer(containers []internal.Container) (string, string, error) {
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf(`%s {{ .Name | cyan | bold }} ({{ .Task }})`, promptui.IconSelect),
		Inactive: `   {{ .Name | cyan }} ({{ .ID }})`,
		Selected: fmt.Sprintf(`%s {{ "Container" | bold }}: {{ .Name | cyan }} ({{ .Task }})`, promptui.IconGood),
	}

	searcher := func(input string, index int) bool {
		container := containers[index]
		name := strings.Replace(strings.ToLower(container.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose a container",
		Items:     containers,
		Templates: templates,
		Size:      15,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", "", err
	}

	return containers[i].Task, containers[i].Name, nil
}
