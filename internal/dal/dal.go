package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

func dbID2ModelID(id *db.ID) *m.ID {
	if id == nil {
		return nil
	}
	a := m.ID(*id)
	return &a
}

func modelID2DBID(id *m.ID) *db.ID {
	if id == nil {
		return nil
	}
	a := db.ID(id.ToInt64())
	return &a
}

func dbID2ModelIDSlice(ids *[]db.ID) *[]m.ID {
	if ids == nil {
		return nil
	}
	var res []m.ID
	for _, id := range *ids {
		res = append(res, *dbID2ModelID(&id))
	}
	return &res
}
func modelID2DBIDSlice(ids *[]m.ID) *[]db.ID {
	if ids == nil {
		return nil
	}
	var res []db.ID
	for _, id := range *ids {
		res = append(res, *modelID2DBID(&id))
	}
	return &res
}

// DAL is a data access layer interface
type DAL struct {
	User       UserDAL
	Doc        DocDAL
	Event      EventDAL
	JP         JPDAL
	Permission PermissionDAL
}

// Connect to the database and implement DAL for PostgreSQL. The first argument is
// connection details of psql database.
// If autoMigrate be true, run auto migration schema to database
func NewPostgresDAL(ConnDetails db.PsqlConnDetails, logger l.Logger, autoMigrate bool) DAL {
	db := db.NewPsqlConn(&ConnDetails, autoMigrate, logger)
	return DAL{
		User:       newPsqlUserDAL(&db, logger),
		Doc:        newPsqlDocDAL(&db, logger),
		Event:      newPsqlEventDAL(&db, logger),
		JP:         newPsqlJPDAL(&db, logger),
		Permission: newPsqlPermissionDAL(&db, logger),
	}
}
