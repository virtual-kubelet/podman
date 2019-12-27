package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/virtual-kubelet/podman/pkg/iopodman"
)

func BuildKeyFromNames(namespace string, name string) (string, error) {
	return fmt.Sprintf("%s-%s", namespace, name), nil
}

// BuildKey is a helper for building the "key" for the providers pod store.
func BuildKey(pod *v1.Pod) string {
	if pod.Namespace == "" {
		return pod.Name
	}
	return fmt.Sprintf("%s-%s", pod.Namespace, pod.Name)
}

func SplitPodName(key string) (namespace, name string) {
	keys := strings.Split(key, "-")
	return keys[0], keys[1]
}

// KubeSpecToPodmanContainer converts v1.Container to podman.Create spec. pod
// argument is used to configure volumes and external configuration to container
func KubeSpecToPodmanContainer(pod v1.Pod, container v1.Container, podName string) iopodman.Create {
	// TODO: Extend this to match most of the fields
	var args []string
	args = append(args, container.Image)
	args = append(args, container.Command...)
	args = append(args, container.Args...)
	containerName := fmt.Sprintf("%s-%s", podName, container.Name)

	// construct hostPath pairs for mount
	var volumes []string
	for _, podVolume := range pod.Spec.Volumes {
		for _, containerVolume := range container.VolumeMounts {
			if podVolume.HostPath != nil && podVolume.Name == containerVolume.Name {
				volumes = append(volumes, fmt.Sprintf("%s:%s", podVolume.HostPath.Path, containerVolume.MountPath))
			}
		}
	}

	podmanPod := iopodman.Create{
		Args:    args,
		Command: &container.Command,
		Name:    &containerName,
		Pod:     &podName,
		Volume:  &volumes,
	}

	if container.SecurityContext != nil {
		podmanPod.Privileged = container.SecurityContext.Privileged
		// if we want to interact with undelaying hostPath
		// TODO: make it separate configurable
		podmanPod.Tty = container.SecurityContext.Privileged
	}

	if pod.Spec.HostNetwork {
		podmanPod.Net = StringPtr("host")
	}

	var vars []string
	for _, e := range pod.Spec.Containers[0].Env {
		vars = append(vars, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}
	podmanPod.Env = &vars

	return podmanPod
}

// GetPodmanPod return podmanPod with v1.Pod metadata in the label
func GetPodmanPod(key string, p *v1.Pod) (*iopodman.PodCreate, error) {
	// preserve original pod spec into lables
	pod := p.DeepCopy()
	data, err := yaml.Marshal(pod)
	if err != nil {
		return nil, err
	}
	podSpecBase := base64.StdEncoding.EncodeToString(data)
	if pod.Labels == nil {
		pod.Labels = make(map[string]string, 1)
	}
	pod.Labels["pod"] = podSpecBase

	podmanPod := iopodman.PodCreate{
		Name:   key,
		Labels: pod.Labels,
	}

	return &podmanPod, nil
}

// GetKubePod returns v1.Pod from podman pod json
// Kuberentes spec is cached in the podman labels
func GetKubePod(podmanJSON string) (*v1.Pod, error) {
	var pPod PodmanPod
	err := json.Unmarshal([]byte(podmanJSON), &pPod)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(pPod.Config.Labels["pod"])
	if err != nil {
		return nil, err
	}
	var kpod v1.Pod
	err = yaml.Unmarshal(data, &kpod)
	if err != nil {
		return nil, err
	}

	// configure status for the kubePod
	kpod.Status, err = GetPodStatus(pPod)
	if err != nil {
		return nil, err
	}

	return &kpod, nil
}

// MarshalPodPod marshals podmanPod json into PodmanPod struct
func MarshalPodPod(podmanJSON string) (*PodmanPod, error) {
	var pPod *PodmanPod
	err := json.Unmarshal([]byte(podmanJSON), pPod)
	if err != nil {
		return nil, err
	}
	return pPod, nil
}

// GetPodStatus returns v1.PodStatus from PodmanPod spec
func GetPodStatus(pPod PodmanPod) (v1.PodStatus, error) {
	now := metav1.NewTime(time.Now())
	status := v1.PodStatus{}
	status.StartTime = &now
	status.HostIP = "1.2.3.4"
	status.PodIP = "5.6.7.8"
	status.Conditions = []v1.PodCondition{
		{
			Type:   v1.PodInitialized,
			Status: v1.ConditionTrue,
		},
		{
			Type:   v1.PodReady,
			Status: v1.ConditionTrue,
		},
		{
			Type:   v1.PodScheduled,
			Status: v1.ConditionTrue,
		},
	}

	for _, c := range pPod.Containers {
		containerStatus := v1.ContainerStatus{}
		containerStatus.Name = c.ID
		containerStatus.Image = c.ID
		var state v1.ContainerState
		switch c.State {
		case "running":
			state = v1.ContainerState{
				Running: &v1.ContainerStateRunning{
					StartedAt: metav1.Time{
						Time: pPod.Config.Created,
					},
				},
			}
			containerStatus.Ready = true
			status.Phase = v1.PodRunning
		case "exited":
			state = v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{
					StartedAt: metav1.Time{
						Time: time.Now(),
					},
				},
			}
			status.Phase = v1.PodFailed
			containerStatus.Ready = false
		default:
			state = v1.ContainerState{
				Terminated: &v1.ContainerStateTerminated{
					StartedAt: metav1.Time{
						Time: time.Now(),
					},
				},
			}
			status.Phase = v1.PodUnknown
			containerStatus.Ready = false
		}
		containerStatus.State = state
		status.ContainerStatuses = append(status.ContainerStatuses, containerStatus)
	}

	return status, nil
}

// StringPtr returns pointer string
func StringPtr(s string) *string {
	return &s
}
