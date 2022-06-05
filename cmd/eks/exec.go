package eks

import (
	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/eks"
	"github.com/spf13/cobra"
)

type execOptions struct {
	clusterName string
	role        string
	namespace   string
	pod         string
	container   string
}

func newExecCmd() *cobra.Command {
	opts := &execOptions{}
	cmd := &cobra.Command{
		Use:           "exec [flags] -- COMMAND [args...]",
		Short:         "Execute a command in a container",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `gotoaws eks exec --cluster gotoaws --role cluster-admin
gotoaws eks exec --cluster gotoaws --role cluster-admin -- /bin/sh
gotoaws eks exec --cluster gotoaws --role cluster-admin -- cat /etc/passwd
gotoaws eks exec --cluster gotoaws --role cluster-admin --namespace default --pod nginx -- date`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := internal.NewConfigFromFlags()
			if err != nil {
				return err
			}

			command := []string{"/bin/sh"}
			if i := cmd.ArgsLenAtDash(); i != -1 {
				command = args[i:]
			}

			cluster, err := findCluster(cfg, opts.clusterName)
			if err != nil {
				return err
			}

			client, err := eks.NewKubeclient(cfg, cluster, opts.role)
			if err != nil {
				return err
			}

			pod, err := findPod(cfg, cluster, opts.role, opts.namespace, opts.pod, opts.container)
			if err != nil {
				return err
			}

			return client.Exec(&eks.ExecInput{
				Namespace: pod.Namespace,
				PodName:   pod.Name,
				Container: pod.Container,
				Command:   command,
			})
		},
	}

	cmd.Flags().StringVarP(&opts.clusterName, "cluster", "", "", "arn or name of the cluster (required)")

	if err := cmd.MarkFlagRequired("cluster"); err != nil {
		panic(err)
	}

	cmd.Flags().StringVarP(&opts.role, "role", "r", "", "arn or name of the role")
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", "namespace of the pod (default \"all namespaces\"")
	cmd.Flags().StringVarP(&opts.pod, "pod", "p", "", "name of the pod")
	cmd.Flags().StringVarP(&opts.container, "container", "c", "", "name of the container")

	return cmd
}
