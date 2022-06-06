package eks

import (
	"github.com/hupe1980/gotoaws/pkg/config"
	"github.com/hupe1980/gotoaws/pkg/iam"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	// At the time EKS no longer supports Kubernetes v1.21 (probably ~Dec 2023),
	// this can be safely changed to default to writing "v1"
	apiVersion = "client.authentication.k8s.io/v1beta1"
)

func NewExecConfig(cfg *config.Config, clusterName, role string) *api.ExecConfig {
	execConfig := &api.ExecConfig{
		APIVersion:      apiVersion,
		Command:         "gotoaws",
		Args:            []string{"--region", cfg.Region, "--silent", "eks", "get-token", "--cluster", clusterName},
		InteractiveMode: api.NeverExecInteractiveMode,
	}

	if role != "" {
		execConfig.Args = append(execConfig.Args, "--role", iam.RoleARN(cfg.Account, role))
	}

	if cfg.Profile != "" {
		execConfig.Env = []api.ExecEnvVar{
			{
				Name:  "AWS_PROFILE",
				Value: cfg.Profile,
			},
		}
	}

	return execConfig
}
