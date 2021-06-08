package internal

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

type MockECSClient struct {
	DescribeTasksOutput *ecs.DescribeTasksOutput
	DescribeTasksError  error
	ListTasksOutput     *ecs.ListTasksOutput
	ListTasksError      error
}

func (m *MockECSClient) DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error) {
	return m.DescribeTasksOutput, m.DescribeTasksError
}

func (m *MockECSClient) ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error) {
	return m.ListTasksOutput, m.ListTasksError
}

func TestContainerFinder(t *testing.T) {
	t.Run("FindByIdentifier", func(t *testing.T) {
		t.Run("no ssm with identifier", func(t *testing.T) {
			finder := &containerFinder{
				timeout: time.Second * 15,
				ecs: &MockECSClient{
					DescribeTasksOutput: &ecs.DescribeTasksOutput{
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
				ecs: &MockECSClient{
					DescribeTasksOutput: &ecs.DescribeTasksOutput{
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
			instances, err := finder.FindByIdentifier("cluster", "task", "")
			assert.Nil(t, err)
			assert.Equal(t, "container", instances[0].Name)
			assert.Equal(t, "1234567890123456789", instances[0].Task)
		})
	})

	t.Run("FindByIdentifier", func(t *testing.T) {
		t.Run("multi container", func(t *testing.T) {
			finder := &containerFinder{
				timeout: time.Second * 15,
				ecs: &MockECSClient{
					DescribeTasksOutput: &ecs.DescribeTasksOutput{
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
			instances, err := finder.FindByIdentifier("cluster", "task", "container1")
			assert.Nil(t, err)
			assert.Equal(t, "container1", instances[0].Name)
			assert.Equal(t, "1234567890123456789", instances[0].Task)
		})
	})
}

func TestTaskID(t *testing.T) {
	arn := "arn:aws:ecs:us-west-2:123456789012:task/MyCluster/1234567890123456789"
	assert.Equal(t, "1234567890123456789", taskID(arn))
}
