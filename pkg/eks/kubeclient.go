package eks

import (
	"context"
	"net/http"
	"os"

	"github.com/hupe1980/gotoaws/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/remotecommand"
)

type Kubeclient struct {
	clientset *kubernetes.Clientset
	restCfg   *rest.Config
}

func NewKubeclient(cfg *config.Config, cluster *Cluster, role string) (*Kubeclient, error) {
	execConfig := NewExecConfig(cfg, cluster.Name, role)

	execConfig.InteractiveMode = api.NeverExecInteractiveMode

	config := &rest.Config{
		Host: cluster.Endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cluster.CAData,
		},
		ExecProvider: execConfig,
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Kubeclient{
		clientset: clientset,
		restCfg:   config,
	}, nil
}

func (k *Kubeclient) GetPod(namespace, podName string) (*v1.Pod, error) {
	podClient := k.clientset.CoreV1().Pods(namespace)

	return podClient.Get(context.TODO(), podName, metav1.GetOptions{})
}

func (k *Kubeclient) ListPods(namespace, labelSelector string) (*v1.PodList, error) {
	podClient := k.clientset.CoreV1().Pods(namespace)

	return podClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: "status.phase=Running",
	})
}

type ExecInput struct {
	Namespace string
	PodName   string
	Container string
	Command   []string
}

func (k *Kubeclient) Exec(input *ExecInput) error {
	req := k.clientset.CoreV1().RESTClient().
		Post().
		Namespace(input.Namespace).
		Resource("pods").
		Name(input.PodName).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command:   input.Command,
			Container: input.Container,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k.restCfg, http.MethodPost, req.URL())
	if err != nil {
		return err
	}

	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             os.Stdin,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		Tty:               true,
		TerminalSizeQueue: nil,
	})
}
