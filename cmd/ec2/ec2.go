package ec2

import (
	"fmt"
	"strings"

	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/hupe1980/gotoaws/pkg/ec2"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func NewEC2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ec2",
		Short:        "Connect to ec2",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newRunCmd(),
		newFwdCmd(),
		newSCPCmd(),
		newSSHCmd(),
		newSessionCmd(),
	)

	return cmd
}

func findInstance(cfg *config.Config, identifier string) (*ec2.Instance, error) {
	finder := ec2.NewInstanceFinder(cfg)
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

func chooseInstance(instances []ec2.Instance) (*ec2.Instance, error) {
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
