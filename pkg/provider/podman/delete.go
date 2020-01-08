package podman

import (
	"context"

	"github.com/virtual-kubelet/virtual-kubelet/log"
	v1 "k8s.io/api/core/v1"
)

// DeletePod deletes the specified pod out of memory.
func (p *PodmanV0Provider) DeletePod(ctx context.Context, pod *v1.Pod) (err error) {
	log.G(ctx).Infof("receive DeletePod %s", pod.Namespace, pod.Name)
	return p.c.Delete(ctx, pod)
}
