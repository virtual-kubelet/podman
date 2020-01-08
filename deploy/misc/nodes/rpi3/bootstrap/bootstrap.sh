#!/bin/bash

FILE=/root/bootstrap.done
if [[ ! -f "$FILE" ]]; then
    echo "Bootstrap starting"
    # remove init startup
    rm /etc/systemd/system/multi-user.target.wants/initial-setup.service
    # extend disk
    xfs_growfs /
    # update
    yum update -y
    # install git and podman
    yum install --enablerepo=updates-testing git podman libvarlink-util libvarlink containers-common -y
    # enable podman
    systemctl enable io.podman.service
    systemctl start io.podman.service

    # install graphical interface
    dnf groupinstall -y "LXDE Desktop" -y
    systemctl set-default graphical.target
    systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target

    rm -rf /usr/lib/systemd/system/initial-setup.service

    # configure users
    useradd rpi
    passwd -d rpi
    # auto-login
    sed -i "s/# autologin=dgod/autologin=rpi/g" /etc/lxdm/lxdm.conf

    # boostrap node
    mkdir -p /etc/kubernetes/ /etc/vkubelet/

    # create podman config
    cat <<EOF >/etc/vkubelet/podman-cfg.json
    {
      "podman": {
        "cpu": "1",
        "memory": "2Gi",
        "pods": "10",
        "socket": "unix:/run/podman/io.podman",
        "daemonSetDisabled": "true"
      }
    }
EOF

    cat <<EOF >/usr/lib/systemd/system/vkubelet-podman.service
[Unit]
Description=vkubelet-podman
Requires=io.podman.service
[Service]
Environment=KUBECONFIG=/etc/kubernetes/admin.conf
ExecStart=/usr/local/bin/virtual-kubelet --provider podman --nodename $HOSTNAME --provider-config /etc/vkubelet/podman-cfg.json
[Install]
WantedBy=multi-user.target
EOF

  # download vk podman binary
  curl https://raw.githubusercontent.com/mjudeikis/podman/master/bin/virtual-kubelet-arm -o /usr/local/bin/virtual-kubelet
  chmod 755 /usr/local/bin/virtual-kubelet

    touch /root/bootstrap.done
    exit 0
else
    echo "Bootstrap already done"
fi
