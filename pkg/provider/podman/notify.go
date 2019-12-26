package podman

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

func (p *PodmanV0Provider) NotifyPods(ctx context.Context, notifier func(pod *v1.Pod)) {
	p.notifier = notifier
}
