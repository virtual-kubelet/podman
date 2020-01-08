package podman

import (
	"context"

	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	v1 "k8s.io/api/core/v1"
)

// CreatePod accepts a Pod definition and stores it in memory.
func (p *PodmanV0Provider) CreatePod(ctx context.Context, pod *v1.Pod) error {
	// if DS is disabled, fail eary
	if p.config.DaemonSetDisabled == "true" {
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "DaemonSet" {
				pod.Status.Phase = v1.PodFailed
				for i := range pod.Status.ContainerStatuses {
					pod.Status.ContainerStatuses[i].State.Terminated = &v1.ContainerStateTerminated{
						ExitCode: 1,
						Message:  "DaemonSetDisabled is disabled on this node",
					}
				}
				pod.Status.Message = "DaemonSetDisabled is disabled on this node"
				p.notifier(pod)
				return errdefs.InvalidInput("DaemonSetDisabled is disabled on this node")
			}
		}

	}

	log.G(ctx).Infof("receive CreatePod %q", pod.Name)
	err := p.c.Create(ctx, pod)
	if err != nil {
		return err
	}

	pod, err = p.c.Get(ctx, pod)
	if err != nil {
		return err
	}

	p.notifier(pod)
	return nil
}
