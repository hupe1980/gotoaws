package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hupe1980/ec2connect/internal"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"
)

var (
	cfg        *internal.Config
	profile    string
	region     string
	instanceID string

	// rootCmd represents the base command when called without any sub-commands
	rootCmd = &cobra.Command{
		Use:   "ec2connect",
		Short: "",
		Long:  `ec2connect is an interactive CLI tool that you can use to connect to your EC2 instances using the [AWS Systems Manager Session Manager.`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "", "default", "AWS profile (optional)")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "", "", "AWS region (optional)")

	rootCmd.Version = version
	rootCmd.InitDefaultVersionFlag()
}

func initConfig() {
	var err error
	cfg, err = internal.NewConfig(profile, region)
	if err != nil {
		panic(err)
	}
}

func findInstance(args []string) (string, error) {
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

func preRun(cmd *cobra.Command, args []string) {
	var err error
	instanceID, err = findInstance(args)
	if err != nil {
		panic(err)
	}
}

func chooseInstance(instances []internal.Instance) (string, error) {
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf("%s {{ .Name | cyan | bold }} ({{ .ID }})", promptui.IconSelect),
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
