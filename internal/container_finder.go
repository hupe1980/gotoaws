package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type Container struct {
	Task string
	Name string
}

type ECSClient interface {
	DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
}

type ContainerFinder interface {
	Find(cluster string) ([]Container, error)
	FindByIdentifier(cluster string, task string, container string) ([]Container, error)
}

type containerFinder struct {
	timeout time.Duration
	ecs     ECSClient
}

func NewContainerFinder(cfg *Config) ContainerFinder {
	return &containerFinder{
		timeout: cfg.timeout,
		ecs:     ecs.NewFromConfig(cfg.awsCfg),
	}
}

func (f *containerFinder) Find(cluster string) ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	output, err := f.ecs.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster: &cluster,
	})
	if err != nil {
		return nil, err
	}
	if len(output.TaskArns) == 0 {
		return nil, fmt.Errorf("no ssm managed containers found")
	}

	tasks, err := f.ecs.DescribeTasks(ctx, &ecs.DescribeTasksInput{
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

func (f *containerFinder) FindByIdentifier(cluster string, task string, container string) ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	tasks, err := f.ecs.DescribeTasks(ctx, &ecs.DescribeTasksInput{
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
