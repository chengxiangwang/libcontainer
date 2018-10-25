package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/chengxiangwang/libcontainer/utils"
)

type Volumn struct {
	HostPath      string
	ContainerPath string
}

func NewVolumns(volumnStrs []string) []*Volumn {
	ret := []*Volumn{}
	for _, volumnStr := range volumnStrs {
		volumnPair := strings.Split(volumnStr, ":")
		if volumnPair == nil || len(volumnPair) != 2 {
			continue
		}
		volumn := &Volumn{
			HostPath:      volumnPair[0],
			ContainerPath: volumnPair[1],
		}
		if !volumn.isLegal() {
			continue
		}
		ret = append(ret, volumn)
	}
	return ret
}

func (v *Volumn) isLegal() bool {
	return utils.PathExists(v.HostPath)
}

func MountVolumn(c *Container) error {
	volumns := NewVolumns(c.Volumns)
	for _, volumn := range volumns {
		fmt.Println(volumn)
		if err := doMountVolumn(c, volumn); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func doMountVolumn(c *Container, volumn *Volumn) error {
	if err := utils.MkDir(volumn.HostPath); err != nil {
		return err
	}
	containerVolumnPath := path.Join(c.Root, volumn.ContainerPath)
	if err := utils.MkDir(containerVolumnPath); err != nil {
		return err
	}
	dirs := "dirs=" + volumn.HostPath
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumnPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func UnmountVolumn(c *Container) error {
	volumns := NewVolumns(c.Volumns)
	for _, volumn := range volumns {
		if err := doUnmountVolumn(c, volumn); err != nil {
			return err
		}
	}
	return nil

}

func doUnmountVolumn(c *Container, volumn *Volumn) error {
	containerVolumnPath := path.Join(c.Root, volumn.ContainerPath)
	if err := syscall.Unmount(containerVolumnPath, syscall.MNT_DETACH); err != nil {
		return err
	}
	return nil
}
