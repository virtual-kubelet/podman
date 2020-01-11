package podman

import (
	"context"
	"time"

	"github.com/virtual-kubelet/podman/pkg/manager"
	"github.com/virtual-kubelet/podman/pkg/podman"
	v1 "k8s.io/api/core/v1"
)

const (
	// Provider configuration defaults.
	defaultCPUCapacity       = "5"
	defaultMemoryCapacity    = "2Gi"
	defaultPodCapacity       = "10"
	defaultSocket            = "unix:/run/podman/io.podman"
	defaultDaemonSetDisabled = "true"
)

// PodmanV0Provider implements the virtual-kubelet provider interface and stores pods in memory.
type PodmanV0Provider struct {
	nodeName           string
	operatingSystem    string
	config             PodmanConfig
	startTime          time.Time
	notifier           func(*v1.Pod)
	internalIP         string
	daemonEndpointPort int32
	c                  podman.Podman
	resourceManager    *manager.ResourceManager
}

// PodmanProvider is like PodmanV0Provider, but implements the PodNotifier interface
type PodmanProvider struct {
	*PodmanV0Provider
}

// PodmanConfig contains a podman virtual-kubelet's configurable parameters.
type PodmanConfig struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	Pods   string `json:"pods,omitempty"`

	Socket string `json:"socket,omitempty"`

	DaemonSetDisabled string `json:"daemonSetDisabled,omitempty"`
}

// NewPodmanV0ProviderPodmanConfig creates a new PodmanV0Provider. podman legacy provider does not implement the new asynchronous podnotifier interface
func NewPodmanV0ProviderPodmanConfig(config PodmanConfig, nodeName, operatingSystem string, resourceManager *manager.ResourceManager) (*PodmanV0Provider, error) {
	client, err := podman.New(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	provider := PodmanV0Provider{
		nodeName:        nodeName,
		operatingSystem: operatingSystem,
		config:          config,
		startTime:       time.Now(),
		c:               client,
		resourceManager: resourceManager,
		// By default notifier is set to a function which is a no-op. In the event we've implemented the PodNotifier interface,
		// it will be set, and then we'll call a real underlying implementation.
		// This makes it easier in the sense we don't need to wrap each method.
		notifier: func(pod *v1.Pod) {},
	}

	go provider.reconcile()
	return &provider, nil
}

// NewPodmanV0Provider creates a new PodmanV0Provider
func NewPodmanV0Provider(providerConfig, nodeName, operatingSystem string, resourceManager *manager.ResourceManager) (*PodmanV0Provider, error) {
	config, err := loadConfig(providerConfig, nodeName)
	if err != nil {
		return nil, err
	}

	return NewPodmanV0ProviderPodmanConfig(config, nodeName, operatingSystem, resourceManager)
}

// NewPodmanProviderPodmanConfig creates a new PodmanProvider with the given config
func NewPodmanProviderPodmanConfig(config PodmanConfig, nodeName, operatingSystem string, resourceManager *manager.ResourceManager) (*PodmanProvider, error) {
	p, err := NewPodmanV0ProviderPodmanConfig(config, nodeName, operatingSystem, resourceManager)

	return &PodmanProvider{PodmanV0Provider: p}, err
}

// NewPodmanProvider creates a new PodmanProvider, which implements the PodNotifier interface
func NewPodmanProvider(providerConfig, nodeName, operatingSystem string, resourceManager *manager.ResourceManager) (*PodmanProvider, error) {
	config, err := loadConfig(providerConfig, nodeName)
	if err != nil {
		return nil, err
	}

	return NewPodmanProviderPodmanConfig(config, nodeName, operatingSystem, resourceManager)
}
