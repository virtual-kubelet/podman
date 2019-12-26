package main

import (
	"github.com/virtual-kubelet/podman/pkg/provider"
	"github.com/virtual-kubelet/podman/pkg/provider/podman"
)

func registerPodman(s *provider.Store) {
	s.Register("podman", func(cfg provider.InitConfig) (provider.Provider, error) { //nolint:errcheck
		return podman.NewPodmanProvider(
			cfg.ConfigPath,
			cfg.NodeName,
			cfg.OperatingSystem,
			cfg.ResourceManager,
		)
	})
}
