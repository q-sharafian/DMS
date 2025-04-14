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

type UserDAL interface {
	// If the user successfully created, return created user's id.
	// If createdByID be nil, the user will be created as admin.
	CreateUser(name string, phoneNumber m.PhoneNumber, createdByID *m.ID) (*m.ID, error)
	GetUserByID(id m.ID) (*m.User, error)
	// If both user and error be empty, means there's not any matched user.
	GetUserByPhone(phoneNumber m.PhoneNumber) (*m.User, error)
	// Returns true if the user is disabled.
	IsDisabledByID(id m.ID) (bool, error)
	IsExistUserByPhone(phoneNumber string) (bool, error)
	// Returns true if the user exists with the given job position.
	// Returns false if the user or job doesn't exist or one of them is deleted.
	// (whether hard or soft delete)
	IsExistsUserWithJP(userID, jpID m.ID) (bool, error)
}

// It's an implementation of UserDAL interface
type psqlUserDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlUserDAL(db *db.PSQLDB, logger l.Logger) *psqlUserDAL {
	return &psqlUserDAL{db, logger}
}

func (d *psqlUserDAL) CreateUser(name string, phoneNumber m.PhoneNumber, createdByID *m.ID) (*m.ID, error) {
	user := db.User{
		Name:        name,
		PhoneNumber: phoneNumber.ToString(),
		IsDisabled:  db.IsNotDisabled,
		CreatedByID: modelID2DBID(createdByID),
	}
	d.logger.Debugf("Trying to create user with details: %+v", user)
	result := d.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, e.NewSError("couldn't create user")
	}
	return dbID2ModelID(&user.ID), nil
}

func (d *psqlUserDAL) GetUserByID(id m.ID) (*m.User, error) {
	var user db.User
	result := d.db.Where(&db.User{BaseModel: db.BaseModel{ID: *modelID2DBID(&id)}}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, fmt.Errorf("there's not any matched user with id %s", id.String())
	}
	return dbUser2ModelUser(&user), nil
}

func (d *psqlUserDAL) GetUserByPhone(phoneNumber m.PhoneNumber) (*m.User, error) {
	var user db.User
	result := d.db.Where(&db.User{PhoneNumber: phoneNumber.ToString()}).Limit(1).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		return nil, nil
	}
	return dbUser2ModelUser(&user), nil
}

func (d *psqlUserDAL) IsDisabledByID(id m.ID) (bool, error) {
	var user db.User
	result := d.db.Where(&db.User{BaseModel: db.BaseModel{ID: *modelID2DBID(&id)}}).First(&user)
	if result.Error != nil {
		return false, result.Error
	}
	d.logger.Debugf("Number of rows found during checking if the user is disabled: %d\nfetched user: %+v", result.RowsAffected, user)
	if result.RowsAffected < 1 {
		return false, fmt.Errorf("there's not any matched user with id %s", id.String())
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

func (p *psqlUserDAL) IsExistsUserWithJP(userID, jpID m.ID) (bool, error) {
	var dest any
	result := p.db.Joins("INNER JOIN users ON users.id = job_positions.user_id").
		Where("users.id = ? AND job_positions.id = ? AND users.deleted_at IS NULL AND job_positions.deleted_at IS NULL",
			userID, jpID).Find(&dest)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected < 1 {
		return false, nil
	} else if result.Error != nil {
		return false, fmt.Errorf("database error in check Existant of user and JP. %s", result.Error)
	}
	return true, nil
}

func dbUser2ModelUser(user *db.User) *m.User {
	return &m.User{
		ID:          *dbID2ModelID(&user.ID),
		Name:        user.Name,
		PhoneNumber: dbPhone2ModelPhone(user.PhoneNumber),
		IsDisabled:  dbDisability2ModelDisability(user.IsDisabled),
		CreatedBy:   dbID2ModelID(user.CreatedByID),
	}
}

func modelUser2DBUser(user *m.User) *db.User {
	return &db.User{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber.ToString(),
		IsDisabled:  modelDisability2DBDisability(user.IsDisabled),
		CreatedByID: modelID2DBID(user.CreatedBy),
	}
}
