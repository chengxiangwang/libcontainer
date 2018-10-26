package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/chengxiangwang/libcontainer/config"
	"github.com/chengxiangwang/libcontainer/utils"
)

func GetContainerLogFile(c *Container) (*os.File, error) {
	if !utils.PathExists(c.Path) {
		if err := os.MkdirAll(c.Path, os.ModeDir); err != nil {
			return nil, err
		}
	}
	containerLogFile := path.Join(c.Path, "container.log")
	stdLogFile, err := os.Create(containerLogFile)
	if err != nil {
		return nil, err
	}
	return stdLogFile, nil
}

func PrintContainerLog(containerId string) {
	containerLogFile := path.Join(config.CONTAINER_ROOT, containerId, "container.log")
	if !utils.PathExists(containerLogFile) {
		return
	}
	file, err := os.Open(containerLogFile)
	defer file.Close()
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	fmt.Fprint(os.Stdout, string(content))
}
