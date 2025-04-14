package dal

import (
	"DMS/internal/db"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"encoding/json"
	"fmt"
)

// If input be nil, return nil
func dbID2ModelID(id *db.ID) *m.ID {
	if id == nil {
		return &m.NilID
	}
	a := m.ID(*id)
	return &a
}

// If input be nil, return nil
func modelID2DBID(id *m.ID) *db.ID {
	if id == nil || id.IsNil() {
		return nil
	}
	a := db.ID(*id)
	return &a
}

func dbIDs2ModelIDs(ids *[]db.ID) *[]m.ID {
	if ids == nil {
		return nil
	}
	var res []m.ID
	for _, id := range *ids {
		res = append(res, *dbID2ModelID(&id))
	}
	return &res
}
func modelIDs2DBIDs(ids *[]m.ID) *[]db.ID {
	if ids == nil {
		return nil
	}
	var res []db.ID
	for _, id := range *ids {
		res = append(res, *modelID2DBID(&id))
	}
	return &res
}

func dbDisability2ModelDisability(userStatus db.Disability) m.Disability {
	if userStatus == db.IsDisabled {
		return m.IsDisabled
	} else if userStatus == db.IsNotDisabled {
		return m.IsNotDisabled
	}
	panic(fmt.Sprintf("unknown user status: %d", userStatus))
}

func modelDisability2DBDisability(userStatus m.Disability) db.Disability {
	if userStatus == m.IsDisabled {
		return db.IsDisabled
	} else if userStatus == m.IsNotDisabled {
		return db.IsNotDisabled
	}
	panic(fmt.Sprintf("unknown user status: %d", userStatus))
}

func dbPhone2ModelPhone(phone string) m.PhoneNumber {
	return m.PhoneNumber(phone)
}

// DAL is a data access layer interface
type DAL struct {
	User       UserDAL
	Doc        DocDAL
	Event      EventDAL
	JP         JPDAL
	Permission PermissionDAL
	Session    SessionDAL
}

// Connect to the database and implement DAL for PostgreSQL. The first argument is
// connection details of psql database.
// If autoMigrate be true, run auto migration schema to database
func NewPostgresDAL(ConnDetails db.PsqlConnDetails, cache InMemoryDAL, logger l.Logger, autoMigrate bool) DAL {
	c := initCache(cache)
	db := db.NewPsqlConn(&ConnDetails, autoMigrate, logger)
	return DAL{
		User:       newPsqlUserDAL(&db, logger),
		Doc:        newPsqlDocDAL(&db, c, logger),
		Event:      newPsqlEventDAL(&db, c, logger),
		JP:         newPsqlJPDAL(&db, c, logger),
		Permission: newPsqlPermissionDAL(&db, logger),
		Session:    newPsqlSessionDAL(&db, logger),
	}
}

type cache struct {
	cache InMemoryDAL
}

// Create a new cache
func initCache(inMemoeyCache InMemoryDAL) *cache {
	return &cache{inMemoeyCache}
}

// If such key doesn't exists, return "e.ErrNotFound" error
func (c *cache) read(key string, dest any) error {
	value, err := c.cache.Get(key)
	if err != nil {
		return err
	} else if value == nil {
		return e.ErrNotFound
	}
	err = json.Unmarshal([]byte(*value), dest)
	if err != nil {
		return fmt.Errorf("can't unmarshal the stored value in the cache: %s", err.Error())
	}
	return nil
}

func (c *cache) write(key string, value any) error {
	stringVal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err.Error())
	}
	err = c.cache.Set(key, string(stringVal))
	if err != nil {
		return fmt.Errorf("faled to set key \"%s\" in the cache: %s", key, err.Error())
	}
	return nil
}
