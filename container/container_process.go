package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/chengxiangwang/libcontainer/cgroups"
	"github.com/chengxiangwang/libcontainer/config"
	"github.com/chengxiangwang/libcontainer/reexec"
	"github.com/chengxiangwang/libcontainer/utils"
)

func init() {
	reexec.Register("initContainer", RunContainerInitProcess)
	if reexec.Init() {
		os.Exit(0)
	}
}

func RunContainer(name string, tty bool, initArgs []string, limits map[string]string, volumns []string) error {
	if CheckContainerByName(name) {
		return fmt.Errorf("The container named %s already exists", name)
	}
	if !CheckInitArgs(initArgs) {
		return fmt.Errorf("init args can not be nil")
	}
	c := &Container{}
	c.ID = NewUUID()
	c.Name = name
	c.TTY = tty
	c.Root = path.Join(config.CONTAINER_ROOT, c.ID.String(), "rootfs")
	c.Path = path.Join(config.CONTAINER_ROOT, c.ID.String())
	c.ReadOnlyLayer = path.Join(config.CONTAINER_ROOT, c.ID.String(), "readOnlyLayer")
	c.WriteLayer = path.Join(config.CONTAINER_ROOT, c.ID.String(), "writeLayer")
	c.Args = initArgs
	c.Limits = limits
	c.Volumns = volumns
	err := WriteContainerConfigFile(c)
	if err != nil {
		return err
	}
	if err := doRun(c); err != nil {
		return err
	}
	return nil
}
func doRun(c *Container) error {
	cmd, writePipe, err := createContainerProcess(c)
	if err != nil || cmd == nil || writePipe == nil {
		return fmt.Errorf("create containerd process error %v", err)
	}
	if err := InitContainerFS(c); err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start containerd process error %v", err)
	}
	//createCgroup(cmd, c)
	if err := sendInitCommand(c, writePipe); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	if err := ReleaseContainerFS(c); err != nil {
		return nil
	}
	return nil
}
func createContainerProcess(c *Container) (*exec.Cmd, *os.File, error) {
	readPipe, writePipe, err := utils.CreatePipePair()
	if err != nil {
		return nil, nil, err
	}
	cmd := reexec.Command("initContainer")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	if c.TTY {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	//cmd.Dir = c.Root
	return cmd, writePipe, nil
}

func sendInitCommand(c *Container, writePipe *os.File) error {
	arr := []string{c.ID.String()}
	arr = append(arr, c.Args...)
	command := strings.Join(arr, " ")
	_, err := writePipe.WriteString(command)
	if err != nil {
		return err
	}
	if err := writePipe.Close(); err != nil {
		return err
	}
	return nil
}

func createCgroup(cmd *exec.Cmd, c *Container) error {
	cgroupManager := cgroups.NewCgroupManager(c.ID.String(), c.Limits)
	defer cgroupManager.Destory()
	cgroupManager.SetLimit()
	cgroupManager.ApplyPid(cmd.Process.Pid)
	return nil
}
