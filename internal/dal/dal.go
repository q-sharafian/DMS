package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
)

// DAL is a data access layer interface
type DAL struct {
	User  UserDAL
	Doc   DocDAL
	Event EventDAL
	JP    JPDAL
}

// Connect to the database and implement DAL for PostgreSQL. The first argument is connection details of psql database.
func NewPostgresDAL(ConnDetails db.PsqlConnDetails, logger l.Logger) DAL {
	db := db.NewPsqlConn(&ConnDetails, &logger)
	return DAL{
		User:  newPsqlUserDAL(&db, &logger),
		Doc:   newPsqlDocDAL(&db, &logger),
		Event: newPsqlEventDAL(&db, &logger),
		JP:    newpsqlJPDAL(&db, &logger),
	}
}
