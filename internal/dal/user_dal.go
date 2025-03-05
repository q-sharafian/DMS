package dal

import (
	"DMS/internal/db"
	m "DMS/internal/models"
)

type UserDAL interface {
	// If the user successfully created, return created user's id
	CreateUser(name, phoneNumber string, createdByID m.ID) (m.ID, error)
	GetUserByID(id m.ID) (*m.User, error)
	// Returns true if the user is disabled.
	IsDisabledByID(id m.ID) (bool, error)
	IsExistUserByPhone(phoneNumber string) (bool, error)
}

// It's an implementation of UserDAL interface
type psqlUserDAL struct {
	db *db.PSQLDB
}

func newPsqlUserDAL(db *db.PSQLDB) *psqlUserDAL {
	return &psqlUserDAL{db}
}

func (d *psqlUserDAL) CreateUser(name, phoneNumber string, id m.ID) (m.ID, error) {
	return m.NilID, nil
}

func (d *psqlUserDAL) GetUserByID(id m.ID) (*m.User, error) {
	return nil, nil
}

// TODO: Complete it
func (d *psqlUserDAL) IsDisabledByID(id m.ID) (bool, error) {
	return false, nil
}

func (d *psqlUserDAL) IsExistUserByPhone(phoneNumber string) (bool, error) {
	return true, nil
}
