package podman

import (
	"context"
	"io"
	"io/ioutil"
	"strings"

	"github.com/virtual-kubelet/podman/pkg/converter"

	//"github.com/davecgh/go-spew/spew"

	"github.com/virtual-kubelet/virtual-kubelet/log"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	v1 "k8s.io/api/core/v1"
)

// GetPod returns a pod by name that is stored in memory.
// TODO impelment pod status fields in the struct we return for the data
func (p *PodmanV0Provider) GetPod(ctx context.Context, namespace, name string) (pod *v1.Pod, err error) {
	log.G(ctx).Infof("receive GetPod %s", namespace, name)
	podName, err := converter.BuildKeyFromNames(namespace, name)
	if err != nil {
		return nil, err
	}
	return p.c.GetByName(ctx, podName)
}

// GetContainerLogs retrieves the logs of a container by name from the provider.
// TODO: To implement
func (p *PodmanV0Provider) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts api.ContainerLogOpts) (io.ReadCloser, error) {
	log.G(ctx).Infof("receive GetContainerLogs %q", podName)
	return ioutil.NopCloser(strings.NewReader("")), nil
}

// RunInContainer executes a command in a container in the pod, copying data
// between in/out/err and the container's stdin/stdout/stderr.
func (p *PodmanV0Provider) RunInContainer(ctx context.Context, namespace, name, container string, cmd []string, attach api.AttachIO) error {
	log.G(context.TODO()).Infof("receive ExecInContainer %q", container)
	return nil
}

// GetPodStatus returns the status of a pod by name that is "running".
// returns nil if a pod by that name is not found.
func (p *PodmanV0Provider) GetPodStatus(ctx context.Context, namespace, name string) (*v1.PodStatus, error) {
	log.G(ctx).Infof("receive GetPodStatus %q", name)

	pod, err := p.GetPod(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	return &pod.Status, nil
}

// GetPods returns a list of all pods known to be "running".
func (p *PodmanV0Provider) GetPods(ctx context.Context) ([]*v1.Pod, error) {
	log.G(ctx).Info("receive GetPods")
	list, err := p.c.List(ctx)
	if err != nil {
		return nil, err
	}
	result := []*v1.Pod{}
	for _, i := range list.Items {
		result = append(result, &i)
	}
	return result, nil
}
