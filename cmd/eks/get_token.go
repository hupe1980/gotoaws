package eks

import (
	"fmt"
	"os"

	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/eks"
	"github.com/spf13/cobra"
)

type getTokenOptions struct {
	cluster   string
	role      string
	tokenOnly bool
}

func newGetTokenCmd() *cobra.Command {
	opts := &getTokenOptions{}
	cmd := &cobra.Command{
		Use:           "get-token",
		Short:         "Get a token for authentication with an Amazon EKS cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := internal.NewConfigFromFlags()
			if err != nil {
				return err
			}

			gen := eks.NewTokenGen(cfg)

			var t *eks.Token
			if opts.role != "" {
				t, err = gen.GetWithRole(opts.cluster, opts.role)
				if err != nil {
					return nil
				}
			} else {
				t, err = gen.Get(opts.cluster)
				if err != nil {
					return nil
				}
			}

			if opts.tokenOnly {
				fmt.Fprintln(os.Stdout, t.Token)
			} else {
				out := gen.FormatJSON(*t)

				fmt.Fprintln(os.Stdout, out)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.cluster, "cluster", "", "", "arn or name of the cluster")
	cmd.Flags().StringVarP(&opts.role, "role", "", "", "arn or name of the role")
	cmd.Flags().BoolVarP(&opts.tokenOnly, "token-only", "", false, "Return only the token for use with Bearer token based tools")

	return cmd
}
