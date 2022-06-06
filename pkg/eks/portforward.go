package eks

import (
	"fmt"
	"net/http"
	"os"

	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type PortForwardInput struct {
	// Namespace of the pod
	Namespace string

	// Name of the pod
	PodName string

	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort int32

	// ContainerPort is the target port for the pod
	ContainerPort int32

	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}

	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
}

func (k *Kubeclient) RunPortForward(input *PortForwardInput) error {
	req := k.clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(input.Namespace).
		Name(input.PodName).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(k.restCfg)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())

	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", input.LocalPort, input.ContainerPort)}, input.StopCh, input.ReadyCh, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}

	return fw.ForwardPorts()
}
