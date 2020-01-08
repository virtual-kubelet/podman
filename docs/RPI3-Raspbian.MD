# Raspberry PI 3 - raspbian

This document explains how to setup RPI3 to be used as edge node for
podman virtual kubelet node

## Install

1. [Download](https://www.raspberrypi.org/downloads/raspbian/) Raspbian ARM image.
2. Install [Balena Etcher](https://www.balena.io/etcher/)
3. Bake ARM image into RPI SD Card:

4. Mount root filesystem of the SD card with R/W and change static files to boot
without promts

```
Add wifi settings:
cat <<EOF >>/run/media/$USER/rootfs/etc/wpa_supplicant/wpa_supplicant.conf
network={
    ssid="SSID_NAME"
    psk="Password"
}
EOF
```

Set Hostname
```
echo "rpi-smallscreen" >> /run/media/$USER/rootfs/etc/hostname
```

Enable SSH
```
touch /run/media/mjudeiki/rootfs/ssh
```
SSH user: pi, password: raspberry


Add startup script to rc.local

cat <<EOF >/run/media/mjudeiki/rootfs/root/bootstrap.sh
#!/bin/bash
curl curl https://gist.githubusercontent.com/mjudeikis/868dbd82bd738e1e926b7c1b25c9555d/raw/637d7d984d94eb249a8ed616a6f166bc5906442b/gistfile1.txt | sh
EOF

Append script to rc.local:

sudo sed -i "s/exit 0//g" /etc/rc.local
