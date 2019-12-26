# Raspberry PI 3 - Fedora

This document explains how to setup RPI3 to be used as edge node for
podman virtual kubelet node

## Install

1. [Download](https://arm.fedoraproject.org/) Fedora ARM minial image.
   * Note: Fedora 31 is not supported on RPI3 and RPI4. Only valid configuration currently is RPI3 & Fedora 30*
2. Install [ARM installer](https://fedoraproject.org/wiki/Architectures/ARM/Installation)
3. Bake ARM image into RPI SD Card:

```
arm-image-installer --addkey=/home/$USER/.ssh/id_rsa.pub \
--target=rpi3 --media=/dev/sda \
--image=/home/$USER/Documents/arm/Fedora-Server-armhfp-30-1.2-sda.raw.xz \
--resizefs
```

4. Mount root filesystem of the SD card with R/W and add WIFI settings

```
# copy network scripts, so WIFI would be "like home"
cp ./deploy/misc/nodes/rpi3/* /run/media/$USER/__/etc/sysconfig/network-scripts/
#cp /etc/sysconfig/network-scripts/* /run/media/$USER/__/etc/sysconfig/network-scripts/
```

Set Hostname
```
echo "rpi3-walldisplay" > /run/media/$USER/__/etc/hostname

```

Copy bootstrap script and systemd file. This is automated way to bootstrap the node.
Important: script is doing `dnf update` so it might take LONG time on the first run.
```
cat <<EOF >/run/media/mjudeiki/__/usr/lib/systemd/system/bootstrap.service
[Unit]

[Service]
ExecStart=curl https://gist.githubusercontent.com/mjudeikis/868dbd82bd738e1e926b7c1b25c9555d/raw/637d7d984d94eb249a8ed616a6f166bc5906442b/gistfile1.txt | sh

[Install]
WantedBy=default.target
EOF
```

Enable service:
```
mkdir -p /run/media/$USER/__/etc/systemd/system/default.target.wants/
ln -s /usr/lib/systemd/system/bootstrap.service  /run/media/$USER/__/etc/systemd/system/default.target.wants/bootstrap.service

```

Startup will install updates, setup virtual-kubelet, rpi user, and make it to
autologin
