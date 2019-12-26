package podman

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/varlink/go/varlink"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/virtual-kubelet/podman/pkg/converter"
	"github.com/virtual-kubelet/podman/pkg/iopodman"
	"github.com/virtual-kubelet/podman/pkg/util/errors"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

var (
	// Provider configuration defaults.
	defaultSocket = "unix:/run/podman/io.podman"
	defaultSleep  = time.Millisecond * 100
)

// Config defines podman configurables
type Config struct {
	Socket *string
	Log    *zap.SugaredLogger
}

type conn struct {
	varlink.Connection
	sync.Mutex
}

type podman struct {
	c   *conn
	log *zap.SugaredLogger
}

// Podman is an simplified interface to interfact with
// podman varlink api
type Podman interface {
	// Methdods locking the connection
	Create(ctx context.Context, pod *corev1.Pod) error
	Delete(ctx context.Context, pod *corev1.Pod) error
	GetByName(ctx context.Context, name string) (*corev1.Pod, error)
	List(ctx context.Context) (*corev1.PodList, error)
	GetPodStats(ctx context.Context, pod *corev1.Pod) (*stats.PodStats, error)
	// Methods using above methods
	Update(ctx context.Context, pod *corev1.Pod) error
	CreateOrUpdate(ctx context.Context, pod *corev1.Pod) error
	Get(ctx context.Context, pod *corev1.Pod) (*corev1.Pod, error)
}

// New created new instance of podman interface
func New(ctx context.Context, c *Config) (Podman, error) {
	podman := podman{}
	cfg := getConfig(c)
	var err error
	vConn, err := varlink.NewConnection(ctx, *cfg.Socket)
	if err != nil {
		return nil, err
	}

	conn := conn{
		Connection: *vConn,
	}
	podman.c = &conn
	podman.log = cfg.Log

	return podman, nil
}

func getConfig(c *Config) *Config {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	if c != nil {
		if c.Socket == nil {
			c.Socket = &defaultSocket
		}
		if c.Log == nil {
			c.Log = log
		}
		return c
	}

	return &Config{
		Socket: &defaultSocket,
		Log:    log,
	}
}

// Create creates podman pod and containers within
func (p podman) Create(ctx context.Context, pod *corev1.Pod) error {
	if pod == nil {
		return fmt.Errorf("create pod can't be nil")
	}

	key := converter.BuildKey(pod)
	podmanPod, err := converter.GetPodmanPod(key, pod)
	if err != nil {
		p.log.Error("getPodmanPod failed", "err", err.Error())
		return err
	}
	p.c.Lock()
	podmanPodName, err := iopodman.CreatePod().Call(ctx, &p.c.Connection, *podmanPod)
	p.c.Unlock()
	if err != nil {
		p.log.Error("create pod failed", "err", err.Error())
		return errors.VKError(err)
	}

	p.log.Info("pod created ", "podName ", podmanPodName)
	// Create hostPath volumes if does not exist
	for _, volume := range pod.Spec.Volumes {
		if volume.HostPath != nil {
			switch *volume.HostPath.Type {
			case v1.HostPathDirectoryOrCreate:
				err := os.MkdirAll(volume.HostPath.Path, os.FileMode(0755))
				if err != nil {
					return err
				}
			case v1.HostPathDirectory:
				if _, err := os.Stat(volume.HostPath.Path); os.IsNotExist(err) {
					return fmt.Errorf("volume %s does not exist", volume.Name)
				}
			default:
				p.log.Debug("hostPath volume type %s is not supported", volume.HostPath.Type)
			}
		} else {
			p.log.Debug("volume provider %s is not supported", volume.String())
		}
	}

	// add containers in the pod
	for _, c := range pod.Spec.Containers {
		p.log.Info("create container ", "pod ", podmanPodName, " container ", c.Name)
		container := converter.KubeSpecToPodmanContainer(*pod, c, podmanPodName)

		// pull image
		p.c.Lock()
		_, err := iopodman.PullImage().Call(ctx, &p.c.Connection, c.Image)
		p.c.Unlock()
		if err != nil {
			p.log.Error("error pullImage", "err", err.Error())
			return errors.VKError(err)
		}

		p.c.Lock()
		_, err = iopodman.CreateContainer().Call(ctx, &p.c.Connection, container)
		p.c.Unlock()
		if err != nil {
			p.log.Error("error createContainer", "err", err.Error())
			return errors.VKError(err)
		}
	}

	// start pod
	p.c.Lock()
	_, err = iopodman.StartPod().Call(ctx, &p.c.Connection, podmanPodName)
	p.c.Unlock()
	if err != nil {
		p.log.Error("error startPod", "err", err.Error())
		return errors.VKError(err)
	}

	// check pod status
	retry := 1
	for retry < 5 {
		retry++
		p.c.Lock()
		podmanPodStatus, err := iopodman.InspectPod().Call(ctx, &p.c.Connection, podmanPodName)
		p.c.Unlock()
		if err != nil {
			p.log.Error("error GetPod.InspectPod ", "err ", err.Error())
			return errors.VKError(err)
		}

		var status PodmanPod
		err = json.Unmarshal([]byte(podmanPodStatus), &status)
		if err != nil {
			return errors.VKError(err)
		}
		healty := true
		for _, c := range status.Containers {
			if c.State != "running" {
				healty = false
			}
		}
		if healty {
			continue
		}
	}

	return nil
}

func (p podman) CreateOrUpdate(ctx context.Context, pod *corev1.Pod) error {
	if pod == nil {
		return fmt.Errorf("create pod can't be nil")
	}

	// for logging only
	key := converter.BuildKey(pod)

	pp, err := p.Get(ctx, pod)
	if err != nil {
		if _, ok := err.(*iopodman.PodNotFound); ok {
			p.log.Debugf("pod not found, creating", " pod ", key)
			return p.Create(ctx, pod)
		}
		if pp != nil && err == nil {
			p.log.Debugf("pod exist, update", " pod ", key)
			return p.Update(ctx, pod)
		}
	}

	return nil
}

func (p podman) Delete(ctx context.Context, pod *corev1.Pod) error {
	if pod == nil {
		p.log.Error("pod can't be nil")
		return fmt.Errorf("pod can't be nil")
	}

	key := converter.BuildKey(pod)
	p.c.Lock()
	_, err := iopodman.RemovePod().Call(ctx, &p.c.Connection, key, true)
	p.c.Unlock()
	if err != nil {
		p.log.Error("error while deleting pod", " pod ", key, " err ", err.Error())
		return errors.VKError(err)
	}

	return nil
}

func (p podman) Update(ctx context.Context, pod *corev1.Pod) error {
	err := p.Delete(ctx, pod)
	if err != nil {
		return errors.VKError(err)
	}
	return p.Create(ctx, pod)
}

func (p podman) Get(ctx context.Context, input *corev1.Pod) (pod *v1.Pod, err error) {
	key := converter.BuildKey(input)
	return p.GetByName(ctx, key)
}

func (p podman) GetByName(ctx context.Context, name string) (pod *v1.Pod, err error) {
	p.c.Lock()
	_, err = iopodman.GetPod().Call(ctx, &p.c.Connection, name)
	p.c.Unlock()
	if err != nil {
		return nil, errors.VKError(err)
	}

	p.c.Lock()
	pPod, err := iopodman.InspectPod().Call(ctx, &p.c.Connection, name)
	p.c.Unlock()
	if err != nil {
		return nil, err
	}

	if len(pPod) > 0 {
		kpod, err := converter.GetKubePod(pPod)
		if err != nil {
			return nil, errors.VKError(err)
		}
		return kpod, nil
	}
	return nil, errors.VKError(err)

}

func (p podman) List(ctx context.Context) (podList *corev1.PodList, err error) {
	p.c.Lock()
	pPods, err := iopodman.ListPods().Call(ctx, &p.c.Connection)
	p.c.Unlock()
	if err != nil {
		return nil, errors.VKError(err)
	}

	kpodsList := &corev1.PodList{}
	for _, podData := range pPods {
		kpod, err := p.GetByName(ctx, podData.Name)
		if err != nil {
			return nil, errors.VKError(err)
		}
		kpodsList.Items = append(kpodsList.Items, *kpod)
	}

	return kpodsList, nil
}

// GetContainerStats return container status from pod name and namespace
// TODO: Implement sum of rss
func (p podman) GetPodStats(ctx context.Context, kPod *v1.Pod) (*stats.PodStats, error) {
	name := converter.BuildKey(kPod)
	p.c.Lock()
	podmanJSON, err := iopodman.InspectPod().Call(ctx, &p.c.Connection, name)
	p.c.Unlock()
	if err != nil {
		return nil, err
	}

	pPod, err := converter.MarshalPodPod(podmanJSON)
	if err != nil {
		return nil, err
	}

	// TODO: When supporting multiple containers in the pod
	// iterate here and return aggregate of rss
	p.c.Lock()
	stat, err := iopodman.GetContainerStats().Call(ctx, &p.c.Connection, pPod.Containers[0].ID)
	p.c.Unlock()
	if err != nil {
		return nil, errors.VKError(err)
	}

	time := metav1.NewTime(time.Now())
	cpuUint64 := uint64(stat.Cpu_nano)
	memUint64 := uint64(stat.Mem_usage)

	pss := &stats.PodStats{
		PodRef: stats.PodReference{
			Name:      kPod.Name,
			Namespace: kPod.Namespace,
			UID:       string(kPod.UID),
		},
		StartTime: kPod.CreationTimestamp,
		// TODO: These should be aggregate all pods stats
		CPU: &stats.CPUStats{
			Time:           time,
			UsageNanoCores: &cpuUint64,
		},
		Memory: &stats.MemoryStats{
			Time:       time,
			UsageBytes: &memUint64,
		},
	}
	// TODO: These should be individual pods stats
	pss.Containers[0] = stats.ContainerStats{
		Name:      kPod.Spec.Containers[0].Name,
		StartTime: kPod.CreationTimestamp,
		CPU: &stats.CPUStats{
			Time:           time,
			UsageNanoCores: &cpuUint64,
		},
		Memory: &stats.MemoryStats{
			Time:       time,
			UsageBytes: &memUint64,
		},
	}

	return pss, nil
}
