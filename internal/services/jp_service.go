package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"fmt"
)

type JPService interface {
	// Return all job positions the user have
	//
	// Possible error codes the function could returns:
	// SEDBError
	GetUserJPs(userID m.ID) ([]m.JobPosotion, *e.Error)
	// Create job position for the given user and details then, reutrn its id.
	//
	// Possible error codes the function could returns:
	// SEDBError
	CreateJP(jp *m.JobPosotion, permissions *m.Permission) (*m.ID, *e.Error)
}

// It's a simple implementation of JPService interface.
// This implementation has minimum functionalities.
type sJPService struct {
	jp     dal.JPDAL
	logger l.Logger
}

func (s *sJPService) GetUserJPs(userID m.ID) ([]m.JobPosotion, *e.Error) {
	jps, err := s.jp.GetJPsByUserID(userID)
	if err != nil {
		return []m.JobPosotion{}, e.NewErrorP(err.Error(), SEDBError)
	}
	return jps, nil
}

// Note that in this implementation, createdTime value doesn't matter and createdTime
// is always the current time.
func (s *sJPService) CreateJP(jp *m.JobPosotion, permissions *m.Permission) (*m.ID, *e.Error) {
	jpID, err := s.jp.CreateJP(jp)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError).
			AppendBegin(
				fmt.Sprintf(
					"failed to create job position for userID %s",
					jp.UserID.ToString(),
				),
			)
	}
	return jpID, nil
}

// Create an instance of sJPService struct
func newSJPService(jp dal.JPDAL, logger l.Logger) JPService {
	return &sJPService{jp, logger}
}
