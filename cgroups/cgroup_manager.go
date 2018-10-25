package cgroups

type CgroupManager struct {
	Name       string
	Limits     map[string]string
	Subsystems map[string]*Subsystem
}

func initSubsystems(cgroupName string, limits map[string]string) map[string]*Subsystem {
	subsystems := make(map[string]*Subsystem)
	subsystems["cpu"], _ = InitCgroupSubSystem("cpu", cgroupName, limits["cpu"], "cpu.shares")
	subsystems["cpuset"], _ = InitCgroupSubSystem("cpuset", cgroupName, limits["cpuset"], "cpuset.cpus")
	subsystems["memory"], _ = InitCgroupSubSystem("memory", cgroupName, limits["memory"], "memory.limit_in_bytes")
	return subsystems
}
func NewCgroupManager(name string, limits map[string]string) *CgroupManager {
	return &CgroupManager{
		Name:       name,
		Limits:     limits,
		Subsystems: initSubsystems(name, limits),
	}
}

func (c *CgroupManager) SetLimit() {
	for _, subsystem := range c.Subsystems {
		subsystem.SetLimit()
	}
}

func (c *CgroupManager) ApplyPid(pid int) {
	for _, subsystem := range c.Subsystems {
		subsystem.ApplyPid(pid)
	}
}

func (c *CgroupManager) Destory() {
	for _, subsystem := range c.Subsystems {
		subsystem.Remove()
	}
}
