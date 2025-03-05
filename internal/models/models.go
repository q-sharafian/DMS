package models

import (
	"strconv"
)

// Use this across the models as the ID type.
type ID int64

// Null value of ID
const NilID ID = -1

// TODO: What happend if the length of id was greater than 8 bytes?
func (i *ID) ToInt64() int64 {
	return int64(*i)
}
func (i *ID) ToString() string {
	return strconv.FormatInt(int64(*i), 10)
}
func (i *ID) FromInt64(id int64) ID {
	return ID(id)
}
