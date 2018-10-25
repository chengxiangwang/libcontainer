package cgroups

import (
	"fmt"
	"os"
	"testing"
)

func TestInitCgroupSubSystem(t *testing.T) {
	subsystem, err := InitCgroupSubSystem("cpu", "test-cgroup1", "512", "cpu.shares")
	if err != nil {
		t.FailNow()
	} else {
		fmt.Println(subsystem)
	}
}
func TestSetLimit(t *testing.T) {
	subsystem, err := InitCgroupSubSystem("cpuset", "test-cgroup1", "1", "cpuset.cpus")
	if err == nil {
		if err = subsystem.SetLimit(); err != nil {
			t.FailNow()
		} else {
			fmt.Println("set limit success")
		}
	} else {
		t.FailNow()
	}
}

func TestApplyPid(t *testing.T) {
	subsystem, err := InitCgroupSubSystem("cpuset", "test-cgroup1", "1", "cpuset.cpus")
	if err == nil {
		if err = subsystem.ApplyPid(os.Getpid()); err != nil {
			fmt.Println("appy pid succes")
		} else {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}

func TestRemove(t *testing.T) {
	subsystem, err := InitCgroupSubSystem("cpu", "test-cgroup1", "512", "cpu.shares")
	if err == nil {
		if err = subsystem.Remove(); err == nil {
			fmt.Println("remove subsystem cgroup  succes")
		} else {
			fmt.Println(err)
			t.FailNow()
		}
	} else {
		fmt.Println(err)
		t.FailNow()
	}
}
