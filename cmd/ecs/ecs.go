package ecs

import (
	"fmt"
	"strings"

	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/hupe1980/gotoaws/pkg/ecs"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func NewECSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ecs",
		Short:        "Connect to ecs",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newExecCmd(),
	)

	return cmd
}

func findContainer(cfg *config.Config, cluster string, task string, cname string) (string, string, error) {
	finder := ecs.NewContainerFinder(cfg)
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

func chooseContainer(containers []ecs.Container) (string, string, error) {
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf(`%s {{ .Name | cyan | bold }} ({{ .Task }})`, promptui.IconSelect),
		Inactive: `   {{ .Name | cyan }} ({{ .Task }})`,
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
