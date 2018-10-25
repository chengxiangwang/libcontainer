package reexec

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

func Self() string {
	return "/proc/self/exe"
}

func Command(args ...string) *exec.Cmd {
	return &exec.Cmd{
		Path: Self(),
		Args: args,
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig: unix.SIGTERM,
		},
	}
}
