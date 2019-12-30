package podman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"k8s.io/apimachinery/pkg/api/resource"
)

// loadConfig loads the given json configuration files.
func loadConfig(providerConfig, nodeName string) (config PodmanConfig, err error) {
	data, err := ioutil.ReadFile(providerConfig)
	if err != nil {
		return config, err
	}
	configMap := map[string]PodmanConfig{}
	err = json.Unmarshal(data, &configMap)
	if err != nil {
		return config, err
	}

	if _, exist := configMap[nodeName]; exist {
		config = configMap[nodeName]
		// set defaults
		if config.CPU == "" {
			config.CPU = defaultCPUCapacity
		}
		if config.Memory == "" {
			config.Memory = defaultMemoryCapacity
		}
		if config.Pods == "" {
			config.Pods = defaultPodCapacity
		}
		if config.Socket == "" {
			config.Socket = defaultSocket
		}
		if config.DaemonSetDisabled == "" {
			config.DaemonSetDisabled = defaultDaemonSetDisabled
		}
		if _, err = resource.ParseQuantity(config.CPU); err != nil {
			return config, fmt.Errorf("Invalid CPU value %v, %v", config.CPU, err)
		}
		if _, err = resource.ParseQuantity(config.Memory); err != nil {
			return config, fmt.Errorf("Invalid memory value %v", config.Memory)
		}
		if _, err = resource.ParseQuantity(config.Pods); err != nil {
			return config, fmt.Errorf("Invalid pods value %v", config.Pods)
		}
		if _, err = strconv.ParseBool(config.DaemonSetDisabled); err != nil {
			return config, fmt.Errorf("Invalid daemonSetDisabled value %v", config.DaemonSetDisabled)
		}
	} else {
		return config, fmt.Errorf("Node config not found %v", nodeName)
	}

	return config, nil
}
