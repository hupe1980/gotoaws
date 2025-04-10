package ec2

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/stretchr/testify/assert"
)

type MockEC2Client struct {
	DescribeInstancesOutput *aws_ec2.DescribeInstancesOutput
	DescribeInstancesError  error
}

func (m *MockEC2Client) DescribeInstances(_ context.Context, _ *aws_ec2.DescribeInstancesInput, _ ...func(*aws_ec2.Options)) (*aws_ec2.DescribeInstancesOutput, error) {
	return m.DescribeInstancesOutput, m.DescribeInstancesError
}

type MockSSMClient struct {
	DescribeInstanceInformationOutput *ssm.DescribeInstanceInformationOutput
	DescribeInstanceInformationError  error
}

func (m *MockSSMClient) DescribeInstanceInformation(_ context.Context, _ *ssm.DescribeInstanceInformationInput, _ ...func(*ssm.Options)) (*ssm.DescribeInstanceInformationOutput, error) {
	return m.DescribeInstanceInformationOutput, m.DescribeInstanceInformationError
}

func TestInstanceFinder(t *testing.T) {
	t.Run("FindByIdentifier", func(t *testing.T) {
		t.Run("no ssm with identifier", func(t *testing.T) {
			finder := &instanceFinder{
				timeout: time.Second * 15,
				ec2Client: &MockEC2Client{
					DescribeInstancesOutput: &aws_ec2.DescribeInstancesOutput{
						Reservations: []types.Reservation{},
					},
					DescribeInstancesError: nil,
				},
			}
			instances, err := finder.FindByIdentifier("XYZ")
			assert.Error(t, err)
			assert.Equal(t, "no ssm managed instances found", err.Error())
			assert.Nil(t, instances)
		})
	})
}

func TestParseIdentifier(t *testing.T) {
	t.Run("instance-id", func(t *testing.T) {
		identifier := "i-08d0906a5bd77e96f"
		expected := types.Filter{
			Name:   aws.String("instance-id"),
			Values: []string{identifier},
		}
		actual := parseIdentifier(identifier)
		assert.Equal(t, expected, actual)
	})

	t.Run("name tag", func(t *testing.T) {
		identifier := "myserver"
		expected := types.Filter{
			Name:   aws.String("tag:Name"),
			Values: []string{identifier},
		}
		actual := parseIdentifier(identifier)
		assert.Equal(t, expected, actual)
	})

	t.Run("public ip", func(t *testing.T) {
		identifier := "80.1.2.3"
		expected := types.Filter{
			Name:   aws.String("ip-address"),
			Values: []string{identifier},
		}
		actual := parseIdentifier(identifier)
		assert.Equal(t, expected, actual)
	})

	t.Run("private ip", func(t *testing.T) {
		identifier := "10.1.0.12"
		expected := types.Filter{
			Name:   aws.String("private-ip-address"),
			Values: []string{identifier},
		}
		actual := parseIdentifier(identifier)
		assert.Equal(t, expected, actual)
	})

	t.Run("public dns", func(t *testing.T) {
		identifier := "ec2-52-58-189-12.eu-central-1.compute.amazonaws.com"
		expected := types.Filter{
			Name:   aws.String("dns-name"),
			Values: []string{identifier},
		}
		actual := parseIdentifier(identifier)
		assert.Equal(t, expected, actual)
	})

	t.Run("private dns", func(t *testing.T) {
		identifier := "ip-172-31-33-49.eu-central-1.compute.internal"
		expected := types.Filter{
			Name:   aws.String("private-dns-name"),
			Values: []string{identifier},
		}
		actual := parseIdentifier(identifier)
		assert.Equal(t, expected, actual)
	})
}
