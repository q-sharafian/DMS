package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Use this across the models as the ID type.
type ID uuid.UUID

var NilID = ID(uuid.Nil)

func (i ID) ToString() string {
	return uuid.UUID(i).String()
}

// If the ID is nil, return an empty string
func (i *ID) ToStringP() string {
	if i == nil {
		return ""
	}
	return uuid.UUID(*i).String()
}
func (i *ID) FromUUID(id uuid.UUID) {
	*i = ID(id)
}

func (i *ID) FromString(id string) error {
	if id == "''" || id == "\"\"" {
		*i = NilID
		return nil
	}
	uuid, err := uuid.Parse(id)
	if err == nil {
		*i = ID(uuid)
		return nil
	}
	// fmt.Printf("err string: '%s' %d\n\n", id, len(id))
	*i = NilID
	return err
}

// If the input id be empty, return nil id and error
func (ID) FromString2(id string) (ID, error) {
	if id == "" || id == "''" || id == "\"\"" {
		return NilID, nil
	}
	uuid, err := uuid.Parse(id)
	if err == nil {
		resultID := ID(uuid)
		return resultID, nil
	}
	return NilID, err
}

// Return true if the ID is nil. (means ID is zero value)
func (i ID) IsNil() bool {
	return uuid.UUID(i) == uuid.Nil
}

// Convert input value to the ID data type.
func (i *ID) UnmarshalJSON(data []byte) error {
	str := string(data[:])
	return i.FromString(str)
}

func ValidateUUIDv4(fl validator.FieldLevel) bool {
	if str, ok := fl.Field().Interface().(string); ok {
		_, err := ID{}.FromString2(str)
		return err == nil
	} else if _, ok := fl.Field().Interface().(ID); ok {
		return true
	} else if _, ok := fl.Field().Interface().(*ID); ok {
		return true
	} else if _, ok := fl.Field().Interface().(uuid.UUID); ok {
		return true
	}
	return false
}

type Disability int8

const (
	IsNotDisabled Disability = 0
	IsDisabled    Disability = 1
)
