package dal

import (
	"DMS/internal/db"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type JPDAL interface {
	// Create a job position and its permissions for specified user and return job position id
	CreateJPWithPermissions(jp *m.JobPosotion, permission *m.Permission) (*m.ID, error)
	// Create a job position for specified user id and return its id
	CreateJP(jp *m.JobPosotion) (*m.ID, error)
	// Create Permission for specified job position and return its id
	CreatePermission(JPID m.ID, permission *m.Permission) (*m.ID, error)
	// Get all job positions of specified user
	GetJPsByUserID(user m.ID) ([]m.JobPosotion, error)
	GetAllJPCount() (uint64, error)
}

type psqlJPDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlJPDAL(db *db.PSQLDB, logger l.Logger) *psqlJPDAL {
	return &psqlJPDAL{db, logger}
}

func (d *psqlJPDAL) CreateJP(jp *m.JobPosotion) (*m.ID, error) {
	newJP := db.JobPosition{
		UserID:   *modelID2DBID(&jp.UserID),
		Title:    jp.Title,
		RegionID: *modelID2DBID(&jp.RegionID),
		ParentID: modelID2DBID(jp.ParentID),
	}
	result := d.db.Create(&newJP)

	if result.Error != nil {
		d.logger.Debugf("Failed to create job position for user-id %s (%s)", newJP.UserID.ToInt64(), result.Error.Error())
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		d.logger.Debugf(`It seems can't create job position for user-id %s. Total rows 
        created are %d"`, newJP.UserID.ToInt64(), result.RowsAffected)
		return nil, e.NewSError("couldn't create job position")
	}
	return dbID2ModelID(&newJP.ID), nil
}

func (d *psqlJPDAL) CreatePermission(JPID m.ID, permission *m.Permission) (*m.ID, error) {
	newPermission := db.JPPermission{
		JPID:            *modelID2DBID(&JPID),
		IsAllowCreateJP: permission.IsAllowCreateJP,
	}
	result := d.db.Create(&newPermission)

	if result.Error != nil {
		d.logger.Debugf("Failed to create permission for job position-id %s (%s)", newPermission.JPID.ToInt64(), result.Error.Error())
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		d.logger.Debugf(`It seems can't create permission for job position-id %s. Total rows 
        created are %d"`, newPermission.JPID.ToInt64(), result.RowsAffected)
		return nil, e.NewSError("couldn't create permission")
	}
	return dbID2ModelID(&newPermission.ID), nil
}

func (d *psqlJPDAL) CreateJPWithPermissions(jp *m.JobPosotion, permission *m.Permission) (*m.ID, error) {
	var jpID *m.ID
	result := d.db.Transaction(func(tx *db.PSQLDB) error {
		var err error
		jpID, err = d.CreateJP(jp)
		if err != nil {
			return err
		}

		_, err = d.CreatePermission(*jpID, permission)
		if err != nil {
			return err
		}
		return nil
	})

	if result != nil {
		d.logger.Debugf("Failed to run transaction to create job position with permissions (%s)", result.Error())
		return nil, result
	}
	return jpID, nil
}

func (d *psqlJPDAL) GetAllJPCount() (uint64, error) {
	var count int64
	result := d.db.Model(&db.JobPosition{}).Count(&count)
	if result.Error != nil {
		d.logger.Debugf("Failed to get all of job position count (%s)", result.Error.Error())
		return 0, result.Error
	}
	ucount := uint64(count)
	return ucount, nil
}

func (d *psqlJPDAL) GetJPsByUserID(user m.ID) ([]m.JobPosotion, error) {
	return []m.JobPosotion{}, nil
}
