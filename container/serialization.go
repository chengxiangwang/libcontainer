package container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/chengxiangwang/libcontainer/config"
)

func ContainerAsJson(c *Container) ([]byte, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func JsonAsContainer(jsonBytes []byte) (*Container, error) {
	ret := Container{}
	err := json.Unmarshal(jsonBytes, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func CheckContainerByName(name string) bool {
	_, err := os.Stat(config.CONTAINER_ROOT)
	if os.IsNotExist(err) {
		err := os.MkdirAll(config.CONTAINER_ROOT, os.ModeDir)
		if err != nil {
			logrus.Errorf("create container root dir error %v", err)
			return false
		}
	}
	rootDir, err := ioutil.ReadDir(config.CONTAINER_ROOT)
	if err != nil {
		logrus.Errorf("read container root dir error %v", err)
		return false
	}
	for _, chDir := range rootDir {
		if chDir.IsDir() {
			configFile := path.Join(config.CONTAINER_ROOT, chDir.Name(), "config.json")
			data, _ := ioutil.ReadFile(configFile)
			c, _ := JsonAsContainer(data)
			if c.Name == name {
				return true
			}
		}
	}
	return false
}

func CheckInitArgs(initArgs []string) bool {
	if initArgs == nil || len(initArgs) < 1 {
		return false
	}
	return true
}

func WriteContainerConfigFile(c *Container) error {
	jsonBytes, err := ContainerAsJson(c)
	if err != nil {
		return fmt.Errorf("container as json error %v", err)
	}
	containerRoot := path.Join(config.CONTAINER_ROOT, c.ID.String())
	err = os.Mkdir(containerRoot, os.ModeDir)
	if err != nil {
		return fmt.Errorf("create container dir %s error %v", containerRoot, err)
	}
	configFile := path.Join(containerRoot, "config.json")
	err = ioutil.WriteFile(configFile, jsonBytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write container config file %s error %v", configFile, err)
	}
	return nil
}
