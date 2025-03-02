package dal

import (
	"DMS/internal/db"
	m "DMS/internal/models"
)

type UserDAL interface {
	CreateUser(name, phoneNumber string) error
	GetUserByID(id m.ID) (*m.User, error)
	// Returns true if the user is disabled.
	IsDisabledByID(id m.ID) (bool, error)
}

// It's an implementation of UserDAL interface
type psqlUserDAL struct {
	db *db.PSQLDB
}

func newPsqlUserDAL(db *db.PSQLDB) *psqlUserDAL {
	return &psqlUserDAL{db}
}

func (d *psqlUserDAL) CreateUser(name, phoneNumber string) error {
	return nil
}

func (d *psqlUserDAL) GetUserByID(id int) (*m.User, error) {
	return nil, nil
}

// TODO: Complete it
func (d *psqlUserDAL) IsDisabledByID(id m.ID) (bool, error) {
	return false, nil
}
