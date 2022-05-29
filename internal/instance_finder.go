package internal

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type Instance struct {
	Name     string
	ID       string
	Platform string
}

type InstanceFinder interface {
	Find() ([]Instance, error)
	FindByIdentifier(identifier string) ([]Instance, error)
}

type instanceFinder struct {
	timeout time.Duration
	ec2     ec2.DescribeInstancesAPIClient
	ssm     ssm.DescribeInstanceInformationAPIClient
}

func NewInstanceFinder(cfg *Config) InstanceFinder {
	return &instanceFinder{
		timeout: cfg.timeout,
		ec2:     ec2.NewFromConfig(cfg.awsCfg),
		ssm:     ssm.NewFromConfig(cfg.awsCfg),
	}
}

func (f *instanceFinder) Find() ([]Instance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	ec2Instances, managedInstances, err := f.findSSMManagedInstances()
	if err != nil {
		return nil, err
	}

	instanceIDs := []string{}
	for _, i := range ec2Instances {
		instanceIDs = append(instanceIDs, *i.InstanceId)
	}

	instanceIDFilter := types.Filter{
		Name:   aws.String("instance-id"),
		Values: instanceIDs,
	}
	instanceRunningFilter := types.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []string{"running"},
	}
	input := &ec2.DescribeInstancesInput{
		Filters:    []types.Filter{instanceRunningFilter, instanceIDFilter},
		MaxResults: aws.Int32(100),
	}

	p := ec2.NewDescribeInstancesPaginator(f.ec2, input)

	instances := []Instance{}

	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, r := range page.Reservations {
			for _, inst := range r.Instances {
				name := ""

				for _, tag := range inst.Tags {
					if *tag.Key == "Name" {
						name = *tag.Value
						break
					}
				}

				instances = append(instances, Instance{Name: name, ID: *inst.InstanceId, Platform: platform(inst)})
			}
		}
	}

	for _, mi := range managedInstances {
		instances = append(instances, Instance{Name: *mi.Name, ID: *mi.InstanceId, Platform: string(mi.PlatformType)})
	}

	return instances, nil
}

func (f *instanceFinder) FindByIdentifier(identifier string) ([]Instance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	input := &ec2.DescribeInstancesInput{
		Filters:    []types.Filter{parseIdentifier(identifier)},
		MaxResults: aws.Int32(100),
	}

	p := ec2.NewDescribeInstancesPaginator(f.ec2, input)

	var instances []Instance

	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, r := range page.Reservations {
			for _, inst := range r.Instances {
				name := ""

				for _, tag := range inst.Tags {
					if *tag.Key == "Name" {
						name = *tag.Value
						break
					}
				}

				instances = append(instances, Instance{Name: name, ID: *inst.InstanceId, Platform: platform(inst)})
			}
		}
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("no ssm managed instances found")
	}

	return instances, nil
}

func (f *instanceFinder) findSSMManagedInstances() ([]ssmTypes.InstanceInformation, []ssmTypes.InstanceInformation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	onlineFilter := ssmTypes.InstanceInformationStringFilter{
		Key:    aws.String("PingStatus"),
		Values: []string{"Online"},
	}
	input := &ssm.DescribeInstanceInformationInput{
		Filters:    []ssmTypes.InstanceInformationStringFilter{onlineFilter},
		MaxResults: 50,
	}

	p := ssm.NewDescribeInstanceInformationPaginator(f.ssm, input)

	var ec2Instances []ssmTypes.InstanceInformation

	var managedInstances []ssmTypes.InstanceInformation

	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, nil, err
		}

		for _, i := range page.InstanceInformationList {
			if i.ResourceType != ssmTypes.ResourceTypeEc2Instance {
				managedInstances = append(managedInstances, i)
			} else {
				ec2Instances = append(ec2Instances, i)
			}
		}
	}

	if len(ec2Instances) == 0 && len(managedInstances) == 0 {
		return nil, nil, fmt.Errorf("no ssm managed instances found")
	}

	return ec2Instances, managedInstances, nil
}

func platform(inst types.Instance) string {
	if inst.Platform != "" {
		return string(inst.Platform)
	}

	return "Linux" // TODO MacOS
}

func parseIdentifier(identifier string) types.Filter {
	if strings.HasPrefix(identifier, "i-") || strings.HasPrefix(identifier, "mi-") {
		return types.Filter{
			Name:   aws.String("instance-id"),
			Values: []string{identifier},
		}
	}

	ip := net.ParseIP(identifier)

	if ip != nil {
		_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
		_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
		_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")

		if private24BitBlock.Contains(ip) || private20BitBlock.Contains(ip) || private16BitBlock.Contains(ip) {
			return types.Filter{
				Name:   aws.String("private-ip-address"),
				Values: []string{identifier},
			}
		}

		return types.Filter{
			Name:   aws.String("ip-address"),
			Values: []string{identifier},
		}
	}

	if strings.HasSuffix(identifier, "compute.amazonaws.com") {
		return types.Filter{
			Name:   aws.String("dns-name"),
			Values: []string{identifier},
		}
	}

	if strings.HasSuffix(identifier, "compute.internal") {
		return types.Filter{
			Name:   aws.String("private-dns-name"),
			Values: []string{identifier},
		}
	}

	return types.Filter{
		Name:   aws.String("tag:Name"),
		Values: []string{identifier},
	}
}
