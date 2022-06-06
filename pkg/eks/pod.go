package eks

import (
	"fmt"

	"github.com/hupe1980/gotoaws/pkg/config"
)

// ContainerPort represents a network port in a single container.
type ContainerPort struct {
	// Number of port to expose on the pod's IP address
	Port int32

	// Protocol for port. Must be UDP, TCP, or SCTP.
	Protocol string
}

// An object representing a pod.
type Pod struct {
	// The name of the pod.
	Name string

	// The namespace of the pod.
	Namespace string

	// The name of the container
	Container string

	// List of ports to expose from the container
	ContainerPorts []ContainerPort
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
			ports := []ContainerPort{}

			for _, p := range c.Ports {
				protocol := "TCP"
				if p.Protocol != "" {
					protocol = string(p.Protocol)
				}

				ports = append(ports, ContainerPort{Port: p.ContainerPort, Protocol: protocol})
			}

			pods = append(pods, Pod{
				Name:           p.Name,
				Namespace:      p.Namespace,
				Container:      c.Name,
				ContainerPorts: ports,
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
			ports := []ContainerPort{}

			for _, p := range c.Ports {
				protocol := "TCP"
				if p.Protocol != "" {
					protocol = string(p.Protocol)
				}

				ports = append(ports, ContainerPort{Port: p.ContainerPort, Protocol: protocol})
			}

			pods = append(pods, Pod{
				Name:           pod.Name,
				Namespace:      pod.Namespace,
				Container:      c.Name,
				ContainerPorts: ports,
			})
		}
	}

	if len(pods) == 0 {
		return nil, fmt.Errorf("no pods found")
	}

	return pods, nil
}
