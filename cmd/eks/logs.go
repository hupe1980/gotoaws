package eks

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/eks"
	"github.com/spf13/cobra"
)

type logsOptions struct {
	clusterName string
	role        string
	namespace   string
	pod         string
	container   string
}

func newLogsCmd() *cobra.Command {
	opts := &logsOptions{}
	cmd := &cobra.Command{
		Use:           "logs",
		Short:         "Print the logs for a container in a pod",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `gotoaws eks logs --cluster gotoaws --role cluster-admin --pod nginx
gotoaws eks logs --cluster gotoaws --role cluster-admin --pod nginx --container nginx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := internal.NewConfigFromFlags()
			if err != nil {
				return err
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

			ctx, cancel := context.WithCancel(context.Background())

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigs
				cancel()
			}()

			return client.PodLogs(ctx, &eks.PodLogsInput{
				Namespace: pod.Namespace,
				PodName:   pod.Name,
				Container: pod.Container,
				Writer: func(line string) {
					prefix := fmt.Sprintf("[pod/%s/%s] ", pod.Name, pod.Container)
					fmt.Fprintln(os.Stdout, prefix, line)
				},
			})
		},
	}

	cmd.Flags().StringVarP(&opts.clusterName, "cluster", "", "", "arn or name of the cluster")
	cmd.Flags().StringVarP(&opts.role, "role", "", "", "arn or name of the role")
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", "namespace of the pod (default for finder \"all namespaces\"")
	cmd.Flags().StringVarP(&opts.pod, "pod", "p", "", "name of the pod")
	cmd.Flags().StringVarP(&opts.container, "container", "c", "", "name of the container")

	return cmd
}
