package dal

import (
	"DMS/internal/db"
)

// DAL is a data access layer interface
type DAL struct {
	User  UserDAL
	Doc   DocDAL
	Event EventDAL
	JP    JPDAL
}

// Implements DAL for PostgreSQL. The first argument is connection details of psql database.
func NewPostgresDAL(ConnDetails db.PsqlConnDetails) DAL {
	db := db.NewPsqlConn(&ConnDetails)
	return DAL{
		User:  newPsqlUserDAL(&db),
		Doc:   newPsqlDocDAL(&db),
		Event: newPsqlEventDAL(&db),
		JP:    newpsqlJPDAL(&db),
	}
}
