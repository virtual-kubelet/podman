# Virtual Kubelet Podman Provider - Beta

Virtual Kubelet Provider implementation to manages containers in
podman environments.

*Important: Project is under development and currently only basic functionality is available*

## Purpose

The whole point of the Virtual Kubelet project is to provide an interface for
container runtimes that don't conform to the standard node-based model. The
[Kubelet](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubelet)
codebase is the comprehensive standard CRI node agent and this Provider is not
attempting to recreate that.

## Dependencies

[Podman](https://podman.io/getting-started/installation.html) must be installed
on the provider node. Podman provider uses [varlink](https://www.projectatomic.io/blog/2018/05/podman-varlink/)
podman [API](https://podman.io/blogs/2019/01/16/podman-varlink.html) to communicate.
This must be enabled for podman provider to work.

## Running

### Production

To run podman provider in production configuration requires minimal configuration

```bash
# create vkubelet and kubernetes directories
mkdir -p /etc/kubernetes/ /etc/vkubelet/
# copy kubeconfig into kubernetes directory
cp $KUBECONFIG /etc/kubernetes/admin.conf
# Copy systemd file into the destination node.
cp ./deploy/systemd/vkubelet-podman.service /usr/lib/systemd/system/vkubelet-podman.service
# Copy vkubelet configuration file. Modify it based on your requirments
cp ./deploy/systemd/podman-cfg.json /etc/vkubelet/podman-cfg.json
# Copy vkubelet binary
cp ./bin/virtual-kubelet .usr/local/bin/virtual-kubelet
# Reload and start vkubelet daemon
systemctl daemon-reload
systemctl start vkubelet-podman
systemctl status vkubelet-podman
```

### Development

For local development it is easiest way to iterate is to use use `[minikube](https://github.com/kubernetes/minikube)`

```bash
#start minikube
minikube start

#start varlink podman socket
systemctl start io.podman.service

# start provider
make run
kubectl create -f deploy/example/pod-example.yaml
# you should see pod running in podman (you might need root to check it)
podman ps
```

## Limitations

* Only `hostPath` volume provider is supported
* Only one container per pod is supported
* No `Secrets` or `ConfigMaps` is supported yet

## Podman install & configuration

```
yum distro-sync --enablerepo=updates-testing install podman containers-common sudo
systemctl enable --now io.podman.socket
systemctl status io.podman.socket
```

I have enabled my local dev user to access socket for development:
```
[Socket]
ListenStream=/run/podman/io.podman
SocketMode=0770
SocketGroup=root
SocketUser=$USER
```

`vi /etc/tmpfiles.d/podman.conf`:
```
d /run/podman 0770 root $USER
```

## Roadmap

[Provider Roadmap](ROADMAP.md)
