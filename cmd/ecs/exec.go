package ecs

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/hupe1980/gotoaws/internal"
	"github.com/hupe1980/gotoaws/pkg/ecs"
	"github.com/spf13/cobra"
)

type execOptions struct {
	cluster   string
	task      string
	container string
}

func newExecCmd() *cobra.Command {
	opts := &execOptions{}
	cmd := &cobra.Command{
		Use:           "exec [flags] -- COMMAND [args...]",
		Short:         "Execute a command in a container",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "gotoaws ecs exec --cluster demo-cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := internal.NewConfigFromFlags()
			if err != nil {
				return err
			}

			task, container, err := findContainer(cfg, opts.cluster, opts.task, opts.container)
			if err != nil {
				return err
			}

			command := []string{"/bin/sh"}
			if i := cmd.ArgsLenAtDash(); i != -1 {
				command = args[i:]
			}

			input := &aws_ecs.ExecuteCommandInput{
				Interactive: true,
				Command:     aws.String(strings.Join(command, " ")),
				Cluster:     &opts.cluster,
				Task:        &task,
				Container:   &container,
			}
			session, err := ecs.NewSession(cfg, input)
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

	cmd.Flags().StringVarP(&opts.cluster, "cluster", "", "default", "arn or name of the cluster")
	cmd.Flags().StringVarP(&opts.task, "task", "", "", "arn or id of the task")
	cmd.Flags().StringVarP(&opts.container, "container", "", "", "name of the container. A container name only needs to be specified for tasks containing multiple containers")

	return cmd
}
