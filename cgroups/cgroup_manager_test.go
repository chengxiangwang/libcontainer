package cgroups

import (
	"fmt"
	"testing"
)

func TestCgroupManager(t *testing.T) {
	limits := map[string]string{
		"cpu":    "512",
		"cpuset": "1",
		"memory": "100m",
	}
	cgroupManager := NewCgroupManager("test-cgroup2", limits)
	for name, subsystem := range cgroupManager.Subsystems {
		fmt.Printf("name %s,path %s\n", name, subsystem.Path)
	}

	cgroupManager.SetLimit()
	cgroupManager.ApplyPid(1000000)
	cgroupManager.Destory()
}
