package eks

import (
	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/eks"
	"github.com/spf13/cobra"
)

type updateKubeconfigOptions struct {
	clusterName string
	role        string
	alias       string
}

func newUpdateKubeconfigCmd() *cobra.Command {
	opts := &updateKubeconfigOptions{}
	cmd := &cobra.Command{
		Use:           "update-kubeconfig",
		Short:         "Configures kubectl so that you can connect to an Amazon EKS cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := internal.NewConfigFromFlags()
			if err != nil {
				return err
			}

			cluster, err := findCluster(cfg, opts.clusterName)
			if err != nil {
				return err
			}

			kubeconfig, err := eks.NewKubeconfig("")
			if err != nil {
				return err
			}

			if opts.alias == "" {
				opts.alias = cluster.ARN
			}

			if err := kubeconfig.Update(cfg, cluster, opts.role, opts.alias); err != nil {
				return err
			}

			if err := kubeconfig.WriteToDisk(); err != nil {
				return err
			}

			internal.PrintInfof("Updated context %s in %s", opts.alias, kubeconfig.Filename())

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.clusterName, "cluster", "", "", "arn or name of the cluster")
	cmd.Flags().StringVarP(&opts.role, "role", "", "", "arn or name of the role")
	cmd.Flags().StringVarP(&opts.alias, "alias", "", "", "alias for the cluster context name (default \"arn of the cluster\"")

	return cmd
}
