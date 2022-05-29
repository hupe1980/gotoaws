package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type Container struct {
	Task string
	Name string
}

type ECSClient interface {
	ecs.ListTasksAPIClient
	ecs.DescribeTasksAPIClient
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

	p := ecs.NewListTasksPaginator(f.ecs, &ecs.ListTasksInput{
		Cluster:    &cluster,
		MaxResults: aws.Int32(100),
	})

	var containers []Container

	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		tasks, err := f.ecs.DescribeTasks(ctx, &ecs.DescribeTasksInput{
			Cluster: &cluster,
			Tasks:   page.TaskArns,
		})
		if err != nil {
			return nil, err
		}

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
				if container == "" || container == *c.Name {
					containers = append(containers, Container{
						Task: taskID(*c.TaskArn),
						Name: *c.Name,
					})
				}
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
