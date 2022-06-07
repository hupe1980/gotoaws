package eks

import (
	"net/http"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

type ExecInput struct {
	// Namespace of the pod
	Namespace string

	// Name of the pod
	PodName string

	// Name of the container
	Container string

	// Command to run
	Command []string
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
