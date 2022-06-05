package eks

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hupe1980/gotoaws/pkg/config"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	kubeConfigEnvName = "KUBECONFIG"
)

type Kubeconfig struct {
	filename string
	config   *clientcmdapi.Config
}

func NewKubeconfig(filename string) (*Kubeconfig, error) {
	kubeConfigFilename := filename

	if kubeConfigFilename == "" {
		kubeConfigFilename = os.Getenv(kubeConfigEnvName)
		// Fallback to default kubeconfig file location if no env variable is set
		if kubeConfigFilename == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}

			kubeConfigFilename = filepath.Join(home, ".kube", "config")
		}
	}

	kubeconfigBytes, err := ioutil.ReadFile(kubeConfigFilename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	kubeconfig, err := clientcmd.Load(kubeconfigBytes)
	if err != nil {
		return nil, err
	}

	return &Kubeconfig{
		filename: kubeConfigFilename,
		config:   kubeconfig,
	}, nil
}

func (k *Kubeconfig) Update(cfg *config.Config, cluster *Cluster, role, alias string) error {
	if alias == "" {
		alias = cluster.ARN
	}

	caCert, err := base64.StdEncoding.DecodeString(cluster.CABase64)
	if err != nil {
		return err
	}

	k.config.Clusters[alias] = &clientcmdapi.Cluster{
		Server:                   cluster.Endpoint,
		CertificateAuthorityData: caCert,
	}

	k.config.Contexts[alias] = &clientcmdapi.Context{
		Cluster:  cluster.ARN,
		AuthInfo: cluster.ARN,
	}

	exec := NewExecConfig(cfg, cluster.Name, role)

	k.config.AuthInfos[alias] = &clientcmdapi.AuthInfo{
		Exec: exec,
	}

	k.config.CurrentContext = alias

	return nil
}

// WriteToDisk writes a KubeConfig object down to disk with mode 0600
func (k *Kubeconfig) WriteToDisk() error {
	err := clientcmd.WriteToFile(*k.config, k.filename)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubeconfig) Filename() string {
	return k.filename
}
