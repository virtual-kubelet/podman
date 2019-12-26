package errors

import (
	"github.com/virtual-kubelet/podman/pkg/iopodman"
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
)

// VKError takes in varlink error and returns Virtual kubelet error
func VKError(err error) error {
	switch err {
	case err.(*iopodman.PodNotFound):
		return errdefs.NotFound("ImageNotFound")
	default:
		return errdefs.AsNotFound(err)
	}
}
