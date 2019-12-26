# Roadmap for Podman Virtual Kubelet provider

Roadmap for podman provider consist of items, needed to be implemented to move
this provider from beta stage to GA.

## Easy

1. Move all errors handling to errdefs implementation
2. Sometimes podman operations hangs on API calls.
   Add context to all operations with timeouts and retries. In addition need to
   incestigate what is making podman to hang.
3. Currently provide does not support restart on container failure.
   This need to be fixed
4. Onboard project to github workers for unit tests
5. Create `curl www.example.com/bootstrap | sh` script to bootsrap different type
   of end devices, like RPI configuration.

## Not so easy

0. Tests tests tests!
   Unit
   E2E
1. Implement `GetContainerLogs`, `RunInContainer`.
2. Add better kube feature parity support. In example to  enable volumes, secrets, configMaps.
3. Add ability to check if podman is alive and update node status on time intervals.
4. Configure node and schedule pods based on configures limits.
   Currently pod limits are not being translated into the pod limits.
5. Add support for "remote vkubelet podman" where vkubelet is running in the
   cluster as a pod and it reaches to podman node via remote varlink api via SSH.
   This might require ssh to be running in the container (yes, its nasty), but
   it would make managing fleet of ssh enabled devices much easier.
   Running kube node in as pod :)
6. Add support for better pod status. Like "CrashBackLoop", "Failed", "ImageNotFound", etc.
7. Add support for 2 containers in the podman pod. Currently only pod 1 is running.
8. Add support for init containers
9.
