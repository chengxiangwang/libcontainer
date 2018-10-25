package cgroups

import (
	"github.com/chengxiangwang/libcontainer/config"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type Subsystem struct {
	Name          string
	CgroupName    string
	Path          string
	LimitFilePath string
	TaskFilePath  string
	Limit         string
}

//create subsystem cgroup dir
func InitCgroupSubSystem(name string, cgroupName string, limit string, limitFileName string) (*Subsystem, error) {
	subsystemCgroupDir := path.Join(config.CGROUP_ROOT, name, cgroupName)
	if err := os.Mkdir(subsystemCgroupDir, os.ModeDir); err != nil && os.IsNotExist(err) {
		return nil, err
	}
	subsystem := &Subsystem{
		Name:          name,
		CgroupName:    cgroupName,
		Path:          subsystemCgroupDir,
		LimitFilePath: path.Join(subsystemCgroupDir, limitFileName),
		TaskFilePath:  path.Join(subsystemCgroupDir, "tasks"),
		Limit:         limit,
	}
	return subsystem, nil
}

func (s *Subsystem) SetLimit() error {
	if s.Limit != "" {
		data := []byte(s.Limit)
		if err := ioutil.WriteFile(s.LimitFilePath, data, 0644); err != nil {
			return err
		}
	}
	return nil

}

func (s *Subsystem) ApplyPid(pid int) error {
	data := []byte(strconv.Itoa(pid))
	if err := ioutil.WriteFile(s.TaskFilePath, data, 0644); err != nil {
		return err
	}
	return nil
}

func (s *Subsystem) Remove() error {
	return os.Remove(s.Path)
}
