package eks

import (
	"fmt"

	"github.com/hupe1980/gotoaws/pkg/config"
)

// An object representing a pod.
type Pod struct {
	// The name of the pod.
	Name string

	// The namespace of the pod.
	Namespace string

	// The name of the container
	Container string
}

type PodFinder interface {
	Find(namespace, labelSelector string) ([]Pod, error)
	FindByIdentifier(namespace, podName, container string) ([]Pod, error)
}

type podFinder struct {
	kubeclient *Kubeclient
}

func NewPodFinder(cfg *config.Config, cluster *Cluster, role string) (PodFinder, error) {
	kubeclient, err := NewKubeclient(cfg, cluster, role)
	if err != nil {
		return nil, err
	}

	return &podFinder{
		kubeclient: kubeclient,
	}, nil
}

func (p *podFinder) Find(namespace, labelSelector string) ([]Pod, error) {
	podList, err := p.kubeclient.ListPods(namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	pods := []Pod{}

	for _, p := range podList.Items {
		for _, c := range p.Spec.Containers {
			pods = append(pods, Pod{
				Name:      p.Name,
				Namespace: p.Namespace,
				Container: c.Name,
			})
		}
	}

	if len(pods) == 0 {
		return nil, fmt.Errorf("no pods found")
	}

	return pods, nil
}

func (p *podFinder) FindByIdentifier(namespace, podName, container string) ([]Pod, error) {
	pod, err := p.kubeclient.GetPod(namespace, podName)
	if err != nil {
		return nil, err
	}

	pods := []Pod{}

	for _, c := range pod.Spec.Containers {
		if container == "" || container == c.Name {
			pods = append(pods, Pod{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Container: c.Name,
			})
		}
	}

	if len(pods) == 0 {
		return nil, fmt.Errorf("no pods found")
	}

	return pods, nil
}
