package container

import (
	"github.com/satori/go.uuid"
)

type UUID struct {
	Original uuid.UUID
}

func NewUUID() *UUID {
	uuid, _ := uuid.NewV4()
	return &UUID{
		Original: uuid,
	}
}

func (u *UUID) String() string {
	return u.Original.String()
}
