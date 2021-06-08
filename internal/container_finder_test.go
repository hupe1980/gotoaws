package internal

import (
	"context"
	"testing"
	"time"

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
			instances, err := finder.FindByIdentifier("cluster", "tssk", "container")
			assert.Error(t, err)
			assert.Equal(t, "no ssm managed containers found", err.Error())
			assert.Nil(t, instances)
		})
	})
}
