package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/chengxiangwang/libcontainer/config"
)

func RunContainerInitProcess() {
	args, err := readInitCommand()
	if err != nil {
		fmt.Printf("read init args from pipe error %v", err)
		os.Exit(1)
	}
	containerId := args[0]
	c, err := readContainerFromConfig(containerId)
	if err != nil {
		fmt.Printf("read container from config file error %v", err)
		os.Exit(1)
	}
	if err := preRunInit(c); err != nil {
		fmt.Printf("pre run container init process error %v", err)
		os.Exit(1)
	}
	if err := postRunInit(c); err != nil {
		fmt.Printf("post run container init process error %v", err)
		os.Exit(1)
	}
}

func readInitCommand() ([]string, error) {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		return nil, err
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " "), nil
}
func readContainerFromConfig(containerId string) (*Container, error) {
	configFile := path.Join(config.CONTAINER_ROOT, containerId, "config.json")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	c, err := JsonAsContainer(data)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func preRunInit(c *Container) error {
	if err := MountProc(c.Root); err != nil {
		return fmt.Errorf("mount proc error %v", err)
	}
	if err := PivotRoot(c.Root); err != nil {
		return fmt.Errorf("pivot root error %v", err)
	}
	if err := syscall.Sethostname([]byte(c.ID.String())); err != nil {
		return fmt.Errorf("set hostname error %v", err)
	}
	return nil
}

func postRunInit(c *Container) error {
	if err := execInitCommand(c); err != nil {
		return err
	}
	return nil
}

func execInitCommand(c *Container) error {
	initArgs := c.Args
	if initArgs == nil || len(initArgs) < 1 {
		return fmt.Errorf("user init command args is nil")
	}
	path, err := exec.LookPath(initArgs[0])
	if err != nil {
		return fmt.Errorf("command %s not found error", initArgs[0])
	}
	if err := syscall.Exec(path, initArgs[0:], os.Environ()); err != nil {
		return fmt.Errorf("exec command %s error %v", initArgs[0], err)
	}
	return nil

}
