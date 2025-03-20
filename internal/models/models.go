package models

import (
	"github.com/google/uuid"
)

// Use this across the models as the ID type.
type ID uuid.UUID

func (i *ID) ToString() string {
	return uuid.UUID(*i).String()
}
func (i ID) FromUUID(id uuid.UUID) ID {
	return ID(id)
}
func (i ID) FromString(id string) (ID, error) {
	uuid, err := uuid.Parse(id)
	return ID(uuid), err
}

type Disability int8

const (
	IsNotDisabled Disability = 0
	IsDisabled    Disability = 1
)
