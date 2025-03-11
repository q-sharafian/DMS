package dal

import (
	"DMS/internal/db"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type UserDAL interface {
	// If the user successfully created, return created user's id.
	// If createdByID be nil, the user will be created as admin.
	CreateUser(name, phoneNumber string, createdByID *m.ID) (m.ID, error)
	GetUserByID(id m.ID) (*m.User, error)
	// Returns true if the user is disabled.
	IsDisabledByID(id m.ID) (bool, error)
	IsExistUserByPhone(phoneNumber string) (bool, error)
}

// It's an implementation of UserDAL interface
type psqlUserDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlUserDAL(db *db.PSQLDB, logger l.Logger) *psqlUserDAL {
	return &psqlUserDAL{db, logger}
}

func (d *psqlUserDAL) CreateUser(name, phoneNumber string, createdByID *m.ID) (m.ID, error) {
	user := db.User{
		Name:        name,
		PhoneNumber: phoneNumber,
		IsDisabled:  db.IsNotDisabled,
		CreatedByID: modelID2DBID(createdByID),
	}
	d.logger.Debugf("Trying to create user with details: %+v", user)
	result := d.db.Create(&user)
	if result.Error != nil {
		return m.NilID, result.Error
	}
	if result.RowsAffected < 1 {
		return m.NilID, e.NewSError("couldn't create user")
	}
	return *dbID2ModelID(&user.ID), nil
}

func (d *psqlUserDAL) GetUserByID(id m.ID) (*m.User, error) {
	return nil, nil
}

// TODO: Complete it
func (d *psqlUserDAL) IsDisabledByID(id m.ID) (bool, error) {
	var user db.User
	result := d.db.Preload("CreatedBy").Where(&db.User{BaseModel: db.BaseModel{ID: *modelID2DBID(&id)}}).First(&user)
	if result.Error != nil {
		return false, result.Error
	}
	d.logger.Debugf("Number of rows found during checking if the user is disabled: %d\nfetched user: %+v", result.RowsAffected, user)
	if result.RowsAffected < 1 {
		return false, e.NewSError("there's not any matched user")
	}
	return user.IsDisabled == db.IsDisabled, nil
}

func (d *psqlUserDAL) IsExistUserByPhone(phoneNumber string) (bool, error) {
	var count int64
	result := d.db.Model(&db.User{}).Where(&db.User{PhoneNumber: phoneNumber}).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	d.logger.Debugf("Number of retrieved users by phone number: %d", count)
	if count >= 1 {
		return true, nil
	}
	return false, nil
}
