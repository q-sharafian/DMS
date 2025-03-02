package repository

import (
	"DMS/internal/dal"
	"DMS/internal/models"
)

type JPRepo struct {
	jpDAL dal.JPDAL
}

func newJPRepo(jpDAL dal.JPDAL) JPRepo {
	return JPRepo{
		jpDAL,
	}
}

// List of some permission the job position could have
type JobPosition struct {
	// Is the user allowed to create a job position as child of himself
	IsAllowCreateJP bool `json:"is_allow_create_jp"`
}

// Create a job position with its permission for a user and return its ID.
func (r *JPRepo) CreateJP(jp *JobPosition) (models.ID, error) {

}
