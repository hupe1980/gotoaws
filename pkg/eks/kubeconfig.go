package eks

import (
	"io/ioutil"
	"os"
	"path/filepath"

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

func (k *Kubeconfig) Update(alias string, cluster *Cluster, exec *clientcmdapi.ExecConfig) {
	if alias == "" {
		alias = cluster.ARN
	}

	k.config.Clusters[alias] = &clientcmdapi.Cluster{
		Server:                   cluster.Endpoint,
		CertificateAuthorityData: cluster.CAData,
	}

	k.config.Contexts[alias] = &clientcmdapi.Context{
		Cluster:  cluster.ARN,
		AuthInfo: cluster.ARN,
	}

	k.config.AuthInfos[alias] = &clientcmdapi.AuthInfo{
		Exec: exec,
	}

	k.config.CurrentContext = alias
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
