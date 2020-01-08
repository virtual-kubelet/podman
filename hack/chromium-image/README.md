# Chromium ARM browser image

Cromium image, used to run Chromium in kiosk mode on `ARM`

```
podman run -it -v /var/run/dbus:/var/run/dbus -e XAUTHORITY=/root/.Xauthority -e DISPLAY=:0.0 -v /home/rpi:/root -v /tmp:/tmp --net host --privileged quay.io/mangirdas/chromium:arm32v7 --app=http://redhat.com --start-fullscreen
```
