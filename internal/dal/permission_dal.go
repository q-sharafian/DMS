package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"fmt"
)

type PermissionDAL interface {
	// Get list of permissions of the given job position.
	// If both returned errror and permission be nil, means there's not any matched job position.
	//
	// Possible error codes
	// SEDBError
	GetPermissionsByJPID(jpID m.ID) (*m.Permission, error)
	// Create a permission for specified job position.
	//
	// Possible error codes:
	// SEDBError
	CreateJPPermission(*m.Permission) error
}

type psqlPermissionDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlPermissionDAL(db *db.PSQLDB, logger l.Logger) *psqlPermissionDAL {
	return &psqlPermissionDAL{db, logger}
}

func (d *psqlPermissionDAL) GetPermissionsByJPID(jpID m.ID) (*m.Permission, error) {
	var permission *db.JPPermission
	result := d.db.Where(&db.JPPermission{BaseModel: db.BaseModel{ID: *modelID2DBID(&jpID)}}).
		Limit(1).First(permission)

	if result.Error != nil {
		d.logger.Debugf("Failed to get permission for job position-id %s: %s", jpID.ToString(), result.Error.Error())
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		d.logger.Warnf(`It seems can't get permission for job position-id %s.`, jpID.ToString())
		return nil, nil
	}
	return dbPermission2Model(permission), nil
}

func (d *psqlPermissionDAL) CreateJPPermission(permission *m.Permission) error {
	dbPermission := modelPermission2DB(permission)
	result := d.db.Create(dbPermission)
	if result.Error != nil {
		err := fmt.Errorf("failed to create permission for job position-id %s: %s", permission.JPID.ToString(), result.Error.Error())
		return err
	}
	return nil
}

func dbPermission2Model(permission *db.JPPermission) *m.Permission {
	return &m.Permission{
		JPID:            *dbID2ModelID(&permission.JpID),
		IsAllowCreateJP: permission.IsAllowCreateJP,
	}
}

func modelPermission2DB(permission *m.Permission) *db.JPPermission {
	return &db.JPPermission{
		JpID:            *modelID2DBID(&permission.JPID),
		IsAllowCreateJP: permission.IsAllowCreateJP,
	}
}
