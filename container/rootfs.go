package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/chengxiangwang/libcontainer/config"
)

func PivotRoot(newroot string) error {
	putold := filepath.Join(newroot, "/.pivot_root")

	// bind mount newroot to itself - this is a slight hack
	// needed to work around a pivot_root requirement
	if err := syscall.Mount(
		newroot,
		newroot,
		"",
		syscall.MS_BIND|syscall.MS_REC,
		"",
	); err != nil {
		return err
	}

	// create putold directory
	if err := os.MkdirAll(putold, 0700); err != nil {
		return err
	}

	// call pivot_root
	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return err
	}

	// ensure current working directory is set to new root
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount putold, which now lives at /.pivot_root
	putold = "/.pivot_root"
	if err := syscall.Unmount(
		putold,
		syscall.MNT_DETACH,
	); err != nil {
		return err
	}

	// remove putold
	if err := os.RemoveAll(putold); err != nil {
		return err
	}

	return nil
}

func MountProc(rootPath string) error {
	source := "proc"
	target := filepath.Join(rootPath, "proc")
	fstype := "proc"
	flags := 0
	data := ""
	if err := syscall.Mount(source, target, fstype, uintptr(flags), data); err != nil {
		return fmt.Errorf("mount proc to root path %s error %v", rootPath, err)
	}
	return nil
}

func InitContainerFS(c *Container) error {
	if err := initReadOnlyLayer(c); err != nil {
		return err
	}
	if err := initWriteLayer(c); err != nil {
		return err
	}
	if err := mountReadWriteLayerToRootFs(c); err != nil {
		return err
	}
	if err := MountVolumn(c); err != nil {
		return err
	}
	return nil
}
func initReadOnlyLayer(c *Container) error {
	_, err := os.Stat(c.ReadOnlyLayer)
	if os.IsNotExist(err) {
		err := os.MkdirAll(c.ReadOnlyLayer, os.ModeDir)
		if err != nil {
			return fmt.Errorf("create container readonlylayer error %v", err)
		}
	}
	if _, err := exec.Command("tar", "-xvf", config.TEST_ROOTFS, "-C", c.ReadOnlyLayer).CombinedOutput(); err != nil {
		return fmt.Errorf("unzip image to contaienr readonlylayer error %v", err)
	}
	return nil
}

func initWriteLayer(c *Container) error {
	_, err := os.Stat(c.WriteLayer)
	if os.IsNotExist(err) {
		err := os.MkdirAll(c.WriteLayer, os.ModeDir)
		if err != nil {
			return fmt.Errorf("create container write layer  error %v", err)
		}
	}
	return nil
}

func mountReadWriteLayerToRootFs(c *Container) error {
	_, err := os.Stat(c.Root)
	if os.IsNotExist(err) {
		err := os.MkdirAll(c.Root, os.ModeDir)
		if err != nil {
			return fmt.Errorf("create container rootfs  error %v", err)
		}
	}

	dirs := "dirs=" + c.WriteLayer + ":" + c.ReadOnlyLayer
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", c.Root)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func ReleaseContainerFS(c *Container) error {
	if err := UnmountVolumn(c); err != nil {
		return err
	}
	if err := syscall.Unmount(c.Root, syscall.MNT_DETACH); err != nil {
		return err
	}
	if err := os.RemoveAll(c.Root); err != nil {
		return err
	}
	return nil
}
