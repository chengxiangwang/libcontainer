package container

type State struct {
	ID   int
	Name string
}

func ContainerCreated() *State {
	return &State{
		ID:   0,
		Name: "Created",
	}
}

func ContainerRunning() *State {
	return &State{
		ID:   1,
		Name: "Running",
	}
}

func ContainerStopped() *State {
	return &State{
		ID:   3,
		Name: "Stopped",
	}
}

func ContainerExited() *State {
	return &State{
		ID:   -1,
		Name: "Exited",
	}
}
