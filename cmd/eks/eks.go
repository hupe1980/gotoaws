package eks

import (
	"fmt"
	"strings"

	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/hupe1980/gotoaws/pkg/eks"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func NewEKSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "eks",
		Short:        "Connect to eks",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newUpdateKubeconfigCmd(),
		newGetTokenCmd(),
		newExecCmd(),
	)

	return cmd
}

func findCluster(cfg *config.Config, clusterName string) (*eks.Cluster, error) {
	finder := eks.NewClusterFinder(cfg)

	clusters, err := finder.Find(clusterName)
	if err != nil {
		return nil, err
	}

	if len(clusters) > 1 {
		return chooseCluster(clusters)
	}

	return &clusters[0], nil
}

// nolint: dupl // ok
func chooseCluster(clusters []eks.Cluster) (*eks.Cluster, error) {
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf(`%s {{ .Name | cyan | bold }} ({{ .Version }})`, promptui.IconSelect),
		Inactive: `   {{ .Name | cyan }} ({{ .Version }})`,
		Selected: fmt.Sprintf(`%s {{ "Cluster" | bold }}: {{ .Name | cyan }} ({{ .Version }})`, promptui.IconGood),
	}

	searcher := func(input string, index int) bool {
		cluster := clusters[index]
		name := strings.Replace(strings.ToLower(cluster.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose a cluster",
		Items:     clusters,
		Templates: templates,
		Size:      15,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &clusters[i], nil
}

func findPod(cfg *config.Config, cluster *eks.Cluster, role, namespace, podName, container string) (*eks.Pod, error) {
	finder, err := eks.NewPodFinder(cfg, cluster, role)
	if err != nil {
		return nil, err
	}

	if podName != "" {
		var pods []eks.Pod

		pods, err = finder.FindByIdentifier(namespace, podName, container)
		if err != nil {
			return nil, err
		}

		if len(pods) > 1 {
			return choosePod(pods)
		}

		return &pods[0], nil
	}

	pods, err := finder.Find(namespace, "")
	if err != nil {
		return nil, err
	}

	fmt.Println(pods)

	return choosePod(pods)
}

// nolint: dupl // ok
func choosePod(pods []eks.Pod) (*eks.Pod, error) {
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf(`%s {{ .Name | cyan | bold }} [{{ .Container }} ({{ .Namespace }})]`, promptui.IconSelect),
		Inactive: `   {{ .Name | cyan }} [{{ .Container }} ({{ .Namespace }})]`,
		Selected: fmt.Sprintf(`%s {{ "Pod" | bold }}: {{ .Name | cyan }} [{{ .Container }} ({{ .Namespace }})]`, promptui.IconGood),
	}

	searcher := func(input string, index int) bool {
		pod := pods[index]
		name := strings.Replace(strings.ToLower(pod.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose a pod",
		Items:     pods,
		Templates: templates,
		Size:      15,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &pods[i], nil
}
