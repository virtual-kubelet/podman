package provider

import (
	"context"

	"github.com/virtual-kubelet/virtual-kubelet/node"
	v1 "k8s.io/api/core/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

// Provider contains the methods required to implement a virtual-kubelet provider.
//
// Errors produced by these methods should implement an interface from
// github.com/virtual-kubelet/virtual-kubelet/errdefs package in order for the
// core logic to be able to understand the type of failure.
type Provider interface {
	node.PodLifecycleHandler

	// ConfigureNode enables a provider to configure the node object that
	// will be used for Kubernetes.
	ConfigureNode(context.Context, *v1.Node)
}

// PodMetricsProvider is an optional interface that providers can implement to expose pod stats
type PodMetricsProvider interface {
	GetStatsSummary(context.Context) (*stats.Summary, error)
}
