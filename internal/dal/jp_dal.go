package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type JPDAL interface {
	// Create a job position for specified user and return its id
	CreateJP(jp *m.JobPosotion) (m.ID, error)
	// Get all job positions of specified user
	GetJPsByUserID(user m.ID) ([]m.JobPosotion, error)
}

type psqlJPDAL struct {
	db     *db.PSQLDB
	logger *l.Logger
}

func newpsqlJPDAL(db *db.PSQLDB, logger *l.Logger) *psqlJPDAL {
	return &psqlJPDAL{db, logger}
}

func (dal *psqlJPDAL) CreateJP(jp *m.JobPosotion) (m.ID, error) {
	return 0, nil
}

func (dal *psqlJPDAL) GetJPsByUserID(user m.ID) ([]m.JobPosotion, error) {
	return []m.JobPosotion{}, nil
}
