package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type Container struct {
	Task string
	Name string
}

func FindPossibleContainers(cfg *Config, cluster string) ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	client := ecs.NewFromConfig(cfg.awsCfg)

	output, err := client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster: &cluster,
	})
	if err != nil {
		return nil, err
	}
	if len(output.TaskArns) == 0 {
		return nil, fmt.Errorf("no ssm managed containers found")
	}

	tasks, err := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &cluster,
		Tasks:   output.TaskArns,
	})
	if err != nil {
		return nil, err
	}

	var containers []Container
	for _, t := range tasks.Tasks {
		if t.EnableExecuteCommand {
			for _, c := range t.Containers {
				containers = append(containers, Container{
					Task: taskID(*c.TaskArn),
					Name: *c.Name,
				})
			}
		}
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("no ssm managed containers found")
	}
	return containers, nil
}

func FindPossibleContainerByIdentifier(cfg *Config, cluster string, task string, container string) ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	client := ecs.NewFromConfig(cfg.awsCfg)
	tasks, err := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &cluster,
		Tasks:   []string{task},
	})
	if err != nil {
		return nil, err
	}
	var containers []Container
	for _, t := range tasks.Tasks {
		if t.EnableExecuteCommand {
			for _, c := range t.Containers {
				containers = append(containers, Container{
					Task: taskID(*c.TaskArn),
					Name: *c.Name,
				})
			}
		}
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("no ssm managed containers found")
	}
	return containers, nil
}

func taskID(a string) string {
	taskARN, _ := arn.Parse(a)
	res := strings.Split(taskARN.Resource, "/")
	return res[len(res)-1]
}
