package dal

import (
	"DMS/internal/db"
	m "DMS/internal/models"
)

type JPDAL interface {
	// Create a job position for specified user
	CreateJP(jp *m.JobPosotion) (m.ID, error)
	// Get all job positions of specified user
	GetUserJPs(user m.ID) (m.JobPosotion, error)
}

type psqlJPDAL struct {
	db *db.PSQLDB
}

func newpsqlJPDAL(db *db.PSQLDB) *psqlJPDAL {
	return &psqlJPDAL{db}
}

func (dal *psqlJPDAL) CreateJP(jp *m.JobPosotion) (m.ID, error) {
	return 0, nil
}

func (dal *psqlJPDAL) GetUserJPs(user m.ID) (m.JobPosotion, error) {
	return m.JobPosotion{}, nil
}
