package eks

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_eks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/hupe1980/gotoaws/pkg/config"
)

// An object representing an Amazon EKS cluster.
type Cluster struct {
	// The Amazon Resource Name (ARN) of the cluster.
	ARN string

	// The name of the cluster.
	Name string

	// The Kubernetes server version for the cluster.
	Version string

	// The endpoint for your Kubernetes API server.
	Endpoint string

	// CAData contains PEM-encoded certificate authority certificates
	CAData []byte
}

type Client interface {
	aws_eks.ListClustersAPIClient
	aws_eks.DescribeClusterAPIClient
}

type ClusterFinder interface {
	Find(name string) ([]Cluster, error)
}

type clusterFinder struct {
	timeout   time.Duration
	eksClient Client
}

func NewClusterFinder(cfg *config.Config) ClusterFinder {
	return &clusterFinder{
		timeout:   cfg.Timeout,
		eksClient: aws_eks.NewFromConfig(cfg.AWSConfig),
	}
}

func (f *clusterFinder) Find(name string) ([]Cluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	var names []string

	if name == "" {
		p := aws_eks.NewListClustersPaginator(f.eksClient, &aws_eks.ListClustersInput{
			MaxResults: aws.Int32(100),
		})

		for p.HasMorePages() {
			page, err := p.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			names = append(names, page.Clusters...)
		}
	} else {
		names = append(names, name)
	}

	clusters := []Cluster{}

	for _, n := range names {
		out, err := f.eksClient.DescribeCluster(ctx, &aws_eks.DescribeClusterInput{
			Name: aws.String(n),
		})

		if err != nil {
			return nil, err
		}

		if out.Cluster.Status == types.ClusterStatusActive || out.Cluster.Status == types.ClusterStatusUpdating {
			caData, err := base64.StdEncoding.DecodeString(*out.Cluster.CertificateAuthority.Data)
			if err != nil {
				return nil, err
			}

			clusters = append(clusters, Cluster{
				ARN:      *out.Cluster.Arn,
				Name:     *out.Cluster.Name,
				Version:  *out.Cluster.Version,
				Endpoint: *out.Cluster.Endpoint,
				CAData:   caData,
			})
		}
	}

	if len(clusters) == 0 {
		return nil, fmt.Errorf("no eks clusters found")
	}

	return clusters, nil
}
