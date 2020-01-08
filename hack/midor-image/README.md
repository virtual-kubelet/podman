# Midori ARM browser image

Midori image is minimal image, used to run browser in Kiosk mode.

Manual start command:
```
cp /home/rpi/.Xauthority /root/.Xauthority
podman run -it -e NO_AT_BRIDGE=1 -e XAUTHORITY=/home/rpi/.Xauthority -e DISPLAY=:0.0 -v /home/rpi:/home/rpi  -v /etc/localtime:/etc/localtime:ro --user rpi -v /usr:/usr  -v /tmp:/tmp -v /etc/passwd:/etc/passwd --net host --privileged quay.io/mangirdas/midori:arm32v7 midori -e Fullscreen -a https://www.containers.ninja

```
