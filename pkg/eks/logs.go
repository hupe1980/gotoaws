package eks

import (
	"bufio"
	"context"

	v1 "k8s.io/api/core/v1"
)

type PodLogsInput struct {
	// Namespace of the pod
	Namespace string

	// Name of the pod
	PodName string

	// Name of the container
	Container string

	Writer func(line string)
}

func (k *Kubeclient) PodLogs(ctx context.Context, input *PodLogsInput) error {
	req := k.clientset.CoreV1().Pods(input.Namespace).GetLogs(input.PodName, &v1.PodLogOptions{
		Container: input.Container,
		Follow:    true,
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}

	defer stream.Close()

	reader := bufio.NewScanner(stream)
	for reader.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
			line := reader.Text()
			input.Writer(line)
		}
	}

	return nil
}
