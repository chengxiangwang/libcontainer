package container

import "time"

type Container struct {
	ID            *UUID
	Name          string
	TTY           bool
	Root          string
	Path          string
	ReadOnlyLayer string
	WriteLayer    string
	Args          []string
	Limits        map[string]string
	Volumns       []string
	Status        string
	ImageID       string
	Created       time.Time
}
