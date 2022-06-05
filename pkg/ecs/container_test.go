package ecs

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	DescribeTasksOutput *aws_ecs.DescribeTasksOutput
	DescribeTasksError  error
	ListTasksOutput     *aws_ecs.ListTasksOutput
	ListTasksError      error
}

func (m *MockClient) DescribeTasks(ctx context.Context, params *aws_ecs.DescribeTasksInput, optFns ...func(*aws_ecs.Options)) (*aws_ecs.DescribeTasksOutput, error) {
	return m.DescribeTasksOutput, m.DescribeTasksError
}

func (m *MockClient) ListTasks(ctx context.Context, params *aws_ecs.ListTasksInput, optFns ...func(*aws_ecs.Options)) (*aws_ecs.ListTasksOutput, error) {
	return m.ListTasksOutput, m.ListTasksError
}

func TestContainerFinder(t *testing.T) {
	t.Run("FindByIdentifier", func(t *testing.T) {
		t.Run("no ssm with identifier", func(t *testing.T) {
			finder := &containerFinder{
				timeout: time.Second * 15,
				ecsClient: &MockClient{
					DescribeTasksOutput: &aws_ecs.DescribeTasksOutput{
						Tasks: []types.Task{},
					},
					DescribeTasksError: nil,
				},
			}
			instances, err := finder.FindByIdentifier("cluster", "task", "container")
			assert.Error(t, err)
			assert.Equal(t, "no ssm managed containers found", err.Error())
			assert.Nil(t, instances)
		})
	})

	t.Run("FindByIdentifier", func(t *testing.T) {
		t.Run("single container", func(t *testing.T) {
			finder := &containerFinder{
				timeout: time.Second * 15,
				ecsClient: &MockClient{
					DescribeTasksOutput: &aws_ecs.DescribeTasksOutput{
						Tasks: []types.Task{
							{
								EnableExecuteCommand: true,
								Containers: []types.Container{{
									TaskArn: aws.String("arn:aws:ecs:us-west-2:123456789012:task/MyCluster/1234567890123456789"),
									Name:    aws.String("container"),
								}},
							},
						},
					},
					DescribeTasksError: nil,
				},
			}
			expected := []Container{
				{Name: "container", Task: "1234567890123456789"},
			}
			instances, err := finder.FindByIdentifier("cluster", "task", "")
			assert.Nil(t, err)
			assert.Equal(t, expected, instances)
		})
	})

	t.Run("FindByIdentifier", func(t *testing.T) {
		t.Run("multi container", func(t *testing.T) {
			finder := &containerFinder{
				timeout: time.Second * 15,
				ecsClient: &MockClient{
					DescribeTasksOutput: &aws_ecs.DescribeTasksOutput{
						Tasks: []types.Task{
							{
								EnableExecuteCommand: true,
								Containers: []types.Container{
									{
										TaskArn: aws.String("arn:aws:ecs:us-west-2:123456789012:task/MyCluster/1234567890123456789"),
										Name:    aws.String("container1"),
									},
									{
										TaskArn: aws.String("arn:aws:ecs:us-west-2:123456789012:task/MyCluster/1234567890123456789"),
										Name:    aws.String("container2"),
									},
								},
							},
						},
					},
					DescribeTasksError: nil,
				},
			}
			expected := []Container{
				{Name: "container1", Task: "1234567890123456789"},
			}
			instances, err := finder.FindByIdentifier("cluster", "task", "container1")
			assert.Nil(t, err)
			assert.Equal(t, expected, instances)
		})
	})
}

func TestTaskID(t *testing.T) {
	arn := "arn:aws:ecs:us-west-2:123456789012:task/MyCluster/1234567890123456789"
	assert.Equal(t, "1234567890123456789", taskID(arn))
}
