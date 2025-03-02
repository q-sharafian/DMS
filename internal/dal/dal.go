package dal

import (
	"DMS/internal/db"
	"strconv"
)

// ID is used as the ID field in database tables
type ID int64

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

// DAL is a data access layer interface
type DAL struct {
	User  UserDAL
	Doc   DocDAL
	Event EventDAL
}

// Implements DAL for PostgreSQL. The first argument is connection details of psql database.
func NewPostgresDAL(ConnDetails db.PsqlConnDetails) DAL {
	db := db.NewPsqlConn(&ConnDetails)
	return DAL{
		User:  newPsqlUserDAL(&db),
		Doc:   newPsqlDocDAL(&db),
		Event: newPsqlEventDAL(&db),
	}
}
