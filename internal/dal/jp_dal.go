package dal

import (
	"DMS/internal/db"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type JPDAL interface {
	// Create a job position and its permissions for specified user and return job position id
	CreateUserJPWithPermissions(jp *m.UserJobPosition, permission *m.Permission) (*m.ID, error)
	// Create an admin job position and its permissions for specified user and return job position id
	CreateAdminJPWithPermissions(jp *m.AdminJobPosition, permission *m.Permission) (*m.ID, error)
	// Create a job position for specified user id and return its id
	CreateUserJP(jp *m.UserJobPosition) (*m.ID, error)
	// Create a admin job position for specified user id and return its id
	CreateAdminJP(jp *m.AdminJobPosition) (*m.ID, error)
	// Create Permission for specified job position and return its id
	CreatePermission(JPID m.ID, permission *m.Permission) (*m.ID, error)
	// Get all job positions of the specified user
	// If both array and error be nil, it means there's not any matched job position.
	GetJPsByUser(user *m.User) (*[]m.UserJobPosition, error)
	// Return true if a job position with given ID belongs to a user with given ID.
	IsExistsUserWithJP(userID, jpID m.ID) (bool, error)
	GetAllJPCount() (uint64, error)
}

type psqlJPDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlJPDAL(db *db.PSQLDB, logger l.Logger) *psqlJPDAL {
	return &psqlJPDAL{db, logger}
}

func (d *psqlJPDAL) CreateUserJP(jp *m.UserJobPosition) (*m.ID, error) {
	newJP := db.JobPosition{
		UserID:   *modelID2DBID(&jp.UserID),
		Title:    jp.Title,
		RegionID: *modelID2DBID(&jp.RegionID),
		ParentID: modelID2DBID(&jp.ParentID),
	}
	result := d.db.Create(&newJP)

	if result.Error != nil {
		d.logger.Debugf("Failed to create job position for user-id %s (%s)", newJP.UserID.ToString(), result.Error.Error())
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		d.logger.Debugf(`It seems can't create job position for user-id %s. Total rows 
        created are %d"`, newJP.UserID.ToString(), result.RowsAffected)
		return nil, e.NewSError("couldn't create job position")
	}
	return dbID2ModelID(&newJP.ID), nil
}

func (d *psqlJPDAL) CreateAdminJP(jp *m.AdminJobPosition) (*m.ID, error) {
	newJP := db.JobPosition{
		UserID:   *modelID2DBID(&jp.UserID),
		Title:    jp.Title,
		RegionID: *modelID2DBID(&jp.RegionID),
		ParentID: nil,
	}
	result := d.db.Create(&newJP)

	if result.Error != nil {
		d.logger.Debugf("Failed to create job position for user-id %s (%s)", newJP.UserID.ToString(), result.Error.Error())
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		d.logger.Debugf(`It seems can't create job position for user-id %s. Total rows 
        created are %d"`, newJP.UserID.ToString(), result.RowsAffected)
		return nil, e.NewSError("couldn't create job position")
	}
	return dbID2ModelID(&newJP.ID), nil
}

func (d *psqlJPDAL) CreatePermission(JPID m.ID, permission *m.Permission) (*m.ID, error) {
	newPermission := db.JPPermission{
		JpID:            *modelID2DBID(&JPID),
		IsAllowCreateJP: permission.IsAllowCreateJP,
	}
	result := d.db.Create(&newPermission)

	if result.Error != nil {
		d.logger.Debugf("Failed to create permission for job position-id %s (%s)", newPermission.JpID.ToString(), result.Error.Error())
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		d.logger.Debugf(`It seems can't create permission for job position-id %s. Total rows 
        created are %d"`, newPermission.JpID.ToString(), result.RowsAffected)
		return nil, e.NewSError("couldn't create permission")
	}
	return dbID2ModelID(&newPermission.ID), nil
}

func (d *psqlJPDAL) CreateUserJPWithPermissions(jp *m.UserJobPosition, permission *m.Permission) (*m.ID, error) {
	var jpID *m.ID
	result := d.db.Transaction(func(tx *db.PSQLDB) error {
		var err error
		jpID, err = d.CreateUserJP(jp)
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

func (d *psqlJPDAL) CreateAdminJPWithPermissions(jp *m.AdminJobPosition, permission *m.Permission) (*m.ID, error) {
	var jpID *m.ID
	result := d.db.Transaction(func(tx *db.PSQLDB) error {
		var err error
		jpID, err = d.CreateAdminJP(jp)
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

func (d *psqlJPDAL) GetJPsByUser(user *m.User) (*[]m.UserJobPosition, error) {
	var jps []db.JobPosition
	dbUser := modelUser2DBUser(user)
	// result := d.db.Select("users.*", "job_positions.*").Joins("INNER JOIN job_positions ON job_positions.user_id = users.id").Where(dbUser).Find(&jps)
	result := d.db.Select("users.*").Where(dbUser).Find(&db.User{}).Select("job_positions.*").Joins("INNER JOIN job_positions ON users.id = job_positions.user_id").Find(&jps)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected < 1 {
		return nil, nil
	} else if result.Error != nil {
		return nil, fmt.Errorf("failed to get job positions of user %+v (%s)", dbUser, result.Error.Error())
	}

	modelJPs := dbJPs2ModelJPs(jps)
	return &modelJPs, nil
}

func (d *psqlJPDAL) IsExistsUserWithJP(userID, jpID m.ID) (bool, error) {
	var jp db.JobPosition
	result := d.db.Where("user_id = ? AND id = ?", userID, jpID).Limit(1).Find(&jp)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check if user with id %s has job position with id %s: %s",
			userID.ToString(), jpID.ToString(), result.Error.Error())
	} else if result.RowsAffected < 1 {
		return false, nil
	} else {
		return true, nil
	}
}

func dbJP2ModelJP(jp *db.JobPosition) *m.UserJobPosition {
	var mParentID m.ID
	dParentID := dbID2ModelID(jp.ParentID)
	if dParentID != nil {
		mParentID = *dParentID
	} else {
		mParentID = m.NilID
	}

	return &m.UserJobPosition{
		CommonJobPosition: m.CommonJobPosition{
			ID:        *dbID2ModelID(&jp.ID),
			UserID:    *dbID2ModelID(&jp.UserID),
			Title:     jp.Title,
			RegionID:  *dbID2ModelID(&jp.RegionID),
			CreatedAt: jp.CreatedAt.UTC().Unix(),
		},
		ParentID: mParentID,
	}
}

func dbJPs2ModelJPs(jps []db.JobPosition) []m.UserJobPosition {
	var modelJPs []m.UserJobPosition
	for _, jp := range jps {
		modelJPs = append(modelJPs, *dbJP2ModelJP(&jp))
	}
	return modelJPs
}
