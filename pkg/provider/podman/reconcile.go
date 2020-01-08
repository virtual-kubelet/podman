package podman

import (
	"context"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/log"
)

func (p *PodmanV0Provider) reconcile() error {
	for {
		ctx := context.Background()
		time.Sleep(10 * time.Second)
		log.G(ctx).Infof("reconcile all pods status")
		pods := p.resourceManager.GetPods()

		if pods != nil {
			for _, pod := range pods {
				updatePod := pod.DeepCopy()
				currentPod, err := p.c.Get(ctx, updatePod)
				if err != nil {
					log.G(ctx).Debugf("error while reconcile pod %s/%s", pod.Namespace, pod.Name)
					continue
				}
				if updatePod != nil {
					updatePod.Status = currentPod.Status
					p.notifier(updatePod)
				}
			}
		}
	}
}
