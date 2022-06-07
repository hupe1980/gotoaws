package eks

import (
	"context"

	"github.com/hupe1980/gotoaws/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Kubeclient struct {
	clientset *kubernetes.Clientset
	restCfg   *rest.Config
}

func NewKubeclient(cfg *config.Config, cluster *Cluster, role string) (*Kubeclient, error) {
	token, err := getToken(cfg, cluster.Name, role)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{
		Host:        cluster.Endpoint,
		BearerToken: token.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cluster.CAData,
		},
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
