package dal

import "DMS/internal/db"

// Representation of user entity
type user struct {
	ID          ID
	Name        string
	PhoneNumber string
	// value 0 means it's not disabled and value 1 means it's disabled.
	IsDisabled  uint8
	CreatedByID ID
}

type UserDAL interface {
	CreateUser(name, phoneNumber string) error
	GetUserByID(id int) (*user, error)
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

func (d *psqlUserDAL) GetUserByID(id int) (*user, error) {
	return nil, nil
}
