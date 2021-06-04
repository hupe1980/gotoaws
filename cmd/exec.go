package cmd

import (
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/spf13/cobra"
)

type execOptions struct {
	cmd       string
	cluster   string
	task      string
	container string
}

func newExecCmd() *cobra.Command {
	opts := &execOptions{}
	cmd := &cobra.Command{
		Use:           "exec",
		Short:         "Exec into container",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "gotoaws ecs exec --cluster demo-cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := newConfig(cmd)
			if err != nil {
				return err
			}

			task, container, err := findContainer(cfg, opts.cluster, opts.task, opts.container)
			if err != nil {
				return err
			}

			input := &ecs.ExecuteCommandInput{
				Interactive: true,
				Command:     &opts.cmd,
				Cluster:     &opts.cluster,
				Task:        &task,
				Container:   &container,
			}
			session, err := internal.NewECSSession(cfg, input)
			if err != nil {
				return err
			}
			defer session.Close()

			if err := session.RunPlugin(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.cluster, "cluster", "", "default", "arn or name of the cluster (optional)")
	cmd.Flags().StringVarP(&opts.task, "task", "", "", "arn or id of the task (optional)")
	cmd.Flags().StringVarP(&opts.container, "container", "", "", "name of the container. A container name only needs to be specified for tasks containing multiple containers. (optional)")
	cmd.Flags().StringVarP(&opts.cmd, "cmd", "c", "/bin/sh", "command to exceute (optional)")

	return cmd
}
