package internal

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type Instance struct {
	Name string
	ID   string
}

func FindPossibleInstances(cfg *Config) ([]Instance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	ec2Instances, managedInstances, err := findSSMManagedInstances(cfg)
	if err != nil {
		return nil, err
	}

	var instanceIDs []string
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
		Filters: []types.Filter{instanceRunningFilter, instanceIDFilter},
	}

	client := ec2.NewFromConfig(cfg.awsCfg)
	output, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	var instances []Instance
	for _, r := range output.Reservations {
		for _, inst := range r.Instances {
			name := ""
			for _, tag := range inst.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
			instances = append(instances, Instance{name, *inst.InstanceId})
		}
	}
	for _, mi := range managedInstances {
		instances = append(instances, Instance{*mi.Name, *mi.InstanceId})
	}
	return instances, nil
}

func findSSMManagedInstances(cfg *Config) ([]ssmTypes.InstanceInformation, []ssmTypes.InstanceInformation, error) {
	client := ssm.NewFromConfig(cfg.awsCfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	onlineFilter := ssmTypes.InstanceInformationStringFilter{
		Key:    aws.String("PingStatus"),
		Values: []string{"Online"},
	}
	input := &ssm.DescribeInstanceInformationInput{
		Filters: []ssmTypes.InstanceInformationStringFilter{onlineFilter},
	}
	out, err := client.DescribeInstanceInformation(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	if len(out.InstanceInformationList) == 0 {
		return nil, nil, fmt.Errorf("no ssm managed instances found")
	}

	var ec2Instances []ssmTypes.InstanceInformation
	var managedInstances []ssmTypes.InstanceInformation
	for _, i := range out.InstanceInformationList {
		if i.ResourceType != ssmTypes.ResourceTypeEc2Instance {
			managedInstances = append(managedInstances, i)
		} else {
			ec2Instances = append(ec2Instances, i)
		}
	}
	return ec2Instances, managedInstances, nil
}

func FindInstanceByIdentifier(cfg *Config, identifier string) ([]Instance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{parseIdentifier(identifier)},
	}

	client := ec2.NewFromConfig(cfg.awsCfg)
	output, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	var instances []Instance
	for _, r := range output.Reservations {
		for _, inst := range r.Instances {
			name := ""
			for _, tag := range inst.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
			instances = append(instances, Instance{name, *inst.InstanceId})
		}
	}
	return instances, nil
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
