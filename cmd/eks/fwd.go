package eks

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/eks"
	"github.com/spf13/cobra"
)

type fwdOptions struct {
	clusterName string
	role        string
	namespace   string
	pod         string
	container   string
	remotePort  int32
	localPort   int32
}

func newFwdCmd() *cobra.Command {
	opts := &fwdOptions{}
	cmd := &cobra.Command{
		Use:           "fwd",
		Short:         "Port forwarding",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: `gotoaws eks fwd --cluster gotoaws --role cluster-admin --pod nginx
gotoaws eks fwd --cluster gotoaws --role cluster-admin --pod nginx --local 8000 --remote 80`,
		RunE: func(_ *cobra.Command, _ []string) error {
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

			containerPort := opts.remotePort
			if containerPort == 0 {
				if len(pod.ContainerPorts) > 0 {
					containerPort = pod.ContainerPorts[0].Port
				} else {
					return errors.New("container port cannot be determined")
				}
			}

			stopCh := make(chan struct{}, 1)
			readyCh := make(chan struct{})

			go func() {
				<-readyCh
				internal.PrintInfo("Port forwarding is ready")
			}()

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigs
				if stopCh != nil {
					close(stopCh)
				}
			}()

			return client.RunPortForward(&eks.PortForwardInput{
				Namespace:     pod.Namespace,
				PodName:       pod.Name,
				LocalPort:     opts.localPort,
				ContainerPort: containerPort,
				StopCh:        stopCh,
				ReadyCh:       readyCh,
			})
		},
	}

	cmd.Flags().StringVarP(&opts.clusterName, "cluster", "", "", "arn or name of the cluster")
	cmd.Flags().StringVarP(&opts.role, "role", "", "", "arn or name of the role")
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", "namespace of the pod (default \"all namespaces\"")
	cmd.Flags().StringVarP(&opts.pod, "pod", "p", "", "name of the pod")
	cmd.Flags().Int32VarP(&opts.remotePort, "remote", "r", 0, "the container port")
	cmd.Flags().Int32VarP(&opts.localPort, "local", "l", 0, "the local port")

	return cmd
}
